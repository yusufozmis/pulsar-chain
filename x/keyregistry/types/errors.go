package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/keyregistry module sentinel errors
var (
	ErrInvalidSigner        = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidSignature     = errors.Register(ModuleName, 1101, "Invalid Signature")
	ErrSecondaryKeyExists   = errors.Register(ModuleName, 1102, "Secondary key already exists")
	ErrInvalidCreatorAddres = errors.Register(ModuleName, 1103, "Invalid creator address")
	ErrInvalidPublicKey     = errors.Register(ModuleName, 1104, "Invalid Public key")
)
