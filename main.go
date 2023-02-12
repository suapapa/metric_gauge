package main

import (
	"context"
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
	fSerialPort  string
	fIntervalStr string

	duration time.Duration
)

func main() {
	flag.StringVar(&fSerialPort, "s", "/dev/ttyUSB0", "Serial port to use")
	flag.StringVar(&fIntervalStr, "i", "1s", "Interval to update gauges")
	flag.Parse()

	log.Println("Opening serial port...")
	gauge, err := serial.OpenPort(&serial.Config{Name: fSerialPort, Baud: 9600})
	if err != nil {
		log.Fatal(err)
	}

	duration, err = time.ParseDuration(fIntervalStr)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancelF := context.WithCancel(context.Background())
	defer cancelF()
	memCh := getMemoryUsage(ctx)
	cpuCh := getCpuPercent(ctx)

	log.Println("Updating gauges...")
	go func(ctx context.Context) {
		for {
			select {
			case mem := <-memCh:
				log.Printf("mem: %f", mem)
				if _, err := fmt.Fprintf(gauge, "A2=%d\r\n", int(mem+0.5)); err != nil {
					log.Println(err)
				}
			case cpu := <-cpuCh:
				log.Printf("cpu: %f", cpu)
				if _, err := fmt.Fprintf(gauge, "A1=%d\r\n", int(cpu+0.5)); err != nil {
					log.Println(err)
				}
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("Shutting down...")
	cancelF()
	if _, err := fmt.Fprintf(gauge, "A=0,0\r\n"); err != nil {
		log.Println(err)
	}
}

// func scale(v, fromLow, fromHigh, toLow, toHigh float64) float64 {
// 	ret := (v-fromLow)/(fromHigh-fromLow)*(toHigh-toLow) + toLow
// 	if ret < toLow {
// 		ret = toLow
// 	}
// 	if ret > toHigh {
// 		ret = toHigh
// 	}
// 	return ret
// }
