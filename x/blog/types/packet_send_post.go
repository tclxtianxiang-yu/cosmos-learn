package types

// GetBytes 提供序列化辅助方法。
func (p SendPostPacketData) GetBytes() ([]byte, error) {
	var modulePacket BlogPacketData

	modulePacket.Packet = &BlogPacketData_SendPostPacket{&p}

	return modulePacket.Marshal()
}
