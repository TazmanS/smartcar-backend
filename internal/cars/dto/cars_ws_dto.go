package dto

import "time"

type CVDetectionEvent struct {
	Type      string     `json:"type"`
	Timestamp time.Time  `json:"timestamp"`
	Objects   []CVObject `json:"objects"`
}

type CVObject struct {
	ClassName  string  `json:"class_name"`
	Confidence float64 `json:"confidence"`
	Left       float64 `json:"left"`
	Top        float64 `json:"top"`
	Right      float64 `json:"right"`
	Bottom     float64 `json:"bottom"`
}
