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
	br rest.BaseReq, msgs []sdk.Msg, req setPrizeReq) {
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
	signedStdTx, err := signTx(cdc, transactionStdTx, req.Account, req.Account, req.Passphrase, req.Sequence, req.AccountNumber)
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
func signTx(cdc *amino.Codec, transactionStdTx types.StdTx, key, account, passphrase string, sequence, accountNumber int64) (signedStdTx types.StdTx, err error) {
	fmt.Println("@@@@@ Start signTx")
	offline := true
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	viper.Set(flags.FlagFrom, account)
	viper.Set(flags.FlagSequence, sequence)
	viper.Set(flags.FlagAccountNumber, accountNumber)
	fmt.Printf("viper.From=%s\n", viper.GetString(flags.FlagFrom))
	fmt.Printf("viper.Sequence=%s\n", viper.GetInt64(flags.FlagSequence))
	fmt.Printf("viper.AccountNumber=%s\n", viper.GetInt64(flags.FlagAccountNumber))
	txBldr := types.NewTxBuilderFromCLI()

	if viper.GetBool(flagValidateSigs) {
		if !printAndValidateSigs(cliCtx, txBldr.ChainID(), transactionStdTx, offline) {
			fmt.Println("@@@@@ Return signTx printAndValidateSigs")
			return signedStdTx, fmt.Errorf("signatures validation failed")
		}
		fmt.Println("@@@@@ Return signTx flagValidateSigs")
		return
	}

	// if --signature-only is on, then override --append
	generateSignatureOnly := viper.GetBool(flagSigOnly)
	multisigAddrStr := viper.GetString(flagMultisig)

	if multisigAddrStr != "" {
		var multisigAddr sdk.AccAddress

		multisigAddr, err = sdk.AccAddressFromBech32(multisigAddrStr)
		if err != nil {
			fmt.Println("@@@@@ Return signTx AccAddressFromBech32")
			return
		}

		signedStdTx, err = utils.SignStdTxWithSignerAddress(
			txBldr, cliCtx, multisigAddr, cliCtx.GetFromName(), transactionStdTx, offline,
		)
		generateSignatureOnly = true
	} else {
		appendSig := viper.GetBool(flagAppend) && !generateSignatureOnly
		signedStdTx, err = SignStdTxAndGivingPassphrase(
			txBldr, cliCtx, key,
			transactionStdTx, appendSig, offline,
			passphrase)
	}
	fmt.Printf("@@@@@ End signTx %v\n", err)
	return
}

func broadcastTx(cdc *amino.Codec, signedStdTx types.StdTx) (err error) {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	cliCtx.BroadcastMode = flags.BroadcastAsync
	txBytes, err := cliCtx.Codec.MarshalBinaryLengthPrefixed(signedStdTx)
	if err != nil {
		fmt.Printf("@@@@@ Return broadcastTx %v\n", err)
		return
	}

	_, err = cliCtx.BroadcastTx(txBytes)
	fmt.Printf("@@@@@ End broadcastTx %v\n", err)
	return
}
