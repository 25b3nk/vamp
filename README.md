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


## To do

- [x] Test the frames received from webcam
- [x] Use serializer for golang to send frame data along with other metadata in pulsar
- [ ] Populate all the serdes into common utils package
- [ ] Finish `detector`, which takes input from `source` and outputs to new stream, with id to source message
- [ ] Use the detection data from stream in `postprocess` and fetch the corresponding image from `source` stream
- [ ] Finally dump dummy events to new stream from `postprocess`


## References

1. https://pulsar.apache.org/docs/
1. https://pkg.go.dev/github.com/apache/pulsar-client-go/pulsar#section-documentation
1. https://pulsar.apache.org/docs/client-libraries-go/
1. https://pulsar.apache.org/docs/schema-understand
