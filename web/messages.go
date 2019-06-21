package web

import "leistungsnachweis-graphiker/problem"

const (
	PostImage   = "PostImage"
	Coordinates = "Coordinates"
	Status      = "Status"
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

type StatusMessageData struct {
	Status problem.Status `json:"status"`
}
