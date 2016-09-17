package main

const (
	TYPE_DATA      = iota
	TYPE_ADVERTISE = iota
)

type IPv4 [4]byte

type TBPacket struct {
	PacketType uint8
	Data       []byte
}

func (packet *TBPacket) Marshal() *[]byte {
	buffer := make([]byte, 1+len(packet.Data))
	buffer[0] = packet.PacketType
	copy(buffer[1:], packet.Data)
	return &buffer

}

func UnmarshalTBPacket(data *[]byte) *TBPacket {
	return &TBPacket{PacketType: (*data)[0], Data: (*data)[1:]}
}

type TBAdvertisePacket struct {
	IPAddress IPv4
}
