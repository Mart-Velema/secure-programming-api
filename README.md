# secure-programming-api

## Server environment

### .env file

The environment variables are stored in a .env file in the project root. It contains the following entries:

```env
PORT="3000" # The port to run the API on

# Database
ENCRYPTION_KEY="Fy2mcmCCOS6LsLRpkoeJSFCNVaixHN9oLEdoU+6lu2A=" # A Base64 encoded 32 byte sequence used as the master key to encrypting and decrypting the database entries
ENCRYPTION_PASSCODE="" # A Generic lenght string that will be used as the master password. Overrides the value of ENCRYPTION_KEY
SQLITE_FILE_LOCATION="./guineatrade.db" # The filepath where the SQLite3 database file is located

```

### Database

The database is a generic SQLite3 database. it is managed by [GORM](https://gorm.io/docs/), which automatically adds in the requried columns and 