package main

import (
	"crypto/ecdsa"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

type TransactionParts struct {
	Context    []byte
	Signature  []byte
	RawContext string
}

// return TransactionParts (decoded context, signature, rawContext), error objects
func parse_transaction(data string, sep string) (*TransactionParts, error) {
	pos := strings.Index(data, sep)
	if pos == -1 {
		return nil, ErrTxNoSeparator
	}

	context, err := base64.StdEncoding.DecodeString(data[:pos])
	if err != nil {
		return nil, errors.New(ErrInvalidTxEncoding.Error() + err.Error())
	}

	return &TransactionParts{
		Context:    context,
		Signature:  []byte(data[pos+len(sep):]),
		RawContext: data[:pos],
	}, nil
}

type InterfaceTimestamp interface {
	GetNanos() int32
	GetSeconds() int64
}

func make_timestamp_string(t InterfaceTimestamp) string {
	return fmt.Sprintf("%d.%d", t.GetSeconds(), t.GetNanos())
}

func parse_pem_public(publicKey []byte) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, ErrPemInvalidPubkey
	}

	return parse_hex_public(block.Bytes)
}

// pass 65 byte -> [0]: type, [1:33]: x. [33:]: y
func parse_pem_with_btcec_public_uncompressed(publicKey []byte) (*ecdsa.PublicKey, error) {
	block, rest := pem.Decode(publicKey)
	fmt.Println(string(rest))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, ErrPemInvalidPubkey
	}

	if block.Bytes[0] != 0x04 {
		return nil, ErrPubkeyTypeNotMatch
	} else if len(block.Bytes) != 65 {
		return nil, ErrInvalidPublicKey
	}

	x, y := &big.Int{}, &big.Int{}
	x.SetBytes(block.Bytes[1:33])
	y.SetBytes(block.Bytes[33:])

	curve := S256()
	if curve.IsOnCurve(x, y) != true {
		return nil, ErrPointsNotOnCurve
	}

	pub := &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}
	return pub, nil

}
func parse_pem_with_der_public(publicKey []byte) (*ecdsa.PublicKey, error) {
	// TODO: need change tx key parser, flow switch
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, ErrPemInvalidPubkey
	} else if len(block.Bytes) == 65 {
		return parse_pem_with_btcec_public_uncompressed(publicKey)
	}

	pub, err := ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.New(ErrDERInvalidPubkey.Error() + ": " + err.Error())
	}

	switch pub := pub.(type) {
	case *ecdsa.PublicKey:
		return pub, nil
	}

	return nil, ErrNotPubkey
}

func parse_hex_public(src []byte) (*ecdsa.PublicKey, error) {
	// if hex.DecodedLen(len(src)) != 64 { // pubkey length is 64byte
	// 	logger.Errorf("Key Length Invalid: cur-%d, req-%d", len(src), 64)
	// 	return nil, ErrInvalidKeyLength
	// }

	x, y := big.Int{}, big.Int{}
	dst := make([]byte, hex.DecodedLen(len(src)))

	n, err := hex.Decode(dst, src)
	if err != nil {
		return nil, ErrInvalidPublicKey
	}

	x.SetBytes(dst[:n/2])
	y.SetBytes(dst[n/2:])

	return &ecdsa.PublicKey{S256(), &x, &y}, nil
}

func parse_signature(src []byte) (*big.Int, *big.Int, error) {
	if hex.DecodedLen(len(src)) != 64 {
		r, s, err := parse_der_signature(src)
		if err != nil {
			return parse_hex_signature(src)
		}
		return r, s, nil
	}
	return parse_hex_signature(src)
}

func parse_hex_signature(src []byte) (*big.Int, *big.Int, error) {
	r, s := big.Int{}, big.Int{}
	dst := make([]byte, hex.DecodedLen(len(src)))

	n, err := hex.Decode(dst, src)
	if err != nil {
		return nil, nil, ErrInvalidSignatureFormat
	}

	r.SetBytes(dst[:n/2])
	s.SetBytes(dst[n/2:])

	return &r, &s, nil
}

type DerSignature struct {
	R *big.Int
	S *big.Int
}

func parse_der_signature(src []byte) (*big.Int, *big.Int, error) {
	dst := make([]byte, hex.DecodedLen(len(src)))
	n, err := hex.Decode(dst, src)
	if err != nil {
		return nil, nil, ErrInvalidSignatureFormat
	}

	signature := DerSignature{}
	rest, err := asn1.Unmarshal(dst[:n], &signature)
	if err != nil {
		return nil, nil, ErrInvalidSignatureFormat
	} else if len(rest) != 0 {
		return nil, nil, ErrInvalidSignatureFormat
	}

	return signature.R, signature.S, nil
}
