package web

const (
	PostImage   = "PostImage"
	Coordinates = "Coordinates"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type ImageMessageData struct {
	Image string `json:"image"`
}

type CoordinatesMessageData struct {
	Coordinates []int `json:"coordinates"`
}
