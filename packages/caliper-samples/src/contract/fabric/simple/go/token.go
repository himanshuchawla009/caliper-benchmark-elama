package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"

	//"./secp256k1"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	//hfc "github.com/soo/sucrose/hfc_contract"
)

func (t *TokenChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	logger.Info("Invoke Run: " + function)
	cache.CacheInit()
	// middleware.MW.Init()

	switch function {
	case "Transaction":
		return t.Transaction(stub, args)
	case "Balance":
		return t.Balance(stub, args)
	case "Exchange":
		return t.Exchange(stub, args)
	case "History":
		return t.History(stub, args)
	case "CreateAccount":
		return t.CreateAccount(stub, args)
	case "Query":
		return t.Query(stub, args)
	case "TXID":
		return t.queryTXID(stub, args)
	case "WriteAccount":
		return t.WriteAccount(stub, args)
	case "Mint":
		return t.Minting(stub, args)
	case "CACHEINIT":
		return t.CacheInit(stub)
	case "CACHESTATUS":
		return t.CacheStatus(stub)
	// case "CreateHTLC":
	// 	contract := new(hfc.ContractChaincode)
	// 	return contract.CreateHTLC(stub, args)
	// case "SettlePayment":
	// 	/*
	// 		type ChaincodeInvoke func(stub shim.ChaincodeStubInterface, args []string) pb.Response
	// 		type ChaincodeQueryBalance func(stub shim.ChaincodeStubInterface, addr string) uint64
	// 	*/
	// 	contract := new(hfc.ContractChaincode)
	// 	return contract.SettlePayment(stub, args, t.Balance_Temp, t.Transaction)
	// case "Refund":
	// 	contract := new(hfc.ContractChaincode)
	// 	return contract.Refund(stub, args, t.Balance_Temp, t.Transaction)
	// case "HTLC_ACCOUNT_QUERY":
	// 	contract := new(hfc.ContractChaincode)
	// 	return contract.Query(stub, args)
	default:
		logger.Error("Unknown function: " + function)
		return shim.Error(ErrUnknownFunction.Error())
	}
}

func (t *TokenChaincode) Balance_Temp(stub shim.ChaincodeStubInterface, addr string) uint64 {
	logger.Info("-- Query Balance --")
	defer logger.Info("-- End Query Balance --")

	if len(addr) != 42 {
		return 0
	}

	balance, balanceErr := func(addr string) (uint64, error) {
		iter, iterErr := stub.GetStateByPartialCompositeKey(TRANSACTION_OBJECT, []string{addr, "0"})
		if iterErr != nil {
			logger.Error(CombinationError(ErrInvalidIterator, iterErr).Error())
			return 0, ErrNodata
		}
		defer iter.Close()

		if !iter.HasNext() {
			logger.Error("Transaction Key Exist but No Values")
			return 0, ErrNodataDisposable
		}

		var i int
		var balance uint64
		for i = 0; iter.HasNext(); i++ {
			data, nextErr := iter.Next()
			if nextErr != nil {
				logger.Error(CombinationError(ErrInvalidIterator, nextErr).Error())
				return 0, CombinationError(ErrInvalidIterator, nextErr, CODE1003)
			}

			tx := PrototypeTransactionOpt{}
			if err := tx.Deserialize(data.Value); err != nil {
				logger.Error(CombinationError(ErrDeserializeFail, err, CODE1003).Error())
				return 0, CombinationError(ErrDeserializeFail, err, CODE1003)
			}

			val, parseErr := strconv.ParseUint(tx.Amount, 10, 64)
			if parseErr != nil {
				return 0, CombinationError(ErrParseFail, parseErr)
			}

			if addr == tx.To {
				balance += val
			} else {
				balance -= val
			}
		}
		return balance, nil
	}(addr)

	if balanceErr == ErrNodataDisposable {
		return 0
	}

	if balanceErr != nil {
		return 0
	}

	return balance
}

func make_pubkey_to_addr(vk *ecdsa.PublicKey) []byte {
	addr := sha256.Sum256(append(vk.X.Bytes(), vk.Y.Bytes()...))
	return addr[:]
}

func (t *TokenChaincode) queryTXID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(CombinationError(ErrInvalidArgLength, "Expecting 1").Error())
	}

	logger.Info("-- query TXID --")
	defer logger.Info("-- End query TXID --")
	txid := args[0]
	if len(txid) != 64 {
		return shim.Error(ErrInvalidTXIDLength.Error() + ": Expecting 64 byte")
	}

	key, keyErr := stub.CreateCompositeKey(TXID_OBJECT, []string{txid, "TRANSACTION", "0"})
	if keyErr != nil {
		return shim.Error(ErrCompositeKeyError.Error())
	}

	if data, err := stub.GetState(key); err != nil {
		return shim.Error(ErrNodata.Error() + ": Cannot Retrive from ledger, Check your Txid")
	} else if data == nil {
		return shim.Error("No Data")
	} else {
		// TODO: make return format
		tx, _ := stub.GetState(string(data))
		return shim.Success([]byte(fmt.Sprintf("Result: %s", string(tx))))
	}
}

