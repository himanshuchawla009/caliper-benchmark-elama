package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strconv"
)

type TokenChaincode struct{}

type BaseTransaction struct {
	Original  string `json:"original"`
	Signature string `json:"signature"`
	Nonce     string `json:"nonce"`
}

type PrototypeTransaction struct {
	BaseTransaction

	Action string `json:"action"`
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
	Pubkey string `json:"pubkey"`
}

// for optimized operation purpose
// TODO: need validation process
type PrototypeTransactionOpt struct {
	BaseTransaction

	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
}

type PrototypeCreateAccount struct {
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature"`
}

type PrototypeAccount struct {
	CreatedAt     string `json:"created,omitempty"`
	Status        string `json:"status,omitempty"`
	LastShrinkKey string `json:"-"`
	ShrinkSize    uint64 `json:"shrinkSize,string,omitempty"`
	Page          uint64 `json:"page,string,omitempty"`
}

type PrototypeHistorystruct struct {
	To        string `json:"to,omitempty"`
	Amount    string `json:"amount,omitempty"`
	Txid      string `json:"txid,omitempty"`
	Action    string `json:"action,omitempty"`
	Page      string `json:"page,omitempty"`
	Timestamp string `json:"Timestamp,omitempty"`
}

type PrototypeResponsestruct struct {
	Txid        string `json:"txid,omitempty"`
	From        string `json:"from,omitempty"`
	To          string `json:"to,omitempty"`
	BlockNumber string `json:"blockNumber,omitempty"`
	Amount      string `json:"amount,omitempty"`
}

var SEPARATOR string = ".ELAMA."
var (
	TRANSACTION_OBJECT        string = "varACCOUNT~PAGE~OPERATION~TXID"
	TXID_OBJECT               string = "varTXID~OBJECTCTYPE~REFCOUNT"
	ACCOUNT_OBJECT            string = "varACCOUNT~TXID"
	COLLISION_MITIGATE_OBJECT string = "varSIGNATURE"
	COLLISION_GUARD_OBJECT    string = "varACCOUNT~TYPE"
)

func (acc PrototypeAccount) PageString() string {
	return strconv.FormatUint(acc.Page, 10)
}

func (tx *PrototypeHistorystruct) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "")
	if err := enc.Encode(tx); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (tx *PrototypeResponsestruct) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "")
	if err := enc.Encode(tx); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (tx *PrototypeTransaction) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "")
	if err := enc.Encode(tx); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (tx *PrototypeTransaction) Deserialize(data []byte) error {
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

func (tx *PrototypeAccount) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "")
	if err := enc.Encode(tx); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (tx *PrototypeAccount) Deserialize(data []byte) error {
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

func (tx *PrototypeTransactionOpt) Deserialize(data []byte) error {
	if json.Valid(data) != true {
		return errors.New("Invalid JSON")
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(tx); err != io.EOF && err != nil {
		return err
	}
	return nil
}
