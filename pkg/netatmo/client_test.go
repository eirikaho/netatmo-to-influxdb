package netatmo

import (
	"fmt"
	"testing"
)

func TestGetStationData(t *testing.T) {
	c, err := GetNewClient()
	if err != nil {
		t.Fatal(err)
	}

	sd, err := c.GetStationData()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v", sd)
}
