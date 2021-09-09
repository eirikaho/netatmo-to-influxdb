package influxdb

import (
	"time"
	"os"
	"fmt"
	"github.com/eirikaho/netatmo-influxdb/pkg/netatmo"
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/sirupsen/logrus"
)

func Write(devices []netatmo.Device) error {
	protocol := os.Getenv("INFLUXDB_PROTOCOL")
	host := os.Getenv("INFLUXDB_HOST")
	port := os.Getenv("INFLUXDB_PORT")

	c, err := newClient(fmt.Sprintf("%s://%s:%s", protocol, host, port))
	if err != nil {
		return err
	}
	defer c.Close()

	bps, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: "netatmo",
	})
	if err != nil {
		return err
	}

	bps.AddPoints(pointsFromDevices(devices))

	return c.Write(bps)
}

func pointsFromDevices(devices []netatmo.Device) ([]*client.Point) {
	var pts []*client.Point
	for _, d := range devices {
		pts = append(pts, pointsFromDevice(d.Name, d.DataType, d.Data, d.Modules...)...)
	}
	return pts
}

func pointsFromDevice(name string, datatypes []string, data map[string]interface{}, modules ...netatmo.Module) []*client.Point {
	var pts []*client.Point
	for _, dt := range datatypes {
		p, err := client.NewPoint(
			dt,
			map[string]string{},
			map[string]interface{}{
				"module": name,
				"value":  data[dt],
			},
			time.Now(),
		)
		if err != nil {
			logrus.Error(err)
			continue
		}
		pts = append(pts, p)
	}
	pts = append(pts, pointsFromModules(modules)...)
	return pts
}

func pointsFromModules(modules []netatmo.Module) []*client.Point {
	var pts []*client.Point
	for _, m := range modules {
		pts = append(pts, pointsFromDevice(m.Name, m.DataType, m.Data)...)
	}
	return pts
}
