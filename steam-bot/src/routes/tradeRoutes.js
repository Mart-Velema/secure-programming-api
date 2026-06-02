const express = require("express");

const { parseTradeUrl } = require("../utils/tradeUrl");
const { getSteamClientStatus } = require("../services/steamClient");

const {
  sendTradeOffer,
  getTradeOffer,
} = require("../services/tradeService");

const router = express.Router();

router.post("/dry-run", (req, res) => {
  const { tradeUrl, items } = req.body;

  if (!tradeUrl) {
    return res.status(400).json({
      ok: false,
      error: "tradeUrl is required",
    });
  }

  const tradeUrlResult = parseTradeUrl(tradeUrl);

  if (!tradeUrlResult.valid) {
    return res.status(400).json({
      ok: false,
      error: tradeUrlResult.error,
    });
  }

  const steamStatus = getSteamClientStatus();

  if (!steamStatus.loggedOn) {
    return res.status(400).json({
      ok: false,
      error: "Steam bot is not logged in",
    });
  }

  res.json({
    ok: true,
    message: "Dry run successful. No trade was sent.",
    partner: tradeUrlResult.partner,
    token: tradeUrlResult.token,
    itemCount: Array.isArray(items) ? items.length : 0,
  });
});

router.post("/", async (req, res) => {
  try {
    const { tradeUrl, itemsToGive, message } = req.body;

    if (!tradeUrl) {
      return res.status(400).json({
        ok: false,
        error: "tradeUrl is required",
      });
    }

    const tradeUrlResult = parseTradeUrl(tradeUrl);

    if (!tradeUrlResult.valid) {
      return res.status(400).json({
        ok: false,
        error: tradeUrlResult.error,
      });
    }

    if (!Array.isArray(itemsToGive) || itemsToGive.length === 0) {
      return res.status(400).json({
        ok: false,
        error: "itemsToGive must contain at least one item",
      });
    }

    const result = await sendTradeOffer(
      tradeUrl,
      itemsToGive,
      message || "GuineaTrade test offer"
    );

    res.json({
      ok: true,
      ...result,
    });
  } catch (error) {
    res.status(500).json({
      ok: false,
      error: error.message,
    });
  }
});

router.get("/:tradeOfferId", async (req, res) => {
  try {
    const offer = await getTradeOffer(req.params.tradeOfferId);

    res.json({
      ok: true,
      offer,
    });
  } catch (error) {
    res.status(500).json({
      ok: false,
      error: error.message,
    });
  }
});

module.exports = router;