package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (t *TokenChaincode) Transaction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(ErrInvalidArgLength.Error() + ": Expecting 1")
	}

	logger.Info("-- New Transaction --")
	defer logger.Info("-- End Transaction --")

	// need transaction length check
	txParts, err := parse_transaction(args[0], SEPARATOR)

	if err != nil {
		return shim.Error(fmt.Sprintf("%s : CODE%d", err.Error(), CODE118))
	}

	tx := PrototypeTransaction{}
	if err := tx.Deserialize(txParts.Context); err != nil {
		logger.Errorf(ErrInvalidTxEncoding.Error() + ": JSON decode fail, Reason: " + err.Error())
		return shim.Error(ErrInvalidTxFormat.Error())
	}
	tx.Original = txParts.RawContext
	tx.Signature = string(txParts.Signature)

	pk, err := parse_pem_with_der_public([]byte(tx.Pubkey))
	if err != nil {
		logger.Error(CombinationError(err.Error(), CODE117).Error())
		return shim.Error(CombinationError(err.Error(), CODE117).Error())
	}

	addr := make_pubkey_to_addr(pk)

	hash := sha256.Sum256([]byte(tx.Original))
	txid := stub.GetTxID()

	logger.Warningf("tx.Original: %s\n\n, tx.Signature: %s\n\n", tx.Original, txParts.Signature)
	r, s, err := parse_signature(txParts.Signature)
	if err != nil {
		logger.Error(CombinationError(err.Error(), CODE117).Error())
		return shim.Error(CombinationError(err.Error(), CODE117).Error())
	}
	logger.Warningf("pk: %+v, hash: 0x%x, r: 0x%x, s: 0x%x", pk, hash[:], r, s)

	result := ecdsa.Verify(pk, hash[:], r, s)

	// TODO
	if tx.From != fmt.Sprintf("%x", addr) && tx.From != fmt.Sprintf("0x%x", addr[:20]) {
		logger.Error(fmt.Sprintf("Sender Address Mismatch: ReqAddr: %s, PrcAddr: %s, PrcAddr2: %s", tx.From, fmt.Sprintf("%x", addr), fmt.Sprintf("0x%x", addr[:20])))
		return shim.Error(fmt.Sprintf("Sender Address Mismatch: Calculated Address - %s", fmt.Sprintf("%x", addr)))
	}

	logger.Info("-- transaction data --",tx)
	logger.Info("-- transaction to data --",tx.To)
	//result = true // REMOVE
	if result == true {
		logger.Info("-- Verified Transaction --")
		defer logger.Info("-- Verified Transaction End--")

		// mitik, mitikErr := stub.CreateCompositeKey(COLLISION_MITIGATE_OBJECT, []string{Sig})
		if coll, mitisErr := stub.GetState(string(txParts.Signature)); mitisErr != nil { // 2^128
			logger.Error(ErrGetStateFail.Error() + ": Unable to Retrieve COLLISION_MITIGATE_OBJECT, State: " + mitisErr.Error())
			return shim.Error(ErrGetStateFail.Error() + ": Unable to Retrieve COLLISION_MITIGATE_OBJECT, State: " + mitisErr.Error())
		} else if coll != nil {
			logger.Error(ErrSignatureCollision.Error(), ": "+string(coll))
			return shim.Error(ErrSignatureCollision.Error() + ": " + string(coll))
		}
		if data, err := tx.Serialize(); err != nil {
			logger.Error(ErrSerializeFail.Error() + ": " + err.Error())
			return shim.Error(ErrSerializeFail.Error() + ": " + err.Error())
		} else if data == nil {
			logger.Error(ErrSerializeFail.Error() + ": data is nil")
			return shim.Error(ErrSerializeFail.Error() + ":data is nil")
		} else {

		
			logger.Info("-- Retrieving From Account --")

			acc := PrototypeAccount{}

			if tx.From != fmt.Sprintf("0x%x", addr[:20]) {
				guard, guardErr := stub.CreateCompositeKey(COLLISION_GUARD_OBJECT, []string{fmt.Sprintf("%x", addr), "CREATE"})
				if guardErr != nil {
					return shim.Error(CombinationError(ErrCompositeKeyError, CODE101).Error())
				}

				if txid, err := stub.GetState(guard); err != nil {
					return shim.Error(CombinationError(ErrGetStateFail, CODE105).Error())
				} else if txid == nil {
					return shim.Error(CombinationError(ErrNodata, CODE106, errors.New("Your Account Data Not Exist")).Error())
				} else {
					key, keyErr := stub.CreateCompositeKey(ACCOUNT_OBJECT, []string{fmt.Sprintf("%x", addr), string(txid)})
					if keyErr != nil {
						return shim.Error(CombinationError(ErrCompositeKeyError, keyErr).Error())
					}
					if data, err := stub.GetState(key); err != nil {
						return shim.Error(CombinationError(ErrGetStateFail, err, CODE107).Error())
					} else {
						if data == nil {
							return shim.Error(CombinationError(ErrInternalPanic).Error())
						}
						if accErr := acc.Deserialize(data); accErr != nil {
							return shim.Error(CombinationError(ErrDeserializeFail, accErr, CODE107).Error())
						}
					}
				}
			}

			logger.Info("-- Retrieving To Account --")
			toacc := PrototypeAccount{}
			toguard, toguardErr := stub.CreateCompositeKey(COLLISION_GUARD_OBJECT, []string{tx.To, "CREATE"})
			if toguardErr != nil {
				return shim.Error(CombinationError(ErrCompositeKeyError, toguardErr, CODE108).Error())
			}
			if txid, err := stub.GetState(toguard); err != nil {
				return shim.Error(CombinationError(ErrGetStateFail, err, CODE107).Error())
			} else if txid == nil {
				logger.Info("-- To Account is not exist. Transaction Forwarding... --")
				// wait policy....
				toacc.Page = 0
			} else {
				key, keyErr := stub.CreateCompositeKey(ACCOUNT_OBJECT, []string{tx.To, string(txid)})
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

			balancer := func(addr string) (uint64, uint64, error) {
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

						if addr == tx.To {
							balance += val
						} else {
							balance -= val
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

						if addr == tx.To {
							balance += val
						} else {
							balance -= val
						}
					}
				}
				return balance, i, nil
			}

			val, parseErr := strconv.ParseUint(tx.Amount, 10, 64)
			if parseErr != nil {
				logger.Error(CombinationError(ErrParseFail, parseErr, CODE206).Error())
				// Todo: do we print this?
				return shim.Error(CombinationError(ErrParseFail, parseErr, CODE206).Error())
			}

			var balance uint64
			var count uint64
			var balanceErr error
			cacheHit := false

			if cacheHit = cache.Query(fmt.Sprintf("%x", addr[:])); cacheHit {
				logger.Info("-- Cache Hit --")
				balance = cache.Determination(fmt.Sprintf("%x", addr[:]))
				count = cache.TransactionCount(fmt.Sprintf("%x", addr[:]))
				logger.Noticef("-- Cache Balance: %d", balance)
				balanceErr = nil
			}
			if cacheHit == false || balance < val {
				balance, count, balanceErr = balancer(fmt.Sprintf("%x", addr[:]))
				logger.Noticef("Hit: %t, Req(b < v): %t, Balance: %d", cacheHit, balance < val, balance)
			}

			if tx.From == fmt.Sprintf("0x%x", addr[:20]) {
				balance = t.Balance_Temp(stub, fmt.Sprintf("0x%x", addr[:20]))
				logger.Infof("VA Balanacer Run: return %d\n", balance)
			}

			if balance < val {
				logger.Noticef("Transaction Balance Not Enough: REQ: %d, CUR: %d", val, balance)
				return shim.Error(CombinationError(ErrBalanceLow, fmt.Sprintf("Request: %d, Current: %d", val, balance)).Error())
			}

			//// SHRINK
			if ((count >= acc.ShrinkSize) && (acc.ShrinkSize != 0) && cacheHit) || (count >= 255) {
				logger.Info("-- Shrink Transactions --")

				stx := PrototypeTransaction{}
				stx.Amount = strconv.FormatUint(balance, 10)
				stx.Action = "SHRINK"
				stx.To = fmt.Sprintf("%x", addr[:])
				stx.From = "TXID CALL - " + stub.GetTxID()
				stx.Original = stx.From
				stx.Pubkey = tx.Pubkey

				acc.Page++
				logger.Info("-- Retrieving From Account TXID --")

				sguard, sguardErr := stub.CreateCompositeKey(COLLISION_GUARD_OBJECT, []string{fmt.Sprintf("%x", addr[:]), "CREATE"})
				if sguardErr != nil {
					return shim.Error(CombinationError(ErrGetStateFail, sguardErr, CODE113).Error())
				}
				if txid, err := stub.GetState(sguard); err != nil {
					return shim.Error(CombinationError(ErrGetStateFail, CODE114).Error())
				} else if txid == nil {
					return shim.Error(CombinationError(ErrInternalPanic, CODE114).Error())
				} else {
					sakey, sakeyErr := stub.CreateCompositeKey(ACCOUNT_OBJECT, []string{fmt.Sprintf("%x", addr[:]), string(txid)})
					if sakeyErr != nil {
						return shim.Error(CombinationError(ErrCompositeKeyError, sakeyErr).Error())
					}
					if sadata, saErr := acc.Serialize(); saErr != nil {
						return shim.Error(CombinationError(ErrSerializeFail, CODE115).Error())
					} else {
						if sastateErr := stub.PutState(sakey, sadata); sastateErr != nil {
							return shim.Error(CombinationError(ErrStateUpdateFail, sastateErr).Error())
						}
					}
				}

				sdata, sErr := stx.Serialize()
				if sErr != nil {
					logger.Error(CombinationError(ErrSerializeFail, sErr, CODE201).Error())
					return shim.Error(CombinationError(ErrSerializeFail, sErr, CODE201).Error())
				}

				skey, skeyErr := stub.CreateCompositeKey(TRANSACTION_OBJECT, []string{tx.From, acc.PageString(), "SHRINK", txid})
				if skeyErr != nil {
					logger.Error(CombinationError(ErrCompositeKeyError, skeyErr, CODE201).Error())
					return shim.Error(CombinationError(ErrCompositeKeyError, skeyErr, CODE201).Error())
				}
				if stateErr := stub.PutState(skey, sdata); stateErr != nil {
					logger.Error(CombinationError(ErrStateUpdateFail, stateErr, CODE202).Error())
					return shim.Error(CombinationError(ErrStateUpdateFail, stateErr, CODE202).Error())
				}

			}
			//// SHRINK END

			if balanceErr == ErrNodata {
				logger.Error(CombinationError(ErrBalanceQuery, balanceErr, CODE203).Error())
				// return shim.Error(CombinationError(ErrBalanceQuery, balanceErr, CODE203).Error())
			}
			if balanceErr != nil && balanceErr != ErrNodata {
				logger.Error(CombinationError(ErrBalanceQuery, balanceErr, CODE209).Error())
				return shim.Error(CombinationError(ErrBalanceQuery, balanceErr, CODE209).Error())
			}

			if balance < val {
				logger.Noticef("Transaction Balance Not Enough: REQ: %d, CUR: %d", val, balance)
				return shim.Error(CombinationError(ErrBalanceLow, fmt.Sprintf("Request: %d, Current: %d", val, balance)).Error())
			}

			fkey, fkeyErr := stub.CreateCompositeKey(TRANSACTION_OBJECT, []string{tx.From, acc.PageString(), "SUB", txid})
			tkey, tkeyErr := stub.CreateCompositeKey(TRANSACTION_OBJECT, []string{tx.To, toacc.PageString(), "ADD", txid})
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
			if stateErr := stub.PutState(string(txParts.Signature), []byte(txid)); stateErr != nil {
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

			cTxResult := &TransactionResult{balance, val, balance - val}
			cache.Insert(fmt.Sprintf("%x", addr[:]), &acc, count, &tx, cTxResult)
			logger.Info("Transaction State Added.")
		}
		logger.Info("-- Valid Transaction End --")
	}

	resp := PrototypeResponsestruct{
		Txid:        stub.GetTxID(),
		From:        tx.From,
		To:          tx.To,
		BlockNumber: stub.GetTxID(),
		Amount:      tx.Amount,
	}

	respRaw, respRawErr := resp.Serialize()
	if respRawErr != nil {
		logger.Error(CombinationError(ErrResponseSerializeFail, respRawErr).Error())
		return shim.Error(CombinationError(ErrResponseSerializeFail, respRawErr).Error())
	}

	if result == false {
		return shim.Error(fmt.Sprintf("Result: %t, Response: %s, Invalid Sign", result, respRaw))
	}

	// JSON?
	return shim.Success([]byte(fmt.Sprintf("Result: %t, Response: %s", result, respRaw)))
}
