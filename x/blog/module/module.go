package blog

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	porttypes "github.com/cosmos/ibc-go/v10/modules/core/05-port/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"cosmos-learn/x/blog/client/cli"
	"cosmos-learn/x/blog/keeper"
	"cosmos-learn/x/blog/types"
)

var (
	_ module.AppModuleBasic = (*AppModule)(nil)
	_ module.AppModule      = (*AppModule)(nil)
	_ module.HasGenesis     = (*AppModule)(nil)

	_ appmodule.AppModule       = (*AppModule)(nil)
	_ appmodule.HasBeginBlocker = (*AppModule)(nil)
	_ appmodule.HasEndBlocker   = (*AppModule)(nil)
	_ porttypes.IBCModule       = (*IBCModule)(nil)
)

// AppModule 实现了 AppModule 接口，提供模块之间所需的互相依赖方法。
type AppModule struct {
	cdc        codec.Codec
	keeper     keeper.Keeper
	authKeeper types.AuthKeeper
	bankKeeper types.BankKeeper
}

func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	authKeeper types.AuthKeeper,
	bankKeeper types.BankKeeper,
) AppModule {
	return AppModule{
		cdc:        cdc,
		keeper:     keeper,
		authKeeper: authKeeper,
		bankKeeper: bankKeeper,
	}
}

// IsAppModule 实现 appmodule.AppModule 接口的标记方法。
func (AppModule) IsAppModule() {}

// Name 以字符串形式返回模块名称。
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec 注册 Amino 编解码器。
func (AppModule) RegisterLegacyAminoCodec(*codec.LegacyAmino) {}

// RegisterGRPCGatewayRoutes 为该模块注册 gRPC 网关路由。
func (AppModule) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(clientCtx.CmdContext, mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// RegisterInterfaces 将模块的接口类型及其具体实现注册为 proto.Message。
func (AppModule) RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registrar)
}

// RegisterServices 注册 gRPC 查询服务，以响应模块特有的 gRPC 查询请求。
func (am AppModule) RegisterServices(registrar grpc.ServiceRegistrar) error {
	types.RegisterMsgServer(registrar, keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(registrar, keeper.NewQueryServerImpl(am.keeper))

	return nil
}

// DefaultGenesis 返回模块的默认 GenesisState，并序列化为 json.RawMessage。
// 默认的 GenesisState 由模块开发者定义，主要用于测试场景。
func (am AppModule) DefaultGenesis(codec.JSONCodec) json.RawMessage {
	return am.cdc.MustMarshalJSON(types.DefaultGenesis())
}

// ValidateGenesis 用于校验 json.RawMessage 形式传入的 GenesisState。
func (am AppModule) ValidateGenesis(_ codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var genState types.GenesisState
	if err := am.cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}

	return genState.Validate()
}

// InitGenesis 执行模块的创世初始化，不会返回验证人更新信息。
func (am AppModule) InitGenesis(ctx sdk.Context, _ codec.JSONCodec, gs json.RawMessage) {
	var genState types.GenesisState
	// 初始化创世状态中的全局索引。
	if err := am.cdc.UnmarshalJSON(gs, &genState); err != nil {
		panic(fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err))
	}

	if err := am.keeper.InitGenesis(ctx, genState); err != nil {
		panic(fmt.Errorf("failed to initialize %s genesis state: %w", types.ModuleName, err))
	}
}

// ExportGenesis 以原始 JSON 字节形式返回模块导出的创世状态。
func (am AppModule) ExportGenesis(ctx sdk.Context, _ codec.JSONCodec) json.RawMessage {
	genState, err := am.keeper.ExportGenesis(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to export %s genesis state: %w", types.ModuleName, err))
	}

	bz, err := am.cdc.MarshalJSON(genState)
	if err != nil {
		panic(fmt.Errorf("failed to marshal %s genesis state: %w", types.ModuleName, err))
	}

	return bz
}

// ConsensusVersion 是模块发生状态破坏性变更时的序列号。
// 每次模块引入破坏共识的变更都需要递增该值。
// 为避免错误或空版本，初始值应设为 1。
func (AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock 承载每个区块开头自动触发的逻辑。
// BeginBlock 的实现是可选的。
func (am AppModule) BeginBlock(_ context.Context) error {
	return nil
}

// EndBlock 承载每个区块结束时自动触发的逻辑。
// EndBlock 的实现是可选的。
func (am AppModule) EndBlock(_ context.Context) error {
	return nil
}

// GetTxCmd 返回模块的根交易命令。
// 这些命令会补充 AutoCLI 生成的交易命令。
func (AppModule) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}
