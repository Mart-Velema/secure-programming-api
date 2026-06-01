const express = require("express");
require("dotenv").config();

const { getSteamClientStatus, loginToSteam } = require("./steamClient");

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