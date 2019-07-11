package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("TokenChaincode")

// Mint Account Settings
const (
	MintAddr      = "586f1462d8cba7572d842002e0bcf63f057d8a6c0d42274d40b74b9ce323cdd7"

	MintPublicKey = "40da80c28a2248e4d07bb5c4829cbdae7551d8fd62b57881a293e07350c2966c9359a2306b42bc62840ff7fb0e0004dfd4d9bccfbc2a742c17f3e222a302a4c4"
)


func main() {
	logger.SetLevel(shim.LogInfo)
	cache.CacheInit()
	err := shim.Start(new(TokenChaincode))
	if err != nil {
		logger.Noticef("[!] Cache Init!: %t", cache.L != nil)
		logger.Noticef("[!] Error starting TokenChaincode: %s", err)
	}
}

	func (t *TokenChaincode) CacheInit(stub shim.ChaincodeStubInterface) pb.Response {
		cache.CacheInit()
		return shim.Success([]byte("Cache Init Success"))
	}

func (t *TokenChaincode) CacheStatus(stub shim.ChaincodeStubInterface) pb.Response {
	if cache.L != nil {
		cache.L.RLock()
		defer cache.L.RUnlock()
		  /// this line was giving error while testing so we have fixed it plz do review
	 	//ret := fmt.Sprintf("[Cache Use: %t]Cache Usage: %d, Cache Data: %+v", len(cache.Mat), cache.Mat, !cache.Disabled)

		//added by quillhash
		ret := fmt.Sprintf("[Cache Use: %+v]Cache Usage: %T, Cache Data: %+v", len(cache.Mat), cache.Mat, !cache.Disabled)
		return shim.Success([]byte(ret))
	}
	return shim.Error("Cache Not Initalized.")
}

func (t *TokenChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	cache.CacheInit()
	return shim.Success(nil)
}

func (t *TokenChaincode) Minting(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(ErrInvalidArgLength.Error() + ": Expecting 1")
	}

	logger.Info("-- New Minting --")
	defer logger.Info("-- End Minting --")

	// need transaction length check
	txParts, err := parse_transaction(args[0], SEPARATOR)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s : CODE%d", err.Error(), CODE110))
	}

	if coll, mitisErr := stub.GetState(string(txParts.Signature)); mitisErr != nil {
		return shim.Error(ErrGetStateFail.Error() + ": Unable to Retrieve COLLISION_MITIGATE_OBJECT, State: " + mitisErr.Error())
	} else if coll != nil {
		return shim.Error(ErrSignatureCollision.Error() + ": " + string(coll))
	}

	tx := PrototypeTransaction{}
	if err := tx.Deserialize(txParts.Context); err != nil {
		return shim.Error(ErrInvalidTxFormat.Error())
	}
	tx.Original = txParts.RawContext
	tx.Signature = string(txParts.Signature)

	hash := sha256.Sum256([]byte(tx.Original))

	r, s, err := parse_signature(txParts.Signature)
	if err != nil {
		return shim.Error(CombinationError(err.Error(), CODE117).Error())
	}

	pk, _ := parse_hex_public([]byte(MintPublicKey))
	result := ecdsa.Verify(pk, hash[:], r, s)
  
	if result {
		logger.Info("-- Retrieving To Account --")
		toacc := PrototypeAccount{}
		toguard, toguardErr := stub.CreateCompositeKey(COLLISION_GUARD_OBJECT, []string{MintAddr, "CREATE"})
		if toguardErr != nil {
			return shim.Error(CombinationError(ErrCompositeKeyError, toguardErr, CODE100).Error())
		}
		if txid, err := stub.GetState(toguard); err != nil {
			return shim.Error(CombinationError(ErrGetStateFail, err, CODE107).Error())
		} else if txid == nil {
			logger.Info("-- To Account is not exist. Transaction Forwarding... --")
			// wait policy....
			toacc.Page = 0
		} else {
			key, keyErr := stub.CreateCompositeKey(ACCOUNT_OBJECT, []string{MintAddr, string(txid)})
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

		logger.Noticef("Mint Account: %+v", toacc)

		tx.To = MintAddr
		tx.From = "MINTING"
		tx.Action = "MINT"

		data, _ := tx.Serialize()
		key, _ := stub.CreateCompositeKey(TRANSACTION_OBJECT, []string{MintAddr, toacc.PageString(), "ADD", stub.GetTxID()})
		if err := stub.PutState(key, data); err != nil {
			return shim.Error(err.Error())
		}

		if stateErr := stub.PutState(string(txParts.Signature), []byte(stub.GetTxID())); stateErr != nil {
			return shim.Error(CombinationError(ErrStateUpdateFail, "Mitigate State", stateErr).Error())
		}
	}
	return shim.Success([]byte(fmt.Sprintf("Minting Done. - [%s]", tx.Amount)))
}
