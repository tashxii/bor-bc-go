package borservice

import (
	"bor-bc-go/x/borservice/internal/keeper"
	"bor-bc-go/x/borservice/internal/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

var (
	NewKeeper         = keeper.NewKeeper
	NewQuerier        = keeper.NewQuerier
	NewMsgBuyPrize    = types.NewMsgBuyPrize
	NewMsgSetPrize    = types.NewMsgSetPrize
	NewMsgDeletePrize = types.NewMsgDeletePrize
	NewWhois          = types.NewWhois
	ModuleCdc         = types.ModuleCdc
	RegisterCodec     = types.RegisterCodec
)

type (
	Keeper          = keeper.Keeper
	MsgSetPrize     = types.MsgSetPrize
	MsgBuyPrize     = types.MsgBuyPrize
	MsgDeletePrize  = types.MsgDeletePrize
	QueryResResolve = types.QueryResResolve
	QueryResNames   = types.QueryResNames
	Whois           = types.Whois
)
