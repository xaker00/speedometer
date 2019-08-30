// Package speedometer provides a way to measure network tx speed in realtime.
// This work only on Linux.
package speedometer

import (
	"errors"
	"io/ioutil"
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Speedometer struct {
	stopchan   chan struct{}
	Device     string
	speed      float64
	err        error
	signalChan chan float64
}

// Stop the internal goroutine from collecting speed information.
func (s *Speedometer) Stop() {
	close(s.stopchan)
}

// Return the last measured speed and any errors.
func (s *Speedometer) GetSpeed() (r float64, err error) {
	return s.speed, s.err
}

// Starts the internal goroutine. The Speedometer struct needs to have the
// Device field set to a valid interface. A float64 channel may be passed
// to recive speed updates. Pass nil to not send data on a channel.
// It takes 4 seconds to read the first value and reports the speed in MBits/sec.
// Returns nill if no error.
func (s *Speedometer) Start(signalChan chan float64) (err error) {
	s.err = errors.New("No data")

	f := filepath.Join("/", "sys", "class", "net", s.Device, "statistics", "tx_bytes")

	_, err = ioutil.ReadFile(f)

	if err != nil {
		return err
	}

	s.stopchan = make(chan struct{})

	// no errors detected - setup the loop
	go func() {
		last_value := -1
		bytes_tx_calc := -1.0
		for {

			select {
			default:
				// TODO: do a bit of the work
				bytes_tx, err := ioutil.ReadFile(f)

				if err != nil {
					s.err = err
					return
				}

				bytes_tx_int, err := strconv.Atoi(strings.Trim(string(bytes_tx), "\n\r"))
				if last_value >= 0 {
					bytes_tx_calc = float64(bytes_tx_int - last_value)
				}
				last_value = bytes_tx_int

				if err != nil {
					s.err = err
					return
				}

				if bytes_tx_calc >= 0 {
					s.speed = math.Round((bytes_tx_calc*2.0)/(1024*1024)*100) / 100
					s.err = nil

					if signalChan != nil {
						select {
						case signalChan <- s.speed:
						default:
							// buffer full
						}
					}
				}

			case <-s.stopchan:
				// stop)
				s.err = errors.New("Speedometer stopped")
				return
			}
			time.Sleep(4 * time.Second)
		}
	}()

	return nil
}
