package influxdb

import (
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/sirupsen/logrus"
)

func newClient(serverUrl string) (client.Client, error) {
	logrus.Info("Connecting to influxdb ...")
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: serverUrl,
	})
	return c, err
}
