package app

import (
	_ "cosmos-learn/x/blog/module"
	blogmoduletypes "cosmos-learn/x/blog/types"
	"time"

	runtimev1alpha1 "cosmossdk.io/api/cosmos/app/runtime/v1alpha1"
	appv1alpha1 "cosmossdk.io/api/cosmos/app/v1alpha1"
	authmodulev1 "cosmossdk.io/api/cosmos/auth/module/v1"
	authzmodulev1 "cosmossdk.io/api/cosmos/authz/module/v1"
	bankmodulev1 "cosmossdk.io/api/cosmos/bank/module/v1"
	circuitmodulev1 "cosmossdk.io/api/cosmos/circuit/module/v1"
	consensusmodulev1 "cosmossdk.io/api/cosmos/consensus/module/v1"
	distrmodulev1 "cosmossdk.io/api/cosmos/distribution/module/v1"
	epochsmodulev1 "cosmossdk.io/api/cosmos/epochs/module/v1"
	evidencemodulev1 "cosmossdk.io/api/cosmos/evidence/module/v1"
	feegrantmodulev1 "cosmossdk.io/api/cosmos/feegrant/module/v1"
	genutilmodulev1 "cosmossdk.io/api/cosmos/genutil/module/v1"
	govmodulev1 "cosmossdk.io/api/cosmos/gov/module/v1"
	groupmodulev1 "cosmossdk.io/api/cosmos/group/module/v1"
	mintmodulev1 "cosmossdk.io/api/cosmos/mint/module/v1"
	nftmodulev1 "cosmossdk.io/api/cosmos/nft/module/v1"
	paramsmodulev1 "cosmossdk.io/api/cosmos/params/module/v1"
	slashingmodulev1 "cosmossdk.io/api/cosmos/slashing/module/v1"
	stakingmodulev1 "cosmossdk.io/api/cosmos/staking/module/v1"
	txconfigv1 "cosmossdk.io/api/cosmos/tx/config/v1"
	upgrademodulev1 "cosmossdk.io/api/cosmos/upgrade/module/v1"
	vestingmodulev1 "cosmossdk.io/api/cosmos/vesting/module/v1"
	"cosmossdk.io/depinject/appconfig"
	_ "cosmossdk.io/x/circuit" // 为副作用而导入
	circuittypes "cosmossdk.io/x/circuit/types"
	_ "cosmossdk.io/x/evidence" // 为副作用而导入
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	_ "cosmossdk.io/x/feegrant/module" // 为副作用而导入
	"cosmossdk.io/x/nft"
	_ "cosmossdk.io/x/nft/module" // 为副作用而导入
	_ "cosmossdk.io/x/upgrade"    // 为副作用而导入
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	_ "github.com/cosmos/cosmos-sdk/x/auth/tx/config" // 为副作用而导入
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	_ "github.com/cosmos/cosmos-sdk/x/auth/vesting" // 为副作用而导入
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	_ "github.com/cosmos/cosmos-sdk/x/authz/module" // 为副作用而导入
	_ "github.com/cosmos/cosmos-sdk/x/bank"         // 为副作用而导入
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	_ "github.com/cosmos/cosmos-sdk/x/consensus" // 为副作用而导入
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	_ "github.com/cosmos/cosmos-sdk/x/distribution" // 为副作用而导入
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	_ "github.com/cosmos/cosmos-sdk/x/epochs" // 为副作用而导入
	epochstypes "github.com/cosmos/cosmos-sdk/x/epochs/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	_ "github.com/cosmos/cosmos-sdk/x/gov" // 为副作用而导入
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	_ "github.com/cosmos/cosmos-sdk/x/group/module" // 为副作用而导入
	_ "github.com/cosmos/cosmos-sdk/x/mint"         // 为副作用而导入
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	_ "github.com/cosmos/cosmos-sdk/x/params" // 为副作用而导入
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	_ "github.com/cosmos/cosmos-sdk/x/slashing" // 为副作用而导入
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	_ "github.com/cosmos/cosmos-sdk/x/staking" // 为副作用而导入
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icatypes "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	"google.golang.org/protobuf/types/known/durationpb"
)

