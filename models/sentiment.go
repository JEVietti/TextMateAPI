package models

import (
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

//Ratings iso
type Ratings struct {
	RatingScore float32        `json:"score"`
	Length      int            `json:"length"`
	Sentiment   SentimentScore `json:"sentimentScore"`
}

//Sentiment Score is
type SentimentScore struct {
	SentimentStruct *languagepb.AnalyzeSentimentResponse `json: googleStruct`
	Score           float32                              `json: sentimentScore`
	BestMessage     string                               `json: bestMessage`
	WorstMessage    string                               `json: worstMessage`
}

type RatingInput struct {
	Input []string `schema:"Input"`
	Count int      `schema:"Count"`
}
