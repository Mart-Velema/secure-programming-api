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
          classId: item.classid,
          instanceId: item.instanceid,
          appId: item.appid,
          contextId: item.contextid,

          marketHashName: item.market_hash_name,
          name: item.name,
          type: item.type,
          iconUrl: item.icon_url,
          iconUrlLarge: item.icon_url_large,

          tradable: item.tradable,
          marketable: item.marketable,

          descriptions: item.descriptions,
          actions: item.actions,
          marketActions: item.market_actions,
          tags: item.tags,
        }))
      );
    });
  });
}

module.exports = {
  getBotInventory,
};