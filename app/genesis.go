package app

import (
	"encoding/json"
)

// 区块链的 GenesisState 在此表示为由标识字符串映射到原始 JSON 消息的字典。
// 该标识符用于确定创世信息属于哪个模块，以便在初始化链时正确路由。
// 在本应用中，默认的创世信息由 ModuleBasicManager 提供，
// 它会在初始化过程中收集每个 BasicModule 的 JSON 数据。
type GenesisState map[string]json.RawMessage
