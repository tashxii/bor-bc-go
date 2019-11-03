package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultCodespace is the Module Name
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodePrizeDoesNotExist sdk.CodeType = 101
)

// ErrPrizeDoesNotExist is the error for name not existing
func ErrPrizeDoesNotExist(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodePrizeDoesNotExist, "Prize does not exist")
}
