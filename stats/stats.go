package stats

import "github.com/Pietot/Gonnect-4/utils"

type Stats struct {
	TotalTimeMicroseconds float64
	NodeCount             uint64
	MeanTimePerNode       float64
	NodesPerSecond        uint64
}

func (s *Stats) String() string {
	return "Total Time: " + utils.GetTime(s.TotalTimeMicroseconds) +
		", \nNumber of Nodes: " + utils.FormatUint64(s.NodeCount) +
		", \nMean Time per Node (µs): " + utils.FormatFloat(s.MeanTimePerNode) +
		", \nNodes per Second: " + utils.FormatUint64(s.NodesPerSecond)
}
