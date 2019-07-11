

package main

import (
  "crypto/ecdsa"
  "bufio"
  "crypto/elliptic"
  // "crypto/md5"
  "crypto/x509"
  "encoding/asn1"
  
  "fmt"
  // "hash"
  "encoding/hex"
  "encoding/pem"
  // "io"
  "math/big"
  "os"
  "errors"
)

const ecPrivKeyVersion = 1
var oidNamedCurveS256 = asn1.ObjectIdentifier{1, 3, 142, 0, 33}
// type DerSignature struct {
// 	R *big.Int
// 	S *big.Int
// }


// Serialize returns the ECDSA signature in the more strict DER format.  Note
// that the serialized bytes returned do not include the appended hash type
// used in Bitcoin signature scripts.
//
// encoding/asn1 is broken so we hand roll this output:
//
// 0x30 <length> 0x02 <length r> r 0x02 <length s> s
func (sig *DerSignature) Serialize() []byte {
	// low 'S' malleability breaker
	sigS := sig.S
	// if sigS.Cmp(S256().halfOrder) == 1 {
	// 	sigS = new(big.Int).Sub(S256().N, sigS)
	// }
	// Ensure the encoded bytes for the r and s values are canonical and
	// thus suitable for DER encoding.
	rb := canonicalizeInt(sig.R)
	sb := canonicalizeInt(sigS)

	// total length of returned signature is 1 byte for each magic and
	// length (6 total), plus lengths of r and s
	length := 6 + len(rb) + len(sb)
	b := make([]byte, length)

	b[0] = 0x30
	b[1] = byte(length - 2)
	b[2] = 0x02
	b[3] = byte(len(rb))
	offset := copy(b[4:], rb) + 4
	b[offset] = 0x02
	b[offset+1] = byte(len(sb))
	copy(b[offset+2:], sb)
	return b
}


// canonicalizeInt returns the bytes for the passed big integer adjusted as
// necessary to ensure that a big-endian encoded integer can't possibly be
// misinterpreted as a negative number.  This can happen when the most
// significant bit is set, so it is padded by a leading zero byte in this case.
// Also, the returned bytes will have at least a single byte when the passed
// value is 0.  This is required for DER encoding.
func canonicalizeInt(val *big.Int) []byte {
	b := val.Bytes()
	if len(b) == 0 {
		b = []byte{0x00}
	}
	if b[0]&0x80 != 0 {
		paddedBytes := make([]byte, len(b)+1)
		copy(paddedBytes[1:], b)
		b = paddedBytes
	}
	return b
}

// func make_pubkey_to_addr(vk *ecdsa.PublicKey) []byte {
// 	addr := sha256.Sum256(append(vk.X.Bytes(), vk.Y.Bytes()...))
// 	return addr[:]
// }


func writePublicPemToDisk(pubkey *ecdsa.PublicKey,fileName string)  (error,bool) {
  pemPublicFile, err := os.Create(fileName+".pem")
  if err != nil {
	  fmt.Println(err)
	  os.Exit(1)
  }

  pk, err := x509.MarshalPKIXPublicKey(pubkey)
  if err != nil {
	  fmt.Println(err,"parsing pem error")
  
  }

  fmt.Println(pk,"der encoded public key")
  var pemPublicBlock = &pem.Block{
    Type:  "ECDSA PUBLIC KEY",
    Bytes: pk,
  }


  err = pem.Encode(pemPublicFile, pemPublicBlock)
  if err != nil {
      fmt.Println(err)
      return err,false
  }
  pemPublicFile.Close()
  return nil,true
}




func writePrivatePemToDisk(prikey *ecdsa.PrivateKey,fileName string)  (error,bool) {
  pemPrivateFile, err := os.Create(fileName+".pem")
  if err != nil {
	  fmt.Println(err)
	  os.Exit(1)
  }

  pk, err := x509.MarshalECPrivateKey(prikey)
  if err != nil {
	  fmt.Println(err,"parsing pem error")
  
  }

  fmt.Println(pk,"der encoded public key")
  var pemPrivateBlock = &pem.Block{
    Type:  "ECDSA PRIVATE KEY",
    Bytes: pk,
  }


  err = pem.Encode(pemPrivateFile, pemPrivateBlock)
  if err != nil {
      fmt.Println(err)
      return err,false
  }
  pemPrivateFile.Close()
  return nil,true
}

