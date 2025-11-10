package app

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

var _ authtypes.GenesisAccount = (*GenesisAccount)(nil)

// GenesisAccount 定义了一个实现 GenesisAccount 接口的类型，
// 用于在创世状态中承载模拟账户信息。
type GenesisAccount struct {
	*authtypes.BaseAccount

	// 归属账户字段
	OriginalVesting  sdk.Coins `json:"original_vesting" yaml:"original_vesting"`   // 初始化时的总归属代币数
	DelegatedFree    sdk.Coins `json:"delegated_free" yaml:"delegated_free"`       // 在委托时刻已归属并被委托的代币
	DelegatedVesting sdk.Coins `json:"delegated_vesting" yaml:"delegated_vesting"` // 在委托时刻尚未归属但已委托的代币
	StartTime        int64     `json:"start_time" yaml:"start_time"`               // 归属开始时间（UNIX 时间戳）
	EndTime          int64     `json:"end_time" yaml:"end_time"`                   // 归属结束时间（UNIX 时间戳）

	// 模块账户字段
	ModuleName        string   `json:"module_name" yaml:"module_name"`               // 模块账户名称
	ModulePermissions []string `json:"module_permissions" yaml:"module_permissions"` // 模块账户权限列表
}

// Validate 会检查归属账户与模块账户相关参数是否存在错误。
func (sga GenesisAccount) Validate() error {
	if !sga.OriginalVesting.IsZero() {
		if sga.StartTime >= sga.EndTime {
			return errors.New("vesting start-time cannot be before end-time")
		}
	}

	if sga.ModuleName != "" {
		ma := authtypes.ModuleAccount{
			BaseAccount: sga.BaseAccount, Name: sga.ModuleName, Permissions: sga.ModulePermissions,
		}
		if err := ma.Validate(); err != nil {
			return err
		}
	}

	return sga.BaseAccount.Validate()
}
