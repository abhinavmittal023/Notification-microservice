package serializers

import "time"

// GraphData struct is a serializer for the graph data
type GraphData struct {
	UpdatedAt  []time.Time `json:"updated_at"`
	Successful []uint64    `json:"successful"`
	Failed     []uint64    `json:"failed"`
	Total      []uint64    `json:"total"`
}

// SuccessFailedData is a serializer to send notifications data for graph
type SuccessFailedData struct {
	UpdatedAt  time.Time `json:"updated_at"`
	Successful uint64    `json:"successful"`
	Failed     uint64    `json:"failed"`
	Total      uint64    `json:"total"`
}
