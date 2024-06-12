# Netatmo-influxdb

Simple go app that polls netatmo API for data points and persists to influxdb
Requires a running instance of influxdb

## Environment variables
NB! refresh token are read from `/var/lib/netatmo/refresh_token.txt` and should be mounted as a volume.
If the file does not exist, the app will use the environment variables and write the token to the file.
To persist the token, mount the volume to a location on the host machine

| Variable name         | example  |
|-----------------------|----------|
| INFLUXDB_HOST         | influxdb |
| INFLUXDB_PROTOCOL     | http     |
| INFLUXDB_PORT         | 8086     |
| NETATMO_CLIENT_ID     |          |
| NETATMO_CLIENT_SECRET |          |
| NETATMO_REFRESH_TOKEN |          |
| INTERVAL_SECONDS      | 600      |

See [Netatmo Connect Docs](https://dev.netatmo.com/apidocumentation/oauth) for documentation regarding netatmo auth variables