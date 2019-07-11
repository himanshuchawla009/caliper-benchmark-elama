package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrUnknown                = errors.New("Unknown")
	ErrInternalPanic          = errors.New("Internal Error")
	ErrUnknownFunction        = errors.New("Unknown Function or Command")
	ErrPemInvalidPubkey       = errors.New("Failed to decode PEM block containing public key")
	ErrDERInvalidPubkey       = errors.New("Failed to parse DER encoded public key")
	ErrNotPubkey              = errors.New("Failed to parse DER encoded public key")
	ErrInvalidArgLength       = errors.New("Incorrect number of arguments")
	ErrInvalidTXIDLength      = errors.New("Invalid TXID Length")
	ErrInvalidAddrLength      = errors.New("Invalid Address Length")
	ErrInvalidTxEncoding      = errors.New("Invalid Transaction: Encoding Error")
	ErrInvalidTxFormat        = errors.New("Invalid Transaction Format")
	ErrInvalidSignatureFormat = errors.New("Invalid Signature Format")
	ErrInvalidPublicKey       = errors.New("Invalid Public Key")
	ErrSignatureCollision     = errors.New("Signature Collision")
	ErrTxNoSeparator          = errors.New("Invalid Transaction: Separator not found")
	ErrGetStateFail           = errors.New("GetState Fail")
	ErrStateUpdateFail        = errors.New("Data Update Fail")
	ErrCompositeKeyError      = errors.New("Create Composite Key Error")
	ErrNodata                 = errors.New("No Data")
	ErrNodataDisposable       = errors.New("No Data, Forward")
	ErrSerializeFail          = errors.New("Serialize Error")
	ErrDeserializeFail        = errors.New("Deserialize Error")
	ErrInvalidIterator        = errors.New("Iterator Error")
	ErrInvalidNext            = errors.New("Iterator Next Error")
	ErrBalanceQuery           = errors.New("Balance Query Error")
	ErrParseFail              = errors.New("Parse Error")
	ErrBalanceLow             = errors.New("Not Enough Balance")
	ErrNoAccount              = errors.New("No Account")
	ErrAccountExists          = errors.New("Account Already Exists")
	ErrResponseSerializeFail  = errors.New("Response Generation Error")
	ErrInvalidTimestamp       = errors.New("Transaction Timestamp Error")
	ErrInvalidKeyLength       = errors.New("Invalid Public Key Length")
	ErrPubkeyTypeNotMatch     = errors.New("Uncompressed Public Key Type Invalid")
	ErrPointsNotOnCurve       = errors.New("X, Y not on lies bitcurve")

	ErrInvalidMsgLen     = errors.New("invalid message length, need 32 bytes")
	ErrInvalidRecoveryID = errors.New("invalid signature recovery id")
	ErrInvalidKey        = errors.New("invalid private key")
	ErrRecoverFailed     = errors.New("recovery failed")
)

type CompositeError struct {
	msg string
}

func NewCompositeError(text string) CompositeError {
	return CompositeError{
		msg: text,
	}
}

func (err CompositeError) Error() string {
	return err.msg
}

func (err CompositeError) ToString() string {
	return err.msg
}

func CombinationError(any ...interface{}) error {
	var lines []string
	for _, v := range any {
		switch v.(type) {
		case error:
			lines = append(lines, v.(error).Error())
		case string:
			lines = append(lines, v.(string))
		case int:
			lines = append(lines, fmt.Sprintf("CODE%d", v.(int)))
			//... etc
		default:
			lines = append(lines, fmt.Sprintf("Error [Type: %s, Contents: %+v]", reflect.TypeOf(v), v))
		}
	}

	return NewCompositeError(strings.Join(lines, " - "))
}

// pre-transaction
const (
	CODE100 int = iota + 100
	CODE101
	CODE102
	CODE103
	CODE104
	CODE105
	CODE106
	CODE107
	CODE108
	CODE109
	CODE110
	CODE111
	CODE112
	CODE113
	CODE114
	CODE115
	CODE116
	CODE117
	CODE118
)

// after shrink
const (
	CODE200 int = iota + 200
	CODE201
	CODE202
	CODE203
	CODE204
	CODE205
	CODE206
	CODE207
	CODE208
	CODE209
)

// state update
const (
	CODE300 int = iota + 300
	CODE301
	CODE302
	CODE303
)

// query
const (
	CODE400 int = iota + 400
	CODE401
	CODE402
	CODE403
)

// History
const (
	CODE500 int = iota + 500
	CODE501
	CODE502
	CODE503
	CODE504
)

// CreateAccount
const (
	CODE600 int = iota + 600
	CODE601
	CODE602
	CODE603
)

// Exchange
const (
	CODE700 int = iota + 700
	CODE701
	CODE702
	CODE703
	CODE704
)

// balance update
const (
	CODE1000 int = iota + 1000
	CODE1001
	CODE1002
	CODE1003
	CODE1004
	CODE1005
	CODE1006
)
