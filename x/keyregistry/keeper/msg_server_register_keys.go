package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func DeriveAddressFromPubkey(cosmosPublicKey []byte) string {
	pubKey := secp256k1.PubKey{
		Key: cosmosPublicKey,
	}
	addr := sdk.AccAddress(pubKey.Address())
	return addr.String()
}

func (k msgServer) RegisterKeys(ctx context.Context, msg *types.MsgRegisterKeys) (*types.MsgRegisterKeysResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidCreatorAddres, "")
	}

	if len(msg.CosmosPublicKey) != secp256k1.PubKeySize {
		return nil, errorsmod.Wrap(types.ErrInvalidPublicKey, "cosmos pubkey must be compressed (33 bytes)")
	}

	derivedAddress := DeriveAddressFromPubkey(msg.CosmosPublicKey)

	if derivedAddress != msg.Creator {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "creator does not match provided cosmos public key")
	}

	cosmosKeyExists, err := k.Keeper.CosmosToMina.Has(ctx, msg.CosmosPublicKey)

	minaKeyExists, err := k.Keeper.MinaToCosmos.Has(ctx, msg.MinaPublicKey)

	if cosmosKeyExists || minaKeyExists {
		return nil, errorsmod.Wrap(types.ErrSecondaryKeyExists, "")
	}

	minaSigValidity := VerifyMinaSig(msg.MinaSignature, msg.CosmosPublicKey, msg.MinaPublicKey)

	cosmosSigValidity := VerifyCosmosSig(msg.CosmosSignature, msg.MinaPublicKey, msg.CosmosPublicKey)

	if !minaSigValidity || !cosmosSigValidity {
		return nil, errorsmod.Wrap(types.ErrInvalidSignature, "invalid cosmos or mina signature")
	}

	err = k.Keeper.CosmosToMina.Set(ctx, msg.CosmosPublicKey, msg.MinaPublicKey)
	if err != nil {
		return nil, err
	}
	err = k.Keeper.MinaToCosmos.Set(ctx, msg.MinaPublicKey, msg.CosmosPublicKey)
	if err != nil {
		return nil, err
	}

	return &types.MsgRegisterKeysResponse{}, nil
}