func (t *TokenChaincode) Balance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	addr := args[0] // Validation needed
	logger.Info("-- Query Balance --")
	defer logger.Info("-- End Query Balance --")

	if len(addr) != 64 {
		return shim.Error(ErrInvalidAddrLength.Error() + ": Expecting 64 byte")
	}

	balance, balanceErr := func(addr string) (uint64, error) {
		logger.Info("-- Retrieving Account --")

		acc := PrototypeAccount{}
		guard, guardErr := stub.CreateCompositeKey(COLLISION_GUARD_OBJECT, []string{addr, "CREATE"})
		if guardErr != nil {
			return 0, CombinationError(ErrCompositeKeyError, guardErr, CODE1001)
		}
		if txid, err := stub.GetState(guard); err != nil {
			return 0, CombinationError(ErrGetStateFail, CODE1001)
		} else if txid == nil {
			logger.Info("No Account. Forwarding...")
			acc.Page = 0
			// return 0, CombinationError(ErrNoAccount, CODE1001)
		} else {
			key, keyErr := stub.CreateCompositeKey(ACCOUNT_OBJECT, []string{addr, string(txid)})
			if keyErr != nil {
				return 0, CombinationError(ErrCompositeKeyError, keyErr, CODE1002)
			}
			if data, err := stub.GetState(key); err != nil {
				return 0, CombinationError(ErrGetStateFail, err, CODE1002)
			} else {
				if accErr := acc.Deserialize(data); accErr != nil {
					return 0, CombinationError(ErrDeserializeFail, accErr, CODE1002)
				}
			}
		}

		iter, iterErr := stub.GetStateByPartialCompositeKey(TRANSACTION_OBJECT, []string{addr, acc.PageString()})
		if iterErr != nil {
			logger.Error(CombinationError(ErrInvalidIterator, iterErr).Error())
			return 0, ErrNodata
		}
		defer iter.Close()

		if !iter.HasNext() {
			logger.Error("Transaction Key Exist but No Values")
			return 0, ErrNodataDisposable
		}

		var i int
		var balance uint64
		for i = 0; iter.HasNext(); i++ {
			data, nextErr := iter.Next()
			if nextErr != nil {
				logger.Error(CombinationError(ErrInvalidIterator, nextErr).Error())
				return 0, CombinationError(ErrInvalidIterator, nextErr, CODE1003)
			}

			tx := PrototypeTransactionOpt{}
			if err := tx.Deserialize(data.Value); err != nil {
				logger.Error(CombinationError(ErrDeserializeFail, err, CODE1003).Error())
				return 0, CombinationError(ErrDeserializeFail, err, CODE1003)
			}

			val, parseErr := strconv.ParseUint(tx.Amount, 10, 64)
			if parseErr != nil {
				return 0, CombinationError(ErrParseFail, parseErr)
			}

			if addr == tx.To {
				balance += val
			} else {
				balance -= val
			}
		}
		return balance, nil
	}(addr)

	if balanceErr == ErrNodataDisposable {
		return shim.Success([]byte(fmt.Sprintf("%d", 0)))
	}

	if balanceErr != nil {
		return shim.Error(CombinationError(balanceErr, CODE1000).Error())
	}

	return shim.Success([]byte(fmt.Sprintf("%d", balance)))
}

// ConfirmedBalance -- Deprecated...
func (t *TokenChaincode) ConfirmedBalance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	account := args[0]
	logger.Info("-- Retrieving Confirmed Balance Data --")

	guard, guardErr := stub.CreateCompositeKey(COLLISION_GUARD_OBJECT, []string{account, "CREATE"})
	if guardErr != nil {
		return shim.Error("Query Key Composite Error: " + guardErr.Error())
	}
	if txid, err := stub.GetState(guard); err != nil {
		return shim.Error(CombinationError(ErrGetStateFail, CODE105).Error())
	} else if txid == nil {
		return shim.Error("Internal Error: CODE106")
	} else {
		key, keyErr := stub.CreateCompositeKey(ACCOUNT_OBJECT, []string{account, string(txid)})
		if keyErr != nil {
			return shim.Error("Query Key Composite Error: " + keyErr.Error())
		}
		if data, err := stub.GetState(key); err != nil {
			return shim.Error("Internal Error: CODE107")
		} else {
			acc := PrototypeAccount{}
			if accErr := acc.Deserialize(data); accErr != nil {
				return shim.Error("Internal Error: CODE107")
			}
			// if page, parseErr := strconv.ParseUint(acc.Page, 10, 64); parseErr != nil {
			// 	return shim.Error("Not Impl")
			// }

		}
	}

	return shim.Error("No Data")
}

