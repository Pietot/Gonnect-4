package stats

import "github.com/Pietot/Gonnect-4/utils"

type Stats struct {
	TotalTimeMicroseconds float64
	NodeCount             int64
	MeanTimePerPosition   float64
	NodesPerSecond        float64
}

func (s *Stats) String() string {
	return "Total Time: " + utils.GetTime(s.TotalTimeMicroseconds) +
		", \nNumber of Nodes: " + utils.FormatInt(s.NodeCount) +
		", \nMean Time per Position (µs): " + utils.FormatFloat(s.MeanTimePerPosition) +
		", \nPositions per Second: " + utils.FormatFloat(s.NodesPerSecond)
}
