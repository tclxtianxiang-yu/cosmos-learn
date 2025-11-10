package keeper

import (
	"context"
	"fmt"

	"cosmos-learn/x/blog/types"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
)

func (k msgServer) SendSendPost(ctx context.Context, msg *types.MsgSendSendPost) (*types.MsgSendSendPostResponse, error) {
	// 校验传入的消息。
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	if msg.Port == "" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid packet port")
	}

	if msg.ChannelID == "" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid packet channel")
	}

	if msg.TimeoutTimestamp == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid packet timeout")
	}

	// TODO: 在发送数据包前补充必要的业务逻辑。

	// 构造待发送的数据包。
	var packet types.SendPostPacketData

	packet.Title = msg.Title
	packet.Content = msg.Content

	// 发送数据包。
	_, err := k.TransmitSendPostPacket(
		ctx,
		packet,
		msg.Port,
		msg.ChannelID,
		clienttypes.ZeroHeight(),
		msg.TimeoutTimestamp,
	)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendSendPostResponse{}, nil
}
