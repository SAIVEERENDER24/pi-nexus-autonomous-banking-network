const Web3 = require('web3');

const piNetworkContractAbi = require('../contracts/PiNetworkContract.json').abi;
const piNetworkContractAddress = '0x...'; // Replace with the deployed Pi Network contract address

const web3 = new Web3(new Web3.providers.HttpProvider('https://mainnet.infura.io/v3/YOUR_PROJECT_ID'));

const piNetworkContract = new web3.eth.Contract(piNetworkContractAbi, piNetworkContractAddress);

async function getNodeCount() {
  try {
    const nodeCount = await piNetworkContract.methods.getNodeCount().call();
    return nodeCount;
  } catch (error) {
    console.error(error);
    return null;
  }
}

async function getNodeByAddress(nodeAddress) {
  try {
    const node = await piNetworkContract.methods.getNodeByAddress(nodeAddress).call();
    return node;
  } catch (error) {
    console.error(error);
    return null;
  }
}

async function getPiBalance(userAddress) {
  try {
    const balance = await piNetworkContract.methods.getPiBalance(userAddress).call();
    return balance;
  } catch (error) {
    console.error(error);
    return null;
  }
}

async function isNodeActive(nodeAddress) {
  try {
    const isActive = await piNetworkContract.methods.isNodeActive(nodeAddress).call();
    return isActive;
  } catch (error) {
    console.error(error);
    return null;
  }
}

module.exports = {
  getNodeCount,
  getNodeByAddress,
  getPiBalance,
  isNodeActive
};
