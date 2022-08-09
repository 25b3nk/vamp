# Video Analytics Mock Pipeline (VAMP)

## Introduction
Mock pipeline which will mimic a real-time video analytics pipeline, where we have four services,
1. Source, which fetches frames from video source
2. Detector, which performs objection detection on the frame
3. Postprocess, which uses object detection information to perform further analytics and events
4. Business logic, where we apply business logic on the events generated from tracker


## Setup

Fetch the submodules
```bash
git submodule init
git submodule sync
git submodule update
```

## Using pulsar as the message bus

### Notes

1. Pulsar message batch size is by default 128KB, need to increase it if we are sending image data
1. Pulsar consumer is by default in blocking mode, to set a timeout, we can set timeout in the context we send


## GoCV

### Installation

> Reference: https://gocv.io/getting-started/linux/

```bash
go get -u -d gocv.io/x/gocv
cd $GOPATH/pkg/mod/gocv.io/x/gocv@v0.31.0/
make install
```


## To do

- [x] Test the frames received from webcam
- [x] Use serializer for golang to send frame data along with other metadata in pulsar
- [ ] Populate all the serdes into common utils package
- [ ] Finish `detector`, which takes input from `source` and outputs to new stream, with id to source message
- [ ] Use the detection data from stream in `postprocess` and fetch the corresponding image from `source` stream
- [ ] Finally dump dummy events to new stream from `postprocess`
- [x] Separate out each module to different repo and use git modules to keep in this repo
- [ ] Create a docker-compose to run the pipeline


## References

1. https://pulsar.apache.org/docs/
1. https://pkg.go.dev/github.com/apache/pulsar-client-go/pulsar#section-documentation
1. https://pulsar.apache.org/docs/client-libraries-go/
1. https://pulsar.apache.org/docs/schema-understand
