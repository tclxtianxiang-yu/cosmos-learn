package app

import (
	"encoding/json"
	"fmt"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"

	storetypes "cosmossdk.io/store/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ExportAppStateAndValidators 将应用状态导出为创世文件所需的数据。
func (app *App) ExportAppStateAndValidators(forZeroHeight bool, jailAllowedAddrs, modulesToExport []string) (servertypes.ExportedApp, error) {
	// 假设验证人可以从下一个区块开始提款
	ctx := app.NewContextLegacy(true, cmtproto.Header{Height: app.LastBlockHeight()})

	// 我们以最后区块高度加一的高度导出，因为 CometBFT 会在该高度启动 InitChain。
	height := app.LastBlockHeight() + 1
	if forZeroHeight {
		height = 0
		app.prepForZeroHeightGenesis(ctx, jailAllowedAddrs)
	}

	genState, err := app.ModuleManager.ExportGenesisForModules(ctx, app.appCodec, modulesToExport)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	appState, err := json.MarshalIndent(genState, "", "  ")
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	validators, err := staking.WriteValidators(ctx, app.StakingKeeper)

	return servertypes.ExportedApp{
		AppState:        appState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: app.BaseApp.GetConsensusParams(ctx),
	}, err
}

// prepForZeroHeightGenesis 用于在高度为零时重新启动链。
// 注意：零高度创世是一个即将废弃的临时功能，
//
//	后续将以按区块高度导出方案取而代之。
func (app *App) prepForZeroHeightGenesis(ctx sdk.Context, jailAllowedAddrs []string) {
	applyAllowedAddrs := false

	// 检查是否传入允许豁免的地址列表
	if len(jailAllowedAddrs) > 0 {
		applyAllowedAddrs = true
	}

	allowedAddrsMap := make(map[string]bool)

	for _, addr := range jailAllowedAddrs {
		_, err := app.InterfaceRegistry().SigningContext().ValidatorAddressCodec().StringToBytes(addr)
		if err != nil {
			panic(err)
		}
		allowedAddrsMap[addr] = true
	}

	/* Handle fee distribution state. */

	// 提取所有验证人的佣金
	err := app.StakingKeeper.IterateValidators(ctx, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
		valBz, err := app.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
		if err != nil {
			panic(err)
		}
		_, _ = app.DistrKeeper.WithdrawValidatorCommission(ctx, valBz)
		return false
	})
	if err != nil {
		panic(err)
	}

	// 提取所有委托人的奖励
	dels, err := app.StakingKeeper.GetAllDelegations(ctx)
	if err != nil {
		panic(err)
	}

	for _, delegation := range dels {
		valAddr, err := app.InterfaceRegistry().SigningContext().ValidatorAddressCodec().StringToBytes(delegation.ValidatorAddress)
		if err != nil {
			panic(err)
		}

		delAddr, err := app.InterfaceRegistry().SigningContext().AddressCodec().StringToBytes(delegation.DelegatorAddress)
		if err != nil {
			panic(err)
		}

		_, _ = app.DistrKeeper.WithdrawDelegationRewards(ctx, delAddr, valAddr)
	}

	// 清理验证人的惩罚事件
	app.DistrKeeper.DeleteAllValidatorSlashEvents(ctx)

	// 清理验证人的历史奖励
	app.DistrKeeper.DeleteAllValidatorHistoricalRewards(ctx)

	// 将上下文高度设置为零
	height := ctx.BlockHeight()
	ctx = ctx.WithBlockHeight(0)

	// 重新初始化所有验证人
	err = app.StakingKeeper.IterateValidators(ctx, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
		valBz, err := app.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
		if err != nil {
			panic(err)
		}
		// 将尚未提取的奖励代币捐赠到社区资金池
		rewards, err := app.DistrKeeper.GetValidatorOutstandingRewardsCoins(ctx, valBz)
		if err != nil {
			panic(err)
		}
		feePool, err := app.DistrKeeper.FeePool.Get(ctx)
		if err != nil {
			panic(err)
		}
		feePool.CommunityPool = feePool.CommunityPool.Add(rewards...)
		if err := app.DistrKeeper.FeePool.Set(ctx, feePool); err != nil {
			panic(err)
		}

		if err := app.DistrKeeper.Hooks().AfterValidatorCreated(ctx, valBz); err != nil {
			panic(err)
		}
		return false
	})
	if err != nil {
		panic(err)
	}

	// 重新初始化所有委托关系
	for _, del := range dels {
		valAddr, err := app.InterfaceRegistry().SigningContext().ValidatorAddressCodec().StringToBytes(del.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		delAddr, err := app.InterfaceRegistry().SigningContext().AddressCodec().StringToBytes(del.DelegatorAddress)
		if err != nil {
			panic(err)
		}

		if err := app.DistrKeeper.Hooks().BeforeDelegationCreated(ctx, delAddr, valAddr); err != nil {
			// 实际不会触发，因为 BeforeDelegationCreated 总是返回 nil
			panic(fmt.Errorf("error while incrementing period: %w", err))
		}

		if err := app.DistrKeeper.Hooks().AfterDelegationModified(ctx, delAddr, valAddr); err != nil {
			// 实际不会触发，因为 AfterDelegationModified 总是返回 nil
			panic(fmt.Errorf("error while creating a new delegation period record: %w", err))
		}
	}

	// 恢复上下文高度
	ctx = ctx.WithBlockHeight(height)

	/* Handle staking state. */

	// 遍历再委托记录并重置创建高度
	err = app.StakingKeeper.IterateRedelegations(ctx, func(_ int64, red stakingtypes.Redelegation) (stop bool) {
		for i := range red.Entries {
			red.Entries[i].CreationHeight = 0
		}
		err = app.StakingKeeper.SetRedelegation(ctx, red)
		if err != nil {
			panic(err)
		}
		return false
	})
	if err != nil {
		panic(err)
	}

	// 遍历解除委托记录并重置创建高度
	err = app.StakingKeeper.IterateUnbondingDelegations(ctx, func(_ int64, ubd stakingtypes.UnbondingDelegation) (stop bool) {
		for i := range ubd.Entries {
			ubd.Entries[i].CreationHeight = 0
		}
		err = app.StakingKeeper.SetUnbondingDelegation(ctx, ubd)
		if err != nil {
			panic(err)
		}
		return false
	})
	if err != nil {
		panic(err)
	}

	// 按投票权从高到低遍历验证人，重置抵押高度，并更新抵押计数器。
	store := ctx.KVStore(app.GetKey(stakingtypes.StoreKey))
	iter := storetypes.KVStoreReversePrefixIterator(store, stakingtypes.ValidatorsKey)

	for ; iter.Valid(); iter.Next() {
		addr := sdk.ValAddress(stakingtypes.AddressFromValidatorsKey(iter.Key()))
		validator, err := app.StakingKeeper.GetValidator(ctx, addr)
		if err != nil {
			panic("expected validator, not found")
		}

		valAddr, err := app.StakingKeeper.ValidatorAddressCodec().BytesToString(addr)
		if err != nil {
			panic(err)
		}

		validator.UnbondingHeight = 0
		if applyAllowedAddrs && !allowedAddrsMap[valAddr] {
			validator.Jailed = true
		}

		if err = app.StakingKeeper.SetValidator(ctx, validator); err != nil {
			panic(err)
		}
	}

	if err := iter.Close(); err != nil {
		app.Logger().Error("error while closing the key-value store reverse prefix iterator: ", err)
		return
	}

	_, err = app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	if err != nil {
		panic(err)
	}

	/* Handle slashing state. */

	// 将签名信息的起始高度重置为零
	if err := app.SlashingKeeper.IterateValidatorSigningInfos(
		ctx,
		func(addr sdk.ConsAddress, info slashingtypes.ValidatorSigningInfo) (stop bool) {
			info.StartHeight = 0
			_ = app.SlashingKeeper.SetValidatorSigningInfo(ctx, addr, info)
			return false
		},
	); err != nil {
		panic(err)
	}

}
