const { manager } = require("./steamClient");

function getBotInventory(appId = 440, contextId = 2) {
  return new Promise((resolve, reject) => {
    manager.getInventoryContents(appId, contextId, true, (err, inventory) => {
      if (err) {
        reject(err);
        return;
      }

      resolve(
        inventory.map((item) => ({
          assetId: item.assetid,
          marketHashName: item.market_hash_name,
          name: item.name,
          tradable: item.tradable,
          marketable: item.marketable,
        }))
      );
    });
  });
}

module.exports = {
  getBotInventory,
};