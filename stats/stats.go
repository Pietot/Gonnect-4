package stats

import "github.com/Pietot/Gonnect-4/utils"

type Stats struct {
	TotalTimeMs         float64
	NumPositions        int64
	MeanTimePerPosition float64
	PositionsPerSecond  float64
}

func (s *Stats) String() string {
	return "Total Time (ms): " + utils.FormatFloat(s.TotalTimeMs) +
		", \nNumber of Positions: " + utils.FormatInt(s.NumPositions) +
		", \nMean Time per Position (ms): " + utils.FormatFloat(s.MeanTimePerPosition) +
		", \nPositions per Second: " + utils.FormatFloat(s.PositionsPerSecond)
}
