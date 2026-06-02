const TradeOfferManager = require("steam-tradeoffer-manager");
const { client, manager } = require("./steamClient");

function sendTradeOffer(tradeUrl, itemsToGive, message = "GuineaTrade test offer") {
  return new Promise((resolve, reject) => {
    if (!client.steamID) {
      reject(new Error("Steam bot is not logged in"));
      return;
    }

    if (!Array.isArray(itemsToGive) || itemsToGive.length === 0) {
      reject(new Error("itemsToGive must contain at least one item"));
      return;
    }

    const offer = manager.createOffer(tradeUrl);

    for (const item of itemsToGive) {
      offer.addMyItem({
        appid: Number(item.appId),
        contextid: String(item.contextId),
        assetid: String(item.assetId),
      });
    }

    offer.setMessage(message);

    offer.send((err, status) => {
      if (err) {
        reject(err);
        return;
      }

      resolve({
        tradeOfferId: offer.id,
        status,
      });
    });
  });
}

function getTradeOffer(tradeOfferId) {
  return new Promise((resolve, reject) => {
    manager.getOffer(tradeOfferId, (err, offer) => {
      if (err) {
        reject(err);
        return;
      }

      resolve({
        id: offer.id,
        state: offer.state,
        stateName: TradeOfferManager.ETradeOfferState[offer.state],
        partner: offer.partner ? offer.partner.getSteamID64() : null,
        message: offer.message,
        created: offer.created,
        updated: offer.updated,
        expires: offer.expires,
        itemsToGive: offer.itemsToGive.map((item) => ({
          assetId: item.assetid,
          appId: item.appid,
          contextId: item.contextid,
          marketHashName: item.market_hash_name,
        })),
        itemsToReceive: offer.itemsToReceive.map((item) => ({
          assetId: item.assetid,
          appId: item.appid,
          contextId: item.contextid,
          marketHashName: item.market_hash_name,
        })),
      });
    });
  });
}

function cancelTradeOffer(tradeOfferId) {
  return new Promise((resolve, reject) => {
    manager.getOffer(tradeOfferId, (err, offer) => {
      if (err) {
        reject(err);
        return;
      }

      offer.cancel((cancelErr) => {
        if (cancelErr) {
          reject(cancelErr);
          return;
        }

        resolve({
          tradeOfferId: offer.id,
          status: "canceled",
        });
      });
    });
  });
}

function getTradeOffers() {
  return new Promise((resolve, reject) => {
    manager.getOffers(
      TradeOfferManager.EOfferFilter.ActiveOnly,
      (err, sent, received) => {
        if (err) {
          reject(err);
          return;
        }

        resolve({
          sent: sent.map((offer) => ({
            id: offer.id,
            state: offer.state,
            stateName: TradeOfferManager.ETradeOfferState[offer.state],
            partner: offer.partner
              ? offer.partner.getSteamID64()
              : null,
            created: offer.created,
            updated: offer.updated,
          })),
          received: received.map((offer) => ({
            id: offer.id,
            state: offer.state,
            stateName: TradeOfferManager.ETradeOfferState[offer.state],
            partner: offer.partner
              ? offer.partner.getSteamID64()
              : null,
            created: offer.created,
            updated: offer.updated,
          })),
        });
      }
    );
  });
}

module.exports = {
  sendTradeOffer,
  getTradeOffer,
  cancelTradeOffer,
  getTradeOffers,
};