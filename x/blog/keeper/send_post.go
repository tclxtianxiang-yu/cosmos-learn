package keeper

import (
	"context"
	"errors"

	"cosmos-learn/x/blog/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
)

// TransmitSendPostPacket transmits the packet over IBC with the specified source port and source channel
func (k Keeper) TransmitSendPostPacket(
	ctx context.Context,
	packetData types.SendPostPacketData,
	sourcePort,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
) (uint64, error) {
	packetBytes, err := packetData.GetBytes()
	if err != nil {
		return 0, errorsmod.Wrapf(sdkerrors.ErrJSONMarshal, "cannot marshal the packet: %s", err)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return k.ibcKeeperFn().ChannelKeeper.SendPacket(sdkCtx, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, packetBytes)
}

// OnRecvSendPostPacket processes packet reception
func (k Keeper) OnRecvSendPostPacket(ctx context.Context, packet channeltypes.Packet, data types.SendPostPacketData) (packetAck types.SendPostPacketAck, err error) {
	// validate packet data upon receiving

	// TODO: packet reception logic

	return packetAck, nil
}

// OnAcknowledgementSendPostPacket responds to the success or failure of a packet
// acknowledgement written on the receiving chain.
func (k Keeper) OnAcknowledgementSendPostPacket(ctx context.Context, packet channeltypes.Packet, data types.SendPostPacketData, ack channeltypes.Acknowledgement) error {
	switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:

		// TODO: failed acknowledgement logic
		_ = dispatchedAck.Error

		return nil
	case *channeltypes.Acknowledgement_Result:
		// Decode the packet acknowledgment
		var packetAck types.SendPostPacketAck

		if err := k.cdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
			// The counter-party module doesn't implement the correct acknowledgment format
			return errors.New("cannot unmarshal acknowledgment")
		}

		// TODO: successful acknowledgement logic

		return nil
	default:
		// The counter-party module doesn't implement the correct acknowledgment format
		return errors.New("invalid acknowledgment format")
	}
}

// OnTimeoutSendPostPacket responds to the case where a packet has not been transmitted because of a timeout
func (k Keeper) OnTimeoutSendPostPacket(ctx context.Context, packet channeltypes.Packet, data types.SendPostPacketData) error {

	// TODO: packet timeout logic

	return nil
}
