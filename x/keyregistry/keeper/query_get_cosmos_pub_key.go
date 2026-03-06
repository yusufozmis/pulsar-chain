package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/node101-io/pulsar-chain/x/keyregistry/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) GetCosmosPubKey(ctx context.Context, req *types.QueryGetCosmosPubKeyRequest) (*types.QueryGetCosmosPubKeyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	exists, err := q.k.MinaToCosmos.Has(sdkCtx, req.MinaPubKey)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal Error")
	}

	if !exists {
		return nil, status.Error(codes.NotFound, "cosmos key not found for given mina key")
	}

	cosmosKey, err := q.k.MinaToCosmos.Get(sdkCtx, req.MinaPubKey)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetCosmosPubKeyResponse{
		CosmosPubKey: cosmosKey,
	}, nil
}
