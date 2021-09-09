package influxdb

import (
	"testing"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

func TestNewClient(t *testing.T) {
	c, err := newClient("http://localhost:8086")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	p, err := client.NewPoint(
		"Temperature",
		map[string]string{"my-tag-key": "my-tag-value"},
		map[string]interface{}{
			"module": "indoor",
			"value":  25.5,
		},
		time.Now(),
	)
	if err != nil {
		t.Fatal(err)
	}

	bpc := client.BatchPointsConfig{
		Database: "netatmo",
		// RetentionPolicy:  "default",
	}

	bps, err := client.NewBatchPoints(bpc)
	if err != nil {
		t.Fatal(err)
	}

	bps.AddPoint(p)

	err = c.Write(bps)
	if err != nil {
		t.Fatal(err)
	}
}
