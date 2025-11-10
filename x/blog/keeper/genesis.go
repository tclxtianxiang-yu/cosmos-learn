package keeper

import (
	"context"
	"errors"

	"cosmos-learn/x/blog/types"

	"cosmossdk.io/collections"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	if err := k.Port.Set(ctx, genState.PortId); err != nil {
		return err
	}
	for _, elem := range genState.PostList {
		if err := k.Post.Set(ctx, elem.Id, elem); err != nil {
			return err
		}
	}

	if err := k.PostSeq.Set(ctx, genState.PostCount); err != nil {
		return err
	}

	return k.Params.Set(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	genesis.PortId, err = k.Port.Get(ctx)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return nil, err
	}
	err = k.Post.Walk(ctx, nil, func(key uint64, elem types.Post) (bool, error) {
		genesis.PostList = append(genesis.PostList, elem)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	genesis.PostCount, err = k.PostSeq.Peek(ctx)
	if err != nil {
		return nil, err
	}

	return genesis, nil
}
