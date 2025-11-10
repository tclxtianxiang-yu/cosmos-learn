package types

import (
	"context"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AuthKeeper 定义了对 Auth 模块的期望接口。
type AuthKeeper interface {
	AddressCodec() address.Codec
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // 仅用于仿真场景
	// 需要引入 account 模块的方法应在此声明。
}

// BankKeeper 定义了对 Bank 模块的期望接口。
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	// 需要引入 bank 模块的方法应在此声明。
}

// ParamSubspace 定义参数子空间的期望接口。
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}
