package main

import (
	"fmt"
	"os"
	"strings"
)

func printType(t Type) string {
	switch t.TypeKind() {
		case "primitive":
			pt := t.(PrimitiveType)
			return pt.Name
		case "struct":
			st := t.(*StructType)
			return toPascalCase(st.Name)
		case "list":
			lt := t.(ListType)
			return fmt.Sprintf("[]%s", printType(lt.ElementType))
		default:
			return "unknown"
	}
}

func printField(f Field) string {
	return fmt.Sprintf(
		"\n\t %s *%s",
		toPascalCase(f.Name),
		printType(f.Type),
	)
}

func printListSerializer(lt ListType) string {
	switch lt.ElementType.TypeKind() {
		case "primitive":
			pt := lt.ElementType.(PrimitiveType)
			return fmt.Sprintf("newListSerializer[%s](serialize%s)", pt.Name, capFirst(pt.Name))
		case "struct":
			st := lt.ElementType.(*StructType)
			return fmt.Sprintf("newListSerializer[%s](serializeStruct[%s])", toPascalCase(st.Name), toPascalCase(st.Name))
		case "list":
			nlt := lt.ElementType.(ListType)
			return fmt.Sprintf("newListSerializer[%s](%s)", printType(nlt), printListSerializer(nlt))
		default:
			return "unknown"
	}
}

func printStruct(sv StructType) string {
	sb := strings.Builder{}

	// print struct
	sb.WriteString(fmt.Sprintf("type %s struct { ",  toPascalCase(sv.Name)))
	for _, f := range sv.Fields {
		sb.WriteString(printField(f))
	}
	sb.WriteString("\n}\n")

	// print TypeID
	sb.WriteString(fmt.Sprintf("func (%s) TypeID() uint16 { return uint16(%d) }\n", toPascalCase(sv.Name), sv.ID))

	// print toBytes
	sb.WriteString(fmt.Sprintf("func (it %s) toBytes(data *bytes.Buffer) {\n", toPascalCase(sv.Name)))
	sb.WriteString(fmt.Sprintf("\tserializeUint16(%d, data)\n", sv.ID))
	sb.WriteString("\tstartLenPos := data.Len()\n")
	sb.WriteString("\tserializeUint32(0, data)\n")
	for i, field := range sv.Fields {
		sb.WriteString(fmt.Sprintf("\tif it.%s != nil {\n", toPascalCase(field.Name)))
		sb.WriteString(fmt.Sprintf("\t\tserializeUint16(%d, data)\n", i))

		switch field.Type.TypeKind() {
			case "primitive":
				pt := field.Type.(PrimitiveType)
				sb.WriteString(fmt.Sprintf("\t\tserialize%s(*it.%s, data)\n", capFirst(pt.Name), toPascalCase(field.Name)))
			case "struct":
				sb.WriteString(fmt.Sprintf("\t\tserializeStruct(*it.%s, data)\n", toPascalCase(field.Name)))
			case "list":
				sb.WriteString(
					fmt.Sprintf("\t\t%s(*it.%s, data)\n",
						printListSerializer(field.Type.(ListType)),
						toPascalCase(field.Name),
					),
				)
		}

		sb.WriteString("\t}\n")
	}
	sb.WriteString("\tbinary.BigEndian.PutUint32(data.Bytes()[startLenPos:], uint32(len(data.Bytes())-(startLenPos+lenSize)))")
	sb.WriteString("\n}\n")

	// print fromBytes
	sb.WriteString(fmt.Sprintf("func (it *%s) fromBytes(data []byte, fieldIndex uint16, offset int) (int, error) {\n", toPascalCase(sv.Name)))
	sb.WriteString("\tswitch fieldIndex {\n")
	for i, field := range sv.Fields {
		var name string
		switch field.Type.TypeKind() {
			case "primitive":
				pt := field.Type.(PrimitiveType)
				name = capFirst(pt.Name)
			case "struct":
				name = "Struct[" + toPascalCase(field.Type.(*StructType).Name) + "]"
			case "list":
				name = "foo"
		}
		sb.WriteString(
			fmt.Sprintf(
				"\tcase %d: val, len, err := deserialize%s(data, offset); it.%s = &val; return len, err\n",
			 	i,
				name,
				toPascalCase(field.Name),
			),
		)
	}
	sb.WriteString("\t}\n")
	sb.WriteString("\treturn 0, UnknownFieldError\n")
	sb.WriteString("}\n")

	return sb.String()
}

func printGo(s Schema) error {
	sb := strings.Builder{}
	sb.WriteString("package schematest\n\n")
	sb.WriteString("import(\n")
	sb.WriteString("\t\"bytes\"\n")
	sb.WriteString("\t\"encoding/binary\"\n")
	sb.WriteString("\t\"errors\"\n")
	sb.WriteString("\t\"fmt\"\n")
	sb.WriteString(")\n\n")

	for _, str := range s.Structs {
		sb.WriteString(printStruct(str) + "\n")
	}

	// Deserialize function
	sb.WriteString("func Deserialize[K any, KT interface {*K; Serializable}](b []byte, out KT) error {\n")
	sb.WriteString("\tif len(b) < idSize+lenSize { return fmt.Errorf(\"data too short to contain message header\") }\n")
	sb.WriteString("\ttypeID := uint16(binary.BigEndian.Uint16(b[0:idSize]))\n")
	sb.WriteString("\tif out.TypeID() != typeID { return fmt.Errorf(\"type ID mismatch: expected %d, got %d\", out.TypeID(), typeID) }\n")
	sb.WriteString("\tswitch v := any(out).(type) {\n")
	for _, str := range s.Structs {
		sb.WriteString(fmt.Sprintf("\tcase *%s: _, err := parse(v.fromBytes, b); return err\n", toPascalCase(str.Name)))
	}
	sb.WriteString("\tdefault: return fmt.Errorf(\"unsupported type for deserialization\")\n\t}\n}\n\n")

	// deserializeStruct function
	sb.WriteString("func deserializeStruct[K Serializable](data []byte, offset int) (K, int, error) {\n")
	sb.WriteString("\tvar val K\n")
	sb.WriteString("\tslice := data[offset:]\n")
	sb.WriteString("\tswitch v := any(val).(type) {\n")
	for _, str := range s.Structs {
		sb.WriteString(fmt.Sprintf("\tcase %s:\n", toPascalCase(str.Name)))
		sb.WriteString("\t\ti, err := parse(v.fromBytes, slice)\n")
		sb.WriteString("\t\treturn any(v).(K), i, err\n")
	}
	sb.WriteString("\tdefault:\n")
	sb.WriteString("\t\treturn val, 0, fmt.Errorf(\"unsupported struct type for deserialization\")\n")
	sb.WriteString("\t}\n}\n\n")


	postlude, err := ReadAfterLine("schema/helpers.go", 9)
	if err != nil {
		return err
	}
	sb.WriteString(postlude)

	os.WriteFile("gen/schema.go", []byte(sb.String()), 0644)
	return nil
}

func capFirst(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}
