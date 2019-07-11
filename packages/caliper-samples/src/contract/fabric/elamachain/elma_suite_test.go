package main

import (
	// "strconv"
	"fmt"
	"testing"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	// "github.com/ethereum/go-ethereum/common/hexutil"
	// "github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"crypto/ecdsa"

	// "bufio"
	"crypto/elliptic"
	// "crypto/md5"
	"crypto/sha256"
	// "crypto/x509"
	"crypto/rand"
  //   "encoding/json"

	// "hash"
	"encoding/hex"
	// "encoding/pem"
	"encoding/base64"
	// "io"
	// "math/big"
	// "os"
)




func TestElma(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Elma Suite")

	
}

var accountAddresses = []string{}
var accountPrivateFiles = []string{"privateZero","privateOne","privateTwo"}
var accountFiles = []string{"accountZero","accountOne","accountTwo"}
var testAccountFiles = []string{"testAccountZero","testAccountOne","testAccountTwo"}

func getAccountAddress(index int) string {
err,pubPEMData := readPemFromDisk(accountFiles[index])
Expect(err).Should(BeNil(),"Failed to read pem public key from file")
//end reading pem public file

//convert pem public key to der format
 pubOne, err := parse_pem_with_der_public([]byte(pubPEMData))
 Expect(err).Should(BeNil(),"Failed to convert pem to der format")
//end converting pem to der format

//convert der format public key to address
addrOne := make_pubkey_to_addr(pubOne)

  formattedAddrOne := fmt.Sprintf("%x", addrOne) 
 return formattedAddrOne
}




