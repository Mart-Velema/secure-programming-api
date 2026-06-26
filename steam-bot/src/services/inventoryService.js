const { manager } = require("./steamClient");

function getBotInventory(appId = 440, contextId = 2) {
  return new Promise((resolve, reject) => {
    manager.getInventoryContents(appId, contextId, true, (err, inventory) => {
      if (err) {
        reject(err);
        return;
      }

      resolve(inventory);
    });
  });
}

module.exports = {
  getBotInventory,
};