const randomxChecker = require('./randomx');
const BigNumber = require('bignumber.js');
console.log("Loaded randomxChecker:", randomxChecker);
console.log("Available methods:", Object.keys(randomxChecker));

// Initialize RandomX with a key
const key = Buffer.from("0410591dc8b3bba89f949212982f05deeb4a1947e939c62679dfc7610c62", "hex");
randomxChecker.InitRandomX(key);

// Convert the nonce to the correct format
const nonce = new BigNumber("65536");
const nonceHex = nonce.toString(16).padStart(8, '0');
const nonceBuffer = Buffer.from(nonceHex, 'hex').reverse(); // Reverse to get little-endian format

console.log("Nonce (hex, big-endian):", nonceHex);
console.log("Nonce (hex, little-endian):", nonceBuffer.toString('hex'));

// Seed hash
const seedHash = "25314901c96d26ff28484bddf315f0a3295f30f13590d056efd65fcb6d8da788";

// Example usage
const blockHeader1 = "101096a5a1b4061274d1d8e13640eff7416062d3366960171731b703b31244d20c252d090c9d97000000008f3f41a03692ea66f71676a3eae82c215be3347b447fd2545b0cfd2c7b850ad837";
const nonce1 = "bf26f10a47750521";
const target1 = "00000000ffff0000000000000000000000000000000000000000000000000000";

const isValid1 = randomxChecker.VerifyEticaRandomXNonce(blockHeader1, nonce1, target1, seedHash);
console.log(`Solution is ${isValid1 ? 'valid' : 'invalid'}`);

console.log('---------- starting Test 2 -------------');

// New values
const blockHeader = "101096a5a1b4061274d1d8e13640eff7416062d3366960171731b703b31244d20c252d090c9d97000000008f3f41a03692ea66f71676a3eae82c215be3347b447fd2545b0cfd2c7b850ad837";
const target = "ff7fffff00000000000000000000000000000000000000000000000000000000"; //new BigNumber(2).pow(248);

// Convert values to the format expected by VerifyEticaRandomXNonce
//const blockHeaderBuffer = Buffer.from(blockHeader.slice(2), 'hex'); // Remove '0x' prefix
const targetBuffer = Buffer.from(target.toString(16).padStart(64, '0'), 'hex');

//console.log("Block Header:", blockHeaderBuffer.toString('hex'));
console.log("Nonce:", nonceBuffer.toString('hex'));
console.log("Target:", targetBuffer.toString('hex'));
console.log("Seed Hash:", seedHash);

//const isValid = randomxChecker.VerifyEticaRandomXNonce(blockHeaderBuffer, nonceBuffer, targetBuffer);
const isValid = randomxChecker.VerifyEticaRandomXNonce(blockHeader, nonceBuffer, targetBuffer, seedHash);
console.log('isValid is: ', isValid);
console.log(`Solution is ${isValid ? 'valid' : 'invalid'}`);