const randomxChecker = require('./randomx');
const BigNumber = require('bignumber.js');
console.log("Loaded randomxChecker:", randomxChecker);
console.log("Available methods:", Object.keys(randomxChecker));

// Initialize RandomX with a key
const key = Buffer.from("0410591dc8b3bba89f949212982f05deeb4a1947e939c62679dfc7610c62", "hex");
randomxChecker.InitRandomX(key);

// Convert the nonce to the correct format
const nonce = new BigNumber("4069499503");
const nonceHex = nonce.toString(16).padStart(8, '0');
const nonceBuffer = Buffer.from(nonceHex, 'hex').reverse(); // Reverse to get little-endian format

console.log("Nonce (hex, big-endian):", nonceHex);
console.log("Nonce (hex, little-endian):", nonceBuffer.toString('hex'));

// Seed hash
const seedHash = "6c04e936f063050f70b86c024b637335dd48c98bd47803a880e2f3bbdaf09642";

// Example usage
const blockHeader1 = "76a3143aab91330383ba1a1ce4f95f326e865146748035f5599e5e6ad163d04dae998839f7b398bb08b36f90d0246a64b0b8ac8376b55e060e8595c1d29556604c2b029812d3b55f4ff84baac290384b";
const nonce1 = "bf26f10a47750521";
const target1 = "00000000ffff0000000000000000000000000000000000000000000000000000";

const isValid1 = randomxChecker.VerifyEticaRandomXNonce(blockHeader1, nonce1, target1, seedHash);
console.log(`Solution is ${isValid1 ? 'valid' : 'invalid'}`);

console.log('---------- starting Test 2 -------------');

// New values
const blockHeader = "76a3143aab91330383ba1a1ce4f95f326e865146748035f5599e5e6ad163d04dae998839f7b398bb08b36f90d0246a64b0b8ac8376b55e060e8595c1d29556604c2b029812d3b55f4ff84baac290384b";
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