package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/node101-io/pulsar-chain/x/keyregistry/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) GetMinaPubKey(ctx context.Context, req *types.QueryGetMinaPubKeyRequest) (*types.QueryGetMinaPubKeyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	exists, err := q.k.CosmosToMina.Has(sdkCtx, req.CosmosPubKey)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal Error")
	}

	if !exists {
		return nil, status.Error(codes.NotFound, "mina key not found for given cosmos key")
	}

	minaKey, err := q.k.CosmosToMina.Get(sdkCtx, req.CosmosPubKey)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetMinaPubKeyResponse{
		MinaPubKey: minaKey,
	}, nil
}
