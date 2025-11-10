package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	channelutils "github.com/cosmos/ibc-go/v10/modules/core/04-channel/client/utils"
	"github.com/cosmos/ibc-go/v10/modules/core/exported"
	"github.com/spf13/cobra"

	"cosmos-learn/x/blog/types"
)

// CmdSendSendPost() returns the SendPost send packet command.
// This command does not use AutoCLI because it gives a better UX to do not.
func CmdSendSendPost() *cobra.Command {
	flagPacketTimeoutTimestamp := "packet-timeout-timestamp"

	cmd := &cobra.Command{
		Use:   "send-send-post [src-port] [src-channel] [title] [content]",
		Short: "Send a sendPost over IBC",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			creator := clientCtx.GetFromAddress().String()
			srcPort := args[0]
			srcChannel := args[1]

			argTitle := args[2]
			argContent := args[3]

			// Get the relative timeout timestamp
			timeoutTimestamp, err := cmd.Flags().GetUint64(flagPacketTimeoutTimestamp)
			if err != nil {
				return err
			}

			clientRes, err := channelutils.QueryChannelClientState(clientCtx, srcPort, srcChannel, false)
			if err != nil {
				return err
			}

			var clientState exported.ClientState
			if err := clientCtx.InterfaceRegistry.UnpackAny(clientRes.IdentifiedClientState.ClientState, &clientState); err != nil {
				return err
			}

			consensusStateAny, err := channelutils.QueryChannelConsensusState(clientCtx, srcPort, srcChannel, clienttypes.Height{}, false)
			if err != nil {
				return err
			}

			var consensusState exported.ConsensusState
			if err := clientCtx.InterfaceRegistry.UnpackAny(consensusStateAny.GetConsensusState(), &consensusState); err != nil {
				return err
			}

			if timeoutTimestamp != 0 {
				timeoutTimestamp = consensusState.GetTimestamp() + timeoutTimestamp //nolint:staticcheck // client side
			}

			msg := types.NewMsgSendSendPost(creator, srcPort, srcChannel, timeoutTimestamp, argTitle, argContent)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Uint64(flagPacketTimeoutTimestamp, DefaultRelativePacketTimeoutTimestamp, "Packet timeout timestamp in nanoseconds. Default is 10 minutes.")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
