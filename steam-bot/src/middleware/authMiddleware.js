function requireBotApiKey(req, res, next) {
  const expectedApiKey = process.env.BOT_API_KEY;
  const providedApiKey = req.header("X-API-Key");

  if (!expectedApiKey) {
    return res.status(500).json({
      ok: false,
      error: "BOT_API_KEY is not configured",
    });
  }

  if (!providedApiKey || providedApiKey !== expectedApiKey) {
    return res.status(401).json({
      ok: false,
      error: "Unauthorized",
    });
  }

  next();
}

module.exports = {
  requireBotApiKey,
};