package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tarm/serial"
)

var (
	fSerialPort     string
	fVraptorAPIHost string
	fIntervalStr    string

	vraptorAPIURL string
	vraptorToken  string
)

func main() {
	flag.StringVar(&fSerialPort, "s", "/dev/ttyUSB0", "Serial port to use")
	flag.StringVar(&fVraptorAPIHost, "h", "http://localhost:5000", "Vraptor API Host")
	flag.StringVar(&fIntervalStr, "i", "5s", "Interval to update gauges")
	flag.Parse()

	var err error

	log.Println("Starting up...")
	vraptorAPIURL = fVraptorAPIHost + "/api"
	vraptorToken, err = getToken(os.Getenv("USER"), os.Getenv("PASS"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Opening serial port...")
	gauge, err := serial.OpenPort(&serial.Config{Name: fSerialPort, Baud: 9600})
	if err != nil {
		log.Fatal(err)
	}

	vraptorRsrc := &VraptorResource{}
	vraptorTemp := &VraptorTemperature{}

	i, err := time.ParseDuration(fIntervalStr)
	if err != nil {
		log.Fatal(err)
	}

	tkr := time.NewTicker(i)
	defer tkr.Stop()

	log.Println("Updating gauges...")
	go func() {
		for range tkr.C {
			if err := vraptorRsrc.Get(); err != nil {
				log.Println(err)
				continue
			}
			if err := vraptorTemp.Get(); err != nil {
				log.Println(err)
				continue
			}

			cpuUsage := scale(vraptorRsrc.CPU.CPUUsage, 0, 99, 0, 100)
			temp := scale(vraptorTemp.Value, 20, 60, 0, 100)
			if _, err := fmt.Fprintf(
				gauge,
				"A=%d,%d\r\n", int(cpuUsage+0.5), int(temp+0.5),
			); err != nil {
				log.Println(err)
				continue
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("Shutting down...")
	if _, err := fmt.Fprintf(gauge, "A=0,0\r\n"); err != nil {
		log.Println(err)
	}
}

func scale(v, fromLow, fromHigh, toLow, toHigh float64) float64 {
	ret := (v-fromLow)/(fromHigh-fromLow)*(toHigh-toLow) + toLow
	if ret < toLow {
		ret = toLow
	}
	if ret > toHigh {
		ret = toHigh
	}
	return ret
}
