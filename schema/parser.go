package main

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"os"

	"strings"
	"unicode"

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
	Structs []StructType
}

type primitiveField struct{}

func parseFile() (Schema, error) {
	bytes, err := os.ReadFile("game.schema")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return Schema{}, err
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
						ID:   uint16(FNV32a(key.String())),
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
		fmt.Println("Struct:", sl)
		structList = append(structList, *sl)
	}
	return Schema{structList}, nil
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

func FNV32a(text string) uint32 {
	algorithm := fnv.New32a()
	algorithm.Write([]byte(text))
	return algorithm.Sum32()
}

func toPascalCase(s string) string {
	if s == "" {
		return ""
	}

	var result strings.Builder
	capitalizeNext := true

	for i, r := range s {
		if r == '_' || r == '-' || r == ' ' || r == '.' {
			// Treat these as word separators
			capitalizeNext = true
			continue
		}

		if unicode.IsUpper(r) && i > 0 {
			// If we encounter an uppercase letter that's not at the start,
			// it might be camelCase, so capitalize it
			result.WriteRune(r)
			capitalizeNext = false
		} else if capitalizeNext {
			result.WriteRune(unicode.ToUpper(r))
			capitalizeNext = false
		} else {
			result.WriteRune(unicode.ToLower(r))
		}
	}

	return result.String()
}

func ReadAfterLine(filename string, lineNum int) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var result strings.Builder
	currentLine := 0

	for scanner.Scan() {
		currentLine++
		if currentLine > lineNum {
			result.WriteString(scanner.Text())
			result.WriteString("\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	return result.String(), nil
}
