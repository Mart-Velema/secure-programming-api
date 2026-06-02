const express = require("express");

const router = express.Router();

router.get("/health", (req, res) => {
  res.json({
    status: "ok",
    service: "steam-bot",
    port: process.env.BOT_PORT || 3001,
  });
});

router.get("/config", (req, res) => {
  res.json({
    steamConfigured: !!process.env.STEAM_USERNAME,
  });
});

module.exports = router;