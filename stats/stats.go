package stats

import "github.com/Pietot/Gonnect-4/utils"

type Stats struct {
	TotalTimeMicroseconds float64
	NodeCount             int64
	MeanTimePerNode       float64
	NodesPerSecond        int64
}

func (s *Stats) String() string {
	return "Total Time: " + utils.GetTime(s.TotalTimeMicroseconds) +
		", \nNumber of Nodes: " + utils.FormatInt(s.NodeCount) +
		", \nMean Time per Node (µs): " + utils.FormatFloat(s.MeanTimePerNode) +
		", \nNodes per Second: " + utils.FormatInt(s.NodesPerSecond)
}
