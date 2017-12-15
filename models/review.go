package models

import (
	"time"
)

type Review struct {
	Reviewer string    `json:"reviewer"`
	Review   string    `json:"review"`
	Posted   time.Time `json:"posted"`
	Edited   bool      `json:"edited"`
	Vote     int16     `json:"vote"`
}

type Reviews []Review
