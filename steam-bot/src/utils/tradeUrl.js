function parseTradeUrl(tradeUrl) {
  try {
    const url = new URL(tradeUrl);

    const partner = url.searchParams.get("partner");
    const token = url.searchParams.get("token");

    if (!partner || !token) {
      return {
        valid: false,
        error: "Trade URL must contain partner and token",
      };
    }

    return {
      valid: true,
      partner,
      token,
    };
  } catch {
    return {
      valid: false,
      error: "Invalid trade URL",
    };
  }
}

module.exports = {
  parseTradeUrl,
};