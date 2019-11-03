package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
)

// WriteGenerateStdTxResponse writes response for the generate only mode.
func WriteGenerateStdTxResponse(w http.ResponseWriter, cliCtx context.CLIContext,
	br rest.BaseReq, msgs []sdk.Msg, account, passphrase string) {
	gasAdj, ok := rest.ParseFloat64OrReturnBadRequest(w, br.GasAdjustment, flags.DefaultGasAdjustment)
	if !ok {
		return
	}

	simAndExec, gas, err := flags.ParseGas(br.Gas)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	txBldr := types.NewTxBuilder(
		utils.GetTxEncoder(cliCtx.Codec), br.AccountNumber, br.Sequence, gas, gasAdj,
		br.Simulate, br.ChainID, br.Memo, br.Fees, br.GasPrices,
	)

	if br.Simulate || simAndExec {
		if gasAdj < 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, errInvalidGasAdjustment.Error())
			return
		}

		txBldr, err = utils.EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if br.Simulate {
			rest.WriteSimulationResponse(w, cliCtx.Codec, txBldr.Gas())
			return
		}
	}

	stdMsg, err := txBldr.BuildSignMsg(msgs)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	transactionStdTx := types.NewStdTx(stdMsg.Msgs, stdMsg.Fee, nil, stdMsg.Memo)

	// Create cdc
	cdc := cliCtx.Codec

	// Sign Tx
	signedStdTx, err := signTx(cdc, transactionStdTx, account, passphrase)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// TODO Broadcast Tx
	err = broadcastTx(cdc, signedStdTx)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	output, err := cliCtx.Codec.MarshalJSON(signedStdTx)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(output); err != nil {
		log.Printf("could not write response: %v", err)
	}

	return
}

// unsignedTx.json --from jack(done) --offline(done) --chain-id prizechain
// --sequence 1 --account-number 0 > signedTx.json
func signTx(cdc *amino.Codec, transactionStdTx types.StdTx, account, passphrase string) (signedStdTx types.StdTx, err error) {
	offline := true
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	viper.Set(flags.FlagSequence, 1)
	viper.Set(flags.FlagAccountNumber, 0)
	txBldr := types.NewTxBuilderFromCLI()

	if viper.GetBool(flagValidateSigs) {
		if !printAndValidateSigs(cliCtx, txBldr.ChainID(), transactionStdTx, offline) {
			return signedStdTx, fmt.Errorf("signatures validation failed")
		}

		return
	}

	// if --signature-only is on, then override --append
	generateSignatureOnly := viper.GetBool(flagSigOnly)
	multisigAddrStr := viper.GetString(flagMultisig)

	if multisigAddrStr != "" {
		var multisigAddr sdk.AccAddress

		multisigAddr, err = sdk.AccAddressFromBech32(multisigAddrStr)
		if err != nil {
			return
		}

		signedStdTx, err = utils.SignStdTxWithSignerAddress(
			txBldr, cliCtx, multisigAddr, cliCtx.GetFromName(), transactionStdTx, offline,
		)
		generateSignatureOnly = true
	} else {
		appendSig := viper.GetBool(flagAppend) && !generateSignatureOnly
		signedStdTx, err = SignStdTxAndGivingPassphrase(txBldr, cliCtx, cliCtx.GetFromName(), transactionStdTx, appendSig, offline, passphrase)
	}
	return
}

func broadcastTx(cdc *amino.Codec, signedStdTx types.StdTx) (err error) {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	txBytes, err := cliCtx.Codec.MarshalBinaryLengthPrefixed(signedStdTx)
	if err != nil {
		return
	}

	_, err = cliCtx.BroadcastTx(txBytes)
	return
}
