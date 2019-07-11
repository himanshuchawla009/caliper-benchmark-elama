  
cat << "EOF"
     
                          
 /$$$$$$$$ /$$       /$$      /$$  /$$$$$$   /$$$$$$  /$$   /$$  /$$$$$$  /$$$$$$ /$$   /$$
| $$_____/| $$      | $$$    /$$$ /$$__  $$ /$$__  $$| $$  | $$ /$$__  $$|_  $$_/| $$$ | $$
| $$      | $$      | $$$$  /$$$$| $$  \ $$| $$  \__/| $$  | $$| $$  \ $$  | $$  | $$$$| $$
| $$$$$   | $$      | $$ $$/$$ $$| $$$$$$$$| $$      | $$$$$$$$| $$$$$$$$  | $$  | $$ $$ $$
| $$__/   | $$      | $$  $$$| $$| $$__  $$| $$      | $$__  $$| $$__  $$  | $$  | $$  $$$$
| $$      | $$      | $$\  $ | $$| $$  | $$| $$    $$| $$  | $$| $$  | $$  | $$  | $$\  $$$
| $$$$$$$$| $$$$$$$$| $$ \/  | $$| $$  | $$|  $$$$$$/| $$  | $$| $$  | $$ /$$$$$$| $$ \  $$
|________/|________/|__/     |__/|__/  |__/ \______/ |__/  |__/|__/  |__/|______/|__/  \__/
                                                                                           
                                                                                                                                                                                    
                                                                                                
"     
EOF

  echo "##########################################################"
  echo "######### Setting channel name in environment ############"
  echo "##########################################################"
export CHANNEL_NAME=mychannel 

# docker exec -it cli bash

  echo "##########################################################"
  echo "######### Checking Cache Status ##########################"
  echo "##########################################################"
docker exec cli peer chaincode query -C $CHANNEL_NAME -n elma -c '{"Args":["CACHESTATUS"]}'



  echo "##########################################################"
  echo "######### Creating MINT Account ##########################"
  echo "##########################################################"

docker exec cli peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n elma --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses peer0.org2.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["WriteAccount","eyJvcmlnaW5hbCI6IiIsInNpZ25hdHVyZSI6IiIsIm5vbmNlIjoiIiwiYWN0aW9uIjoiIiwiZnJvbSI6IiIsInRvIjoiIiwiYW1vdW50IjoiIiwicHVia2V5IjoiNTg2ZjE0NjJkOGNiYTc1NzJkODQyMDAyZTBiY2Y2M2YwNTdkOGE2YzBkNDIyNzRkNDBiNzRiOWNlMzIzY2RkNyJ9Cg==.ELAMA.304402201a8e2854f94d5f75a0d25663a2b7932d2a133860e901c38b39df36722e928fe80220293f999dbb7b308e4ae03467fe59eb9876dda128bad1ed93f70ea52cbdc26d19"]}'

  echo "##########################################################"
  echo "######### Creating MINT Account Done #####################"
  echo "##########################################################"

  echo "##########################################################"
  echo "######### Creating New Account (Account zero) ############"
  echo "##########################################################"

docker exec cli peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n elma --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses peer0.org2.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["CreateAccount","eyJvcmlnaW5hbCI6IiIsInNpZ25hdHVyZSI6IiIsIm5vbmNlIjoiIiwiYWN0aW9uIjoiIiwiZnJvbSI6IjZlYTMzM2NhYjMyNzBiYjBhYTNmNzIzY2QxM2RlNDllZDZlY2NhNjA1M2FkYmY2ZGZjY2ZhMjQyMTcyZTg4YmMiLCJ0byI6IjZlYTMzM2NhYjMyNzBiYjBhYTNmNzIzY2QxM2RlNDllZDZlY2NhNjA1M2FkYmY2ZGZjY2ZhMjQyMTcyZTg4YmMiLCJhbW91bnQiOiIiLCJwdWJrZXkiOiItLS0tLUJFR0lOIEVDRFNBIFBVQkxJQyBLRVktLS0tLVxuTUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFOWFicXQzR2NpSkZobHJWZS82NzZDcWV3Qzd0MlxuOC9ncEJRbGc4VVVUZklXNG9SSW9vSnZOK1JmeFdOVkhyU2FZa1p4QVdzbkxDYTBvV1dBMnQ3V2hTZz09XG4tLS0tLUVORCBFQ0RTQSBQVUJMSUMgS0VZLS0tLS1cbiJ9Cg==.ELAMA.3046022100fb156304c2f9205fc6a31ee929981269ea590989c96d95ad174a6b10027e1d40022100d2c3c221d5021903c7212b612002b5fd5534d801e545c58a8c2335410658782c"]}'

  echo "##########################################################"
  echo "######### Finished Creating New Account ##################"
  echo "##########################################################"

  echo "##########################################################"
  echo "######### Minting tokens to mint account #################"
  echo "##########################################################"

