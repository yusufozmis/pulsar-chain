package keeper

import (
	"context"

	"github.com/node101-io/pulsar-chain/x/keyregistry/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {

	keyPairs := genState.Keys

	for _, keyPair := range keyPairs {
		err := k.CosmosToMina.Set(ctx, keyPair.CosmosKey, keyPair.MinaKey)
		if err != nil {
			return err
		}
		err = k.MinaToCosmos.Set(ctx, keyPair.MinaKey, keyPair.CosmosKey)
		if err != nil {
			return err
		}
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

	return genesis, nil
}
