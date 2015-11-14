package common

import (
	"encoding/binary"
	"math"
)

//==================================
//-------------- Read --------------
//==================================
func ReadBool(buf []byte) bool {
	return bool(buf[0] != 0)
}

func ReadByte(buf []byte) byte {
	return buf[0]
}

func ReadInt8(buf []byte) int8 {
	return int8(buf[0])
}

func ReadInt16(buf []byte) int16 {
	return int16(binary.BigEndian.Uint16(buf))
}

func ReadInt32(buf []byte) int32 {
	return int32(binary.BigEndian.Uint32(buf))
}

func ReadInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func ReadUint8(buf []byte) uint8 {
	return uint8(buf[0])
}

func ReadUint16(buf []byte) uint16 {
	return uint16(binary.BigEndian.Uint16(buf))
}

func ReadUint32(buf []byte) uint32 {
	return uint32(binary.BigEndian.Uint32(buf))
}

func ReadUint64(buf []byte) uint64 {
	return uint64(binary.BigEndian.Uint64(buf))
}

func ReadFloat32(buf []byte) float32 {
	return float32(binary.BigEndian.Uint32(buf))
}

func ReadFloat64(buf []byte) float64 {
	return float64(binary.BigEndian.Uint64(buf))
}

func ReadString(buf []byte) string {
	return string(buf[:])
}

func ReadBytes(buf []byte) []byte {
	return buf
}

//==================================
//-------------- Write -------------
//==================================

func WriteBool(buf []byte, v bool) {
	if v {
		buf[0] = 1
	} else {
		buf[0] = 0
	}
}

func WriteByte(buf []byte, v byte) {
	buf[0] = byte(v)
}

func WriteInt16(buf []byte, v int16) {
	binary.BigEndian.PutUint16(buf, uint16(v))
}

func WriteInt32(buf []byte, v int32) {
	binary.BigEndian.PutUint32(buf, uint32(v))
}

func WriteInt64(buf []byte, v int64) {
	binary.BigEndian.PutUint64(buf, uint64(v))
}

func WriteUint8(buf []byte, v uint8) {
	buf[0] = byte(v)
}

func WriteUint16(buf []byte, v uint16) {
	binary.BigEndian.PutUint16(buf, v)
}

func WriteUint32(buf []byte, v uint32) {
	binary.BigEndian.PutUint32(buf, v)
}

func WriteUint64(buf []byte, v uint64) {
	binary.BigEndian.PutUint64(buf, v)
}

func WriteFloat32(buf []byte, v float32) {
	binary.BigEndian.PutUint32(buf, math.Float32bits(v))
}

func WriteFloat64(buf []byte, v float64) {
	binary.BigEndian.PutUint64(buf, math.Float64bits(v))
}

func WriteString(buf []byte, v string) {
	copy(buf[0:], []byte(v))
}

func WriteBytes(buf []byte, v []byte) {
	copy(buf[0:], v)
}
