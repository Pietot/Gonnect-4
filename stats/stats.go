package stats

import "github.com/Pietot/Gonnect-4/utils"

type Stats struct {
	TotalTimeMicroseconds float64
	NumberPositions       int64
	MeanTimePerPosition   float64
	PositionsPerSecond    float64
}

func (s *Stats) String() string {
	return "Total Time (µs): " + utils.FormatFloat(s.TotalTimeMicroseconds) +
		", \nNumber of Positions: " + utils.FormatInt(s.NumberPositions) +
		", \nMean Time per Position (µs): " + utils.FormatFloat(s.MeanTimePerPosition) +
		", \nPositions per Second: " + utils.FormatFloat(s.PositionsPerSecond)
}
