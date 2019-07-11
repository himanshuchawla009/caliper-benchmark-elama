/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

'use strict';

let shell = require('shelljs')
let wallet= require('./wallet');

module.exports.info  = 'generating user accounts';

//let account_array = [];
let txnPerBatch;
//let initMoney;
let bc, contx;

module.exports.init = function(blockchain, context, args) {
  
    if(!args.hasOwnProperty('txnPerBatch')) {
        args.txnPerBatch = 1;
    }
    // initMoney = args.money;
    txnPerBatch = args.txnPerBatch;
    bc = blockchain;
    contx = context;

    return Promise.resolve();
};



/**
 * Generates simple workload
 * @returns {Object} array of json objects
 */
async function generateWorkload() {
    let workload = [];
    for(let i= 0; i < txnPerBatch; i++) {
        
        const key = wallet.genKey();
        console.log("ecdsa keys",key)
        //getting address from ecdsa keys
        const address = wallet.getWallet(key);
        console.log("address",address)
        //getting private key in pem file form from ecdsa key
        const pem = await wallet.getPrivatekeyAsPem(key);

        console.log("private pem",pem)

        //encrypting private key pem file with pincode
        const cipherString = wallet.getPrivateKeyAsCrypto(pem, '1111');

        console.log("cipher string",cipherString)
        //creating hash of json to be sent to chaincode
        const json_hash = wallet.createAccount(key, ".ELAMA.");
        // console.log(address)
        console.log("json hash",json_hash)
        // const t = new Date().getTime();
        // const keyfile = {
        //   enc: `enc_key_${address}_${t}.${delimiter.replace(/\./g, '').toLowerCase()}`,
        //   pem: `pk-key_${address}_${t}.pem`,
        // };

        //encrypted file key 
        //todo: save this in file system
        //download(cipherString, keyfile.enc);

        //unencrypted file , to do: sent to user
        //download(pem, keyfile.pem);

        //todo: save cipherstring,address in database
        // return {
        //   address,
        //   cipherString,
        //   json_hash,
        //   pem,
        //   keyfile,
        // }; 

        // res.status(200).json({
        //     success:true,
        //     address,
        //     cipherString,
        //     json_hash
        // })
        let args = [json_hash];
        workload.push({
            chaincodeFunction: 'CreateAccount',
            chaincodeArguments: args,
        });

    }
    return workload;
}

module.exports.run = function() {
    let args = generateWorkload();
    return bc.invokeSmartContract(contx, 'simple', 'v0', args, 300);
};

module.exports.end = function() {
    return Promise.resolve();
};
