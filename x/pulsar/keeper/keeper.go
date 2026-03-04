package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/node101-io/pulsar-chain/x/pulsar/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema collections.Schema
	Params collections.Item[types.Params]

	CosmosToMina collections.Map[[]byte, []byte] // Cosmos PubKey --> Mina PubKey
	MinaToCosmos collections.Map[[]byte, []byte] // Mina PubKey --> Cosmos PubKey
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,

) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,

		Params: collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),

		CosmosToMina: collections.NewMap(sb, collections.NewPrefix("pulsar"), "cosmos_to_mina", collections.BytesKey, collections.BytesValue),
		MinaToCosmos: collections.NewMap(sb, collections.NewPrefix("pulsar"), "mina_to_cosmos", collections.BytesKey, collections.BytesValue),
	}
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

func (k Keeper) SetCosmosToMina(ctx context.Context, cosmosPublicKey, minaPublicKey []byte) error {
	return k.CosmosToMina.Set(ctx, cosmosPublicKey, minaPublicKey)
}

func (k Keeper) GetCosmosToMina(ctx context.Context, cosmosPublicKey []byte) ([]byte, error) {
	return k.CosmosToMina.Get(ctx, cosmosPublicKey)
}

func (k Keeper) SetMinaToCosmos(ctx context.Context, minaPublicKey, cosmosPublicKey []byte) error {
	return k.MinaToCosmos.Set(ctx, minaPublicKey, cosmosPublicKey)
}

func (k Keeper) GetMinaToCosmos(ctx context.Context, minaPublicKey []byte) ([]byte, error) {
	return k.MinaToCosmos.Get(ctx, minaPublicKey)
}
