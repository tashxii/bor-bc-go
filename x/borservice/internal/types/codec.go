package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	// TODO: nameservice -> borservice
	// TODO: SetName, BuyName, DeleteName -> SetPrize, BuyPrize, DeletePrize
	cdc.RegisterConcrete(MsgSetPrize{}, "borservice/SetPrize", nil)
	cdc.RegisterConcrete(MsgBuyPrize{}, "borservice/BuyPrize", nil)
	cdc.RegisterConcrete(MsgDeletePrize{}, "borservice/DeletePrize", nil)
}
