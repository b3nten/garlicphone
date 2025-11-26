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

type structValue struct {
	id     uint16
	name   string
	uuid   string
	fields []fieldWriter
}

type fieldWriter interface {
	printField() string
	printToBytes() string
	printFromBytes() string
}

type primitiveField struct {
	id   uint16
	name string
	typ  string
}

func (pfv primitiveField) printField() string {
	return fmt.Sprintf(
		"\n\t %s optional.Optional[%s]",
		toPascalCase(pfv.name),
		pfv.typ,
	)
}

func (pfv primitiveField) printToBytes() string {
	return fmt.Sprintf(
		"\tif it.%s.Exists() {\n\t\tappendUint16(&data, %d)\n\t\tappend%s(&data, it.%s.MustGet())\n\t}\n",
		toPascalCase(pfv.name),
		pfv.id,
		toPascalCase(pfv.typ),
		toPascalCase(pfv.name),
	)
}

func (pfv primitiveField) printFromBytes() string {
	return fmt.Sprintf(
		"\tcase %d: return read%s(data, offset, &it.%s)\n",
		pfv.id,
		toPascalCase(pfv.typ),
		toPascalCase(pfv.name),
	)
}

func (sv structValue) printStruct() string {
	sb := strings.Builder{}

	// print struct
	sb.WriteString(fmt.Sprintf("type %s struct { ",  toPascalCase(sv.name)))
	for _, field := range sv.fields {
		sb.WriteString(field.printField())
	}
	sb.WriteString("\n}\n\n")

	// print to bytes
	sb.WriteString("func (it " + toPascalCase(sv.name) + ") ToBytes() []byte {\n")
	sb.WriteString("\tvar data bytes.Buffer\n\n")
	sb.WriteString(fmt.Sprintf("\tappendUint16(&data, %d)\n", sv.id))
	sb.WriteString("\tappendUint32(&data, 0)\n\n")
	for _, field := range sv.fields {
		sb.WriteString(field.printToBytes())
	}
	sb.WriteString("\n\tbinary.BigEndian.PutUint32(data.Bytes()[idSize:], uint32(len(data.Bytes())-(idSize+lenSize)))\n")
	sb.WriteString("\treturn data.Bytes()\n")
	sb.WriteString("}\n\n")

	// print from bytes
	sb.WriteString("func (it *" + toPascalCase(sv.name) + ") FromBytes(bytes []byte) (int, error) {\n")
	sb.WriteString("\treturn parse(it, bytes)\n")
	sb.WriteString("}\n\n")

	// print parse fields
	sb.WriteString("func (it *" + toPascalCase(sv.name) + ") parseFields(data []byte, fieldIndex uint16, offset int) (int, error) {\n")
	sb.WriteString("\tswitch fieldIndex {\n")
	for _, field := range sv.fields {
		sb.WriteString(field.printFromBytes())
	}
	sb.WriteString("\t}\n")
	sb.WriteString("\treturn 0, UnknownFieldError\n")
	sb.WriteString("}\n\n")
	return sb.String()
}

func printGo(structs []*structValue) error {
	sb := strings.Builder{}
	sb.WriteString("package schematest\n\n")
	sb.WriteString("import(\n")
	sb.WriteString("\t\"bytes\"\n")
	sb.WriteString("\t\"encoding/binary\"\n")
	sb.WriteString("\t\"6enten/garlicphone/optional\"\n")
	sb.WriteString("\t\"errors\"\n")
	sb.WriteString("\t\"fmt\"\n")
	sb.WriteString(")\n\n")

	sb.WriteString("// Type definitions\n\n")
	for _, str := range structs {
		sb.WriteString(str.printStruct() + "\n")
	}

	sb.WriteString("// Utilities\n\n")
	postlude, err := ReadAfterLine("schema/util.go", 10)
	if err != nil {
		return err
	}
	sb.WriteString(postlude)

	os.WriteFile("gen/schema.go", []byte(sb.String()), 0644)
	return nil
}

func parseFile() ([]*structValue, error) {
	bytes, err := os.ReadFile("game.schema")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}
	L := lua.NewState()
	defer L.Close()

	if err := L.DoString(string(bytes)); err != nil {
		panic(err)
	}

	lstructs := []*lua.LTable{}
	structs := []*structValue{}
	ltos := map[*lua.LTable]*structValue{}
	if tbl, ok := L.GetGlobal("_G").(*lua.LTable); ok {
		// Iterate over the global table to find structs defined in the global scope
		tbl.ForEach(func(key lua.LValue, value lua.LValue) {
			if structTbl, ok := value.(*lua.LTable); ok {
				if structTbl.RawGet(lua.LString("type")) == lua.LString("struct") {
					sv := structValue{
						name: key.String(),
						id:   uint16(FNV32a(key.String())),
						uuid: structTbl.RawGet(lua.LString("uuid")).String(),
					}
					structs = append(structs, &sv)
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
						_ = fieldTbl.RawGet(lua.LString("type")).String()
						sv := ltos[lstruct]
						sv.fields = append(sv.fields, primitiveField{
							id:   uint16(i),
							name: key.String(),
							typ:  fieldTbl.RawGet(lua.LString("name")).String(),
						})
						i++
					}
				})
			}
		}
	}
	return structs, nil
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
