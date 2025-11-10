package cmd

import (
	cmtcfg "github.com/cometbft/cometbft/config"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
)

// initCometBFTConfig 用于覆盖默认的 CometBFT 配置值。
// 如果应用不需要自定义配置，则返回 cmtcfg.DefaultConfig。
func initCometBFTConfig() *cmtcfg.Config {
	cfg := cmtcfg.DefaultConfig()

	// 这些配置会让节点内存承受更大压力
	// cfg.P2P.MaxNumInboundPeers = 100 示例：调高入站连接上限
	// cfg.P2P.MaxNumOutboundPeers = 40 示例：调高出站连接上限

	return cfg
}

// initAppConfig 用于覆盖默认的 appConfig 模板和配置。
// 如果应用不需要自定义配置，则返回 "", nil。
func initAppConfig() (string, interface{}) {
	// 以下代码片段仅供参考。
	type CustomAppConfig struct {
		serverconfig.Config `mapstructure:",squash"`
	}

	// 视情况允许链的开发者覆盖 SDK 的默认服务器配置。
	srvCfg := serverconfig.DefaultConfig()
	// SDK 在 app.toml 中将默认最小 gas 价格设为 ""（空值）。
	// 如果验证人保持为空，节点会在启动时停止运行。
	// 链的开发者可以在此为验证人设置 app.toml 的默认值。
	//
	// 总结如下：
	// - 若保持 srvCfg.MinGasPrices = ""，所有验证人必须自行调整各自的 app.toml 配置；
	// - 若将 srvCfg.MinGasPrices 设置为非空，验证人可以通过修改自己的 app.toml 覆盖该值，或者直接使用这里的默认值。
	//
	// 在测试环境中，我们将最小 gas 价格设为 0。
	// srvCfg.MinGasPrices = "0stake" 示例：设置默认 gas 单价

	customAppConfig := CustomAppConfig{
		Config: *srvCfg,
	}

	customAppTemplate := serverconfig.DefaultConfigTemplate
	// 编辑默认模板文件
	//
	// customAppTemplate := serverconfig.DefaultConfigTemplate + `
	// [wasm]
	// # 这是允许任意 x/wasm “智能”查询使用的最大 SDK gas（含 wasm 与存储）
	// query_gas_limit = 300000  # 示例：将查询 gas 上限设为 300000
	// # 这是为了提速而缓存于内存中的 wasm 虚拟机实例数量
	// # 警告：当前此功能不够稳定，可能导致崩溃，除非本地测试建议保持为 0
	// lru_size = 0  # 示例：维持 LRU 缓存大小为 0`

	return customAppTemplate, customAppConfig
}
