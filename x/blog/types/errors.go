package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/blog 模块的哨兵错误定义。
var (
	ErrInvalidSigner        = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidPacketTimeout = errors.Register(ModuleName, 1500, "invalid packet timeout")
	ErrInvalidVersion       = errors.Register(ModuleName, 1501, "invalid version")
)
