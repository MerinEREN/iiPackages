package video

import "time"

// Status is active or deactive
type Video struct {
	Path         string    `json:"path"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Uploaded     time.Time `json:"uploaded"`
	LastModified time.Time `json:"lastModified"`
	Status       string    `json:"status"`
}

type Videos []Video
