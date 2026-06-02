const express = require("express");
require("dotenv").config();

const healthRoutes = require("./routes/healthRoutes");
const steamRoutes = require("./routes/steamRoutes");
const tradeRoutes = require("./routes/tradeRoutes");

const { requireBotApiKey } = require("./middleware/authMiddleware");

const app = express();
const PORT = process.env.BOT_PORT || 3001;

app.use(express.json());

app.use("/", healthRoutes);

app.use("/steam", requireBotApiKey, steamRoutes);
app.use("/steam/trade-offers", requireBotApiKey, tradeRoutes);

app.listen(PORT, () => {
  console.log(`Steam bot service running on port ${PORT}`);
});