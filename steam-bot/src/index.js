const express = require("express");
require("dotenv").config();

const {
  getSteamClientStatus,
  loginToSteam,
  getBotInventory,
} = require("./steamClient");

const { parseTradeUrl } = require("./tradeUrl");

const app = express();

app.use(express.json());

const PORT = process.env.BOT_PORT || 3001;

app.get("/health", (req, res) => {
  res.json({
    status: "ok",
    service: "steam-bot",
    port: PORT,
  });
});

app.get("/config", (req, res) => {
  res.json({
    steamConfigured: !!process.env.STEAM_USERNAME,
  });
});

app.get("/steam/status", (req, res) => {
  res.json(getSteamClientStatus());
});

app.listen(PORT, () => {
  console.log(`Steam bot service running on port ${PORT}`);
});

app.post("/steam/login", (req, res) => {
  const { authCode } = req.body;

  const result = loginToSteam(authCode);

  res.json(result);
});

app.get("/steam/inventory", async (req, res) => {
  try {
    const appId = Number(req.query.appId || 440);
    const contextId = Number(req.query.contextId || 2);

    const inventory = await getBotInventory(appId, contextId);

    res.json({
      appId,
      contextId,
      count: inventory.length,
      inventory,
    });
  } catch (error) {
    res.status(500).json({
      error: error.message,
    });
  }
});

app.post("/steam/trade-url/validate", (req, res) => {
  const { tradeUrl } = req.body;

  if (!tradeUrl) {
    return res.status(400).json({
      valid: false,
      error: "tradeUrl is required",
    });
  }

  const result = parseTradeUrl(tradeUrl);

  res.status(result.valid ? 200 : 400).json(result);
});

app.post("/steam/trade-offers/dry-run", (req, res) => {
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