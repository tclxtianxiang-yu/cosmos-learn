package keeper

import (
	"cosmos-learn/x/blog/types"
)

var _ types.QueryServer = queryServer{}

// NewQueryServerImpl 根据传入的 Keeper 返回 QueryServer 接口实现。
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k Keeper
}
