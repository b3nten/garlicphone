package parser

import (
	"fmt"
	"hash/fnv"
	"os"

	lua "github.com/yuin/gopher-lua"
)

type Type interface {
	TypeKind() string
}

type PrimitiveType struct {
	Name string
}

func (p PrimitiveType) TypeKind() string { return "primitive" }

type Field struct {
	Name string
	ID   uint16
	Type Type
}

type StructType struct {
	ID     uint16
	UUID   string
	Name   string
	Fields []Field
}

func (s StructType) TypeKind() string { return "struct" }

type ListType struct {
	ElementType Type
}

func (l ListType) TypeKind() string { return "list" }

type Schema struct {
	Name string
	Version int
	Structs []StructType
}

func GenerateSchema(file string) (*Schema, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return &Schema{}, err
	}

	L := lua.NewState()
	defer L.Close()

	if err := L.DoString(string(bytes)); err != nil {
		panic(err)
	}

	lstructs := []*lua.LTable{}
	structs := map[string]*StructType{}
	ltos := map[*lua.LTable]*StructType{}
	if tbl, ok := L.GetGlobal("_G").(*lua.LTable); ok {
		// Iterate over the global table to find structs defined in the global scope
		tbl.ForEach(func(key lua.LValue, value lua.LValue) {
			if structTbl, ok := value.(*lua.LTable); ok {
				if structTbl.RawGet(lua.LString("type")) == lua.LString("struct") {
					sv := StructType{
						Name: key.String(),
						ID:   uint16(fnv32a(key.String())),
						UUID: structTbl.RawGet(lua.LString("uuid")).String(),
					}
					structs[sv.UUID] = &sv
					lstructs = append(lstructs, structTbl)
					ltos[structTbl] = &sv
				}
			}
		})
		// gather all fields in structs
		for _, lstruct := range lstructs {
			if fields, ok := lstruct.RawGet(lua.LString("fields")).(*lua.LTable); ok {
				i := 10
				fields.ForEach(func(key lua.LValue, value lua.LValue) {
					if fieldTbl, ok := value.(*lua.LTable); ok {
						sv := ltos[lstruct]
						typ := mapType(fieldTbl, structs)
						sv.Fields = append(sv.Fields, Field{
							ID:   uint16(i),
							Name: key.String(),
							Type: typ,
						})
						i++
					}
				})
			}
		}
	}

	structList := []StructType{}
	for _, sl := range structs {
		structList = append(structList, *sl)
	}

	// get file name without extension or base path
	fileName := file
	if stat, err := os.Stat(file); err == nil {
		fileName = stat.Name()
	}
	extIndex := -1
	for i := len(fileName) - 1; i >= 0; i-- {
		if fileName[i] == '.' {
			extIndex = i
			break
		}
	}
	if extIndex != -1 {
		fileName = fileName[:extIndex]
	}

	return &Schema{fileName, 1, structList}, nil
}

func generateTypeTable(L *lua.LState, lt Type) *lua.LTable {
	tbl := L.NewTable()
	switch lt.TypeKind() {
	case "primitive":
		pt := lt.(PrimitiveType)
		tbl.RawSet(lua.LString("kind"), lua.LString("primitive"))
		tbl.RawSet(lua.LString("name"), lua.LString(pt.Name))
	case "struct":
		st := lt.(*StructType)
		tbl.RawSet(lua.LString("kind"), lua.LString("struct"))
		tbl.RawSet(lua.LString("name"), lua.LString(st.Name))
		tbl.RawSet(lua.LString("uuid"), lua.LString(st.UUID))
		tbl.RawSet(lua.LString("id"), lua.LNumber(st.ID))
	case "list":
		lt := lt.(ListType)
		tbl.RawSet(lua.LString("kind"), lua.LString("list"))
		tbl.RawSet(lua.LString("of"), generateTypeTable(L, lt.ElementType))
	default:
		panic(fmt.Sprintf("unknown type kind: %T", lt))
	}
	return tbl
}

func CreateLuaState(s *Schema) *lua.LState {
	// create state from schema so lua files can generate code
	L := lua.NewState()
	// set schema properties
	schema := L.NewTable()
	schema.RawSet(lua.LString("name"), lua.LString(s.Name))
	schema.RawSet(lua.LString("version"), lua.LNumber(s.Version))

	structsTable := L.NewTable()

	// set structs
	for _, s := range s.Structs {
		structTable := L.NewTable()

		// set struct properties
		structTable.RawSet(lua.LString("id"), lua.LNumber(s.ID))
		structTable.RawSet(lua.LString("uuid"), lua.LString(s.UUID))
		structTable.RawSet(lua.LString("name"), lua.LString(s.Name))

		// set field properties
		fieldsTable := L.NewTable()
		for _, f := range s.Fields {
			fieldTable := L.NewTable()

			fieldTable.RawSet(lua.LString("name"), lua.LString(f.Name))
			fieldTable.RawSet(lua.LString("id"), lua.LNumber(f.ID))

			// set type properties
			typeTable := generateTypeTable(L, f.Type)
			// assign type to field
			fieldTable.RawSet(lua.LString("type"), typeTable)
			// assign field to fields table
			fieldsTable.RawSet(lua.LString(f.Name), fieldTable)
		}
		// assign fields to struct
		structTable.RawSet(lua.LString("fields"), fieldsTable)
		// assign struct to structs table
		structsTable.RawSet(lua.LString(s.Name), structTable)
	}

	schema.RawSet(lua.LString("structs"), structsTable)
	L.SetGlobal("Schema", schema)
	return L
}

func mapType(tbl *lua.LTable, structs map[string]*StructType) Type {
	switch tbl.RawGet(lua.LString("type")).String() {
	case "primitive":
		return PrimitiveType{Name: tbl.RawGet(lua.LString("name")).String()}
	case "struct":
		return structs[tbl.RawGet(lua.LString("uuid")).String()]
	case "list":
		return ListType{
			ElementType: mapType(tbl.RawGet(lua.LString("of")).(*lua.LTable), structs),
		}
	default:
		return nil
	}
}

func fnv32a(text string) uint32 {
	algorithm := fnv.New32a()
	algorithm.Write([]byte(text))
	return algorithm.Sum32()
}
