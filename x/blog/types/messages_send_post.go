package types

func NewMsgSendSendPost(
	creator string,
	port string,
	channelID string,
	timeoutTimestamp uint64,
	title string,
	content string,
) *MsgSendSendPost {
	return &MsgSendSendPost{
		Creator:          creator,
		Port:             port,
		ChannelID:        channelID,
		TimeoutTimestamp: timeoutTimestamp,
		Title:            title,
		Content:          content,
	}
}
