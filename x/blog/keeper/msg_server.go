package keeper

import (
	"cosmos-learn/x/blog/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl 根据传入的 Keeper 返回 MsgServer 接口实现。
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}
