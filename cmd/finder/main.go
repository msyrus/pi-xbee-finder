package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	mhttp "github.com/msyrus/go/http"
	"github.com/msyrus/pi-xbee-finder/version"
	"github.com/msyrus/pi-xbee-finder/xbee"
	serial "go.bug.st/serial.v1"
)

const usage = `
Usage:	finder [options] <tty port device>
Options:
	-b <baudrate> (default 115200)
	-p <web port> (default 8080)
`

func main() {
	fmt.Println("Hello", version.Version)

	var baudR, webP int
	flag.IntVar(&baudR, "b", 115200, "baud rate")
	flag.IntVar(&webP, "p", 8080, "web port")

	flag.Parse()

	serialP := flag.Arg(0)
	if serialP == "" {
		fmt.Println(usage)
		os.Exit(1)
	}

	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	found := false
	for _, port := range ports {
		if port == serialP {
			found = true
		}
	}

	if !found {
		fmt.Println(serialP, "not found")
		fmt.Println("Available serial ports")
		for _, port := range ports {
			fmt.Println(port)
		}
	}

	devSig := map[string]bool{}
	dev := 1
	var devErr error

	port, err := serial.Open(serialP, &serial.Mode{BaudRate: baudR})
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	dataCh := make(chan *xbee.Frame, 100)
	errCh := make(chan error, 100)

	go func() {
		for err := range errCh {
			log.Println(err)
			devErr = err
		}
	}()

	go func() {
		sum := map[string]int{}
		cnt := map[string]int{}
		for data := range dataCh {
			// fmt.Printf("%#v\n", data)
			id := data.DeviceID
			val := 0
			if data.DataD[dev] {
				val = 1
			}
			sum[id] = sum[id] + val
			cnt[id]++
			if cnt[id] >= 10 {
				devSig[id] = (float64(sum[id])/float64(cnt[id]) > 0.5)
				sum[id] = 0
				cnt[id] = 0
			}
		}

	}()

	go func() {
		for {
			buff := make([]byte, 1000)
			n, err := port.Read(buff)
			if err != nil {
				log.Fatal(err)
			}
			// fmt.Println("Read", n, "bytes")
			// fmt.Println("Received:", hex.EncodeToString(buff[:n]))
			f, err := xbee.ParseFrame(buff[:n])
			if err != nil {
				errCh <- err
			} else {
				dataCh <- f
			}
		}
	}()

	webHandler := func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"data":  devSig,
			"error": devErr,
		}
		devErr = nil
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Println(err)
		}
	}

	srvr := http.Server{
		Addr:    ":" + strconv.Itoa(webP),
		Handler: http.HandlerFunc(webHandler),
	}

	if err := mhttp.ManageServer(&srvr, 30*time.Second); err != nil {
		log.Fatal(err)
	}
}
