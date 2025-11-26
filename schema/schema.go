package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func Deserialize[K any, KT interface {*K; Serializable}](b []byte, out KT) error {
	if len(b) < idSize+lenSize {
		return fmt.Errorf("data too short to contain message header")
	}
	typeID := uint16(binary.BigEndian.Uint16(b[0:idSize]))
	if out.TypeID() != typeID {
		return fmt.Errorf("type ID mismatch: expected %d, got %d", out.TypeID(), typeID)
	}
	switch v := any(out).(type) {
	case *MyMessage: _, err := parse(v.fromBytes, b); return err
	case *MyStruct:
		_, err := parse(v.fromBytes, b)
		return err
	default:
		return fmt.Errorf("unsupported type for deserialization")
	}
}

func deserializeStruct[K Serializable](data []byte, offset int) (K, int, error) {
	var val K
	slice := data[offset:]
	switch v := any(val).(type) {
	case MyMessage:
		i, err := parse(v.fromBytes, slice)
		return any(v).(K), i, err
	case MyStruct:
		i, err := parse(v.fromBytes, slice)
		return any(v).(K), i, err
	default:
		return val, 0, fmt.Errorf("unsupported struct type for deserialization")
	}
}

// id = 1
type MyMessage struct {
	Id          *uint32         // index 0
	Name        *string         // index 1
	IsSomething *bool           // index 2
	SomeStruct  *MyStruct       // index 3
	SomeList    *[]string       // index 4
	SomeList2   *[]MyStruct     // index 5
	SomeList3   *[][]string     // index 6
	SomeList4   *[][][]MyStruct // index 7
}

func (MyMessage) TypeID() uint16 { return 1 }

func (msg MyMessage) toBytes(data *bytes.Buffer) {

	serializeInt16(1, data)
	startLenPos := data.Len()
	serializeUint32(0, data)

	if msg.Id != nil {
		serializeUint16(0, data)
		serializeUint32(*msg.Id, data)
	}

	if msg.Name != nil {
		serializeUint16(1, data)
		serializeString(*msg.Name, data)
	}

	if msg.IsSomething != nil {
		serializeUint16(2, data)
		serializeBool(*msg.IsSomething, data)
	}

	if msg.SomeStruct != nil {
		serializeUint16(3, data)
		serializeStruct(*msg.SomeStruct, data)
	}

	if msg.SomeList != nil {
		serializeUint16(4, data)
		newListSerializer(serializeString)(*msg.SomeList, data)
	}

	if msg.SomeList2 != nil {
		serializeUint16(5, data)
		newListSerializer(serializeStruct[MyStruct])(*msg.SomeList2, data)
	}

	if msg.SomeList3 != nil {
		serializeUint16(6, data)
		newListSerializer(newListSerializer(serializeString))(*msg.SomeList3, data)
	}

	if msg.SomeList4 != nil {
		serializeUint16(7, data)
		newListSerializer(newListSerializer(newListSerializer(serializeStruct[MyStruct])))(*msg.SomeList4, data)
	}

	binary.BigEndian.PutUint32(data.Bytes()[startLenPos:], uint32(len(data.Bytes())-(startLenPos+lenSize)))
}

func (msg *MyMessage) fromBytes(data []byte, fieldIndex uint16, offset int) (int, error) {
	switch fieldIndex {
	case 0:
		val, len, err := deserializeUint32(data, offset); msg.Id = &val ;return len, err
	case 1:
		val, len, err := deserializeString(data, offset)
		msg.Name = &val
		return len, err
	case 2:
		val, len, err := deserializeBool(data, offset)
		msg.IsSomething = &val
		return len, err
	case 3:
		val, len, err := deserializeStruct[MyStruct](data, offset)
		msg.SomeStruct = &val
		return len, err
	case 4:
		val, len, err := newListDeserializer(deserializeString)(data, offset)
		msg.SomeList = &val
		return len, err
	case 5:
		val, len, err := newListDeserializer(deserializeStruct[MyStruct])(data, offset)
		msg.SomeList2 = &val
		return len, err
	case 6:
		val, len, err := newListDeserializer(newListDeserializer(deserializeString))(data, offset)
		msg.SomeList3 = &val
		return len, err
	case 7:
		val, len, err := newListDeserializer(newListDeserializer(newListDeserializer(deserializeStruct[MyStruct])))(data, offset)
		msg.SomeList4 = &val
		return len, err
	}
	return 0, UnknownFieldError
}

// id = 2
type MyStruct struct {
	Value *string // index 0
}

func (MyStruct) TypeID() uint16 { return 2 }

func (s MyStruct) toBytes(data *bytes.Buffer) {
	serializeInt16(2, data)
	startLenPos := data.Len()
	serializeUint32(0, data)
	if s.Value != nil {
		serializeUint16(0, data)
		serializeString(*s.Value, data)
	}
	binary.BigEndian.PutUint32(data.Bytes()[startLenPos:], uint32(len(data.Bytes())-(startLenPos+lenSize)))
}

func (s *MyStruct) fromBytes(data []byte, fieldIndex uint16, offset int) (int, error) {
	switch fieldIndex {
	case 0:
		val, len, err := deserializeString(data, offset)
		s.Value = &val
		return len, err
	}
	return 0, UnknownFieldError
}

func main() {
	values, _ := parseFile()
	printGo(values)

	original := MyMessage{
		Id:          Ptr(uint32(42)),
		Name:        Ptr("Example"),
		IsSomething: Ptr(true),
		SomeStruct: &MyStruct{
			Value: Ptr("Nested"),
		},
		SomeList: &[]string{"one", "two", "three"},
		SomeList2: &[]MyStruct{
			{Value: Ptr("Struct1")},
			{Value: Ptr("Struct2")},
		},
		SomeList3: &[][]string{
			[]string{"A1", "A2"},
			[]string{"B1", "B2"},
		},
		SomeList4: &[][][]MyStruct{
			[][]MyStruct{
				[]MyStruct{
					{Value: Ptr("X1")},
					{Value: Ptr("X2")},
				},
				[]MyStruct{
					{Value: Ptr("Y1")},
				},
			},
			[][]MyStruct{
				[]MyStruct{
					{Value: Ptr("Z1")},
				},
			},
		},
	}

	serialized, err := Serialize(&original)
	fmt.Printf("Serialized bytes: %v\n", serialized)
	if err != nil {
		fmt.Println("Error during serialization:", err)
		return
	}

	deserialized := MyMessage{}
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		fmt.Println("Error during deserialization:", err)
		return
	}

	fmt.Println(
		*deserialized.Id,
		*deserialized.Name,
		*deserialized.IsSomething,
		*deserialized.SomeStruct.Value,
		*deserialized.SomeList,
		*(*deserialized.SomeList2)[1].Value,
		*deserialized.SomeList3,
		*(*deserialized.SomeList4)[0][1][0].Value,
	)
}
