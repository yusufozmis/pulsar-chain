package keeper_test

import (
	"testing"

	"github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/node101-io/pulsar-chain/x/keyregistry/keeper"
	"github.com/node101-io/pulsar-chain/x/keyregistry/types"
	"github.com/stretchr/testify/require"
)

var mockCosmosSignature = "cosmosSig"
var mockMinaSignature = "minaSig"

func TestRegisterKeysFail(t *testing.T) {

	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(f.keeper)

	creatorAddr := sdk.AccAddress([]byte("pulsar"))
	_, err := ms.RegisterKeys(f.ctx, &types.MsgRegisterKeys{
		Creator:         creatorAddr.String(),
		CosmosSignature: mockCosmosSignature,
		MinaSignature:   mockMinaSignature,
		CosmosPublicKey: CosmosPubKey,
		MinaPublicKey:   MinaPubKey,
	})
	require.ErrorIs(t, err, types.ErrInvalidPublicKey)
}

func TestRegisterKeysSuccess(t *testing.T) {

	priv := secp256k1.GenPrivKey()

	pub := priv.PubKey()

	addr := sdk.AccAddress(pub.Address())

	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(f.keeper)

	resp, err := ms.RegisterKeys(f.ctx, &types.MsgRegisterKeys{
		Creator:         addr.String(),
		CosmosSignature: mockCosmosSignature,
		MinaSignature:   mockMinaSignature,
		CosmosPublicKey: pub.Bytes(),
		MinaPublicKey:   MinaPubKey,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	exists, err := f.keeper.CosmosToMina.Has(f.ctx, pub.Bytes())
	require.NoError(t, err)
	require.Equal(t, exists, true)

	exists, err = f.keeper.MinaToCosmos.Has(f.ctx, MinaPubKey)
	require.NoError(t, err)
	require.Equal(t, exists, true)

}

func TestInvalidCreatorAddress(t *testing.T) {

	priv := secp256k1.GenPrivKey()

	pub := priv.PubKey()

	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(f.keeper)

	_, err := ms.RegisterKeys(f.ctx, &types.MsgRegisterKeys{
		Creator:         "creator",
		CosmosSignature: mockCosmosSignature,
		MinaSignature:   mockMinaSignature,
		CosmosPublicKey: pub.Bytes(),
		MinaPublicKey:   MinaPubKey,
	})
	require.ErrorIs(t, err, types.ErrInvalidCreatorAddres)
}

func TestInvalidSigner(t *testing.T) {

	priv := secp256k1.GenPrivKey()

	pub := priv.PubKey()

	secondaryPriv := secp256k1.GenPrivKey()

	secondaryPublic := secondaryPriv.PubKey()

	addr := sdk.AccAddress(secondaryPublic.Address())

	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(f.keeper)

	_, err := ms.RegisterKeys(f.ctx, &types.MsgRegisterKeys{
		Creator:         addr.String(),
		CosmosSignature: mockCosmosSignature,
		MinaSignature:   mockMinaSignature,
		CosmosPublicKey: pub.Bytes(),
		MinaPublicKey:   MinaPubKey,
	})

	require.ErrorIs(t, err, types.ErrInvalidSigner)

}

// TODO: Update require.NoError to require.ErrorIs once the VerifyCosmosSig and VerifyMinaSig is implemented
func TestInvalidSignature(t *testing.T) {

	priv := secp256k1.GenPrivKey()

	pub := priv.PubKey()

	addr := sdk.AccAddress(pub.Address())

	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(f.keeper)

	invalidSig := "cosmosSig"

	_, err := ms.RegisterKeys(f.ctx, &types.MsgRegisterKeys{
		Creator:         addr.String(),
		CosmosSignature: invalidSig,
		MinaSignature:   mockMinaSignature,
		CosmosPublicKey: pub.Bytes(),
		MinaPublicKey:   MinaPubKey,
	})

	require.NoError(t, err)

}

func TestInsertSecondaryKeysFail(t *testing.T) {
	f := initFixture(t)

	priv := secp256k1.GenPrivKey()

	pub := priv.PubKey()

	addr := sdk.AccAddress(pub.Address())

	ms := keeper.NewMsgServerImpl(f.keeper)

	resp, err := ms.RegisterKeys(f.ctx, &types.MsgRegisterKeys{
		Creator:         addr.String(),
		CosmosSignature: mockCosmosSignature,
		MinaSignature:   mockMinaSignature,
		CosmosPublicKey: pub.Bytes(),
		MinaPublicKey:   MinaPubKey,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	resp, err = ms.RegisterKeys(f.ctx, &types.MsgRegisterKeys{
		Creator:         addr.String(),
		CosmosSignature: mockCosmosSignature,
		MinaSignature:   mockMinaSignature,
		CosmosPublicKey: pub.Bytes(),
		MinaPublicKey:   MinaPubKey,
	})

	require.ErrorIs(t, err, types.ErrSecondaryKeyExists)
}
