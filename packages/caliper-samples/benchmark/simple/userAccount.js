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

module.exports.info  = 'opening accounts';


//let account_array = [];
let txnPerBatch;
let txSig='eyJvcmlnaW5hbCI6IiIsInNpZ25hdHVyZSI6IiIsIm5vbmNlIjoiIiwiYWN0aW9uIjoiIiwiZnJvbSI6IjZlYTMzM2NhYjMyNzBiYjBhYTNmNzIzY2QxM2RlNDllZDZlY2NhNjA1M2FkYmY2ZGZjY2ZhMjQyMTcyZTg4YmMiLCJ0byI6IjZlYTMzM2NhYjMyNzBiYjBhYTNmNzIzY2QxM2RlNDllZDZlY2NhNjA1M2FkYmY2ZGZjY2ZhMjQyMTcyZTg4YmMiLCJhbW91bnQiOiIiLCJwdWJrZXkiOiItLS0tLUJFR0lOIEVDRFNBIFBVQkxJQyBLRVktLS0tLVxuTUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFOWFicXQzR2NpSkZobHJWZS82NzZDcWV3Qzd0MlxuOC9ncEJRbGc4VVVUZklXNG9SSW9vSnZOK1JmeFdOVkhyU2FZa1p4QVdzbkxDYTBvV1dBMnQ3V2hTZz09XG4tLS0tLUVORCBFQ0RTQSBQVUJMSUMgS0VZLS0tLS1cbiJ9Cg==.ELAMA.3046022100fb156304c2f9205fc6a31ee929981269ea590989c96d95ad174a6b10027e1d40022100d2c3c221d5021903c7212b612002b5fd5534d801e545c58a8c2335410658782c';

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
function generateWorkload() {
    let workload = [];
    for(let i= 0; i < txnPerBatch; i++) {
        
        workload.push({
            chaincodeFunction: 'CreateAccount',
            chaincodeArguments: [txSig],
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
