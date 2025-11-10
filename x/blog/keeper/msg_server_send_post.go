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
	// validate incoming message
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

	// TODO: logic before transmitting the packet

	// Construct the packet
	var packet types.SendPostPacketData

	packet.Title = msg.Title
	packet.Content = msg.Content

	// Transmit the packet
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
