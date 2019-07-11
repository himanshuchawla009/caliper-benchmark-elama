package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type PrototypeExchange struct {
	//PrototypeOrignalTransactionProposal
	Original  string `json:"original"`
	Signature string `json:"Signature"`
	CreatedAt string `json:"created,omitempty"`

	Action string `json:"action,omitempty"`

	From       string `json:"from,omitempty"`
	To         string `json:"to,omitempty"`
	FromAmount string `json:"fromAmount,omitempty"`
	Amount     string `json:"amount,omitempty"`
	Currency   string `json:"currency,omitempty"`
	Hash       string `json:"hash,omitempty"`

	FromExchangeRate string `json:"fromExchangeRate,omitempty"`
	ToExchangeRate   string `json:"toExchangeRate,omitempty"`
	Nonce            string `json:"nonce,omitempty"`
}

func (tx *PrototypeExchange) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "")
	if err := enc.Encode(tx); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (tx *PrototypeExchange) Deserialize(data []byte) error {
	if json.Valid(data) != true {
		return errors.New("Invalid JSON")
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(tx); err != io.EOF && err != nil {
		return err
	}
	return nil
}

func (t *TokenChaincode) Exchange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(CombinationError(ErrInvalidArgLength, "Expecting 1").Error())
	}

	rawData := args[0]
	logger.Info("-- Exchange --")
	defer logger.Info("-- End Exchange --")

	txParts, err := parse_transaction(rawData, SEPARATOR)
	if err != nil {
		return shim.Error(err.Error())
	}

	context := txParts.Context
	signature := txParts.Signature

	var ts string
	if tsObj, err := stub.GetTxTimestamp(); err != nil {
		return shim.Error(ErrInvalidTimestamp.Error())
	} else {
		ts = make_timestamp_string(tsObj)
	}

	txExchange := PrototypeExchange{}
	if err := txExchange.Deserialize(context); err != nil {
		return shim.Error(CombinationError(ErrDeserializeFail, err, CODE700).Error())
	}
	txExchange.From = MintAddr
	txExchange.Original = string(txParts.RawContext)
	txExchange.Signature = string(signature)
	txExchange.CreatedAt = ts

	hash := sha256.Sum256([]byte(txExchange.Original))

	r, s, err := parse_signature(signature)
	if err != nil {
		return shim.Error(CombinationError(ErrInvalidSignatureFormat, ErrParseFail, CODE701).Error())
	}

	pk, _ := parse_hex_public([]byte(MintPublicKey))
	valid := ecdsa.Verify(pk, hash[:], r, s)
	if valid == true {
		logger.Info("-- Verified Exchange --")
		defer logger.Info("-- Verified Exchange End--")

		if coll, mitisErr := stub.GetState(string(signature)); mitisErr != nil { // 2^128
			logger.Error(ErrGetStateFail.Error() + ": Unable to Retrieve COLLISION_MITIGATE_OBJECT, State: " + mitisErr.Error())
			return shim.Error(ErrGetStateFail.Error() + ": Unable to Retrieve COLLISION_MITIGATE_OBJECT, State: " + mitisErr.Error())
		} else if coll != nil {
			logger.Error(ErrSignatureCollision.Error(), ": "+string(coll))
			return shim.Error(ErrSignatureCollision.Error() + ": " + string(coll))
		}

		if data, err := txExchange.Serialize(); err != nil {
			logger.Error(ErrSerializeFail.Error() + ": " + err.Error())
			return shim.Error(ErrSerializeFail.Error() + ": " + err.Error())
		} else if data == nil {
			logger.Error(ErrSerializeFail.Error() + ": data is nil")
			return shim.Error(ErrSerializeFail.Error() + ":data is nil")
		} else {
			txid := stub.GetTxID()

			logger.Info("-- Retrieving Mint Account --")

			acc := PrototypeAccount{}
			guard, guardErr := stub.CreateCompositeKey(COLLISION_GUARD_OBJECT, []string{MintAddr, "CREATE"})
			if guardErr != nil {
				return shim.Error(CombinationError(ErrCompositeKeyError, CODE101).Error())
			}
			if txid, err := stub.GetState(guard); err != nil {
				return shim.Error(CombinationError(ErrGetStateFail, CODE105).Error())
			} else if txid == nil {
				return shim.Error(CombinationError(ErrNodata, CODE106, "Mint Account Data Not Exist", ErrInternalPanic).Error())
			} else {
				key, keyErr := stub.CreateCompositeKey(ACCOUNT_OBJECT, []string{MintAddr, string(txid)})
				if keyErr != nil {
					return shim.Error(CombinationError(ErrCompositeKeyError, keyErr).Error())
				}
				if data, err := stub.GetState(key); err != nil {
					return shim.Error(CombinationError(ErrGetStateFail, CODE107).Error())
				} else {
					if accErr := acc.Deserialize(data); accErr != nil {
						return shim.Error(CombinationError(ErrDeserializeFail, CODE107).Error())
					}
				}
			}

			logger.Info("-- Retrieving To Account --")
			toacc := PrototypeAccount{}
			toguard, toguardErr := stub.CreateCompositeKey(COLLISION_GUARD_OBJECT, []string{txExchange.To, "CREATE"})
			if toguardErr != nil {
				return shim.Error(CombinationError(ErrCompositeKeyError, toguardErr, CODE108).Error())
			}
			if txid, err := stub.GetState(toguard); err != nil {
				return shim.Error(CombinationError(ErrDeserializeFail, CODE107).Error())
			} else if txid == nil {
				logger.Info("-- To Account is not exist. Transaction Forwarding... --")
				// wait policy....
				toacc.Page = 0
			} else {
				key, keyErr := stub.CreateCompositeKey(ACCOUNT_OBJECT, []string{txExchange.To, string(txid)})
				if keyErr != nil {
					return shim.Error(CombinationError(ErrCompositeKeyError, keyErr, CODE108).Error())
				}
				// CODE107 Scope
				if data, err := stub.GetState(key); err != nil {
					return shim.Error(CombinationError(ErrGetStateFail, err, CODE107).Error())
				} else {
					if accErr := toacc.Deserialize(data); accErr != nil {
						return shim.Error(CombinationError(ErrDeserializeFail, accErr, CODE107).Error())
					}
				}
			}

			balance, _, balanceErr := func(addr string) (uint64, uint64, error) {
				logger.Info("-- Sum Transactions --")
				defer logger.Info("-- Sum Transactions End --")

				iter, iterErr := stub.GetStateByPartialCompositeKey(TRANSACTION_OBJECT, []string{addr, acc.PageString()})
				if iterErr != nil {
					logger.Error(CombinationError(ErrInvalidIterator, iterErr, "Unable To Retrieve Transaction History", CODE109).Error())
					return 0, 0, ErrNodata
				}
				defer iter.Close()

				if !iter.HasNext() {
					logger.Error("Empty Transaction History")
					return 0, 0, ErrNodata
				}

				var i uint64
				var balance uint64
				// var list []string
				for i = 0; iter.HasNext(); i++ {
					data, nextErr := iter.Next()
					if nextErr != nil {
						logger.Error(CombinationError(ErrInvalidNext, nextErr, "Unable To Retrieve Transaction History", CODE109).Error())
						return 0, 0, ErrNodata
					}

					if strings.Contains(string(data.Value), "created") {
						tx := PrototypeExchange{}
						if err := tx.Deserialize(data.Value); err != nil {
							logger.Error(CombinationError(ErrDeserializeFail, err).Error())
							return 0, 0, CombinationError(ErrDeserializeFail, err, CODE111)
						}

						val, parseErr := strconv.ParseUint(tx.Amount, 10, 64)
						if parseErr != nil {
							return 0, 0, CombinationError(parseErr, CODE112)
						}

						if addr == tx.From {
							balance -= val
						} else if addr == tx.To {
							balance += val
						}
					} else {
						tx := PrototypeTransaction{}
						if err := tx.Deserialize(data.Value); err != nil {
							logger.Error(CombinationError(ErrDeserializeFail, err).Error())
							return 0, 0, CombinationError(ErrDeserializeFail, err, CODE111)
						}

						val, parseErr := strconv.ParseUint(tx.Amount, 10, 64)
						if parseErr != nil {
							return 0, 0, CombinationError(parseErr, CODE112)
						}

						if addr == tx.From {
							balance -= val
						} else if addr == tx.To {
							balance += val
						}
					}
				}
				return balance, i, nil
			}(MintAddr)

			if balanceErr != nil && balanceErr != ErrNodata {
				logger.Error(CombinationError(ErrBalanceQuery, balanceErr, CODE203).Error())
				return shim.Error(CombinationError(ErrBalanceQuery, balanceErr, CODE203).Error())
			}

			val, parseErr := strconv.ParseUint(txExchange.Amount, 10, 64)
			if parseErr != nil {
				logger.Error(CombinationError(ErrParseFail, parseErr, CODE206).Error())
				// Todo: do we print this?
				return shim.Error(CombinationError(ErrParseFail, parseErr, CODE206).Error())
			}

			if balance < val {
				logger.Notice("Transaction Balance Not Enough: REQ: %d, CUR: %d", val, balance)
				return shim.Error(CombinationError(ErrBalanceLow, fmt.Sprintf("Request: %d, Current: %d", val, balance)).Error())
			}

			fkey, fkeyErr := stub.CreateCompositeKey(TRANSACTION_OBJECT, []string{MintAddr, acc.PageString(), "SUB", txid})
			tkey, tkeyErr := stub.CreateCompositeKey(TRANSACTION_OBJECT, []string{txExchange.To, toacc.PageString(), "ADD", txid})
			rkey, rkeyErr := stub.CreateCompositeKey(TXID_OBJECT, []string{stub.GetTxID(), "TRANSACTION", "0"})
			if fkeyErr != nil {
				logger.Error(CombinationError(ErrCompositeKeyError, fkeyErr).Error())
				return shim.Error(CombinationError(ErrCompositeKeyError, fkeyErr).Error())
			}
			if tkeyErr != nil {
				logger.Error(CombinationError(ErrCompositeKeyError, tkeyErr).Error())
				return shim.Error(CombinationError(ErrCompositeKeyError, tkeyErr).Error())
			}
			if rkeyErr != nil {
				logger.Error(CombinationError(ErrCompositeKeyError, rkeyErr).Error())
				return shim.Error(CombinationError(ErrCompositeKeyError, rkeyErr).Error())
			}
			if stateErr := stub.PutState(string(signature), []byte(stub.GetTxID())); stateErr != nil {
				logger.Error(CombinationError(ErrStateUpdateFail, "Mitigate State", stateErr).Error())
				return shim.Error(CombinationError(ErrStateUpdateFail, "Mitigate State", stateErr).Error())
			}
			if stateErr := stub.PutState(fkey, data); stateErr != nil {
				logger.Error(CombinationError(ErrStateUpdateFail, stateErr, CODE300).Error())
				return shim.Error(CombinationError(ErrStateUpdateFail, stateErr, CODE300).Error())
			}
			if stateErr := stub.PutState(tkey, data); stateErr != nil {
				logger.Error(CombinationError(ErrStateUpdateFail, stateErr, CODE301).Error())
				return shim.Error(CombinationError(ErrStateUpdateFail, stateErr, CODE301).Error())
			}
			// TODO: is it printable?
			if stateErr := stub.PutState(rkey, []byte(fkey)); stateErr != nil {
				logger.Error(CombinationError(ErrStateUpdateFail, "Relation Key", stateErr, CODE303).Error())
				return shim.Error(CombinationError(ErrStateUpdateFail, "Relation Key", stateErr, CODE303).Error())
			}
			logger.Info("Transaction State Added.")
		}
		logger.Info("-- Valid Transaction End --")
	}

	resp := PrototypeResponsestruct{
		Txid:        stub.GetTxID(),
		From:        MintAddr,
		To:          txExchange.To,
		BlockNumber: stub.GetTxID(),
		Amount:      txExchange.Amount,
	}

	respRaw, respRawErr := resp.Serialize()
	if respRawErr != nil {
		logger.Error(CombinationError(ErrResponseSerializeFail, respRawErr).Error())
		return shim.Error(CombinationError(ErrResponseSerializeFail, respRawErr).Error())
	}

	if valid == false {
		return shim.Error(fmt.Sprintf("Result: %t, Response: %s, Invalid Sign", valid, respRaw))
	}

	return shim.Success([]byte(fmt.Sprintf("Result: %t, Response: %s", valid, respRaw)))
}
