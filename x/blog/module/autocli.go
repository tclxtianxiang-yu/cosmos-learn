package blog

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"cosmos-learn/x/blog/types"
)

// AutoCLIOptions 实现了 autocli.HasAutoCLIConfig 接口。
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "ListPost",
					Use:       "list-post",
					Short:     "List all post",
				},
				{
					RpcMethod:      "GetPost",
					Use:            "get-post [id]",
					Short:          "Gets a post by id",
					Alias:          []string{"show-post"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				// 此行供 Ignite 脚手架使用 # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // 仅当需要启用自定义命令时才必需
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // 因为受权限控制而跳过
				},
				{
					RpcMethod:      "CreatePost",
					Use:            "create-post [title] [content]",
					Short:          "Create post",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "title"}, {ProtoField: "content"}},
				},
				{
					RpcMethod:      "UpdatePost",
					Use:            "update-post [id] [title] [content]",
					Short:          "Update post",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}, {ProtoField: "title"}, {ProtoField: "content"}},
				},
				{
					RpcMethod:      "DeletePost",
					Use:            "delete-post [id]",
					Short:          "Delete post",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				// 此行供 Ignite 脚手架使用 # autocli/tx
			},
		},
	}
}
