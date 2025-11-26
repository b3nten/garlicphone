package schematest

import(
	"bytes"
	"encoding/binary"
	"6enten/garlicphone/optional"
	"errors"
	"fmt"
)

// Type definitions

type Player struct { 
	 Id optional.Optional[uint32]
	 Name optional.Optional[string]
}

func (it Player) ToBytes() []byte {
	var data bytes.Buffer

	appendUint16(&data, 49920)
	appendUint32(&data, 0)

	if it.Id.Exists() {
		appendUint16(&data, 10)
		appendUint32(&data, it.Id.MustGet())
	}
	if it.Name.Exists() {
		appendUint16(&data, 11)
		appendString(&data, it.Name.MustGet())
	}

	binary.BigEndian.PutUint32(data.Bytes()[idSize:], uint32(len(data.Bytes())-(idSize+lenSize)))
	return data.Bytes()
}

func (it *Player) FromBytes(bytes []byte) (int, error) {
	return parse(it, bytes)
}

func (it *Player) parseFields(data []byte, fieldIndex uint16, offset int) (int, error) {
	switch fieldIndex {
	case 10: return readUint32(data, offset, &it.Id)
	case 11: return readString(data, offset, &it.Name)
	}
	return 0, UnknownFieldError
}


// Utilities


const lenSize = 4
const idSize = 2

type Serializable interface {
	FromBytes([]byte) (int, error)
	ToBytes() []byte
}

func readBool(data []byte, offset int, out **bool) (int, error) {
	slice := data[offset:]
	if len(slice) < 1 {
		return 0, fmt.Errorf("insufficient data for bool")
	}
	val := slice[0] != 0
	*out = &val
	return 1, nil
}

func appendBool(data *bytes.Buffer, value bool) int {
	var boolByte byte = 0
	if value {
		boolByte = 1
	}
	data.WriteByte(boolByte)
	return 1
}

func readInt8(data []byte, offset int, out **int8) (int, error) {
	slice := data[offset:]
	if len(slice) < 1 {
		return 0, fmt.Errorf("insufficient data for int8")
	}
	val := int8(slice[0])
	*out = &val
	return 1, nil
}

func appendInt8(data *bytes.Buffer, value int8) int {
	data.WriteByte(byte(value))
	return 1
}

func readUint8(data []byte, offset int, out **uint8) (int, error) {
	slice := data[offset:]
	if len(slice) < 1 {
		return 0, fmt.Errorf("insufficient data for uint8")
	}
	val := uint8(slice[0])
	*out = &val
	return 1, nil
}

func appendUint8(data *bytes.Buffer, value uint8) int {
	data.WriteByte(byte(value))
	return 1
}

func readInt16(data []byte, offset int, out **int16) (int, error) {
	slice := data[offset:]
	if len(slice) < 2 {
		return 0, fmt.Errorf("insufficient data for int16")
	}
	val := int16(binary.BigEndian.Uint16(slice))
	*out = &val
	return 2, nil
}

func appendInt16(data *bytes.Buffer, value int16) int {
	binary.Write(data, binary.BigEndian, value)
	return 2
}

func readUint16(data []byte, offset int, out **uint16) (int, error) {
	slice := data[offset:]
	if len(slice) < 2 {
		return 0, fmt.Errorf("insufficient data for uint16")
	}
	val := uint16(binary.BigEndian.Uint16(slice))
	*out = &val
	return 2, nil
}

func appendUint16(data *bytes.Buffer, value uint16) int {
	binary.Write(data, binary.BigEndian, value)
	return 2
}

func readInt32(data []byte, offset int, out **int32) (int, error) {
	slice := data[offset:]
	if len(slice) < 4 {
		return 0, fmt.Errorf("insufficient data for int32")
	}
	val := int32(binary.BigEndian.Uint32(slice))
	*out = &val
	return 4, nil
}

func appendInt32(data *bytes.Buffer, value int32) int {
	binary.Write(data, binary.BigEndian, value)
	return 4
}

func readUint32(data []byte, offset int, out **uint32) (int, error) {
	slice := data[offset:]
	if len(slice) < 4 {
		return 0, fmt.Errorf("insufficient data for uint32")
	}
	val := uint32(binary.BigEndian.Uint32(slice))
	*out = &val
	return 4, nil
}

func appendUint32(data *bytes.Buffer, value uint32) int {
	binary.Write(data, binary.BigEndian, value)
	return 4
}

func readString(data []byte, offset int, out **string) (int, error) {
	slice := data[offset:]
	if len(slice) < 4 {
		return 0, fmt.Errorf("insufficient data for string length")
	}
	strLen := int(binary.BigEndian.Uint32(slice))
	if strLen == 0 {
		*out = new(string)
		return 4, nil
	}
	if len(slice[4:]) < strLen {
		return 0, fmt.Errorf("insufficient data for string content")
	}
	val := string(slice[4 : 4+strLen])
	*out = &val
	return 4 + strLen, nil
}

func appendString(data *bytes.Buffer, str string) int {
	binary.Write(data, binary.BigEndian, int32(len(str)))
	data.Write([]byte(str))
	return len(str) + 4
}

func appendList[T Serializable](list []T, data *bytes.Buffer) {
	appendUint32(data, uint32(0))
	i := len(data.Bytes())
	for _, item := range list {
		data.Write(item.ToBytes())
	}
	binary.BigEndian.PutUint32(data.Bytes()[i-lenSize:], uint32(len(data.Bytes())-i))
}

func appendPrimitiveList[T any](list []T, data *bytes.Buffer, serializer func(*bytes.Buffer, T) int) {
	appendUint32(data, uint32(0))
	i := len(data.Bytes())
	for _, item := range list {
		serializer(data, item)
	}
	binary.BigEndian.PutUint32(data.Bytes()[i-lenSize:], uint32(len(data.Bytes())-i))
}

func readStruct[K Serializable](data []byte, offset int, out K) (int, error) {
	n, err := out.FromBytes(data[offset:])
	if err != nil && !errors.Is(err, UnknownFieldError) {
		return n, err
	}
	return n, nil
}

func readList[K any](data []byte, offset int, list *[]K, reader func([]byte, int, **K) (int, error)) (int, error) {
	listLen := int(binary.BigEndian.Uint32(data[offset:]))
	for i := offset + lenSize; i < offset+listLen; {
		var val = new(K)
		n, _ := reader(data, i, &val)
		*list = append(*list, *val)
		i += n
	}
	return listLen + lenSize, nil
}

func parse(s interface {
	Serializable
	parseFields([]byte, uint16, int) (int, error)
}, bytes []byte) (int, error) {
	if len(bytes) < idSize+lenSize {
		return 0, fmt.Errorf("data too short: need at least 6 bytes for id and length")
	}
	dataSize := binary.BigEndian.Uint32(bytes[idSize:])
	totalSize := int(dataSize) + idSize + lenSize
	if len(bytes) < totalSize {
		return 0, fmt.Errorf("data too short: expected %d bytes, got %d", totalSize, len(bytes))
	}
	for i := idSize + lenSize; i < totalSize; {
		fieldIndex, err := getField(bytes[i:])
		if err != nil {
			return totalSize, err
		}
		i += 2
		next, err := s.parseFields(bytes, fieldIndex, i)
		if err != nil {
			if errors.Is(err, UnknownFieldError) {
				return totalSize, nil
			} else {
				return totalSize, err
			}
		}
		i += next
	}
	return totalSize, nil
}

func getField(b []byte) (uint16, error) {
	if len(b) < 2 {
		return 0, fmt.Errorf("insufficient data for field index")
	}
	return binary.BigEndian.Uint16(b), nil
}
