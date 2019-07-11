package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (t *TokenChaincode) CreateAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	logger.Info("-- CreateAccount --")
	defer logger.Info("-- End CreateAccount --")

	// sync with Transaction format
	txParts, err := parse_transaction(args[0], SEPARATOR)
	if err != nil {
		return shim.Error(err.Error())
	}

	txca := PrototypeTransaction{}

	if err := txca.Deserialize(txParts.Context); err != nil {
		return shim.Error(CombinationError(ErrDeserializeFail, err, CODE601).Error())
	}

	pk, err := parse_pem_with_der_public([]byte(txca.Pubkey))
	if err != nil {
		return shim.Error(CombinationError(ErrPemInvalidPubkey, err, CODE601).Error())
	}

	addr := make_pubkey_to_addr(pk)
	r, s, err := parse_signature(txParts.Signature)
	if err != nil {
		return shim.Error(CombinationError(ErrInvalidSignatureFormat, err, CODE704).Error())
	}

	hash := sha256.Sum256([]byte(txParts.RawContext))
	valid := ecdsa.Verify(pk, hash[:], r, s)
	if valid == false {
		return shim.Error(fmt.Sprintf("Result: %t, Invalid Sign", valid))
	}

	// add from - to - addr validation process
	if fmt.Sprintf("%x", addr) != txca.From && fmt.Sprintf("%x", addr) != txca.To {
		logger.Error(fmt.Sprintf("Sender Address Mismatch: ReqAddr: %s, PrcAddr: %s", txca.From, fmt.Sprintf("%x", addr)))
		return shim.Error(fmt.Sprintf("Sender Address Mismatch: Calculated Address - %s", fmt.Sprintf("%x", addr)))
	}

	if valid == true {
		logger.Info("-- Valid Transaction --")
		guard, guardErr := stub.CreateCompositeKey(COLLISION_GUARD_OBJECT, []string{fmt.Sprintf("%x", addr), "CREATE"})
		if guardErr != nil {
			return shim.Error(CombinationError(ErrCompositeKeyError, guardErr, CODE101).Error())
		}
		if exist, err := stub.GetState(guard); err != nil {
			return shim.Error(CombinationError(ErrGetStateFail, err, CODE603).Error())
		} else if len(exist) != 0 {
			logger.Error(ErrAccountExists.Error())
			return shim.Error(ErrAccountExists.Error())
		}

		ts, tsErr := stub.GetTxTimestamp()
		if tsErr != nil {
			return shim.Error(CombinationError(ErrInvalidTimestamp, "Could not retrieve transaction timestamp.", tsErr, CODE602).Error())
		}
		tsStr := make_timestamp_string(ts)

		tx := PrototypeAccount{
			CreatedAt:     tsStr,
			Page:          0,
			LastShrinkKey: "nil",
			ShrinkSize:    0xFF,
			Status:        "Normal",
		}
		txid := stub.GetTxID()

		if data, err := tx.Serialize(); err != nil {
			return shim.Error(CombinationError(ErrSerializeFail, err).Error())
		} else if data == nil {
			return shim.Error(CombinationError(ErrSerializeFail, ErrNodata).Error())
		} else {
			key, keyErr := stub.CreateCompositeKey(ACCOUNT_OBJECT, []string{fmt.Sprintf("%x", addr), txid})
			if keyErr != nil {
				return shim.Error(CombinationError(ErrCompositeKeyError, keyErr).Error())
			}
			if stateErr := stub.PutState(key, data); stateErr != nil {
				logger.Error(CombinationError(ErrStateUpdateFail, stateErr).Error())
				return shim.Error(CombinationError(ErrStateUpdateFail, stateErr).Error())
			}
			if guardErr := stub.PutState(guard, []byte(txid)); guardErr != nil {
				logger.Error(CombinationError(ErrStateUpdateFail, "Collision Guard Update Error", guardErr).Error())
				return shim.Error(CombinationError(ErrStateUpdateFail, "Collision Guard Update Error", guardErr).Error())
			}
		}
	}

	return shim.Success([]byte(fmt.Sprintf("Result: %t", valid)))
}

func (t *TokenChaincode) WriteAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	logger.Info("-- CreateAccount --")
	defer logger.Info("-- End CreateAccount --")

	// sync with Transaction format
	txParts, err := parse_transaction(args[0], SEPARATOR)
	if err != nil {
		return shim.Error(err.Error())
	}

	txca := PrototypeTransaction{}

	if err := txca.Deserialize(txParts.Context); err != nil {
		return shim.Error(CombinationError(ErrDeserializeFail, err, CODE601).Error())
	}

	pk_mint, err := parse_hex_public([]byte(MintPublicKey))
	if err != nil {
		return shim.Error(CombinationError(ErrPemInvalidPubkey, err, CODE601).Error())
	}

	addr := txca.Pubkey
	if len(addr) != 64 {
		return shim.Error(ErrInvalidAddrLength.Error())
	}
	r, s, err := parse_signature(txParts.Signature)
	if err != nil {
		return shim.Error(CombinationError(ErrInvalidSignatureFormat, err, CODE704).Error())
	}

	hash := sha256.Sum256([]byte(txParts.RawContext))
	valid := ecdsa.Verify(pk_mint, hash[:], r, s)
	if valid == false {
		return shim.Error(fmt.Sprintf("Result: %t, Invalid Sign", valid))
	}

	if valid == true {
		logger.Info("-- Valid Transaction --")
		guard, guardErr := stub.CreateCompositeKey(COLLISION_GUARD_OBJECT, []string{addr, "CREATE"})
		if guardErr != nil {
			return shim.Error(CombinationError(ErrCompositeKeyError, guardErr, CODE101).Error())
		}
		if exist, err := stub.GetState(guard); err != nil {
			return shim.Error(CombinationError(ErrGetStateFail, err, CODE603).Error())
		} else if len(exist) != 0 {
			logger.Error(ErrAccountExists.Error())
			return shim.Error(ErrAccountExists.Error())
		}

		ts, tsErr := stub.GetTxTimestamp()
		if tsErr != nil {
			return shim.Error(CombinationError(ErrInvalidTimestamp, "Could not retrieve transaction timestamp.", tsErr, CODE602).Error())
		}
		tsStr := make_timestamp_string(ts)

		tx := PrototypeAccount{
			CreatedAt:     tsStr,
			Page:          0,
			LastShrinkKey: "nil",
			ShrinkSize:    0xFF,
			Status:        "Normal",
		}
		txid := stub.GetTxID()

		if data, err := tx.Serialize(); err != nil {
			return shim.Error(CombinationError(ErrSerializeFail, err).Error())
		} else if data == nil {
			return shim.Error(CombinationError(ErrSerializeFail, ErrNodata).Error())
		} else {
			key, keyErr := stub.CreateCompositeKey(ACCOUNT_OBJECT, []string{addr, txid})
			if keyErr != nil {
				return shim.Error(CombinationError(ErrCompositeKeyError, keyErr).Error())
			}
			if stateErr := stub.PutState(key, data); stateErr != nil {
				logger.Error(CombinationError(ErrStateUpdateFail, stateErr).Error())
				return shim.Error(CombinationError(ErrStateUpdateFail, stateErr).Error())
			}
			if guardErr := stub.PutState(guard, []byte(txid)); guardErr != nil {
				logger.Error(CombinationError(ErrStateUpdateFail, "Collision Guard Update Error", guardErr).Error())
				return shim.Error(CombinationError(ErrStateUpdateFail, "Collision Guard Update Error", guardErr).Error())
			}
		}
	}

	return shim.Success([]byte(fmt.Sprintf("Result: %t", valid)))
}
