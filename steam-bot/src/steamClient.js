const SteamUser = require("steam-user");

const client = new SteamUser();

function getSteamClientStatus() {
  return {
    clientCreated: !!client,
    loggedOn: client.steamID !== null,
  };
}

module.exports = {
  client,
  getSteamClientStatus,
};