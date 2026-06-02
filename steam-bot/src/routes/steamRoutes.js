const express = require("express");

const {
  getSteamClientStatus,
  loginToSteam,
} = require("../services/steamClient");

const { getBotInventory } = require("../services/inventoryService");

const router = express.Router();

router.get("/status", (req, res) => {
  res.json(getSteamClientStatus());
});

router.post("/login", (req, res) => {
  const { authCode } = req.body;
  const result = loginToSteam(authCode);
  res.json(result);
});

router.get("/inventory", async (req, res) => {
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

module.exports = router;