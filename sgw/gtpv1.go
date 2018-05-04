package main

import (
	"encoding/binary"
)

const (
	GTPV1_PORT = 2152
)

type GTPV1 struct {
	Version                 uint8  // 3bit
	ProtocolType            uint8  // 1bit
	Reserved                uint8  // 1bit
	ExtensionHeaderFlag     uint8  // 1bit
	SequenceNumberFlag      uint8  // 1bit
	N_PDUNumberFlag         uint8  // 1bit
	MessageType             uint8  // 8bit
	MessageLength           uint16 // 16bit
	TEID                    uint32 // 32bit
	SequenceNumber          uint16 // 16bit
	N_PDUNumber             uint8  // 8bit
	NextExtensionFeaderType uint8  // 8bit
	Data                    []byte // Userdata
}

func (g *GTPV1) Marshal(userdata []byte) []byte {
	b := make([]byte, 12+len(userdata))

	b[0] = byte(((g.Version & 0x07) << 5) + ((g.ProtocolType & 0x01) << 4) + ((g.Reserved & 0x01) << 3) + ((g.ExtensionHeaderFlag & 0x01) << 2) + ((g.SequenceNumberFlag & 0x01) << 1) + (g.N_PDUNumberFlag & 0x01))
	b[1] = byte(g.MessageType)
	binary.BigEndian.PutUint16(b[2:4], g.MessageLength)
	binary.BigEndian.PutUint32(b[4:8], g.TEID)
	binary.BigEndian.PutUint16(b[8:10], g.SequenceNumber)
	b[10] = byte(g.N_PDUNumber)
	b[11] = byte(g.NextExtensionFeaderType)
	copy(b[12:], userdata)

	return b
}

// TODO：各フラグ値をみて下位４ビットがない場合の構造体を作る
func (g *GTPV1) Parse(b []byte) error {
	g.Version = uint8((b[0] & 0xE0) >> 5)
	g.ProtocolType = uint8((b[0] & 0x10) >> 4)
	g.Reserved = uint8((b[0] & 0x08) >> 3)
	g.ExtensionHeaderFlag = uint8((b[0] & 0x04) >> 2)
	g.SequenceNumberFlag = uint8((b[0] & 0x02) >> 1)
	g.N_PDUNumberFlag = uint8(b[0] & 0x01)
	g.MessageType = uint8(b[1])
	g.MessageLength = uint16(binary.BigEndian.Uint16(b[2:4]))
	g.TEID = uint32(binary.BigEndian.Uint32(b[4:8]))
	g.SequenceNumber = uint16(binary.BigEndian.Uint16(b[8:10]))
	g.N_PDUNumber = uint8(b[10])
	g.NextExtensionFeaderType = uint8(b[11])
	g.Data = b[12:]

	return nil
}
