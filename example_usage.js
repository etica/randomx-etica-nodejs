const randomxChecker = require('./randomx');
const BigNumber = require('bignumber.js');
console.log("Loaded randomxChecker:", randomxChecker);
console.log("Available methods:", Object.keys(randomxChecker));

// Initialize RandomX with a key
const key = Buffer.from("0410591dc8b3bba89f949212982f05deeb4a1947e939c62679dfc7610c62", "hex");
randomxChecker.InitRandomX(key);

// Convert the nonce to the correct format
/*const nonce = new BigNumber("4069499503");
const nonceHex = nonce.toString(16).padStart(8, '0');
const nonceBuffer = Buffer.from(nonceHex, 'hex').reverse(); // Reverse to get little-endian format*/
const nonceHex = "d43b0000";
const nonceBuffer = Buffer.from(nonceHex, 'hex');
console.log("Nonce (hex, big-endian):", nonceHex);
console.log("Nonce (hex, little-endian):", nonceBuffer.toString('hex'));

// Seed hash
const seedHash = "6c04e936f063050f70b86c024b637335dd48c98bd47803a880e2f3bbdaf09642";
const expectedHash = "a3ffc49f732ffb916475750ef4c69274e61fd954fb8b774bcceef8b951270100";


console.log('-------------------------------------------- starting Test 1  (should be valid)------------------------------------------------------');
console.log('------------------------------------------   TEST 1    ---------------------------------------------------------------------------------------');

// New values
const blockHeader = "918785aa0dba7671f2e9b62078b664bdcb11dfc9ec88ddb7dfd36f507de69d33a89d31c96f7090b271770347526ea4de850896b1870948eb67216d49e3eb213be11cd4c767d4ba87421e83f98ccb3914";
//const target = "ff7fffff00000000000000000000000000000000000000000000000000000000"; //new BigNumber(2).pow(248);
const target = "14f8b588e368f08461f9f01b866e43aa79bbadc0980b242070b8cfbfc6540";

// Convert values to the format expected by VerifyEticaRandomXNonce
//const blockHeaderBuffer = Buffer.from(blockHeader.slice(2), 'hex'); // Remove '0x' prefix
const targetBuffer = Buffer.from(target.toString(16).padStart(64, '0'), 'hex');

//console.log("Block Header:", blockHeaderBuffer.toString('hex'));
console.log("Nonce:", nonceBuffer.toString('hex'));
console.log("Target:", targetBuffer.toString('hex'));
console.log("Seed Hash:", seedHash);

//const isValid = randomxChecker.VerifyEticaRandomXNonce(blockHeaderBuffer, nonceBuffer, targetBuffer);
const isValid = randomxChecker.VerifyEticaRandomXNonce(blockHeader, nonceBuffer, targetBuffer, seedHash, expectedHash);
console.log('isValid is: ', isValid);
console.log(`Solution is ${isValid ? 'valid' : 'invalid'}`);

console.log('------------------------------------------   TEST 1    ---------------------------------------------------------------------------------------');
console.log('-------------------------------------------- ended Test 1  (should be valid)------------------------------------------------------');


console.log('-------------------------------------------- starting Test 2  (should be invalid)------------------------------------------------------');
console.log('------------------------------------------   TEST 2    ---------------------------------------------------------------------------------------');
// Example usage
const blockHeader2 = "76a3143aab91330383ba1a1ce4f95f326e865146748035f5599e5e6ad163d04dae998839f7b398bb08b36f90d0246a64b0b8ac8376b55e060e8595c1d29556604c2b029812d3b55f4ff84baac290384b";


const isValid2 = randomxChecker.VerifyEticaRandomXNonce(blockHeader2, nonceBuffer, targetBuffer, seedHash, expectedHash);
console.log(`Solution is ${isValid2 ? 'valid' : 'invalid'}`);

console.log('------------------------------------------   TEST 2    ---------------------------------------------------------------------------------------');
console.log('-------------------------------------------- ended Test 2  (should be invalid)------------------------------------------------------');


console.log('-------------------------------------------- starting Test 3  (should be invalid)------------------------------------------------------');
console.log('------------------------------------------   TEST 3    ---------------------------------------------------------------------------------------');
// Example usage
const targetBuffer3 = "00000000ffff0000000000000000000000000000000000000000000000000000";

const isValid3 = randomxChecker.VerifyEticaRandomXNonce(blockHeader, nonceBuffer, targetBuffer3, seedHash, expectedHash);
console.log(`Solution is ${isValid3 ? 'valid' : 'invalid'}`);


console.log('------------------------------------------   TEST 3    ---------------------------------------------------------------------------------------');
console.log('-------------------------------------------- ended Test 3  (should be invalid)------------------------------------------------------');

if(isValid && !isValid2 && !isValid3){
      console.log('- - - -  - - - - -  -  - - -  - - - - -  - SUCCESSFUL TESTS --  - -  -  - - -  - - - - -  - - - - - - -  - - - - - - -');
}

else {
    console.log('- - - -  - - - - -  -  - - -  - - - - -  -  AT LEAST ONE TEST HAS FAILED --  - -  -  - - -  - - - - -  - - - - - - -  - - - - - - -');
}