var _ = Describe("Tests for Token Chaincode", func() {
	scc := new(TokenChaincode) 
	// scc.testMode = true 
	stub := shim.NewMockStub("testingStub", scc) 
  

	It("Should initialize the chaincode", func() {
		result:= stub.MockInit("000", nil)
		status200 := int32(200)
		Expect(result.Status).Should(Equal(status200))
	})

	

	It("Should be able to initialize the cache", func() {
		// result:= stub.MockInit("000", nil)
		argsToRead := [][]byte{[]byte("CACHEINIT")}
		result := stub.MockInvoke("000", argsToRead)
		status200 := int32(200)
		expectedPayload:= []byte("Cache Init Success")
		Expect(result.Payload).Should(Equal(expectedPayload))
		Expect(result.Status).Should(Equal(status200))
	})

 
	It("Should be able to check status of cache", func() {
		// result:= stub.MockInit("000", nil)
		 argsToRead := [][]byte{[]byte("CACHESTATUS")}
		 result := stub.MockInvoke("000", argsToRead)
		status200 := int32(200)
		fmt.Println("cache result",result.Payload)
		Expect(result.Status).Should(Equal(status200))
	})

	It("Should be able to create a new account", func() {

		pubkeyCurve := elliptic.P256() //see http://golang.org/pkg/crypto/elliptic/#P256

		privatekey := new(ecdsa.PrivateKey)
		privatekey, err := ecdsa.GenerateKey(pubkeyCurve, rand.Reader) // this generates a public & private key pair
		Expect(err).Should(BeNil(),"Failed to generate key pairs")

		var pubkey ecdsa.PublicKey
		pubkey = privatekey.PublicKey
		
		//end creating key pairs

		err, don := writePrivatePemToDisk(privatekey,accountPrivateFiles[0])
	    fmt.Println("writing private file result",don)
		Expect(err).Should(BeNil(),"Failed to write pem private key to file")
	  
		//convert ecdsa public key to pem format and write it to a file.
		err, bool := writePublicPemToDisk(&pubkey,accountFiles[0])
	    fmt.Println("writing file result",bool)
		Expect(err).Should(BeNil(),"Failed to write pem public key to file")
	  
		//end writing
	   
	   //read public key pem file from disk
		err,pubPEMData := readPemFromDisk(accountFiles[0])
		Expect(err).Should(BeNil(),"Failed to read pem public key from file")
	   //end reading pem public file
	  
		
		//convert pem public key to der format
		 pub, err := parse_pem_with_der_public([]byte(pubPEMData))
		 Expect(err).Should(BeNil(),"Failed to convert pem to der format")
		//end converting pem to der format
	  
		//convert der format public key to address
		addr := make_pubkey_to_addr(pub)
	  
		  formattedAddr := fmt.Sprintf("%x", addr) 
		  accountAddresses:= append(accountAddresses,formattedAddr)
		  fmt.Println(accountAddresses[0],"account number zero address")

		//end der to address conversion
	  
	   //creating transaction struct
		transaction := PrototypeTransaction{}
		transaction.Pubkey = string(pubPEMData)
		transaction.To = formattedAddr
		transaction.From = formattedAddr
	   //end creating transaction Struct
	  
	  
		//serialize transaction
		serializedTx,err := transaction.Serialize()
		Expect(err).Should(BeNil(),"Failed to serialize the transaction")
		//end serializing transaction
	  
		//encode serialize transaction in to base64
		rawContext := base64.StdEncoding.EncodeToString(serializedTx)
		//end encoding
	  
		//create hash of encoded transaction
	   
		  h := sha256.New()
		h.Write([]byte(rawContext))
		signHash := h.Sum(nil)
		  fmt.Printf("%x", h.Sum(nil))
	   //sign the hash of encoded transaction
		r, s, serr := ecdsa.Sign(rand.Reader, privatekey, signHash)
	  
		Expect(serr).Should(BeNil(),"Failed to create signature")
	  
		//end signing
	  
		//convert the signature in to der structure
		signature := DerSignature{}
		signature.R = r
		signature.S = s
	  
	   //serialize the der structure signature
		derForm := signature.Serialize()
	  
	  
	   //end serializing

	  //hex encode signature
	  dst := make([]byte, hex.EncodedLen(len(derForm)))
	  hex.Encode(dst, derForm)
	  
	   //joining all tx data
	  
	   finalTx := rawContext + ".ELAMA." + string(dst)
	  
	   fmt.Println("========================Creatign account zero tx======================",finalTx)
	  
		// result:= stub.MockInit("000", nil)
		 argsToRead := [][]byte{[]byte("CreateAccount"),[]byte(finalTx)}
		 result := stub.MockInvoke("000", argsToRead)
		status200 := int32(200)
		expectedPayload:= []byte(fmt.Sprintf("Result: %t", true))
		fmt.Println("account result",result.Payload)
		Expect(result.Payload).Should(Equal(expectedPayload))
		Expect(result.Status).Should(Equal(status200))

	})




	It("Should be able to check balance of account one which should be zero", func() {
		// result:= stub.MockInit("000", nil)
	  //read public key pem file from disk
	  err,pubPEMData := readPemFromDisk(accountFiles[0])
	  Expect(err).Should(BeNil(),"Failed to read pem public key from file")
	 //end reading pem public file
	
	  
	  //convert pem public key to der format
	   pub, err := parse_pem_with_der_public([]byte(pubPEMData))
	   Expect(err).Should(BeNil(),"Failed to convert pem to der format")
	  //end converting pem to der format
	
	  //convert der format public key to address
	  addr := make_pubkey_to_addr(pub)
	
		formattedAddr := fmt.Sprintf("%x", addr) 
		 argsToRead := [][]byte{[]byte("Balance"),[]byte(formattedAddr)}
		 result := stub.MockInvoke("000", argsToRead)
		status200 := int32(200)
		expectedPayload := []byte(fmt.Sprintf("%d", 0))
		fmt.Println("balance result",result.Payload)
		Expect(result.Status).Should(Equal(status200))
		Expect(result.Payload).Should(Equal(expectedPayload))
	
	})


	It("Should throw err if account creation transaction contains invalid seperator", func() {

		pubkeyCurve := elliptic.P256() //see http://golang.org/pkg/crypto/elliptic/#P256

		privatekey := new(ecdsa.PrivateKey)
		privatekey, err := ecdsa.GenerateKey(pubkeyCurve, rand.Reader) // this generates a public & private key pair
		Expect(err).Should(BeNil(),"Failed to generate key pairs")

		var pubkey ecdsa.PublicKey
		pubkey = privatekey.PublicKey
		
		//end creating key pairs
	  
		//convert ecdsa public key to pem format and write it to a file.
		err, bool := writePublicPemToDisk(&pubkey,testAccountFiles[0])
	    fmt.Println("writing file result",bool)
		Expect(err).Should(BeNil(),"Failed to write pem public key to file")
	  
		//end writing
	   
	   //read public key pem file from disk
		err,pubPEMData := readPemFromDisk(testAccountFiles[0])
		Expect(err).Should(BeNil(),"Failed to read pem public key from file")
	   //end reading pem public file
	  
		
		//convert pem public key to der format
		 pub, err := parse_pem_with_der_public([]byte(pubPEMData))
		 Expect(err).Should(BeNil(),"Failed to convert pem to der format")
		//end converting pem to der format
	  
		//convert der format public key to address
		addr := make_pubkey_to_addr(pub)
	  
		  formattedAddr := fmt.Sprintf("%x", addr) 
		  accountAddresses:= append(accountAddresses,formattedAddr)
		  fmt.Println(accountAddresses[0],"account number two address")

		//end der to address conversion
	  
	   //creating transaction struct
		transaction := PrototypeTransaction{}
		transaction.Pubkey = string(pubPEMData)
		transaction.To = formattedAddr
		transaction.From = formattedAddr
	   //end creating transaction Struct
	  
	  
		//serialize transaction
		serializedTx,err := transaction.Serialize()
		Expect(err).Should(BeNil(),"Failed to serialize the transaction")
		//end serializing transaction
	  
		//encode serialize transaction in to base64
		rawContext := base64.StdEncoding.EncodeToString(serializedTx)
		//end encoding
	  
		//create hash of encoded transaction
	   
		  h := sha256.New()
		h.Write([]byte(rawContext))
		signHash := h.Sum(nil)
		  fmt.Printf("%x", h.Sum(nil))
	   //sign the hash of encoded transaction
		r, s, serr := ecdsa.Sign(rand.Reader, privatekey, signHash)
	  
		Expect(serr).Should(BeNil(),"Failed to create signature")
	  
		//end signing
	  
		//convert the signature in to der structure
		signature := DerSignature{}
		signature.R = r
		signature.S = s
	  
	   //serialize the der structure signature
		derForm := signature.Serialize()
	  
	  
	   //end serializing

	  //hex encode signature
	  dst := make([]byte, hex.EncodedLen(len(derForm)))
	  hex.Encode(dst, derForm)
	  
	//   fmt.Printf("%s\n", dst)
	  
	   //joining all tx data
	  
	   finalTx := rawContext + ".ELAA." + string(dst)
	  
	//    fmt.Println("final tx",finalTx)
	  
		// result:= stub.MockInit("000", nil)
		 argsToRead := [][]byte{[]byte("CreateAccount"),[]byte(finalTx)}
		 result := stub.MockInvoke("000", argsToRead)
		status500 := int32(500)
		// expectedPayload:= []byte(fmt.Sprintf("Result: %t", true))
		fmt.Println("account result",result.Message)
		expectedErrorMessage := "Invalid Transaction: Separator not found"
		Expect(result.Message).Should(Equal(expectedErrorMessage))
		Expect(result.Status).Should(Equal(status500))
		Expect(result.Payload).Should(BeNil(),"payload should be nil because  tx contains invalid seperator")
	})
	



	It("Should throw err if account creation transaction is in invalid format", func() {

		pubkeyCurve := elliptic.P256() //see http://golang.org/pkg/crypto/elliptic/#P256

		privatekey := new(ecdsa.PrivateKey)
		privatekey, err := ecdsa.GenerateKey(pubkeyCurve, rand.Reader) // this generates a public & private key pair
		Expect(err).Should(BeNil(),"Failed to generate key pairs")
  
		var pubkey ecdsa.PublicKey
		pubkey = privatekey.PublicKey
		
		//end creating key pairs
	  
		//convert ecdsa public key to pem format and write it to a file.
		err, bool := writePublicPemToDisk(&pubkey,testAccountFiles[0])
	    fmt.Println("writing file result",bool)
		Expect(err).Should(BeNil(),"Failed to write pem public key to file")
	  
		//end writing
	   
	   //read public key pem file from disk
		err,pubPEMData := readPemFromDisk(testAccountFiles[0])
		Expect(err).Should(BeNil(),"Failed to read pem public key from file")
	   //end reading pem public file
	  
		
		//convert pem public key to der format
		 pub, err := parse_pem_with_der_public([]byte(pubPEMData))
		 Expect(err).Should(BeNil(),"Failed to convert pem to der format")
		//end converting pem to der format
	  
		//convert der format public key to address
		addr := make_pubkey_to_addr(pub)
	  
		  formattedAddr := fmt.Sprintf("%x", addr) 
		  accountAddresses:= append(accountAddresses,formattedAddr)
		  fmt.Println(accountAddresses[0],"account number two address")

		//end der to address conversion
	  
	   //creating transaction struct
		transaction := PrototypeTransaction{}
		transaction.Pubkey = string(pubPEMData)
		transaction.To = formattedAddr
		transaction.From = formattedAddr
	   //end creating transaction Struct
	  
	  
		//serialize transaction
		serializedTx,err := transaction.Serialize()
		Expect(err).Should(BeNil(),"Failed to serialize the transaction")
		//end serializing transaction
	  
		//encode serialize transaction in to base64
		rawContext := base64.StdEncoding.EncodeToString(serializedTx)
		//end encoding
	  
		//create hash of encoded transaction
	   
		  h := sha256.New()
		h.Write([]byte(rawContext))
		signHash := h.Sum(nil)
		  fmt.Printf("%x", h.Sum(nil))
	   //sign the hash of encoded transaction
		r, s, serr := ecdsa.Sign(rand.Reader, privatekey, signHash)
	  
		Expect(serr).Should(BeNil(),"Failed to create signature")
	  
		//end signing
	  
		//convert the signature in to der structure
		signature := DerSignature{}
		signature.R = r
		signature.S = s
	  
	   //serialize the der structure signature
		derForm := signature.Serialize()
	  
	  
	   //end serializing

	  //hex encode signature
	  dst := make([]byte, hex.EncodedLen(len(derForm)))
	  hex.Encode(dst, derForm)
	  
	//   fmt.Printf("%s\n", dst)
	  
	   //joining all tx data
	  
	   finalTx := rawContext + ".ELAMA." + string(dst) + "invalid text to check tx format"
	  
	//    fmt.Println("final tx",finalTx)
	  
		// result:= stub.MockInit("000", nil)
		 argsToRead := [][]byte{[]byte("CreateAccount"),[]byte(finalTx)}
		 result := stub.MockInvoke("000", argsToRead)
		status500 := int32(500)
		// expectedPayload:= []byte(fmt.Sprintf("Result: %t", true))
		fmt.Println("account result",result.Message)
		expectedErrorMessage := "Invalid Signature Format - Invalid Signature Format - CODE704"
		Expect(result.Message).Should(Equal(expectedErrorMessage))
		Expect(result.Status).Should(Equal(status500))
		Expect(result.Payload).Should(BeNil(),"payload should be nil because of invalid tx format")
	})




	It("Should not be able to create account with different address which can be derived from public key", func() {

		pubkeyCurve := elliptic.P256() //see http://golang.org/pkg/crypto/elliptic/#P256

		privatekey := new(ecdsa.PrivateKey)
		privatekey, err := ecdsa.GenerateKey(pubkeyCurve, rand.Reader) // this generates a public & private key pair
		Expect(err).Should(BeNil(),"Failed to generate key pairs")

		var pubkey ecdsa.PublicKey
		pubkey = privatekey.PublicKey
		
		//end creating key pairs
	  
		//convert ecdsa public key to pem format and write it to a file.
		err, bool := writePublicPemToDisk(&pubkey,testAccountFiles[0])
	    fmt.Println("writing file result",bool)
		Expect(err).Should(BeNil(),"Failed to write pem public key to file")
	  
		//end writing
	   
	   //read public key pem file from disk
		err,pubPEMData := readPemFromDisk(testAccountFiles[0])
		Expect(err).Should(BeNil(),"Failed to read pem public key from file")
	   //end reading pem public file
	  
		
		//convert pem public key to der format
		 pub, err := parse_pem_with_der_public([]byte(pubPEMData))
		 Expect(err).Should(BeNil(),"Failed to convert pem to der format")
		//end converting pem to der format
	  
		
	   //convert der format public key to address
	   addr := make_pubkey_to_addr(pub)
	 
		 formattedAddr := fmt.Sprintf("%x", addr) 
	    //trye any dummy address here which is not associated with public key
		  dummyAddr := "8cda5793ab4d9f3b680e16998c2fa1923898578daa01b01588f7d219d35cbe24"
		

		//end der to address conversion
	  
	   //creating transaction struct
		transaction := PrototypeTransaction{}
		transaction.Pubkey = string(pubPEMData)
		transaction.To = dummyAddr //fake address
		transaction.From = dummyAddr
	   //end creating transaction Struct
	  
	  
		//serialize transaction
		serializedTx,err := transaction.Serialize()
		Expect(err).Should(BeNil(),"Failed to serialize the transaction")
		//end serializing transaction
	  
		//encode serialize transaction in to base64
		rawContext := base64.StdEncoding.EncodeToString(serializedTx)
		//end encoding
	  
		//create hash of encoded transaction
	   
		  h := sha256.New()
		h.Write([]byte(rawContext))
		signHash := h.Sum(nil)
		  fmt.Printf("%x", h.Sum(nil))
	   //sign the hash of encoded transaction
		r, s, serr := ecdsa.Sign(rand.Reader, privatekey, signHash)
	  
		Expect(serr).Should(BeNil(),"Failed to create signature")
	  
		//end signing
	  
		//convert the signature in to der structure
		signature := DerSignature{}
		signature.R = r
		signature.S = s
	  
	   //serialize the der structure signature
		derForm := signature.Serialize()
	  
	  
	   //end serializing

	  //hex encode signature
	  dst := make([]byte, hex.EncodedLen(len(derForm)))
	  hex.Encode(dst, derForm)
	  
	//   fmt.Printf("%s\n", dst)
	  
	   //joining all tx data
	  
	   finalTx := rawContext + ".ELAMA." + string(dst)
	  
	//    fmt.Println("final tx",finalTx)
	  
		// result:= stub.MockInit("000", nil)
		 argsToRead := [][]byte{[]byte("CreateAccount"),[]byte(finalTx)}
		 result := stub.MockInvoke("000", argsToRead)
		status500 := int32(500)
		// expectedPayload:= []byte(fmt.Sprintf("Result: %t", true))
		fmt.Println("account result",result.Message)
		expectedErrorMessage := "Sender Address Mismatch: Calculated Address - " + formattedAddr
		Expect(result.Message).Should(Equal(expectedErrorMessage))
		Expect(result.Status).Should(Equal(status500))
		Expect(result.Payload).Should(BeNil(),"payload should be nil because of invalid tx format")
	})





	It("Should not be able to create account with existing account key pairs", func() {

		err,priPEMData := readPemFromDisk(accountPrivateFiles[0])
		Expect(err).Should(BeNil(),"Failed to retrieve privat pem file")
		fmt.Println("private pem file",priPEMData)
		privatekey,err := parse_pem_with_der_private(priPEMData);
		Expect(err).Should(BeNil(),"Failed to parse private file pem")
		fmt.Println("parsed private key",privatekey)
		var pubkey ecdsa.PublicKey
		pubkey = privatekey.PublicKey
		
		//end creating key pairs
	  
		//convert ecdsa public key to pem format and write it to a file.
		err, bool := writePublicPemToDisk(&pubkey,testAccountFiles[0])
	    fmt.Println("writing file result",bool)
		Expect(err).Should(BeNil(),"Failed to write pem public key to file")
	  
		//end writing
	   
	   //read public key pem file from disk
		err,pubPEMData := readPemFromDisk(testAccountFiles[0])
		Expect(err).Should(BeNil(),"Failed to read pem public key from file")
	   //end reading pem public file
	  
		
		//convert pem public key to der format
		 pub, err := parse_pem_with_der_public([]byte(pubPEMData))
		 Expect(err).Should(BeNil(),"Failed to convert pem to der format")
		//end converting pem to der format
	  
		
	   //convert der format public key to address
	   addr := make_pubkey_to_addr(pub)
	 
	   formattedAddr := fmt.Sprintf("%x", addr) 
	   

		//end der to address conversion
	  
	   //creating transaction struct
		transaction := PrototypeTransaction{}
		transaction.Pubkey = string(pubPEMData)
		transaction.To = formattedAddr 
		transaction.From = formattedAddr
	   //end creating transaction Struct
	  
	  
		//serialize transaction
		serializedTx,err := transaction.Serialize()
		Expect(err).Should(BeNil(),"Failed to serialize the transaction")
		//end serializing transaction
	  
		//encode serialize transaction in to base64
		rawContext := base64.StdEncoding.EncodeToString(serializedTx)
		//end encoding
	  
		//create hash of encoded transaction
	   
		  h := sha256.New()
		h.Write([]byte(rawContext))
		signHash := h.Sum(nil)
		  fmt.Printf("%x", h.Sum(nil))
	   //sign the hash of encoded transaction
		r, s, serr := ecdsa.Sign(rand.Reader, privatekey, signHash)
	  
		Expect(serr).Should(BeNil(),"Failed to create signature")
	  
		//end signing
	  
		//convert the signature in to der structure
		signature := DerSignature{}
		signature.R = r
		signature.S = s
	  
	   //serialize the der structure signature
		derForm := signature.Serialize()
	  
	  
	   //end serializing

	  //hex encode signature
	  dst := make([]byte, hex.EncodedLen(len(derForm)))
	  hex.Encode(dst, derForm)
	  
	   //fmt.Printf("%s\n", dst)
	  
	   //joining all tx data
	  
	   finalTx := rawContext + ".ELAMA." + string(dst)
	  
	    // fmt.Println("final tx",finalTx)
 	  
		// result:= stub.MockInit("000", nil)
		 argsToRead := [][]byte{[]byte("CreateAccount"),[]byte(finalTx)}
		 result := stub.MockInvoke("000", argsToRead)
		status500 := int32(500)
		// expectedPayload:= []byte(fmt.Sprintf("Result: %t", true))
		fmt.Println("account result",result.Message)
		expectedErrorMessage := "Account Already Exists"
		Expect(result.Message).Should(Equal(expectedErrorMessage))
		Expect(result.Status).Should(Equal(status500))
		Expect(result.Payload).Should(BeNil(),"payload should be nil because of invalid tx format")
	})




	It("Should be able to write the  mint account", func() {
   
		

		err,priPEMData := readPemFromDisk("mintPrivateKey")
		Expect(err).Should(BeNil(),"Failed to retrieve privat pem file")
		fmt.Println("private pem file",priPEMData)
		private,err := parse_pem_with_hex_private(priPEMData);
	
		Expect(err).Should(BeNil(),"Failed to parse private file pem")
		// fmt.Println("parsed private key",privatekey)
		var pubkey ecdsa.PublicKey
		pubkey = private.PublicKey

		fmt.Println(pubkey,"pub key")
		fmt.Println(private,"admin private key")
		
		//end creating key pairs
	 
	
	  
		hexKey := elliptic.Marshal(S256(), pubkey.X, pubkey.Y)
		fmt.Println("hex key",hexKey)
		hexFormat :=Encode(hexKey)[4:]
		//convert pem public key to der format
		 pub, err :=  parse_hex_public([]byte(hexFormat))
		 Expect(err).Should(BeNil(),"Failed to convert pem to der format")
		//end converting pem to der format
	  
		
	   //convert der format public key to address
	     addr := make_pubkey_to_addr(pub)
	 
		 formattedAddr := fmt.Sprintf("%x", addr) 
		 fmt.Println("mint account address",formattedAddr)
	   

	// 	//end der to address conversion
	  
	//    //creating transaction struct
		transaction := PrototypeTransaction{}
		transaction.Pubkey = formattedAddr
		
	//    //end creating transaction Struct
	  
	  
		//serialize transaction
		serializedTx,err := transaction.Serialize()
		Expect(err).Should(BeNil(),"Failed to serialize the transaction")
		//end serializing transaction
	  
		//encode serialize transaction in to base64
		rawContext := base64.StdEncoding.EncodeToString(serializedTx)
		//end encoding
	  
		//create hash of encoded transaction
	   
		 h := sha256.New()
		h.Write([]byte(rawContext))
		signHash := h.Sum(nil)
		fmt.Printf("%x", h.Sum(nil))
		
		r, s, serr := ecdsa.Sign(rand.Reader, private, signHash)
		fmt.Println("r and s out side", r,s)
		Expect(serr).Should(BeNil(),"Failed to create signature")
	  
		//end signing
	  
		//convert the signature in to der structure
		signature := DerSignature{}
		signature.R = r
		signature.S = s
	  
	   //serialize the der structure signature
		derForm := signature.Serialize()
	  
	  
	   //end serializing

	  //hex encode signature
	  dst := make([]byte, hex.EncodedLen(len(derForm)))
	  hex.Encode(dst, derForm)
	  
	   //fmt.Printf("%s\n", dst)
	  
	   //joining all tx data
	  
	   finalTx := rawContext + ".ELAMA." + string(dst)
	  fmt.Println(finalTx,"mint account creation transaction")
	   
		argsToRead := [][]byte{[]byte("WriteAccount"),[]byte(finalTx)}
		result := stub.MockInvoke("000", argsToRead)
		status200 := int32(200)
		expectedPayload:= []byte(fmt.Sprintf("Result: %t", true))
		fmt.Println("account result",result.Message)
		fmt.Println("account result",result.Payload)
		Expect(result.Payload).Should(Equal(expectedPayload))
		Expect(result.Status).Should(Equal(status200))
	
	})


	// It("Should be able to create a mint account", func() {

	// 	pubkeyCurve := S256() 

	// 	privatekey := new(ecdsa.PrivateKey)
	// 	privatekey, err := ecdsa.GenerateKey(pubkeyCurve, rand.Reader) // this generates a public & private key pair
	// 	Expect(err).Should(BeNil(),"Failed to generate key pairs")
	// 	fmt.Println("private key of mint",privatekey)
		
	// 	var pubkey ecdsa.PublicKey
	// 	pubkey = privatekey.PublicKey
		
	// 	//end creating key pairs

	// 	err, don := writePrivatePemHexToDisk(privatekey,"mintPrivateKey")
	//     fmt.Println("writing private file result",privatekey,don)
	// 	Expect(err).Should(BeNil(),"Failed to write pem private key to file")
	  
	// 	//convert ecdsa public key to pem format and write it to a file.
	// 	err, bool := writePublicPemToDisk(&pubkey,"mintPublicKey")
	//     fmt.Println("writing file result",bool)
	// 	Expect(err).Should(BeNil(),"Failed to write pem public key to file")
	  
	// 	//end writing
	   
	
	// 	hexKey := elliptic.Marshal(S256(), pubkey.X, pubkey.Y)
	// 	hexFormat :=Encode(hexKey)[4:]
	// 	//parse key from hex format
	// 	pub, err := parse_hex_public([]byte(hexFormat))
	// 	 Expect(err).Should(BeNil(),"Failed to convert pem to der format")
	// 	//end converting pem to der format
	  
	// 	//convert der format public key to address
	// 	addr := make_pubkey_to_addr(pub)
	  
	// 	formattedAddr := fmt.Sprintf("%x", addr) 
	// 	// fmt.Printf("%s\n", dst)
	// 	// hexFormat := hex.EncodeToString(dst)
	// 	  fmt.Println("mint public key",hexFormat)
	// 	//   fmt.Printf("%s\n", hexFormat)
	// 	  fmt.Println("mint address",formattedAddr)
	// })




	It("Should be able to mint tokens", func() {
   
		
		err,priPEMData := readPemFromDisk("mintPrivateKey")
		Expect(err).Should(BeNil(),"Failed to retrieve privat pem file")
		fmt.Println("private pem file",priPEMData)
		private,err := parse_pem_with_hex_private(priPEMData);
		Expect(err).Should(BeNil(),"Failed to parse private file pem")
		// fmt.Println("parsed private key",privatekey)
		var pubkey ecdsa.PublicKey
		pubkey = private.PublicKey

		fmt.Println(pubkey,"pub key")
		
		//end creating key pairs
	 
	
	  
	    hexKey := elliptic.Marshal(S256(), pubkey.X, pubkey.Y)
        fmt.Println("hex key",hexKey)
		hexFormat :=Encode(hexKey)[4:]
		//convert pem public key to der format
		 pub, err :=  parse_hex_public([]byte(hexFormat))
		 Expect(err).Should(BeNil(),"Failed to convert pem to der format")
		//end converting pem to der format
	  
		
	   //convert der format public key to address
	     addr := make_pubkey_to_addr(pub)
	 
		 formattedAddr := fmt.Sprintf("%x", addr) 
		 fmt.Println("address",formattedAddr)
	   

	// 	//end der to address conversion
	  
	//    //creating transaction struct
		transaction := PrototypeTransaction{}
		transaction.Amount = "1000"
		
	//    //end creating transaction Struct
	  
	  
		//serialize transaction
		serializedTx,err := transaction.Serialize()
		Expect(err).Should(BeNil(),"Failed to serialize the transaction")
		//end serializing transaction
	  
		//encode serialize transaction in to base64
		rawContext := base64.StdEncoding.EncodeToString(serializedTx)
		//end encoding
	  
		//create hash of encoded transaction
	   
		 h := sha256.New()
		h.Write([]byte(rawContext))
		signHash := h.Sum(nil)
		fmt.Printf("%x", h.Sum(nil))
		
		r, s, serr := ecdsa.Sign(rand.Reader, private, signHash)
		fmt.Println("r and s out side", r,s)
		Expect(serr).Should(BeNil(),"Failed to create signature")
	  
		//end signing
	  
		//convert the signature in to der structure
		signature := DerSignature{}
		signature.R = r
		signature.S = s
	  
	   //serialize the der structure signature
		derForm := signature.Serialize()
	  
	  
	   //end serializing

	  //hex encode signature
	  dst := make([]byte, hex.EncodedLen(len(derForm)))
	  hex.Encode(dst, derForm)
	  
	   //fmt.Printf("%s\n", dst)
	  
	   //joining all tx data
	  
	   finalTx := rawContext + ".ELAMA." + string(dst)
	  fmt.Println(finalTx,"mint tokens transaction")
	   
		argsToRead := [][]byte{[]byte("Mint"),[]byte(finalTx)}
		result := stub.MockInvoke("001", argsToRead)
		status200 := int32(200)
		expectedPayload:= []byte(fmt.Sprintf("Minting Done. - [%s]",transaction.Amount))
		fmt.Println("account result",result.Message)
		fmt.Println("account result",result.Payload)
		Expect(result.Payload).Should(Equal(expectedPayload))
		Expect(result.Status).Should(Equal(status200))
	
	})




// 	It("Should not be able to mint tokens from other account except mint account", func() {
   
		
// 		err,priPEMData := readPemFromDisk(accountPrivateFiles[0])
// 	Expect(err).Should(BeNil(),"Failed to retrieve privat pem file")
// 	fmt.Println("private pem file",priPEMData)
// 	private,err := parse_pem_with_der_private(priPEMData);

// 	err,pubPEMData := readPemFromDisk(accountFiles[0])
// 	Expect(err).Should(BeNil(),"Failed to read pem public key from file")
//    //end reading pem public file
  
	
// 	//convert pem public key to der format
// 	 pub, err := parse_pem_with_der_public([]byte(pubPEMData))
// 	 Expect(err).Should(BeNil(),"Failed to convert pem to der format")
// 	//end converting pem to der format
  
	  
		
// 	   //convert der format public key to address
// 	     addr := make_pubkey_to_addr(pub)
	 
// 		 formattedAddr := fmt.Sprintf("%x", addr) 
// 		 fmt.Println("address",formattedAddr)
	   

// 	// 	//end der to address conversion
	  
// 	//    //creating transaction struct
// 		transaction := PrototypeTransaction{}
// 		transaction.Amount = "1000"
		
// 	//    //end creating transaction Struct
	  
	  
// 		//serialize transaction
// 		serializedTx,err := transaction.Serialize()
// 		Expect(err).Should(BeNil(),"Failed to serialize the transaction")
// 		//end serializing transaction
	  
// 		//encode serialize transaction in to base64
// 		rawContext := base64.StdEncoding.EncodeToString(serializedTx)
// 		//end encoding
	  
// 		//create hash of encoded transaction
	   
// 		 h := sha256.New()
// 		h.Write([]byte(rawContext))
// 		signHash := h.Sum(nil)
// 		fmt.Printf("%x", h.Sum(nil))
		
// 		r, s, serr := ecdsa.Sign(rand.Reader, private, signHash)
// 		fmt.Println("r and s out side", r,s)
// 		Expect(serr).Should(BeNil(),"Failed to create signature")
	  
// 		//end signing
	  
// 		//convert the signature in to der structure
// 		signature := DerSignature{}
// 		signature.R = r
// 		signature.S = s
	  
// 	   //serialize the der structure signature
// 		derForm := signature.Serialize()
	  
	  
// 	   //end serializing

// 	  //hex encode signature
// 	  dst := make([]byte, hex.EncodedLen(len(derForm)))
// 	  hex.Encode(dst, derForm)
	  
// 	   //fmt.Printf("%s\n", dst)
	  
// 	   //joining all tx data
	  
// 	   finalTx := rawContext + ".ELAMA." + string(dst)
// 	  fmt.Println(finalTx,"final")
	   
// 		argsToRead := [][]byte{[]byte("Mint"),[]byte(finalTx)}
// 		result := stub.MockInvoke("001", argsToRead)
// 		status500 := int32(500)
// 	//	expectedPayload:= []byte(fmt.Sprintf("Minting Done. - [%s]",transaction.Amount))
// 		fmt.Println("account result",result.Message)
// 	//	fmt.Println("account result",result.Payload)
// 	//	Expect(result.Payload).Should(Equal(expectedPayload))
// 		Expect(result.Status).Should(Equal(status500))
	
// 	})



	It("Should be able to check balance of minting account which should be 1000", func() {
	
		argsToRead := [][]byte{[]byte("Balance"),[]byte("586f1462d8cba7572d842002e0bcf63f057d8a6c0d42274d40b74b9ce323cdd7")}
		result := stub.MockInvoke("000", argsToRead)
		status200 := int32(200)
		expectedPayload := []byte(fmt.Sprintf("%d", 1000))
		fmt.Println("balance result",result.Payload)
		Expect(result.Status).Should(Equal(status200))
		Expect(result.Payload).Should(Equal(expectedPayload))
	
	})



	It("Should be able to send 500 tokens to account zero from mint account using exchange function", func() {
   
		
		formattedAddrOne := getAccountAddress(0)

		err,priPEMData := readPemFromDisk("mintPrivateKey")
		Expect(err).Should(BeNil(),"Failed to retrieve privat pem file")
		fmt.Println("private pem file",priPEMData)
		private,err := parse_pem_with_hex_private(priPEMData);
		Expect(err).Should(BeNil(),"Failed to parse private file pem")
		// fmt.Println("parsed private key",privatekey)
		var pubkey ecdsa.PublicKey
		pubkey = private.PublicKey

		fmt.Println(pubkey,"pub key")
		
		//end creating key pairs
	 
	
	  
	    hexKey := elliptic.Marshal(S256(), pubkey.X, pubkey.Y)
        fmt.Println("hex key",hexKey)
		hexFormat :=Encode(hexKey)[4:]
		//convert pem public key to der format
		 pub, err :=  parse_hex_public([]byte(hexFormat))
		 Expect(err).Should(BeNil(),"Failed to convert pem to der format")
		//end converting pem to der format
	  
		
	   //convert der format public key to address
	     addr := make_pubkey_to_addr(pub)
	 
		 formattedAddr := fmt.Sprintf("%x", addr) 
		 fmt.Println("address",formattedAddr)
	   

	// 	//end der to address conversion
	  
	//    //creating exchange struct
		transaction := PrototypeExchange{}
		transaction.Amount = "500"
		transaction.To = formattedAddrOne
		
	//    //end creating transaction Struct
	  
	  
		//serialize transaction
		serializedTx,err := transaction.Serialize()
		Expect(err).Should(BeNil(),"Failed to serialize the transaction")
		//end serializing transaction
	  
		//encode serialize transaction in to base64
		rawContext := base64.StdEncoding.EncodeToString(serializedTx)
		//end encoding
	  
		//create hash of encoded transaction
	   
		 h := sha256.New()
		h.Write([]byte(rawContext))
		signHash := h.Sum(nil)
		fmt.Printf("%x", h.Sum(nil))
		
		r, s, serr := ecdsa.Sign(rand.Reader, private, signHash)
		fmt.Println("r and s out side", r,s)
		Expect(serr).Should(BeNil(),"Failed to create signature")
	  
		//end signing
	  
		//convert the signature in to der structure
		signature := DerSignature{}
		signature.R = r
		signature.S = s
	  
	   //serialize the der structure signature
		derForm := signature.Serialize()
	  
	  
	   //end serializing

	  //hex encode signature
	  dst := make([]byte, hex.EncodedLen(len(derForm)))
	  hex.Encode(dst, derForm)
	  
	   //fmt.Printf("%s\n", dst)
	  
	   //joining all tx data
	  
	   finalTx := rawContext + ".ELAMA." + string(dst)
	  fmt.Println(finalTx,"exchange final tx")
	   
		argsToRead := [][]byte{[]byte("Exchange"),[]byte(finalTx)}
		result := stub.MockInvoke("002", argsToRead)
		status200 := int32(200)
		//expectedPayload:= []byte(fmt.Sprintf("Minting Done. - [%s]",transaction.Amount))
		fmt.Println("account result",result.Message)
		fmt.Println("account result",result.Payload)
	//	Expect(result.Payload).Should(Equal(expectedPayload))
		Expect(result.Status).Should(Equal(status200))
	
	})

	

	It("Should be able to check balance of minting account which should be 500", func() {
	
		argsToRead := [][]byte{[]byte("Balance"),[]byte("586f1462d8cba7572d842002e0bcf63f057d8a6c0d42274d40b74b9ce323cdd7")}
		result := stub.MockInvoke("000", argsToRead)
		status200 := int32(200)
		expectedPayload := []byte(fmt.Sprintf("%d", 500))
		fmt.Println("balance result",result.Payload)
		Expect(result.Status).Should(Equal(status200))
		Expect(result.Payload).Should(Equal(expectedPayload))
	
	})

	It("Should be able to check balance of  account zero which should be 500", func() {
	    addr := getAccountAddress(0)
		argsToRead := [][]byte{[]byte("Balance"),[]byte(addr)}
		result := stub.MockInvoke("000", argsToRead)
		status200 := int32(200)
		expectedPayload := []byte(fmt.Sprintf("%d", 500))
		fmt.Println("balance result",result.Payload)
		Expect(result.Status).Should(Equal(status200))
		Expect(result.Payload).Should(Equal(expectedPayload))
	
	})

	 
	It("Should throw error if minting account tries to exchange more than its balance", func() {
		currentArgs := [][]byte{[]byte("Balance"),[]byte("586f1462d8cba7572d842002e0bcf63f057d8a6c0d42274d40b74b9ce323cdd7")}
		currentArgsresult := stub.MockInvoke("000", currentArgs)
		status200 := int32(200)
		currentExpectedPayload := []byte(fmt.Sprintf("%d", 500))
		fmt.Println("balance result",currentArgsresult.Payload)
		Expect(currentArgsresult.Status).Should(Equal(status200))
		Expect(currentArgsresult.Payload).Should(Equal(currentExpectedPayload))
		
	  formattedAddrOne := getAccountAddress(0)
	  err,priPEMData := readPemFromDisk("mintPrivateKey")
	  Expect(err).Should(BeNil(),"Failed to retrieve privat pem file")
	  fmt.Println("private pem file",priPEMData)
	  private,err := parse_pem_with_hex_private(priPEMData);
	  Expect(err).Should(BeNil(),"Failed to parse private file pem")
	  // fmt.Println("parsed private key",privatekey)
	  var pubkey ecdsa.PublicKey
	  pubkey = private.PublicKey

	  fmt.Println(pubkey,"pub key")
	  
	  //end creating key pairs
   
  
	
	  hexKey := elliptic.Marshal(S256(), pubkey.X, pubkey.Y)
	  fmt.Println("hex key",hexKey)
	  hexFormat :=Encode(hexKey)[4:]
	  //convert pem public key to der format
	   pub, err :=  parse_hex_public([]byte(hexFormat))
	   Expect(err).Should(BeNil(),"Failed to convert pem to der format")
	  //end converting pem to der format
	
	  
	 //convert der format public key to address
	   addr := make_pubkey_to_addr(pub)
   
	   formattedAddr := fmt.Sprintf("%x", addr) 
	   fmt.Println("address",formattedAddr)
	 

  // 	//end der to address conversion
	
  //    //creating exchange struct
	  transaction := PrototypeExchange{}
	  transaction.Amount = "700"
	  transaction.To = formattedAddrOne
	  
  //    //end creating transaction Struct
	
	
	  //serialize transaction
	  serializedTx,err := transaction.Serialize()
	  Expect(err).Should(BeNil(),"Failed to serialize the transaction")
	  //end serializing transaction
	
	  //encode serialize transaction in to base64
	  rawContext := base64.StdEncoding.EncodeToString(serializedTx)
	  //end encoding
	
	  //create hash of encoded transaction
	 
	   h := sha256.New()
	  h.Write([]byte(rawContext))
	  signHash := h.Sum(nil)
	  fmt.Printf("%x", h.Sum(nil))
	  
	  r, s, serr := ecdsa.Sign(rand.Reader, private, signHash)
	  fmt.Println("r and s out side", r,s)
	  Expect(serr).Should(BeNil(),"Failed to create signature")
	
	  //end signing
	
	  //convert the signature in to der structure
	  signature := DerSignature{}
	  signature.R = r
	  signature.S = s
	
	 //serialize the der structure signature
	  derForm := signature.Serialize()
	
	
	 //end serializing

	//hex encode signature
	dst := make([]byte, hex.EncodedLen(len(derForm)))
	hex.Encode(dst, derForm)
	
	 //fmt.Printf("%s\n", dst)
	
	 //joining all tx data
	
	 finalTx := rawContext + ".ELAMA." + string(dst)
	fmt.Println(finalTx,"final")
	 
	  argsToRead := [][]byte{[]byte("Exchange"),[]byte(finalTx)}
	  result := stub.MockInvoke("003", argsToRead)
	  status500 := int32(500)
	  //expectedPayload:= []byte(fmt.Sprintf("Minting Done. - [%s]",transaction.Amount))
	  fmt.Println("account result",result.Message)
	  fmt.Println("account result",result.Payload)
  //	Expect(result.Payload).Should(Equal(expectedPayload))
	  Expect(result.Status).Should(Equal(status500))
	
  
  })

  It("Should be able to send 500 tokens to account zero from mint account using exchange function", func() {
   
		
	formattedAddrOne := getAccountAddress(0)

  err,priPEMData := readPemFromDisk("mintPrivateKey")
  Expect(err).Should(BeNil(),"Failed to retrieve privat pem file")
  fmt.Println("private pem file",priPEMData)
  private,err := parse_pem_with_hex_private(priPEMData);
  Expect(err).Should(BeNil(),"Failed to parse private file pem")
  // fmt.Println("parsed private key",privatekey)
  var pubkey ecdsa.PublicKey
  pubkey = private.PublicKey

  fmt.Println(pubkey,"pub key")
  
  //end creating key pairs



  hexKey := elliptic.Marshal(S256(), pubkey.X, pubkey.Y)
  fmt.Println("hex key",hexKey)
  hexFormat :=Encode(hexKey)[4:]
  //convert pem public key to der format
   pub, err :=  parse_hex_public([]byte(hexFormat))
   Expect(err).Should(BeNil(),"Failed to convert pem to der format")
  //end converting pem to der format

  
 //convert der format public key to address
   addr := make_pubkey_to_addr(pub)

   formattedAddr := fmt.Sprintf("%x", addr) 
   fmt.Println("address",formattedAddr)
 

// 	//end der to address conversion

//    //creating exchange struct
  transaction := PrototypeExchange{}
  transaction.Amount = "300"
  transaction.To = formattedAddrOne
  
//    //end creating transaction Struct


  //serialize transaction
  serializedTx,err := transaction.Serialize()
  Expect(err).Should(BeNil(),"Failed to serialize the transaction")
  //end serializing transaction

  //encode serialize transaction in to base64
  rawContext := base64.StdEncoding.EncodeToString(serializedTx)
  //end encoding

  //create hash of encoded transaction
 
   h := sha256.New()
  h.Write([]byte(rawContext))
  signHash := h.Sum(nil)
  fmt.Printf("%x", h.Sum(nil))
  
  r, s, serr := ecdsa.Sign(rand.Reader, private, signHash)
  fmt.Println("r and s out side", r,s)
  Expect(serr).Should(BeNil(),"Failed to create signature")

  //end signing

  //convert the signature in to der structure
  signature := DerSignature{}
  signature.R = r
  signature.S = s

 //serialize the der structure signature
  derForm := signature.Serialize()


 //end serializing

//hex encode signature
dst := make([]byte, hex.EncodedLen(len(derForm)))
hex.Encode(dst, derForm)

 //fmt.Printf("%s\n", dst)

 //joining all tx data

 finalTx := rawContext + ".ELAMA." + string(dst)
fmt.Println(finalTx,"final")
 
  argsToRead := [][]byte{[]byte("Exchange"),[]byte(finalTx)}
  result := stub.MockInvoke("004", argsToRead)
  status200 := int32(200)
  //expectedPayload:= []byte(fmt.Sprintf("Minting Done. - [%s]",transaction.Amount))
  fmt.Println("account result",result.Message)
  fmt.Println("account result",result.Payload)
//	Expect(result.Payload).Should(Equal(expectedPayload))
  Expect(result.Status).Should(Equal(status200))

})



  It("Should be able to check account zero transaction history", func() {
	
	addr := getAccountAddress(0)
	fmt.Println("my address",addr)
	argsToRead := [][]byte{[]byte("History"),[]byte(addr)}
	result := stub.MockInvoke("000", argsToRead)
	status200 := int32(200)
	//expectedPayload := []byte(fmt.Sprintf("%d", 500))
	fmt.Println("balance result",result.Payload)
	Expect(result.Status).Should(Equal(status200))

	txs := string(result.Payload)
	fmt.Println("**********Account Zero Transaction History*********************")
	fmt.Println(txs)
	//Expect(result.Payload).Should(Equal(expectedPayload))

 })


 It("Should be able to check mint account transaction history", func() {
	
	
	argsToRead := [][]byte{[]byte("History"),[]byte("586f1462d8cba7572d842002e0bcf63f057d8a6c0d42274d40b74b9ce323cdd7")}
	result := stub.MockInvoke("000", argsToRead)
	status200 := int32(200)
	//expectedPayload := []byte(fmt.Sprintf("%d", 500))
	fmt.Println("balance result",result.Payload)
	Expect(result.Status).Should(Equal(status200))

	txs := string(result.Payload)
	fmt.Println("**********Minting Zero Transaction History*********************")
	fmt.Println(txs)
	//Expect(result.Payload).Should(Equal(expectedPayload))

 })


 It("Should be able to query mint account", func() {
	
	
	argsToRead := [][]byte{[]byte("Query"),[]byte("586f1462d8cba7572d842002e0bcf63f057d8a6c0d42274d40b74b9ce323cdd7")}
	result := stub.MockInvoke("000", argsToRead)
	status200 := int32(200)
	//expectedPayload := []byte(fmt.Sprintf("%d", 500))
	// fmt.Println("balance result",result.Payload)
	Expect(result.Status).Should(Equal(status200))

	fmt.Println("**********Minting  account details*********************")
	
	acc := PrototypeAccount{}
	acc.Deserialize([]byte(result.Payload)); 
	fmt.Println(string(result.Payload))
	fmt.Println("accout",acc.Status)
	Expect(acc.Status).Should(Equal("Normal"))
	Expect(acc.ShrinkSize).Should(Equal(uint64(255)))
	Expect(acc.Page).Should(Equal(uint64(0)))


	//Expect(result.Payload).Should(Equal(expectedPayload))

 })



 It("Should be able to query  account zero", func() {
	
	addr := getAccountAddress(0)
	argsToRead := [][]byte{[]byte("Query"),[]byte(addr)}
	result := stub.MockInvoke("000", argsToRead)
	status200 := int32(200)
	//expectedPayload := []byte(fmt.Sprintf("%d", 500))
	// fmt.Println("balance result",result.Payload)
	Expect(result.Status).Should(Equal(status200))

	fmt.Println("**********  account  zero details*********************")
	
	acc := PrototypeAccount{}
	acc.Deserialize([]byte(result.Payload)); 
	fmt.Println(string(result.Payload))
	fmt.Println("accout",acc.Status)
	Expect(acc.Status).Should(Equal("Normal"))
	Expect(acc.ShrinkSize).Should(Equal(uint64(255)))
	Expect(acc.Page).Should(Equal(uint64(0)))


	//Expect(result.Payload).Should(Equal(expectedPayload))

 })


 It("Should not be able to make transaction with same transaction id which has been used", func() {
   
		
	formattedAddrOne := getAccountAddress(0)

  err,priPEMData := readPemFromDisk("mintPrivateKey")
  Expect(err).Should(BeNil(),"Failed to retrieve privat pem file")
  fmt.Println("private pem file",priPEMData)
  private,err := parse_pem_with_hex_private(priPEMData);
  Expect(err).Should(BeNil(),"Failed to parse private file pem")
  // fmt.Println("parsed private key",privatekey)
  var pubkey ecdsa.PublicKey
  pubkey = private.PublicKey

  fmt.Println(pubkey,"pub key")
  
  //end creating key pairs



  hexKey := elliptic.Marshal(S256(), pubkey.X, pubkey.Y)
  fmt.Println("hex key",hexKey)
  hexFormat :=Encode(hexKey)[4:]
  //convert pem public key to der format
   pub, err :=  parse_hex_public([]byte(hexFormat))
   Expect(err).Should(BeNil(),"Failed to convert pem to der format")
  //end converting pem to der format

  
 //convert der format public key to address
   addr := make_pubkey_to_addr(pub)

   formattedAddr := fmt.Sprintf("%x", addr) 
   fmt.Println("address",formattedAddr)
 

// 	//end der to address conversion

//    //creating exchange struct
  transaction := PrototypeExchange{}
  transaction.Amount = "10"
  transaction.To = formattedAddrOne
  
//    //end creating transaction Struct


  //serialize transaction
  serializedTx,err := transaction.Serialize()
  Expect(err).Should(BeNil(),"Failed to serialize the transaction")
  //end serializing transaction

  //encode serialize transaction in to base64
  rawContext := base64.StdEncoding.EncodeToString(serializedTx)
  //end encoding

  //create hash of encoded transaction
 
   h := sha256.New()
  h.Write([]byte(rawContext))
  signHash := h.Sum(nil)
  fmt.Printf("%x", h.Sum(nil))
  
  r, s, serr := ecdsa.Sign(rand.Reader, private, signHash)
  fmt.Println("r and s out side", r,s)
  Expect(serr).Should(BeNil(),"Failed to create signature")

  //end signing

  //convert the signature in to der structure
  signature := DerSignature{}
  signature.R = r
  signature.S = s

 //serialize the der structure signature
  derForm := signature.Serialize()


 //end serializing

//hex encode signature
dst := make([]byte, hex.EncodedLen(len(derForm)))
hex.Encode(dst, derForm)

 //fmt.Printf("%s\n", dst)

 //joining all tx data

 finalTx := rawContext + ".ELAMA." + string(dst)
fmt.Println(finalTx,"final")
 
  argsToRead := [][]byte{[]byte("Exchange"),[]byte(finalTx)}
  result := stub.MockInvoke("004", argsToRead)
  status500 := int32(500)
  //expectedPayload:= []byte(fmt.Sprintf("Minting Done. - [%s]",transaction.Amount))
  fmt.Println("account result",result.Message)
  fmt.Println("account result",result.Payload)
//	Expect(result.Payload).Should(Equal(expectedPayload))
  Expect(result.Status).Should(Equal(status500))

})
  

It("Should be able to create second account", func() {
	
	pubkeyCurve := elliptic.P256() //see http://golang.org/pkg/crypto/elliptic/#P256

	privatekey := new(ecdsa.PrivateKey)
	privatekey, err := ecdsa.GenerateKey(pubkeyCurve, rand.Reader) // this generates a public & private key pair
	Expect(err).Should(BeNil(),"Failed to generate key pairs")

	var pubkey ecdsa.PublicKey
	pubkey = privatekey.PublicKey
	
	//end creating key pairs

	err, don := writePrivatePemToDisk(privatekey,accountPrivateFiles[1])
	fmt.Println("writing private file result",don)
	Expect(err).Should(BeNil(),"Failed to write pem private key to file")
  
	//convert ecdsa public key to pem format and write it to a file.
	err, bool := writePublicPemToDisk(&pubkey,accountFiles[1])
	fmt.Println("writing file result",bool)
	Expect(err).Should(BeNil(),"Failed to write pem public key to file")
  
	//end writing
   
   //read public key pem file from disk
	err,pubPEMData := readPemFromDisk(accountFiles[1])
	Expect(err).Should(BeNil(),"Failed to read pem public key from file")
   //end reading pem public file
  
	
	//convert pem public key to der format
	 pub, err := parse_pem_with_der_public([]byte(pubPEMData))
	 Expect(err).Should(BeNil(),"Failed to convert pem to der format")
	//end converting pem to der format
  
	//convert der format public key to address
	addr := make_pubkey_to_addr(pub)
  
	  formattedAddr := fmt.Sprintf("%x", addr) 
	  accountAddresses:= append(accountAddresses,formattedAddr)
	  fmt.Println(accountAddresses[0],"account number two address")

	//end der to address conversion
  
   //creating transaction struct
	transaction := PrototypeTransaction{}
	transaction.Pubkey = string(pubPEMData)
	transaction.To = formattedAddr
	transaction.From = formattedAddr
   //end creating transaction Struct
  
  
	//serialize transaction
	serializedTx,err := transaction.Serialize()
	Expect(err).Should(BeNil(),"Failed to serialize the transaction")
	//end serializing transaction
  
	//encode serialize transaction in to base64
	rawContext := base64.StdEncoding.EncodeToString(serializedTx)
	//end encoding
  
	//create hash of encoded transaction
   
	  h := sha256.New()
	h.Write([]byte(rawContext))
	signHash := h.Sum(nil)
	  fmt.Printf("%x", h.Sum(nil))
   //sign the hash of encoded transaction
	r, s, serr := ecdsa.Sign(rand.Reader, privatekey, signHash)
  
	Expect(serr).Should(BeNil(),"Failed to create signature")
  
	//end signing
  
	//convert the signature in to der structure
	signature := DerSignature{}
	signature.R = r
	signature.S = s
  
   //serialize the der structure signature
	derForm := signature.Serialize()
  
  
   //end serializing

  //hex encode signature
  dst := make([]byte, hex.EncodedLen(len(derForm)))
  hex.Encode(dst, derForm)
  
//   fmt.Printf("%s\n", dst)
  
   //joining all tx data
  
   finalTx := rawContext + ".ELAMA." + string(dst)
  
   fmt.Println("account one",finalTx)
  
	// result:= stub.MockInit("000", nil)
	 argsToRead := [][]byte{[]byte("CreateAccount"),[]byte(finalTx)}
	 result := stub.MockInvoke("000", argsToRead)
	status200 := int32(200)
	expectedPayload:= []byte(fmt.Sprintf("Result: %t", true))
	fmt.Println("account result",result.Payload)
	Expect(result.Payload).Should(Equal(expectedPayload))
	Expect(result.Status).Should(Equal(status200))

})






It("Should be make a transaction with 64 bit txid", func() {

    //sender balance before 
	senderAddr := getAccountAddress(0)
	senderArgs := [][]byte{[]byte("Balance"),[]byte(senderAddr)}
	senderResult := stub.MockInvoke("000", senderArgs)
	status200 := int32(200)
	//expectedPayload := []byte(fmt.Sprintf("%d", 500))
	
	bal:=fmt.Sprintf("%d", senderResult.Payload)
	fmt.Println("sender balance result before",string(bal))
	Expect(senderResult.Status).Should(Equal(status200))
	//Expect(result.Payload).Should(Equal(expectedPayload))
	//sender balance before 
	receiverAddr := getAccountAddress(1)
	receiverArgs := [][]byte{[]byte("Balance"),[]byte(receiverAddr)}
	receiverResult := stub.MockInvoke("000", receiverArgs)

	receiverExpectedPayload := []byte(fmt.Sprintf("%d", 0))
	fmt.Println("receiver balance result before",receiverResult.Payload)
	Expect(receiverResult.Status).Should(Equal(status200))
    Expect(receiverResult.Payload).Should(Equal(receiverExpectedPayload))


	//receiver account address
	formattedAddrOne := getAccountAddress(1)
	//sender private key
	err,priPEMData := readPemFromDisk(accountPrivateFiles[0])
	Expect(err).Should(BeNil(),"Failed to retrieve privat pem file")
	fmt.Println("private pem file",priPEMData)
	privatekey,err := parse_pem_with_der_private(priPEMData);

	// var pubkey ecdsa.PublicKey
	// pubkey = privatekey.PublicKey
	
	//end creating key pairs


   
   //read public key pem file from disk
	err,pubPEMData := readPemFromDisk(accountFiles[0])
	Expect(err).Should(BeNil(),"Failed to read pem public key from file")
   //end reading pem public file
  
	
	//convert pem public key to der format
	 pub, err := parse_pem_with_der_public([]byte(pubPEMData))
	 Expect(err).Should(BeNil(),"Failed to convert pem to der format")
	//end converting pem to der format
  
	//convert der format public key to address
	addr := make_pubkey_to_addr(pub)
  
	  formattedAddr := fmt.Sprintf("%x", addr) 
	  accountAddresses:= append(accountAddresses,formattedAddr)
	  fmt.Println(accountAddresses[0],"account number two address")

	//end der to address conversion
  
   //creating transaction struct
	transaction := PrototypeTransaction{}
	transaction.Pubkey = string(pubPEMData)
	transaction.To = formattedAddrOne
	transaction.Amount = "10"
	transaction.From = formattedAddr
   //end creating transaction Struct
  
  
	//serialize transaction
	serializedTx,err := transaction.Serialize()
	Expect(err).Should(BeNil(),"Failed to serialize the transaction")
	//end serializing transaction
  
	//encode serialize transaction in to base64
	rawContext := base64.StdEncoding.EncodeToString(serializedTx)
	//end encoding
  
	//create hash of encoded transaction
   
	  h := sha256.New()
	h.Write([]byte(rawContext))
	signHash := h.Sum(nil)
	  fmt.Printf("%x", h.Sum(nil))
   //sign the hash of encoded transaction
	r, s, serr := ecdsa.Sign(rand.Reader, privatekey, signHash)
  
	Expect(serr).Should(BeNil(),"Failed to create signature")
  
	//end signing
  
	//convert the signature in to der structure
	signature := DerSignature{}
	signature.R = r
	signature.S = s
  
   //serialize the der structure signature
	derForm := signature.Serialize()
  
  
   //end serializing

  //hex encode signature
  dst := make([]byte, hex.EncodedLen(len(derForm)))
  hex.Encode(dst, derForm)
  
   //fmt.Printf("%s\n", dst)
  
   //joining all tx data
  
   finalTx := rawContext + ".ELAMA." + string(dst)
  
   fmt.Println("account zero to account one",finalTx)
   status500 := int32(200)

	// result:= stub.MockInit("000", nil)
	 argsToRead := [][]byte{[]byte("Transaction"),[]byte(finalTx)}
	 //uuid less than 64 byte : 901928
	 result := stub.MockInvoke("e87b402a-659a-11e9-a923-1681be663d3ee87b402a-65e87b402a-65745343", argsToRead)
	
//	expectedPayload:= []byte(fmt.Sprintf("Result: %t", true))
	fmt.Println("account result",result.Payload)
//	Expect(result.Payload).Should(Equal(expectedPayload))
	Expect(result.Status).Should(Equal(status500))


})



 It("Should be able to check account one transaction history", func() {
	
	addr := getAccountAddress(1)
	fmt.Println("account one address",addr)
	argsToRead := [][]byte{[]byte("History"),[]byte(addr)}
	result := stub.MockInvoke("000", argsToRead)
	status200 := int32(200)
	//expectedPayload := []byte(fmt.Sprintf("%d", 500))
	fmt.Println("balance result",result.Payload)
	Expect(result.Status).Should(Equal(status200))

	txs := string(result.Payload)
	fmt.Println("**********Account one Transaction History*********************")
	fmt.Println(txs)
	//Expect(result.Payload).Should(Equal(expectedPayload))

 })

 It("Should be able to check mint accoount transaction history", func() {
	
	addr := getAccountAddress(1)
	fmt.Println("my address",addr)
	argsToRead := [][]byte{[]byte("History"),[]byte("586f1462d8cba7572d842002e0bcf63f057d8a6c0d42274d40b74b9ce323cdd7")}
	result := stub.MockInvoke("000", argsToRead)
	status200 := int32(200)
	//expectedPayload := []byte(fmt.Sprintf("%d", 500))
	fmt.Println("balance result",result.Payload)
	Expect(result.Status).Should(Equal(status200))

	txs := string(result.Payload)
	fmt.Println("**********Mint Transaction History*********************")
	fmt.Println(txs)
	//Expect(result.Payload).Should(Equal(expectedPayload))

 })


It("Should not be able to able to make transactions with same signatures", func() {
	status500 := int32(500)
    //sender balance before 
	senderAddr := getAccountAddress(0)
	senderArgs := [][]byte{[]byte("Balance"),[]byte(senderAddr)}
	senderResult := stub.MockInvoke("000", senderArgs)
	status200 := int32(200)
	//expectedPayload := []byte(fmt.Sprintf("%d", 500))
	
	bal:=fmt.Sprintf("%d", senderResult.Payload)
	fmt.Println("sender balance result before",string(bal))
	Expect(senderResult.Status).Should(Equal(status200))
	//Expect(result.Payload).Should(Equal(expectedPayload))
	//sender balance before 
	receiverAddr := getAccountAddress(1)
	receiverArgs := [][]byte{[]byte("Balance"),[]byte(receiverAddr)}
	receiverResult := stub.MockInvoke("000", receiverArgs)

	receiverExpectedPayload := []byte(fmt.Sprintf("%d", 10))
	fmt.Println("receiver balance result before",receiverResult.Payload)
	Expect(receiverResult.Status).Should(Equal(status200))
    Expect(receiverResult.Payload).Should(Equal(receiverExpectedPayload))


	
	//sender private key
	err,priPEMData := readPemFromDisk(accountPrivateFiles[0])
	Expect(err).Should(BeNil(),"Failed to retrieve privat pem file")
	fmt.Println("private pem file",priPEMData)
	privatekey,err := parse_pem_with_der_private(priPEMData);

	// var pubkey ecdsa.PublicKey
	// pubkey = privatekey.PublicKey
	
	//end creating key pairs


   
   //read public key pem file from disk
	err,pubPEMData := readPemFromDisk(accountFiles[0])
	Expect(err).Should(BeNil(),"Failed to read pem public key from file")
   //end reading pem public file
  
	
	//convert pem public key to der format
	 pub, err := parse_pem_with_der_public([]byte(pubPEMData))
	 Expect(err).Should(BeNil(),"Failed to convert pem to der format")
	//end converting pem to der format
  
	//convert der format public key to address
	addr := make_pubkey_to_addr(pub)
  
	  formattedAddr := fmt.Sprintf("%x", addr) 
	  accountAddresses:= append(accountAddresses,formattedAddr)
	  fmt.Println(accountAddresses[0],"account number two address")

	//end der to address conversion
  
   //creating transaction struct
	transaction := PrototypeTransaction{}
	transaction.Pubkey = string(pubPEMData)
	transaction.To = receiverAddr
	transaction.Amount = "10"
	transaction.From = senderAddr
   //end creating transaction Struct
  
  
	//serialize transaction
	serializedTx,err := transaction.Serialize()
	Expect(err).Should(BeNil(),"Failed to serialize the transaction")
	//end serializing transaction
  
	//encode serialize transaction in to base64
	rawContext := base64.StdEncoding.EncodeToString(serializedTx)
	//end encoding
  
	//create hash of encoded transaction
   
	  h := sha256.New()
	h.Write([]byte(rawContext))
	signHash := h.Sum(nil)
	  fmt.Printf("%x", h.Sum(nil))
   //sign the hash of encoded transaction
	r, s, serr := ecdsa.Sign(rand.Reader, privatekey, signHash)
  
	Expect(serr).Should(BeNil(),"Failed to create signature")
  
	//end signing
  
	//convert the signature in to der structure
	signature := DerSignature{}
	signature.R = r
	signature.S = s
  
   //serialize the der structure signature
	derForm := signature.Serialize()
  
  
   //end serializing

  //hex encode signature
  dst := make([]byte, hex.EncodedLen(len(derForm)))
  hex.Encode(dst, derForm)
  
   //fmt.Printf("%s\n", dst)
  
   //joining all tx data
  
   finalTx := rawContext + ".ELAMA." + string(dst)
  
   fmt.Println("final tx",finalTx)
  
	// result:= stub.MockInit("000", nil)
	 argsToRead := [][]byte{[]byte("Transaction"),[]byte(finalTx)}
	 result := stub.MockInvoke("006", argsToRead)

	
//	expectedPayload:= []byte(fmt.Sprintf("Result: %t", true))
	fmt.Println("account result",result.Payload)
//	Expect(result.Payload).Should(Equal(expectedPayload))
	Expect(result.Status).Should(Equal(status200))



	   //sender balance after 
	  
	   senderArgsafter := [][]byte{[]byte("Balance"),[]byte(senderAddr)}
	   senderResultafter := stub.MockInvoke("000", senderArgsafter)
	   
	   //expectedPayload := []byte(fmt.Sprintf("%d", 500))
	   fmt.Println("sender balance result after",senderResultafter.Payload)
	   Expect(senderResultafter.Status).Should(Equal(status200))
	   //Expect(result.Payload).Should(Equal(expectedPayload))
	   //sender balance before 

	   receiverArgsafter := [][]byte{[]byte("Balance"),[]byte(receiverAddr)}
	   receiverResultafter := stub.MockInvoke("000", receiverArgsafter)
	  
	   receiverExpectedPayloadafter := []byte(fmt.Sprintf("%d", 20))
	   fmt.Println("receiver balance result after",receiverResultafter.Payload)
	   Expect(receiverResultafter.Status).Should(Equal(status200))
	   Expect(receiverResultafter.Payload).Should(Equal(receiverExpectedPayloadafter))
       
	   replayTx := stub.MockInvoke("007", argsToRead)
	   Expect(replayTx.Status).Should(Equal(status500),"Failed to send a error on replay transaction")
})





// It("Should not be able to able to make bulk transactions", func() {

//     //sender balance before 
// 	senderAddr := getAccountAddress(0)
// 	senderArgs := [][]byte{[]byte("Balance"),[]byte(senderAddr)}
// 	senderResult := stub.MockInvoke("000", senderArgs)
// 	status200 := int32(200)
// 	//expectedPayload := []byte(fmt.Sprintf("%d", 500))
	
// 	bal:=fmt.Sprintf("%d", senderResult.Payload)
// 	fmt.Println("sender balance result before",string(bal))
// 	Expect(senderResult.Status).Should(Equal(status200))
// 	//Expect(result.Payload).Should(Equal(expectedPayload))
// 	//sender balance before 
// 	receiverAddr := getAccountAddress(1)
// 	receiverArgs := [][]byte{[]byte("Balance"),[]byte(receiverAddr)}
// 	receiverResult := stub.MockInvoke("000", receiverArgs)

// 	receiverExpectedPayload := []byte(fmt.Sprintf("%d", 20))
// 	fmt.Println("receiver balance result before",receiverResult.Payload)
// 	Expect(receiverResult.Status).Should(Equal(status200))
//     Expect(receiverResult.Payload).Should(Equal(receiverExpectedPayload))


	
// 	//sender private key
// 	err,priPEMData := readPemFromDisk(accountPrivateFiles[0])
// 	Expect(err).Should(BeNil(),"Failed to retrieve privat pem file")
// 	fmt.Println("private pem file",priPEMData)
// 	privatekey,err := parse_pem_with_der_private(priPEMData);

// 	// var pubkey ecdsa.PublicKey
// 	// pubkey = privatekey.PublicKey
	
// 	//end creating key pairs


   
//    //read public key pem file from disk
// 	err,pubPEMData := readPemFromDisk(accountFiles[0])
// 	Expect(err).Should(BeNil(),"Failed to read pem public key from file")
//    //end reading pem public file
  
	
// 	//convert pem public key to der format
// 	 pub, err := parse_pem_with_der_public([]byte(pubPEMData))
// 	 Expect(err).Should(BeNil(),"Failed to convert pem to der format")
// 	//end converting pem to der format
  
// 	//convert der format public key to address
// 	addr := make_pubkey_to_addr(pub)
  
// 	  formattedAddr := fmt.Sprintf("%x", addr) 
// 	  accountAddresses:= append(accountAddresses,formattedAddr)
// 	  fmt.Println(accountAddresses[0],"account number two address")

// 	//end der to address conversion
  
  
//    for i := 10; i < 500; i++ {
//  //creating transaction struct
//  transaction := PrototypeTransaction{}
//  transaction.BaseTransaction.Nonce = strconv.Itoa(i)
//  transaction.Pubkey = string(pubPEMData)
//  transaction.To = receiverAddr
//  transaction.Amount = "1"
//  transaction.From = senderAddr
// //end creating transaction Struct


//  //serialize transaction
//  serializedTx,err := transaction.Serialize()
//  Expect(err).Should(BeNil(),"Failed to serialize the transaction")
//  //end serializing transaction

//  //encode serialize transaction in to base64
//  rawContext := base64.StdEncoding.EncodeToString(serializedTx)
//  //end encoding

//  //create hash of encoded transaction

//    h := sha256.New()
//  h.Write([]byte(rawContext))
//  signHash := h.Sum(nil)
//    fmt.Printf("%x", h.Sum(nil))
// //sign the hash of encoded transaction
//  r, s, serr := ecdsa.Sign(rand.Reader, privatekey, signHash)

//  Expect(serr).Should(BeNil(),"Failed to create signature")

//  //end signing

//  //convert the signature in to der structure
//  signature := DerSignature{}
//  signature.R = r
//  signature.S = s

// //serialize the der structure signature
//  derForm := signature.Serialize()


// //end serializing

// //hex encode signature
// dst := make([]byte, hex.EncodedLen(len(derForm)))
// hex.Encode(dst, derForm)

// //fmt.Printf("%s\n", dst)

// //joining all tx data

// finalTx := rawContext + ".ELAMA." + string(dst)

//  // result:= stub.MockInit("000", nil)
//   argsToRead := [][]byte{[]byte("Transaction"),[]byte(finalTx)}
//   result := stub.MockInvoke("00"+ strconv.Itoa(i), argsToRead)

 

// //	Expect(result.Payload).Should(Equal(expectedPayload))
//  Expect(result.Status).Should(Equal(status200))

// 	}
  
// 	receiverArgsafter := [][]byte{[]byte("Balance"),[]byte(receiverAddr)}
// 	receiverResultafter := stub.MockInvoke("000", receiverArgsafter)
   
// 	receiverExpectedPayloadafter := []byte(fmt.Sprintf("%d", 510))
// 	fmt.Println("receiver balance result after",receiverResultafter.Payload)
// 	Expect(receiverResultafter.Status).Should(Equal(status200))
// 	Expect(receiverResultafter.Payload).Should(Equal(receiverExpectedPayloadafter))
	   
// })


It("Should be able to check account zero transaction history once again", func() {
	
	addr := getAccountAddress(0)
	fmt.Println("my address",addr)
	argsToRead := [][]byte{[]byte("History"),[]byte(addr)}
	result := stub.MockInvoke("000", argsToRead)
	status200 := int32(200)
	//expectedPayload := []byte(fmt.Sprintf("%d", 500))
	Expect(result.Status).Should(Equal(status200))

	// txs := string(result.Payload)
	// fmt.Println("**********Account Zero Transaction History*********************")
	// fmt.Println(txs)
	//Expect(result.Payload).Should(Equal(expectedPayload))

 })



 It("Should be able to query  account zero", func() {
	
	accountZero := getAccountAddress(0)
	
	argsToRead := [][]byte{[]byte("Query"),[]byte(accountZero)}
	result := stub.MockInvoke("000", argsToRead)
	status200 := int32(200)
	//expectedPayload := []byte(fmt.Sprintf("%d", 500))
	// fmt.Println("balance result",result.Payload)
	Expect(result.Status).Should(Equal(status200))

	fmt.Println("*********Account Zero details*********************")
	
	acc := PrototypeAccount{}
	acc.Deserialize([]byte(result.Payload)); 
	fmt.Println(string(result.Payload))
	fmt.Println("accout",acc.Status)
	Expect(acc.Status).Should(Equal("Normal"))
	Expect(acc.ShrinkSize).Should(Equal(uint64(255)))
	Expect(acc.Page).Should(Equal(uint64(0)))


	//Expect(result.Payload).Should(Equal(expectedPayload))

 })



It("Should not be able to able to make transaction more than its balance", func() {
	
    //sender balance before 
	senderAddr := getAccountAddress(1)
	senderArgs := [][]byte{[]byte("Balance"),[]byte(senderAddr)}
	senderResult := stub.MockInvoke("000", senderArgs)
	status200 := int32(200)
	senderExpectedPayload := []byte(fmt.Sprintf("%d", 510))
	//expectedPayload := []byte(fmt.Sprintf("%d", 500))
	Expect(senderResult.Payload).Should(Equal(senderExpectedPayload))
	bal:=fmt.Sprintf("%d", senderResult.Payload)
	fmt.Println("sender balance result before",string(bal))
	Expect(senderResult.Status).Should(Equal(status200))
	//Expect(result.Payload).Should(Equal(expectedPayload))
	//sender balance before 
	receiverAddr := getAccountAddress(0)
	receiverArgs := [][]byte{[]byte("Balance"),[]byte(receiverAddr)}
	receiverResult := stub.MockInvoke("000", receiverArgs)

	// receiverExpectedPayload := []byte(fmt.Sprintf("%d", 10))
	// fmt.Println("receiver balance result before",receiverResult.Payload)
	Expect(receiverResult.Status).Should(Equal(status200))
    


	
	//sender private key
	err,priPEMData := readPemFromDisk(accountPrivateFiles[1])
	Expect(err).Should(BeNil(),"Failed to retrieve privat pem file")
	fmt.Println("private pem file",priPEMData)
	privatekey,err := parse_pem_with_der_private(priPEMData);

	// var pubkey ecdsa.PublicKey
	// pubkey = privatekey.PublicKey
	
	//end creating key pairs


   
   //read public key pem file from disk
	err,pubPEMData := readPemFromDisk(accountFiles[1])
	Expect(err).Should(BeNil(),"Failed to read pem public key from file")
   //end reading pem public file
  
	
	//convert pem public key to der format
	 pub, err := parse_pem_with_der_public([]byte(pubPEMData))
	 Expect(err).Should(BeNil(),"Failed to convert pem to der format")
	//end converting pem to der format
  
	//convert der format public key to address
	addr := make_pubkey_to_addr(pub)
  
	  formattedAddr := fmt.Sprintf("%x", addr) 
	  accountAddresses:= append(accountAddresses,formattedAddr)
	  fmt.Println(accountAddresses[0],"account number two address")

	//end der to address conversion
  
   //creating transaction struct
	transaction := PrototypeTransaction{}
	transaction.Pubkey = string(pubPEMData)
	transaction.To = receiverAddr
	transaction.Amount = "520"
	transaction.From = senderAddr
   //end creating transaction Struct
  
  
	//serialize transaction
	serializedTx,err := transaction.Serialize()
	Expect(err).Should(BeNil(),"Failed to serialize the transaction")
	//end serializing transaction
  
	//encode serialize transaction in to base64
	rawContext := base64.StdEncoding.EncodeToString(serializedTx)
	//end encoding
  
	//create hash of encoded transaction
   
	  h := sha256.New()
	h.Write([]byte(rawContext))
	signHash := h.Sum(nil)
	  fmt.Printf("%x", h.Sum(nil))
   //sign the hash of encoded transaction
	r, s, serr := ecdsa.Sign(rand.Reader, privatekey, signHash)
  
	Expect(serr).Should(BeNil(),"Failed to create signature")
  
	//end signing
  
	//convert the signature in to der structure
	signature := DerSignature{}
	signature.R = r
	signature.S = s
  
   //serialize the der structure signature
	derForm := signature.Serialize()
  
  
   //end serializing

  //hex encode signature
  dst := make([]byte, hex.EncodedLen(len(derForm)))
  hex.Encode(dst, derForm)

  
   finalTx := rawContext + ".ELAMA." + string(dst)
  
   fmt.Println("final tx",finalTx)
  
	// result:= stub.MockInit("000", nil)
	 argsToRead := [][]byte{[]byte("Transaction"),[]byte(finalTx)}
	 result := stub.MockInvoke("837837", argsToRead)
	 status500 := int32(500)
	
//	expectedPayload:= []byte(fmt.Sprintf("Result: %t", true))
	fmt.Println("account result",result.Payload)
//	Expect(result.Payload).Should(Equal(expectedPayload))
	Expect(result.Status).Should(Equal(status500))


})




It("Should not be able to able to make transaction if private key is not associated with public key", func() {
	
    //sender balance before 
	senderAddr := getAccountAddress(1)
	senderArgs := [][]byte{[]byte("Balance"),[]byte(senderAddr)}
	senderResult := stub.MockInvoke("000", senderArgs)
	status200 := int32(200)
	senderExpectedPayload := []byte(fmt.Sprintf("%d", 510))
	//expectedPayload := []byte(fmt.Sprintf("%d", 500))
	Expect(senderResult.Payload).Should(Equal(senderExpectedPayload))
	bal:=fmt.Sprintf("%d", senderResult.Payload)
	fmt.Println("sender balance result before",string(bal))
	Expect(senderResult.Status).Should(Equal(status200))
	//Expect(result.Payload).Should(Equal(expectedPayload))
	//sender balance before 
	receiverAddr := getAccountAddress(0)
	receiverArgs := [][]byte{[]byte("Balance"),[]byte(receiverAddr)}
	receiverResult := stub.MockInvoke("000", receiverArgs)

	// receiverExpectedPayload := []byte(fmt.Sprintf("%d", 10))
	// fmt.Println("receiver balance result before",receiverResult.Payload)
	Expect(receiverResult.Status).Should(Equal(status200))
    


	
	//sender private key
	err,priPEMData := readPemFromDisk(accountPrivateFiles[1])
	Expect(err).Should(BeNil(),"Failed to retrieve privat pem file")
	fmt.Println("private pem file",priPEMData)
	privatekey,err := parse_pem_with_der_private(priPEMData);

	// var pubkey ecdsa.PublicKey
	// pubkey = privatekey.PublicKey
	
	//end creating key pairs


   
   //wrong public key for testing
	err,pubPEMData := readPemFromDisk(accountFiles[0])
	Expect(err).Should(BeNil(),"Failed to read pem public key from file")
   //end reading pem public file
  
	
	//convert pem public key to der format
	 pub, err := parse_pem_with_der_public([]byte(pubPEMData))
	 Expect(err).Should(BeNil(),"Failed to convert pem to der format")
	//end converting pem to der format
  
	//convert der format public key to address
	addr := make_pubkey_to_addr(pub)
  
	  formattedAddr := fmt.Sprintf("%x", addr) 
	  accountAddresses:= append(accountAddresses,formattedAddr)
	  fmt.Println(accountAddresses[0],"account number two address")

	//end der to address conversion
  
   //creating transaction struct
	transaction := PrototypeTransaction{}
	transaction.Pubkey = string(pubPEMData)
	transaction.To = receiverAddr
	transaction.Amount = "520"
	transaction.From = senderAddr
   //end creating transaction Struct
  
  
	//serialize transaction
	serializedTx,err := transaction.Serialize()
	Expect(err).Should(BeNil(),"Failed to serialize the transaction")
	//end serializing transaction
  
	//encode serialize transaction in to base64
	rawContext := base64.StdEncoding.EncodeToString(serializedTx)
	//end encoding
  
	//create hash of encoded transaction
   
	  h := sha256.New()
	h.Write([]byte(rawContext))
	signHash := h.Sum(nil)
	  fmt.Printf("%x", h.Sum(nil))
   //sign the hash of encoded transaction
	r, s, serr := ecdsa.Sign(rand.Reader, privatekey, signHash)
  
	Expect(serr).Should(BeNil(),"Failed to create signature")
  
	//end signing
  
	//convert the signature in to der structure
	signature := DerSignature{}
	signature.R = r
	signature.S = s
  
   //serialize the der structure signature
	derForm := signature.Serialize()
  
  
   //end serializing

  //hex encode signature
  dst := make([]byte, hex.EncodedLen(len(derForm)))
  hex.Encode(dst, derForm)

  
   finalTx := rawContext + ".ELAMA." + string(dst)
  
   fmt.Println("final tx",finalTx)
  
	// result:= stub.MockInit("000", nil)
	 argsToRead := [][]byte{[]byte("Transaction"),[]byte(finalTx)}
	 result := stub.MockInvoke("837837", argsToRead)
	 status500 := int32(500)
	
//	expectedPayload:= []byte(fmt.Sprintf("Result: %t", true))
	fmt.Println("account result",result.Payload)
//	Expect(result.Payload).Should(Equal(expectedPayload))
	Expect(result.Status).Should(Equal(status500))


})




It("Should not be able to mint tokens from other account except mint account", func() {
   
		
	err,priPEMData := readPemFromDisk(accountPrivateFiles[1])
Expect(err).Should(BeNil(),"Failed to retrieve privat pem file")
fmt.Println("private pem file",priPEMData)
private,err := parse_pem_with_der_private(priPEMData);

err,pubPEMData := readPemFromDisk(accountFiles[1])
Expect(err).Should(BeNil(),"Failed to read pem public key from file")
//end reading pem public file


//convert pem public key to der format
 pub, err := parse_pem_with_der_public([]byte(pubPEMData))
 Expect(err).Should(BeNil(),"Failed to convert pem to der format")
//end converting pem to der format

  
	
   //convert der format public key to address
	 addr := make_pubkey_to_addr(pub)
 
	 formattedAddr := fmt.Sprintf("%x", addr) 
	 fmt.Println("address",formattedAddr)
   

// 	//end der to address conversion
  
//    //creating transaction struct
	transaction := PrototypeTransaction{}
	transaction.Amount = "1000"
	
//    //end creating transaction Struct
  
  
	//serialize transaction
	serializedTx,err := transaction.Serialize()
	Expect(err).Should(BeNil(),"Failed to serialize the transaction")
	//end serializing transaction
  
	//encode serialize transaction in to base64
	rawContext := base64.StdEncoding.EncodeToString(serializedTx)
	//end encoding
  
	//create hash of encoded transaction
   
	 h := sha256.New()
	h.Write([]byte(rawContext))
	signHash := h.Sum(nil)
	fmt.Printf("%x", h.Sum(nil))
	
	r, s, serr := ecdsa.Sign(rand.Reader, private, signHash)
	fmt.Println("r and s out side", r,s)
	Expect(serr).Should(BeNil(),"Failed to create signature")
  
	//end signing
  
	//convert the signature in to der structure
	signature := DerSignature{}
	signature.R = r
	signature.S = s
  
   //serialize the der structure signature
	derForm := signature.Serialize()
  
  
   //end serializing

  //hex encode signature
  dst := make([]byte, hex.EncodedLen(len(derForm)))
  hex.Encode(dst, derForm)
  
   //fmt.Printf("%s\n", dst)
  
   //joining all tx data
  
   finalTx := rawContext + ".ELAMA." + string(dst)
  fmt.Println(finalTx,"final")
   
	argsToRead := [][]byte{[]byte("Mint"),[]byte(finalTx)}
	result := stub.MockInvoke("001", argsToRead)
	status500 := int32(500)
//	expectedPayload:= []byte(fmt.Sprintf("Minting Done. - [%s]",transaction.Amount))
	fmt.Println("account result",result.Message)
//	fmt.Println("account result",result.Payload)
//	Expect(result.Payload).Should(Equal(expectedPayload))
	Expect(result.Status).Should(Equal(status500))

})



 It("Should be able to check account zero transaction history once again", func() {
	
	addr := getAccountAddress(0)
	fmt.Println("my address",addr)
	argsToRead := [][]byte{[]byte("History"),[]byte(addr)}
	result := stub.MockInvoke("000", argsToRead)
	status200 := int32(200)
	//expectedPayload := []byte(fmt.Sprintf("%d", 500))
	fmt.Println("balance result",result.Payload)
	Expect(result.Status).Should(Equal(status200))

	txs := string(result.Payload)
	fmt.Println("**********Account Zero Transaction History*********************")
	fmt.Println(txs)
	//Expect(result.Payload).Should(Equal(expectedPayload))

 })



It("Should be throw error if txid length is less than 64 bit", func() {

    //sender balance before 
	senderAddr := getAccountAddress(0)
	senderArgs := [][]byte{[]byte("Balance"),[]byte(senderAddr)}
	senderResult := stub.MockInvoke("000", senderArgs)
	status200 := int32(200)
	//expectedPayload := []byte(fmt.Sprintf("%d", 500))
	
	bal:=fmt.Sprintf("%d", senderResult.Payload)
	fmt.Println("sender balance result before",string(bal))
	Expect(senderResult.Status).Should(Equal(status200))
	//Expect(result.Payload).Should(Equal(expectedPayload))
	//sender balance before 
	receiverAddr := getAccountAddress(1)
	receiverArgs := [][]byte{[]byte("Balance"),[]byte(receiverAddr)}
	receiverResult := stub.MockInvoke("000", receiverArgs)

	receiverExpectedPayload := []byte(fmt.Sprintf("%d", 0))
	fmt.Println("receiver balance result before",receiverResult.Payload)
	Expect(receiverResult.Status).Should(Equal(status200))
    Expect(receiverResult.Payload).Should(Equal(receiverExpectedPayload))


	//receiver account address
	formattedAddrOne := getAccountAddress(1)
	//sender private key
	err,priPEMData := readPemFromDisk(accountPrivateFiles[0])
	Expect(err).Should(BeNil(),"Failed to retrieve privat pem file")
	fmt.Println("private pem file",priPEMData)
	privatekey,err := parse_pem_with_der_private(priPEMData);

	// var pubkey ecdsa.PublicKey
	// pubkey = privatekey.PublicKey
	
	//end creating key pairs


   
   //read public key pem file from disk
	err,pubPEMData := readPemFromDisk(accountFiles[0])
	Expect(err).Should(BeNil(),"Failed to read pem public key from file")
   //end reading pem public file
  
	
	//convert pem public key to der format
	 pub, err := parse_pem_with_der_public([]byte(pubPEMData))
	 Expect(err).Should(BeNil(),"Failed to convert pem to der format")
	//end converting pem to der format
  
	//convert der format public key to address
	addr := make_pubkey_to_addr(pub)
  
	  formattedAddr := fmt.Sprintf("%x", addr) 
	  accountAddresses:= append(accountAddresses,formattedAddr)
	  fmt.Println(accountAddresses[0],"account number two address")

	//end der to address conversion
  
   //creating transaction struct
	transaction := PrototypeTransaction{}
	transaction.Pubkey = string(pubPEMData)
	transaction.To = formattedAddrOne
	transaction.Amount = "10"
	transaction.From = formattedAddr
   //end creating transaction Struct
  
  
	//serialize transaction
	serializedTx,err := transaction.Serialize()
	Expect(err).Should(BeNil(),"Failed to serialize the transaction")
	//end serializing transaction
  
	//encode serialize transaction in to base64
	rawContext := base64.StdEncoding.EncodeToString(serializedTx)
	//end encoding
  
	//create hash of encoded transaction
   
	  h := sha256.New()
	h.Write([]byte(rawContext))
	signHash := h.Sum(nil)
	  fmt.Printf("%x", h.Sum(nil))
   //sign the hash of encoded transaction
	r, s, serr := ecdsa.Sign(rand.Reader, privatekey, signHash)
  
	Expect(serr).Should(BeNil(),"Failed to create signature")
  
	//end signing
  
	//convert the signature in to der structure
	signature := DerSignature{}
	signature.R = r
	signature.S = s
  
   //serialize the der structure signature
	derForm := signature.Serialize()
  
  
   //end serializing

  //hex encode signature
  dst := make([]byte, hex.EncodedLen(len(derForm)))
  hex.Encode(dst, derForm)
  
   //fmt.Printf("%s\n", dst)
  
   //joining all tx data
  
   finalTx := rawContext + ".ELAMA." + string(dst)
  
   fmt.Println("final tx",finalTx)
   status500 := int32(500)

	// result:= stub.MockInit("000", nil)
	 argsToRead := [][]byte{[]byte("Transaction"),[]byte(finalTx)}
	 //uuid less than 64 byte : 901928
	 result := stub.MockInvoke("901928", argsToRead)
	
//	expectedPayload:= []byte(fmt.Sprintf("Result: %t", true))
	fmt.Println("account result",result.Payload)
//	Expect(result.Payload).Should(Equal(expectedPayload))
	Expect(result.Status).Should(Equal(status500))


})



It("Should be able to query  a transaction id", func() {
	
	// addr := getAccountAddress(0)
	argsToRead := [][]byte{[]byte("TXID"),[]byte("e87b402a-659a-11e9-a923-1681be663d3ee87b402a-65e87b402a-65745343")}
	result := stub.MockInvoke("", argsToRead)
	status200 := int32(200)
	//expectedPayload := []byte(fmt.Sprintf("%d", 500))
	fmt.Println("balance result",string(result.Payload))
	Expect(result.Status).Should(Equal(status200))


 })

 It("Should be able throw error if txid doesn,t exist", func() {
	
	// addr := getAccountAddress(0)
	argsToRead := [][]byte{[]byte("TXID"),[]byte("e87b409a-659a-11e9-a923-1681be663d3ee87b402a-65e87b402a-65745343")}
	result := stub.MockInvoke("", argsToRead)
	status500 := int32(500)
	//expectedPayload := []byte(fmt.Sprintf("%d", 500))
	fmt.Println("balance result",string(result.Message))
	Expect(result.Message).Should(Equal("No Data"))
	Expect(result.Status).Should(Equal(status500))

	addrZero := getAccountAddress(0)
	addrOne := getAccountAddress(1)

	fmt.Println("accout zero",addrZero)
	fmt.Println("accout one",addrOne)


 })

 


})

