package main

import (
	"gocv.io/x/gocv"
)

func main() {
	webcam, _ := gocv.VideoCaptureDevice(0)
	img := gocv.NewMat()
    webcam.Read(&img)
    gocv.IMWrite("tmp.jpg", img)
}
