package common

import "time"

// Neighbor P2P debug neighbor format
type Neighbor struct {
	Address  string `json:"address"`
	OutBound bool   `json:"out_bound"`
}

// ReportMsg P2P debug report message format
type ReportMsg struct {
	ReportPeer string    `json:"report_peer"`
	From       string    `json:"from"`
	Time       time.Time `json:"time"`
}

// ErrorResponse error response format
type ErrorResponse struct {
	Error string `json:"error"`
}

func NewResponse(err error) *ErrorResponse {
	return &ErrorResponse{
		Error: err.Error(),
	}
}
