package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tarm/serial"
)

var (
	fSerialPort     string
	fVraptorAPIHost string

	vraptorAPIURL string
	vraptorToken  string
)

func main() {
	flag.StringVar(&fSerialPort, "serial", "/dev/ttyUSB0", "Serial port to use")
	flag.StringVar(&fVraptorAPIHost, "vraptor", "http://localhost:5000", "Vraptor API Host")
	flag.Parse()

	var err error
	vraptorAPIURL = fVraptorAPIHost + "/api"
	vraptorToken, err = getToken(os.Getenv("USER"), os.Getenv("PASS"))
	if err != nil {
		log.Fatal(err)
	}

	gauge, err := serial.OpenPort(&serial.Config{Name: fSerialPort, Baud: 9600})
	if err != nil {
		log.Fatal(err)
	}

	vraptorRsrc := &VraptorResource{}
	vraptorTemp := &VraptorTemperature{}

	tkr := time.NewTicker(5 * time.Second)
	defer tkr.Stop()
	for range tkr.C {
		if err := vraptorRsrc.Get(); err != nil {
			log.Println(err)
			continue
		}
		if err := vraptorTemp.Get(); err != nil {
			log.Println(err)
			continue
		}

		cpuUsage := int(vraptorRsrc.CPU.CPUUsage + 0.5)
		temp := int(vraptorTemp.Value + 0.5)
		if _, err := fmt.Fprintf(gauge, "A=%d,%d\r\n", cpuUsage, temp); err != nil {
			log.Println(err)
			continue
		}
	}
}
