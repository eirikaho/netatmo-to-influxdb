package main

import (
	"os"
	"strconv"
	"time"

	"github.com/eirikaho/netatmo-influxdb/pkg/influxdb"
	"github.com/eirikaho/netatmo-influxdb/pkg/netatmo"
	"github.com/sirupsen/logrus"
)

func main() {
	i, err := strconv.Atoi(os.Getenv("INTERVAL_SECONDS"))
	if err != nil {
		panic(err)
	}
	ticker := time.NewTicker(time.Duration(i) * time.Second)
	for range ticker.C {
		err := persistDatapointsFromNetatmo()
		if err != nil {
			logrus.Error(err)
		}
		logrus.Infof("ticker sleeps for %d seconds ...", i)
	}
}

func persistDatapointsFromNetatmo() error {
	logrus.Info("Starting routine ...")
	c, err := getNetatmoClient()
	if err != nil {
		return err
	}
	logrus.Info("Obtained netatmo client")

	sd, err := c.GetStationData()
	if err != nil {
		return err
	}
	logrus.Info("Recieved station data")

	err = influxdb.Write(sd.Devices)
	if err != nil {
		return err
	}
	logrus.Info("Wrote data to influxdb")
	logrus.Info("Routine finished OK")
	return nil
}

func getNetatmoClient() (*netatmo.Client, error) {
	clientId := os.Getenv("NETATMO_CLIENT_ID")
	clientSecret := os.Getenv("NETATMO_CLIENT_SECRET")
	user := os.Getenv("NETATMO_USERNAME")
	pass := os.Getenv("NETATMO_PASSWORD")

	return netatmo.NewClient(clientId, clientSecret, user, pass)
}

