const PubNub = require('pubnub');

const credentials = {
  publishKey: 'pub-c-29e982a4-9db5-4f8e-89a3-a1a63b8eeab3',
  subscribeKey: 'sub-c-e5f9c3c3-6562-4762-b15a-e5d2a42ad12a',
  secretKey: 'sec-c-ODczZjNhMjktNTM4YS00NDI1LWEzNjAtNzQ0NzViMTkxYTgx'
};

const CHANNELS = {
  TEST: 'TEST',
  BLOCKCHAIN: 'BLOCKCHAIN',
  TRANSACTION: 'TRANSACTION'
};

class PubSub {
  constructor({ blockchain, transactionPool }) {
    this.blockchain = blockchain;
    this.transactionPool = transactionPool;

    this.pubnub = new PubNub(credentials);
    this.pubnub.subscribe({ channels: Object.values(CHANNELS) });
    this.pubnub.addListener(this.listener());
  }

  handleMessage(channel, message) {
    console.log(`Message received. Channel: ${channel}. Message: ${message}`);

    const parsedMessage = JSON.parse(message);

    switch(channel) {
      case CHANNELS.BLOCKCHAIN:
        this.blockchain.replaceChain(parsedMessage, true, () => {
          this.transactionPool.clearBlockchainTransactions({
             chain: parsedMessage
          });
        });
        break;
      case CHANNELS.TRANSACTION:
        this.transactionPool.setTransaction(parsedMessage);
        break;
      default:
        return;
    }
  }


  listener() {
    return {
      message: messageObject => {
        const { channel, message } = messageObject;

        this.handleMessage(channel, message);
      }
    };
  }

  publish({ channel, message}) {
    this.pubnub.publish({ channel, message });
  }

  broadcastChain() {
    this.publish({
      channel: CHANNELS.BLOCKCHAIN,
      message: JSON.stringify(this.blockchain.chain)
    });
  }

  broadcastTransaction(transaction) {
    this.publish({
      channel: CHANNELS.TRANSACTION,
      message: JSON.stringify(transaction)
    })
  }
}

module.exports = PubSub;
