package speedometer

import (
	"fmt"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	//err := Start("enp0s25")
	// err := Start("lo")
	// if err != nil {
	// 	t.Errorf("%s", err)
	// }

	s := Speedometer{Device: "enp0s25"}

	sc := make(chan float64)

	err := s.Start(sc)
	if err != nil {
		t.Errorf("%s", err)
	}

	// speed := <-s.speed
	// fmt.Printf("Speed of %s %fMB/s\n", s.Device, speed)

	// s.Stop()

	speed := <-sc
	fmt.Printf("Speed of %s %fMB/s\n", s.Device, speed)

	time.Sleep(5 * time.Second)
	speed, err = s.GetSpeed()
	if err == nil {
		fmt.Printf("Speed of %s %fMB/s\n", s.Device, speed)
	} else {
		t.Error(err)
	}

}
