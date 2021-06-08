package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/asch/go-tcmu"
)

type null struct{
}

func (n *null) ReadAt(b []byte, off int64) (int, error) {
	return len(b), nil
}

func (n *null) WriteAt(b []byte, off int64) (int, error) {
	return len(b), nil
}

func main() {
	null := &null{}

	handler := tcmu.BasicSCSIHandler(null)
	handler.VolumeName = "TCMU NULL Device"
	handler.DataSizes.VolumeSize = 1024*1024*1024*1024
	d, err := tcmu.OpenTCMUDevice("/dev/tcmunull", handler)
	if err != nil {
		die("couldn't tcmu: %v", err)
	}
	defer d.Close()

	mainClose := make(chan bool)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		for range signalChan {
			fmt.Println("\nReceived an interrupt, stopping services...")
			close(mainClose)
		}
	}()
	<-mainClose
}

func die(why string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, why+"\n", args...)
	os.Exit(1)
}