var (
	moduleAccPerms = []*authmodulev1.ModuleAccountPermission{
		{Account: authtypes.FeeCollectorName},
		{Account: distrtypes.ModuleName},
		{Account: minttypes.ModuleName, Permissions: []string{authtypes.Minter}},
		{Account: stakingtypes.BondedPoolName, Permissions: []string{authtypes.Burner, stakingtypes.ModuleName}},
		{Account: stakingtypes.NotBondedPoolName, Permissions: []string{authtypes.Burner, stakingtypes.ModuleName}},
		{Account: govtypes.ModuleName, Permissions: []string{authtypes.Burner}},
		{Account: nft.ModuleName},
		{Account: ibctransfertypes.ModuleName, Permissions: []string{authtypes.Minter, authtypes.Burner}},
		{Account: icatypes.ModuleName},
	}

	// 被禁止的模块账户地址
	blockAccAddrs = []string{
		authtypes.FeeCollectorName,
		distrtypes.ModuleName,
		minttypes.ModuleName,
		stakingtypes.BondedPoolName,
		stakingtypes.NotBondedPoolName,
		nft.ModuleName,
		// 允许以下模块账户接收资金：
		// govtypes.ModuleName
	}

	// 应用配置（由 depinject 使用）
	appConfig = appconfig.Compose(&appv1alpha1.Config{
		Modules: []*appv1alpha1.ModuleConfig{
			{
				Name: runtime.ModuleName,
				Config: appconfig.WrapAny(&runtimev1alpha1.Module{
					AppName: Name,
					// 注意：升级模块必须优先执行
					PreBlockers: []string{
						upgradetypes.ModuleName,
						authtypes.ModuleName,
						// starport 脚手架使用的占位行 # stargate/app/preBlockers
					},
					// 在 BeginBlock 中，削减流程发生在 distr.BeginBlocker 之后，
					// 以确保验证人费用池不会残留资金，从而保持 CanWithdrawInvariant 不变式。
					// 注意：当 HistoricalEntries 参数大于 0 时必须启用 staking 模块。
					BeginBlockers: []string{
						minttypes.ModuleName,
						distrtypes.ModuleName,
						slashingtypes.ModuleName,
						evidencetypes.ModuleName,
						stakingtypes.ModuleName,
						authz.ModuleName,
						epochstypes.ModuleName,
						// IBC 模块
						ibcexported.ModuleName,
						// 链上业务模块
						blogmoduletypes.ModuleName,
						// starport 脚手架使用的占位行 # stargate/app/beginBlockers
					},
					EndBlockers: []string{
						govtypes.ModuleName,
						stakingtypes.ModuleName,
						feegrant.ModuleName,
						group.ModuleName,
						// 链上业务模块
						blogmoduletypes.ModuleName,
						// starport 脚手架使用的占位行 # stargate/app/endBlockers
					},
					// 下列配置通常只在 ModuleName 与 StoreKey 不一致时需要。
					OverrideStoreKeys: []*runtimev1alpha1.StoreKeyConfig{
						{
							ModuleName: authtypes.ModuleName,
							KvStoreKey: "acc",
						},
					},
					// 注意：genutils 模块必须位于 staking 之后，
					// 才能用创世账户的代币正确初始化各个池子。
					// 注意：genutils 模块也必须位于 auth 之后，以便读取 auth 暴露的参数。
					InitGenesis: []string{
						consensustypes.ModuleName,
						authtypes.ModuleName,
						banktypes.ModuleName,
						distrtypes.ModuleName,
						stakingtypes.ModuleName,
						slashingtypes.ModuleName,
						govtypes.ModuleName,
						minttypes.ModuleName,
						genutiltypes.ModuleName,
						evidencetypes.ModuleName,
						authz.ModuleName,
						feegrant.ModuleName,
						vestingtypes.ModuleName,
						nft.ModuleName,
						group.ModuleName,
						upgradetypes.ModuleName,
						circuittypes.ModuleName,
						epochstypes.ModuleName,
						// IBC 模块
						ibcexported.ModuleName,
						ibctransfertypes.ModuleName,
						icatypes.ModuleName,
						// 链上业务模块
						blogmoduletypes.ModuleName,
						// starport 脚手架使用的占位行 # stargate/app/initGenesis
					},
				}),
			},
			{
				Name: authtypes.ModuleName,
				Config: appconfig.WrapAny(&authmodulev1.Module{
					Bech32Prefix:                AccountAddressPrefix,
					ModuleAccountPermissions:    moduleAccPerms,
					EnableUnorderedTransactions: true,
					// 默认情况下，模块的权限由治理模块掌控。可以通过以下方式自定义：
					// Authority: "group", // 通过模块名称设置自定义权限控制者
					// Authority: "cosmos1cwwv22j5ca08ggdv9c2uky355k908694z577tv", // 或者指定某个地址
				}),
			},
			{
				Name:   vestingtypes.ModuleName,
				Config: appconfig.WrapAny(&vestingmodulev1.Module{}),
			},
			{
				Name: banktypes.ModuleName,
				Config: appconfig.WrapAny(&bankmodulev1.Module{
					BlockedModuleAccountsOverride: blockAccAddrs,
				}),
			},
			{
				Name:   stakingtypes.ModuleName,
				Config: appconfig.WrapAny(&stakingmodulev1.Module{}),
			},
			{
				Name:   slashingtypes.ModuleName,
				Config: appconfig.WrapAny(&slashingmodulev1.Module{}),
			},
			{
				Name:   "tx",
				Config: appconfig.WrapAny(&txconfigv1.Config{}),
			},
			{
				Name:   genutiltypes.ModuleName,
				Config: appconfig.WrapAny(&genutilmodulev1.Module{}),
			},
			{
				Name:   authz.ModuleName,
				Config: appconfig.WrapAny(&authzmodulev1.Module{}),
			},
			{
				Name:   upgradetypes.ModuleName,
				Config: appconfig.WrapAny(&upgrademodulev1.Module{}),
			},
			{
				Name:   distrtypes.ModuleName,
				Config: appconfig.WrapAny(&distrmodulev1.Module{}),
			},
			{
				Name:   evidencetypes.ModuleName,
				Config: appconfig.WrapAny(&evidencemodulev1.Module{}),
			},
			{
				Name:   minttypes.ModuleName,
				Config: appconfig.WrapAny(&mintmodulev1.Module{}),
			},
			{
				Name: group.ModuleName,
				Config: appconfig.WrapAny(&groupmodulev1.Module{
					MaxExecutionPeriod: durationpb.New(time.Second * 1209600),
					MaxMetadataLen:     255,
				}),
			},
			{
				Name:   nft.ModuleName,
				Config: appconfig.WrapAny(&nftmodulev1.Module{}),
			},
			{
				Name:   feegrant.ModuleName,
				Config: appconfig.WrapAny(&feegrantmodulev1.Module{}),
			},
			{
				Name:   govtypes.ModuleName,
				Config: appconfig.WrapAny(&govmodulev1.Module{}),
			},
			{
				Name:   consensustypes.ModuleName,
				Config: appconfig.WrapAny(&consensusmodulev1.Module{}),
			},
			{
				Name:   circuittypes.ModuleName,
				Config: appconfig.WrapAny(&circuitmodulev1.Module{}),
			},
			{
				Name:   paramstypes.ModuleName,
				Config: appconfig.WrapAny(&paramsmodulev1.Module{}),
			},
			{
				Name:   epochstypes.ModuleName,
				Config: appconfig.WrapAny(&epochsmodulev1.Module{}),
			},
			{
				Name:   blogmoduletypes.ModuleName,
				Config: appconfig.WrapAny(&blogmoduletypes.Module{}),
			},
			// starport 脚手架使用的占位行 # stargate/app/moduleConfig
		},
	})
)
