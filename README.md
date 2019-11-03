# What is this?
This is a REST API server developed during DeFi Hackathon 2019
- Writen by Go lang
- Cosmos SDK is used
- Can store the honor prize numbers for our system players in bitcoin.

# How did we implemenet? 
- Implement a new REST API Server with using Cosmos SDK.
- Add new sign and broadcast for REST API server.

# What are new points?
- Building new REST API server within 2 days.
  - Code-base is "nameservice" in the following;
    - https://github.com/cosmos/sdk-tutorials
    - Details: https://github.com/cosmos/sdk-tutorials/blob/master/nameservice/README.md
-   Extend new two functions for REST API SDK, they don't exist ever.
    - Sign Transaciton
    - Broadcast Transaction 

## Our purpose
- Provide "one click" experience to send "prize" via blockchain.

## Issue
- Cosmos REST API SDK provides "Create new transcion" function, but they don"t have "Sing" and "Broadcast" functions.
  - We need to use Cosmos CLI SDK to execute sign and broadcast.
  - Even with using CLI, user input is required in the promot.

## How we solved
- Investigate Cosmos CLI API to execute sing and broadcast transction.
- Skip to show the promt to user for requiring input passphrase.
- Extend the REST API request to be input "passphrase" and any other information to "one click"
 
## Codes to extend Cosoms SDK
### Additionnal Files
```
  x/borserivcerest/
 +   extend_broadcast.go
 +   extend_rest.go
 +   extend_tx.go
 +   extend_tx_sign.go

 *   tx.go
```

### tx.go
Change 'setPrizeHandler" function to call hacked our "WriteGenerateStdTxResponse" function instead.
This is entry point function to sign and broadcast and prevent to show prompt.
```
// Original code
// utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
// Hacked code:
	WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg}, req)   
```

### "WriteGenerateSiteTxResponse" function
We added to call two hacked functions.
- One is "signTx" function for signing without inputing passphrase
- Second is "broadcastTx" function for broadcast transaction

```
    //
    // Same code of utils.WriteGenerateSiteTxResponse from above
    //


	// Create cdc
	cdc := cliCtx.Codec

	// Sign Tx
	signedStdTx, err := signTx(cdc, transactionStdTx, req.Account, req.Account, req.Passphrase, req.Sequence, req.AccountNumber)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Broadcast Tx
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

```
