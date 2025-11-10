package types

import "cosmossdk.io/collections"

const (
	// ModuleName 定义模块名称。
	ModuleName = "blog"

	// StoreKey 定义模块的主存储键。
	StoreKey = ModuleName

	// GovModuleName 复刻治理模块名称，以避免直接依赖 x/gov。
	// 如若治理模块名称调整，此处也需同步更新。
	// 参见：https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-beta.2/x/gov/types/keys.go#L9
	GovModuleName = "gov"

	// Version 定义 IBC 模块当前支持的版本。
	Version = "blog-1"

	// PortID 是模块绑定的默认端口标识。
	PortID = "blog"
)

var (
	// PortKey 定义在存储中保存端口 ID 的键。
	PortKey = collections.NewPrefix("blog-port-")
)

// ParamsKey 是读取全部 Params 的前缀。
var ParamsKey = collections.NewPrefix("p_blog")

var (
	PostKey      = collections.NewPrefix("post/value/")
	PostCountKey = collections.NewPrefix("post/count/")
)
