package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is the module name router key
const RouterKey = ModuleName // this was defined in your key.go file

// MsgSetPrize defines a SetName message
type MsgSetPrize struct {
	Name  string         `json:"name"`
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
}

// NewMsgSetPrize is a constructor function for MsgSetPrize
func NewMsgSetPrize(name string, value string, owner sdk.AccAddress) MsgSetPrize {
	return MsgSetPrize{
		Name:  name,
		Value: value,
		Owner: owner,
	}
}

// Route should return the name of the module
func (msg MsgSetPrize) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetPrize) Type() string { return "set_prize" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetPrize) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Name) == 0 || len(msg.Value) == 0 {
		return sdk.ErrUnknownRequest("Prize and/or Value cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetPrize) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetPrize) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgBuyPrize defines the BuyName message
type MsgBuyPrize struct {
	Name  string         `json:"name"`
	Bid   sdk.Coins      `json:"bid"`
	Buyer sdk.AccAddress `json:"buyer"`
}

// NewMsgBuyPrize is the constructor function for MsgBuyPrize
func NewMsgBuyPrize(name string, bid sdk.Coins, buyer sdk.AccAddress) MsgBuyPrize {
	return MsgBuyPrize{
		Name:  name,
		Bid:   bid,
		Buyer: buyer,
	}
}

// Route should return the name of the module
func (msg MsgBuyPrize) Route() string { return RouterKey }

// Type should return the action
func (msg MsgBuyPrize) Type() string { return "buy_name" }

// ValidateBasic runs stateless checks on the message
func (msg MsgBuyPrize) ValidateBasic() sdk.Error {
	if msg.Buyer.Empty() {
		return sdk.ErrInvalidAddress(msg.Buyer.String())
	}
	if len(msg.Name) == 0 {
		return sdk.ErrUnknownRequest("Prize cannot be empty")
	}
	if !msg.Bid.IsAllPositive() {
		return sdk.ErrInsufficientCoins("Bids must be positive")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgBuyPrize) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgBuyPrize) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}

// MsgDeletePrize defines a DeleteName message
type MsgDeletePrize struct {
	Name  string         `json:"name"`
	Owner sdk.AccAddress `json:"owner"`
}

// NewMsgDeletePrize is a constructor function for MsgDeletePrize
func NewMsgDeletePrize(name string, owner sdk.AccAddress) MsgDeletePrize {
	return MsgDeletePrize{
		Name:  name,
		Owner: owner,
	}
}

// Route should return the name of the module
func (msg MsgDeletePrize) Route() string { return RouterKey }

// Type should return the action
func (msg MsgDeletePrize) Type() string { return "delete_name" }

// ValidateBasic runs stateless checks on the message
func (msg MsgDeletePrize) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Name) == 0 {
		return sdk.ErrUnknownRequest("Prize cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgDeletePrize) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgDeletePrize) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
