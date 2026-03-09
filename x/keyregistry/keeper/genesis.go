package keeper

import (
	"context"

	"github.com/node101-io/pulsar-chain/x/keyregistry/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {

	keyPairs := genState.KeyPairs

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

	var keyPairs []*types.KeyPair

	var existenceMap = make(map[string]bool)

	cosmosIterator, err := k.CosmosToMina.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer cosmosIterator.Close()

	for cosmosIterator.Valid() {
		cosmosKey, err := cosmosIterator.Key()
		if err != nil {
			return genesis, err
		}
		minaKey, err := cosmosIterator.Value()
		if err != nil {
			return genesis, err
		}
		keyPair := &types.KeyPair{
			MinaKey:   minaKey,
			CosmosKey: cosmosKey,
		}
		keyPairs = append(keyPairs, keyPair)
		existenceMap[keyPair.String()] = true
		cosmosIterator.Next()
	}

	minaIterator, err := k.MinaToCosmos.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer minaIterator.Close()

	for minaIterator.Valid() {
		minaKey, err := minaIterator.Key()
		if err != nil {
			return genesis, err
		}
		cosmosKey, err := minaIterator.Value()
		if err != nil {
			return genesis, err
		}
		keyPair := &types.KeyPair{
			MinaKey:   minaKey,
			CosmosKey: cosmosKey,
		}
		if existenceMap[keyPair.String()] {
			minaIterator.Next()
			continue
		}
		keyPairs = append(keyPairs, keyPair)
		minaIterator.Next()
	}

	genesis.KeyPairs = keyPairs

	return genesis, nil
}