type ecPrivateKey struct {
	Version       int
	PrivateKey    []byte
	NamedCurveOID asn1.ObjectIdentifier `asn1:"optional,explicit,tag:0"`
	PublicKey     asn1.BitString        `asn1:"optional,explicit,tag:1"`
}
// marshalECPrivateKey marshals an EC private key into ASN.1, DER format and
// sets the curve ID to the given OID, or omits it if OID is nil.
func marshalECPrivateKeyWithOID(key *ecdsa.PrivateKey, oid asn1.ObjectIdentifier) ([]byte, error) {
	privateKeyBytes := key.D.Bytes()
	paddedPrivateKey := make([]byte, (key.Curve.Params().N.BitLen()+7)/8)
	copy(paddedPrivateKey[len(paddedPrivateKey)-len(privateKeyBytes):], privateKeyBytes)

	return asn1.Marshal(ecPrivateKey{
		Version:       1,
		PrivateKey:    paddedPrivateKey,
		NamedCurveOID: oid,
		PublicKey:     asn1.BitString{Bytes: elliptic.Marshal(key.Curve, key.X, key.Y)},
	})
}

func writePrivatePemHexToDisk(prikey *ecdsa.PrivateKey,fileName string)  (error,bool) {
  pemPrivateFile, err := os.Create(fileName+".pem")
  if err != nil {
	  fmt.Println(err)
	  os.Exit(1)
  }

  // pk, err := x509.MarshalECPrivateKey(prikey)

  // oid, ok := oidFromNamedCurve(S256())
	// if !ok {
	// 	return nil, false
	// }

  pk, err := marshalECPrivateKeyWithOID(prikey, oidNamedCurveS256)


  if err != nil {
	  fmt.Println(err,"parsing pem error")
  
  }

  fmt.Println(pk,"der encoded public key")
  var pemPrivateBlock = &pem.Block{
    Type:  "ECDSA PRIVATE KEY",
    Bytes: pk,
  }


  err = pem.Encode(pemPrivateFile, pemPrivateBlock)
  if err != nil {
      fmt.Println(err)
      return err,false
  }
  pemPrivateFile.Close()
  return nil,true
}



func readPemFromDisk(filename string)  (error,[]byte) {
  publicKeyFile, err := os.Open(filename+".pem")
  if err != nil {
      fmt.Println(err)
      os.Exit(1)
  }

  pemfileinfo, _ := publicKeyFile.Stat()
  var size int64 = pemfileinfo.Size()
  pembytes := make([]byte, size)
  buffer := bufio.NewReader(publicKeyFile)
  _, err = buffer.Read(pembytes)
  if( err != nil){
    return err,nil
  }
  // data, _ := pem.Decode([]byte(pembytes))
  
  publicKeyFile.Close()
  return nil, pembytes
}



func parse_pem_with_der_private(privateKey []byte) (*ecdsa.PrivateKey, error) {
	// TODO: need change tx key parser, flow switch
	block, _ := pem.Decode(privateKey)

	pri, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil,err
	}

  return pri,nil
}


func parse_pem_with_hex_private(privateKey []byte) (*ecdsa.PrivateKey, error) {
	// TODO: need change tx key parser, flow switch
	block, _ := pem.Decode(privateKey)

	pri, err := parseECPrivateKey(oidNamedCurveS256,block.Bytes)
	if err != nil {
		return nil,err
	}

  return pri,nil
}

func Encode(b []byte) string {
	enc := make([]byte, len(b)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], b)
	return string(enc)
}


func parseECPrivateKey(namedCurveOID asn1.ObjectIdentifier, der []byte) (key *ecdsa.PrivateKey, err error) {
	var privKey ecPrivateKey
	if _, err := asn1.Unmarshal(der, &privKey); err != nil {
		return nil, errors.New("x509: failed to parse EC private key: " + err.Error())
	}
	if privKey.Version != ecPrivKeyVersion {
		return nil, fmt.Errorf("x509: unknown EC private key version %d", privKey.Version)
	}

	var curve elliptic.Curve
	if namedCurveOID != nil {
		curve = S256()
	} else {
		curve = S256()
	}
	if curve == nil {
		return nil, errors.New("x509: unknown elliptic curve")
	}

	k := new(big.Int).SetBytes(privKey.PrivateKey)
	curveOrder := curve.Params().N
	if k.Cmp(curveOrder) >= 0 {
		return nil, errors.New("x509: invalid elliptic curve private key value")
	}
	priv := new(ecdsa.PrivateKey)
	priv.Curve = curve
	priv.D = k

	privateKey := make([]byte, (curveOrder.BitLen()+7)/8)

	// Some private keys have leading zero padding. This is invalid
	// according to [SEC1], but this code will ignore it.
	for len(privKey.PrivateKey) > len(privateKey) {
		if privKey.PrivateKey[0] != 0 {
			return nil, errors.New("x509: invalid private key length")
		}
		privKey.PrivateKey = privKey.PrivateKey[1:]
	}

	// Some private keys remove all leading zeros, this is also invalid
	// according to [SEC1] but since OpenSSL used to do this, we ignore
	// this too.
	copy(privateKey[len(privateKey)-len(privKey.PrivateKey):], privKey.PrivateKey)
	priv.X, priv.Y = curve.ScalarBaseMult(privateKey)

	return priv, nil
}



