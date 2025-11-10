package types

// GetBytes is a helper for serialising
func (p SendPostPacketData) GetBytes() ([]byte, error) {
	var modulePacket BlogPacketData

	modulePacket.Packet = &BlogPacketData_SendPostPacket{&p}

	return modulePacket.Marshal()
}
