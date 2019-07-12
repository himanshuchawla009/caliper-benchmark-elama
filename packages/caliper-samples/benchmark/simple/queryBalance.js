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

module.exports.info  = 'querying balance';



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
function generateWorkload() {
    let workload = [];
    for(let i= 0; i < txnPerBatch; i++) {
 
         
        
        workload.push({
            chaincodeFunction: 'Balance',
            chaincodeArguments: ["960c86d4acd8a0705a38495c37ca48424cd421f387f6913c390af1a57e9fbdc1"],
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
