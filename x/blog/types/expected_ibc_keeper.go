package types

import (
	"context"

	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
)

// ChannelKeeper 定义了期望的 IBC 通道 keeper 接口。
type ChannelKeeper interface {
	GetChannel(ctx context.Context, portID, channelID string) (channeltypes.Channel, bool)
	GetNextSequenceSend(ctx context.Context, portID, channelID string) (uint64, bool)
	SendPacket(
		ctx context.Context,
		sourcePort string,
		sourceChannel string,
		timeoutHeight clienttypes.Height,
		timeoutTimestamp uint64,
		data []byte,
	) (uint64, error)
	ChanCloseInit(ctx context.Context, portID, channelID string) error
}