func (t *TokenChaincode) Query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	account := args[0]
	logger.Info("-- Retrieving Account --")
	defer logger.Info("-- End Retrieving Account --")

	guard, guardErr := stub.CreateCompositeKey(COLLISION_GUARD_OBJECT, []string{account, "CREATE"})
	if guardErr != nil {
		return shim.Error(CombinationError(ErrCompositeKeyError, guardErr).Error())
	}
	if txid, err := stub.GetState(guard); err != nil {
		return shim.Error(CombinationError(ErrGetStateFail, CODE400).Error())
	} else if txid == nil {

		return shim.Error(ErrNoAccount.Error())
	} else {
		key, keyErr := stub.CreateCompositeKey(ACCOUNT_OBJECT, []string{account, string(txid)})
		if keyErr != nil {
			return shim.Error(CombinationError(ErrCompositeKeyError, keyErr).Error())
		}
		if data, err := stub.GetState(key); err != nil {
			return shim.Error(CombinationError(ErrGetStateFail, CODE401).Error())
		} else {
			acc := PrototypeAccount{}
			if accErr := acc.Deserialize(data); accErr != nil {
				return shim.Error(CombinationError(ErrDeserializeFail, CODE402).Error())
			}
			return shim.Success(data)
		}
	}

	return shim.Error(ErrNodata.Error())
}

func (t *TokenChaincode) History(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) > 2 || len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	account := args[0]
	var page string
	if len(args) == 2 {
		page = args[1]
		if _, err := strconv.ParseUint(page, 10, 64); err != nil {
			return shim.Error(CombinationError(ErrParseFail, "page is unsigned number", CODE116).Error())
		}
	}

	logger.Info("-- Retrieving History --")
	defer logger.Info("-- End Retrieving History --")

	var lines []string

	var component []string
	if len(args) == 2 {
		component = []string{account, page}
	} else {
		component = []string{account}
	}

	iter, iterErr := stub.GetStateByPartialCompositeKey(TRANSACTION_OBJECT, component)
	if iterErr != nil {
		logger.Error(CombinationError(ErrInvalidIterator, "Transaction History Not Exist", iterErr).Error())
		return shim.Success([]byte(ErrNodata.Error()))
	}
	defer iter.Close()

	if !iter.HasNext() {
		logger.Error("Transaction Key Exist but No Values")
		return shim.Error(ErrNodata.Error())
	}

	var i int
	for i = 0; iter.HasNext(); i++ {
		data, nextErr := iter.Next()
		if nextErr != nil {
			logger.Error("Transaction Iterator Error: " + nextErr.Error())
			return shim.Error(CombinationError(ErrInvalidIterator.Error(), CODE501).Error())
		}

		_, keyParts, compositeKeyErr := stub.SplitCompositeKey(data.Key)
		if compositeKeyErr != nil {
			logger.Error(CombinationError(ErrCompositeKeyError, compositeKeyErr, CODE501).Error())
			return shim.Error(CombinationError(ErrCompositeKeyError, compositeKeyErr, CODE501).Error())
		}

		if strings.Contains(string(data.Value), "created") {
			tx := PrototypeExchange{}
			if err := tx.Deserialize(data.Value); err != nil {
				logger.Error(CombinationError(ErrDeserializeFail, err, CODE504).Error())
				return shim.Error(CombinationError(ErrDeserializeFail, err, CODE504).Error())
			}

			resp := PrototypeHistorystruct{
				Page:      keyParts[1],
				Action:    keyParts[2],
				Txid:      keyParts[3],
				To:        tx.To,
				Amount:    tx.Amount,
				Timestamp: tx.CreatedAt,
			}
			row, _ := resp.Serialize()

			lines = append(lines, string(row))
		} else {
			tx := PrototypeTransaction{}
			if err := tx.Deserialize(data.Value); err != nil {
				logger.Error(CombinationError(ErrDeserializeFail, err, CODE502).Error())
				return shim.Error(CombinationError(ErrDeserializeFail, err, CODE502).Error())
			}

			resp := PrototypeHistorystruct{
				Page:   keyParts[1],
				Action: keyParts[2],
				Txid:   keyParts[3],
				To:     tx.To,
				Amount: tx.Amount,
			}
			row, _ := resp.Serialize()

			lines = append(lines, string(row))
		}
	}

	return shim.Success([]byte("[" + strings.Join(lines, ", ") + "]"))
}
