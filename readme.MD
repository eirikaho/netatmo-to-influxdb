# Netatmo-influxdb

Simple go app that polls netatmo API for data points and persists to influxdb
Requires a running instance of influxdb

## Environment variables
| Variable name | example |
| ---|---|
|INFLUXDB_HOST|influxdb
|INFLUXDB_PROTOCOL|http
|INFLUXDB_PORT|8086
|NETATMO_CLIENT_ID|
|NETATMO_CLIENT_SECRET|
|NETATMO_USERNAME|
|NETATMO_PASSWORD|
|INTERVAL_SECONDS|600

See [Netatmo Connect Docs](https://dev.netatmo.com/apidocumentation/oauth) for documentation regarding netatmo auth variables