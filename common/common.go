package common

// TODO: Move all the common serializer-deserializer stuff here

type DetectionData struct {
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	ClassId   int    `json:"class_id"`
	ClassName string `json:"class_name"`
}
type Metadata struct {
	ID                 int             `json:"id"`
	NumberOfDetections int             `json:"num_of_detections"`
	Detections         []DetectionData `json:"detections"`
}

type FrameData struct {
	ID        int    `json:"id"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	MatType   int    `json:"mat_type"`
	ImageData []byte `json:"data"`
}

const (
	AvroFrameDataScheme = `
		{
			"type" : "record",
			"namespace" : "VAMP",
			"name" : "framedata",
			"fields" : [
				{ "name" : "id" , "type" : "int" },
				{ "name" : "width" , "type" : "int" },
				{ "name" : "height" , "type" : "int" },
				{ "name" : "mat_type" , "type" : "int" },
				{ "name" : "data" , "type" : "bytes" }
			]
		}
	`

	AvroMetadataScheme = `
        {
            "type" : "record",
            "namespace" : "VAMP",
            "name" : "metadata",
            "fields" : [
                { "name" : "id" , "type" : "int" },
                { "name" : "num_of_detections" , "type" : "int" },
                {
                    "name" : "detections",
                    "type" : {
                        "type" : "array",
                        "items" : {
                            "type" : "record",
                            "name" : "detectiondata",
                            "fields" : [
                                { "name" : "x" , "type" : "int" },
                                { "name" : "y" , "type" : "int" },
                                { "name" : "width" , "type" : "int" },
                                { "name" : "height" , "type" : "int" },
                                { "name" : "class_id" , "type" : "int" },
                                { "name" : "class_name" , "type" : "string" }
                            ]
                        }
                    }
                }
			]
		}
	`
)
