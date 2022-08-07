# Mock pipeline for video analytics

## Introduction
Mock pipeline which will mimic a real-time video analytics pipeline, where we have four services,
1. Source, which fetches frames from video source
2. Detector, which performs objection detection on the frame
3. Postprocess, which uses object detection information to perform further analytics and events
4. Business logic, where we apply business logic on the events generated from tracker

## Using pulsar as the message bus

## GoCV

### Installation

> Reference: https://gocv.io/getting-started/linux/

```bash
$ go get -u -d gocv.io/x/gocv
$ cd $GOPATH/pkg/mod/gocv.io/x/gocv@v0.31.0/
$ make install
```


