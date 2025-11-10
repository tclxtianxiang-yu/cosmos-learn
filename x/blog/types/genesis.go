package types

import (
	"fmt"

	host "github.com/cosmos/ibc-go/v10/modules/core/24-host"
)

// DefaultGenesis 返回模块的默认创世状态。
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		PortId: PortID, PostList: []Post{}}
}

// Validate 执行基础的创世状态校验，发现问题时返回错误，
// 防止初始化失败。
func (gs GenesisState) Validate() error {
	if err := host.PortIdentifierValidator(gs.PortId); err != nil {
		return err
	}
	postIdMap := make(map[uint64]bool)
	postCount := gs.GetPostCount()
	for _, elem := range gs.PostList {
		if _, ok := postIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for post")
		}
		if elem.Id >= postCount {
			return fmt.Errorf("post id should be lower or equal than the last id")
		}
		postIdMap[elem.Id] = true
	}

	return gs.Params.Validate()
}