docker exec cli peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n elma --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses peer0.org2.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["Mint","eyJvcmlnaW5hbCI6IiIsInNpZ25hdHVyZSI6IiIsIm5vbmNlIjoiIiwiYWN0aW9uIjoiIiwiZnJvbSI6IiIsInRvIjoiIiwiYW1vdW50IjoiMTAwMCIsInB1YmtleSI6IiJ9Cg==.ELAMA.3046022100bcf5d209719030b707cdd89e3a99516e345c749a1abd6d191bebe66e2f562ec1022100df82cdf9f92594d186f6b3547a166d7e67dc15dfe30b2f36fd849cb141a24f0f"]}'


  echo "##########################################################"
  echo "######### Finished Minting tokens to mint account ########"
  echo "##########################################################"


  echo "##########################################################"
  echo "## Checking Mint Account Balance #########################"
  echo "##########################################################"
docker exec cli peer chaincode query -C $CHANNEL_NAME -n elma -c '{"Args":["Balance","586f1462d8cba7572d842002e0bcf63f057d8a6c0d42274d40b74b9ce323cdd7"]}'


  echo "##########################################################"
  echo "## Sending tokens from mint account to account zero ######"
  echo "##########################################################"


docker exec cli peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n elma --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses peer0.org2.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["Exchange","eyJvcmlnaW5hbCI6IiIsIlNpZ25hdHVyZSI6IiIsInRvIjoiNmVhMzMzY2FiMzI3MGJiMGFhM2Y3MjNjZDEzZGU0OWVkNmVjY2E2MDUzYWRiZjZkZmNjZmEyNDIxNzJlODhiYyIsImFtb3VudCI6IjUwMCJ9Cg==.ELAMA.304502202283d6febd3da0924b494f26b479d4f6813793bf0b17345e906a1e521c0bff640221008dad431cad4e351eef113835528d99a3f9c1a13cb54ec952b68babfdba26e422"]}'

  echo "##########################################################"
  echo "## Sent tokens from mint account to account zero #########"
  echo "##########################################################"

  echo "##########################################################"
  echo "## Checking Mint Account Balance #########################"
  echo "##########################################################"
docker exec cli peer chaincode query -C $CHANNEL_NAME -n elma -c '{"Args":["Balance","586f1462d8cba7572d842002e0bcf63f057d8a6c0d42274d40b74b9ce323cdd7"]}'


  echo "##########################################################"
  echo "## Checking Account Zero balance #########################"
  echo "##########################################################"

docker exec cli peer chaincode query -C $CHANNEL_NAME -n elma -c '{"Args":["Balance","6ea333cab3270bb0aa3f723cd13de49ed6ecca6053adbf6dfccfa242172e88bc"]}'

  echo "##########################################################"
  echo "## Checking history of account zero ######################"
  echo "##########################################################" 
docker exec cli peer chaincode query -C $CHANNEL_NAME -n elma -c '{"Args":["History","6ea333cab3270bb0aa3f723cd13de49ed6ecca6053adbf6dfccfa242172e88bc"]}'

  echo "##########################################################"
  echo "## Checking history of mint account ######################"
  echo "##########################################################" 
docker exec cli peer chaincode query -C $CHANNEL_NAME -n elma -c '{"Args":["History","586f1462d8cba7572d842002e0bcf63f057d8a6c0d42274d40b74b9ce323cdd7"]}'



  echo "##########################################################"
  echo "## Creating a new account (Account one) ##################"
  echo "##########################################################" 
