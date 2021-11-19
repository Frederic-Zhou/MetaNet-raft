package main

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLanscan(t *testing.T) {

	interfaces, _ := net.Interfaces()

	jsondata, _ := json.Marshal(interfaces)
	logrus.Info(string(jsondata))

}
