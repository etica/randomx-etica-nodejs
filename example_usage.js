const randomxChecker = require('./randomx');
const BigNumber = require('bignumber.js');
console.log("Loaded randomxChecker:", randomxChecker);
console.log("Available methods:", Object.keys(randomxChecker));

// Initialize RandomX with a key
const key = Buffer.from("0410591dc8b3bba89f949212982f05deeb4a1947e939c62679dfc7610c62", "hex");
randomxChecker.InitRandomX(key);

// Example usage
const blockHeader1 = "8139df23abad4be6fcbbbdfcbc5e41d11e9d890432772d9ed1714862e0f5e362";
const nonce1 = "bf26f10a47750521";
const target1 = "00000000ffff0000000000000000000000000000000000000000000000000000";

const isValid1 = randomxChecker.VerifyEticaRandomXNonce(blockHeader1, nonce1, target1);
console.log(`Solution is ${isValid1 ? 'valid' : 'invalid'}`);

console.log('---------- starting Test 2 -------------');

// New values
const blockHeader = "0x15b5584cf95dd4b07e9e2c30c5a3d015527e07e35f2f1a614ce7f5e8f943ae37";
const nonce = new BigNumber("12848617903317052952");
const target = new BigNumber(2).pow(248);

// Convert values to the format expected by VerifyEticaRandomXNonce
const blockHeaderBuffer = Buffer.from(blockHeader.slice(2), 'hex'); // Remove '0x' prefix
const nonceBuffer = Buffer.from(nonce.toString(16).padStart(16, '0'), 'hex');
const targetBuffer = Buffer.from(target.toString(16).padStart(64, '0'), 'hex');

console.log("Block Header:", blockHeaderBuffer.toString('hex'));
console.log("Nonce:", nonceBuffer.toString('hex'));
console.log("Target:", targetBuffer.toString('hex'));

const isValid = randomxChecker.VerifyEticaRandomXNonce(blockHeaderBuffer, nonceBuffer, targetBuffer);
console.log('isValid is: ', isValid);
console.log(`Solution is ${isValid ? 'valid' : 'invalid'}`);