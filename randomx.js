const randomxAddon = require('./build/Release/randomx_addon');

function InitRandomX(key) {
  const keyBuffer = Buffer.from(key, 'hex');
  const randomxVM = randomxAddon.InitRandomX(keyBuffer);
  if (!randomxVM) {
    throw new Error('Failed to initialize RandomX');
  }
  return randomxVM;
}

function VerifyEticaRandomXNonce(blockHeader, nonce, target, seedHash, expectedHash) {
  const blockHeaderBuffer = Buffer.from(blockHeader, 'hex');
  const nonceBuffer = Buffer.from(nonce, 'hex');
  const targetBuffer = Buffer.from(target, 'hex');
  const seedHashBuffer = Buffer.from(seedHash, 'hex');
  const expectedHashBuffer = Buffer.from(expectedHash, 'hex');

  const result = randomxAddon.VerifyEticaRandomXNonce(blockHeaderBuffer, nonceBuffer, targetBuffer, seedHashBuffer, expectedHashBuffer);

  if(result == true){
    return true;
  }
  else {
    return false;
  }
  
}

module.exports = {
  InitRandomX,
  VerifyEticaRandomXNonce
};