docker exec cli peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n elma --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses peer0.org2.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["CreateAccount","eyJvcmlnaW5hbCI6IiIsInNpZ25hdHVyZSI6IiIsIm5vbmNlIjoiIiwiYWN0aW9uIjoiIiwiZnJvbSI6ImRiZGMxZmVjNWViZWNlYzA0YzU3NzJjNjVhYTk0NTVkZTMxYjVhNzM2ODU0NjE4NzU5ZGEyYjBhNjk2ZDk3NWIiLCJ0byI6ImRiZGMxZmVjNWViZWNlYzA0YzU3NzJjNjVhYTk0NTVkZTMxYjVhNzM2ODU0NjE4NzU5ZGEyYjBhNjk2ZDk3NWIiLCJhbW91bnQiOiIiLCJwdWJrZXkiOiItLS0tLUJFR0lOIEVDRFNBIFBVQkxJQyBLRVktLS0tLVxuTUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFQmtKSFUxZ09XM0dySE5QU0pTM2gzY2NNbElmdlxuN2F6a3RWRGF2OHBLY0xGbkIyRkFXUG5LZmtLck9JeXUxWExFNmt0Z1VFQlFOM3VraTlBZlhLY1pNdz09XG4tLS0tLUVORCBFQ0RTQSBQVUJMSUMgS0VZLS0tLS1cbiJ9Cg==.ELAMA.304502202b0fc3c1c71cec0cfda83c6b7013e9c82d57f54b3642f851d5991e4e4ecd37eb0221009af71b8c8a11128b1c493d47081af81fae9eca0a9bc95bc404d0f51df998053a"]}'


  echo "##########################################################"
  echo "## Transferring tokens from account zero to account one ##"
  echo "##########################################################" 
docker exec cli peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n elma --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses peer0.org2.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["Transaction","eyJvcmlnaW5hbCI6IiIsInNpZ25hdHVyZSI6IiIsIm5vbmNlIjoiIiwiYWN0aW9uIjoiIiwiZnJvbSI6IjZlYTMzM2NhYjMyNzBiYjBhYTNmNzIzY2QxM2RlNDllZDZlY2NhNjA1M2FkYmY2ZGZjY2ZhMjQyMTcyZTg4YmMiLCJ0byI6ImRiZGMxZmVjNWViZWNlYzA0YzU3NzJjNjVhYTk0NTVkZTMxYjVhNzM2ODU0NjE4NzU5ZGEyYjBhNjk2ZDk3NWIiLCJhbW91bnQiOiIxMCIsInB1YmtleSI6Ii0tLS0tQkVHSU4gRUNEU0EgUFVCTElDIEtFWS0tLS0tXG5NRmt3RXdZSEtvWkl6ajBDQVFZSUtvWkl6ajBEQVFjRFFnQUU5YWJxdDNHY2lKRmhsclZlLzY3NkNxZXdDN3QyXG44L2dwQlFsZzhVVVRmSVc0b1JJb29Kdk4rUmZ4V05WSHJTYVlrWnhBV3NuTENhMG9XV0EydDdXaFNnPT1cbi0tLS0tRU5EIEVDRFNBIFBVQkxJQyBLRVktLS0tLVxuIn0K.ELAMA.304502207b78ebbbc1f1d196134bda40ca8f02d292bc9130e8e3b147d80081004cc26da60221009c37b0a5651357d50a0913e290218ab89c203961eb01cf1e6d3287faa380c9f3"]}'


  echo "##########################################################"
  echo "## checking balance of account zero #######################"
  echo "##########################################################"
docker exec cli peer chaincode query -C $CHANNEL_NAME -n elma -c '{"Args":["Balance","6ea333cab3270bb0aa3f723cd13de49ed6ecca6053adbf6dfccfa242172e88bc"]}'



  echo "##########################################################"
  echo "## checking balance of account one #######################"
  echo "##########################################################"
docker exec cli peer chaincode query -C $CHANNEL_NAME -n elma -c '{"Args":["Balance","dbdc1fec5ebecec04c5772c65aa9455de31b5a736854618759da2b0a696d975b"]}'


  echo "##########################################################"
  echo "## checking history of account one #######################"
  echo "##########################################################"
docker exec cli peer chaincode query -C $CHANNEL_NAME -n elma -c '{"Args":["History","dbdc1fec5ebecec04c5772c65aa9455de31b5a736854618759da2b0a696d975b"]}'


  echo "##########################################################"
  echo "## checking details of account one #######################"
  echo "##########################################################"
docker exec cli peer chaincode query -C $CHANNEL_NAME -n elma -c '{"Args":["Query","dbdc1fec5ebecec04c5772c65aa9455de31b5a736854618759da2b0a696d975b"]}'

  echo "##########################################################"
  echo "## checking balance of account zero #######################"
  echo "##########################################################"
docker exec cli peer chaincode query -C $CHANNEL_NAME -n elma -c '{"Args":["Query","6ea333cab3270bb0aa3f723cd13de49ed6ecca6053adbf6dfccfa242172e88bc"]}'


  echo "##########################################################"
  echo "## checking balance of mint account ######################"
  echo "##########################################################"
docker exec cli peer chaincode query -C $CHANNEL_NAME -n elma -c '{"Args":["Query","586f1462d8cba7572d842002e0bcf63f057d8a6c0d42274d40b74b9ce323cdd7"]}'


  echo "##########################################################"
  echo "######### END ############################################"
  echo "##########################################################"


