const randomxAddon = require('./build/Release/randomx_addon');

function InitRandomX(key) {
  const keyBuffer = Buffer.from(key, 'hex');
  console.log("Calling InitRandomX with key:", key);
  const randomxVM = randomxAddon.InitRandomX(keyBuffer);
  console.log("InitRandomX randomxVM:", randomxVM);
  if (!randomxVM) {
    throw new Error('Failed to initialize RandomX');
  }
  return randomxVM;
}

function VerifyEticaRandomXNonce(blockHeader, nonce, target) {
  const blockHeaderBuffer = Buffer.from(blockHeader, 'hex');
  const nonceBuffer = Buffer.from(nonce, 'hex');
  const targetBuffer = Buffer.from(target, 'hex');

  const result = randomxAddon.VerifyEticaRandomXNonce(blockHeaderBuffer, nonceBuffer, targetBuffer);
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