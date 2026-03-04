package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/node101-io/pulsar-chain/x/keyregistry/types"
)

// TODO: Implement Mina signature verification
func VerifyMinaSig(sig string, msg, minaPublicKey []byte) bool {
	return true
}

// TODO: Implement Cosmos signature verification
func VerifyCosmosSig(sig string, msg, cosmosPublicKey []byte) bool {
	return true
}

func (k msgServer) RegisterKeys(ctx context.Context, msg *types.MsgRegisterKeys) (*types.MsgRegisterKeysResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	minaSigValidity := VerifyMinaSig(msg.MinaSignature, msg.CosmosPublicKey, msg.MinaPublicKey)

	cosmosSigValidity := VerifyCosmosSig(msg.CosmosSignature, msg.MinaPublicKey, msg.CosmosPublicKey)

	if minaSigValidity && cosmosSigValidity {

		err := k.Keeper.CosmosToMina.Set(ctx, msg.CosmosPublicKey, msg.MinaPublicKey)
		if err != nil {
			return &types.MsgRegisterKeysResponse{}, err
		}
		err = k.Keeper.MinaToCosmos.Set(ctx, msg.MinaPublicKey, msg.CosmosPublicKey)
		if err != nil {
			return &types.MsgRegisterKeysResponse{}, err
		}
	}

	return &types.MsgRegisterKeysResponse{}, nil
}
