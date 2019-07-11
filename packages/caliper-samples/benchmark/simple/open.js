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

require('./wallet.js')

//let account_array = [];
let txnPerBatch;
let txSig='eyJvcmlnaW5hbCI6IiIsInNpZ25hdHVyZSI6IiIsIm5vbmNlIjoiIiwiYWN0aW9uIjoiIiwiZnJvbSI6IiIsInRvIjoiIiwiYW1vdW50IjoiIiwicHVia2V5IjoiNTg2ZjE0NjJkOGNiYTc1NzJkODQyMDAyZTBiY2Y2M2YwNTdkOGE2YzBkNDIyNzRkNDBiNzRiOWNlMzIzY2RkNyJ9Cg==.ELAMA.304502210084d468f1fa8a847afba964795f9eec70d35097f3f2d28b2a3763aff693d9ebc0022051adcc29d5450d926f549bb1ead021b680f3d659a73b05147bbaf197ad36537a';
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
            chaincodeFunction: 'WriteAccount',
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
