package test

import (
	"math"
	"testing"

	"github.com/Pietot/Gonnect-4/evaluation"
)

type compareEvaluation struct {
	eval1    *evaluation.Evaluation
	eval2    *evaluation.Evaluation
	isBetter bool
}

var compareEvaluations = []compareEvaluation{
	{
		eval1: &evaluation.Evaluation{
			Score:         nil,
			BestMove:      nil,
			RemainingMove: nil,
		},
		eval2: &evaluation.Evaluation{
			Score:         nil,
			BestMove:      nil,
			RemainingMove: nil,
		},
		isBetter: false,
	},
	{
		eval1: &evaluation.Evaluation{
			Score:         evaluation.Float64Ptr(math.Inf(1)),
			BestMove:      nil,
			RemainingMove: nil,
		},
		eval2: &evaluation.Evaluation{
			Score:         nil,
			BestMove:      nil,
			RemainingMove: nil,
		},
		isBetter: true,
	},
	{
		eval1: &evaluation.Evaluation{
			Score:         nil,
			BestMove:      nil,
			RemainingMove: nil,
		},
		eval2: &evaluation.Evaluation{
			Score:         evaluation.Float64Ptr(math.Inf(1)),
			BestMove:      nil,
			RemainingMove: nil,
		},
		isBetter: false,
	},
	{
		eval1: &evaluation.Evaluation{
			Score:         evaluation.Float64Ptr(math.Inf(1)),
			BestMove:      nil,
			RemainingMove: nil,
		},
		eval2: &evaluation.Evaluation{
			Score:         evaluation.Float64Ptr(math.Inf(-1)),
			BestMove:      nil,
			RemainingMove: nil,
		},
		isBetter: true,
	},
	{
		eval1: &evaluation.Evaluation{
			Score:         evaluation.Float64Ptr(math.Inf(-1)),
			BestMove:      nil,
			RemainingMove: nil,
		},
		eval2: &evaluation.Evaluation{
			Score:         evaluation.Float64Ptr(math.Inf(1)),
			BestMove:      nil,
			RemainingMove: nil,
		},
		isBetter: false,
	},
	{
		eval1: &evaluation.Evaluation{
			Score:         evaluation.Float64Ptr(math.Inf(1)),
			BestMove:      nil,
			RemainingMove: evaluation.IntPtr(5),
		},
		eval2: &evaluation.Evaluation{
			Score:         evaluation.Float64Ptr(math.Inf(1)),
			BestMove:      nil,
			RemainingMove: evaluation.IntPtr(10),
		},
		isBetter: true,
	},
	{
		eval1: &evaluation.Evaluation{
			Score:         evaluation.Float64Ptr(math.Inf(1)),
			BestMove:      nil,
			RemainingMove: evaluation.IntPtr(10),
		},
		eval2: &evaluation.Evaluation{
			Score:         evaluation.Float64Ptr(math.Inf(1)),
			BestMove:      nil,
			RemainingMove: evaluation.IntPtr(5),
		},
		isBetter: false,
	},
	{
		eval1: &evaluation.Evaluation{
			Score:         evaluation.Float64Ptr(math.Inf(-1)),
			BestMove:      nil,
			RemainingMove: evaluation.IntPtr(5),
		},
		eval2: &evaluation.Evaluation{
			Score:         evaluation.Float64Ptr(math.Inf(-1)),
			BestMove:      nil,
			RemainingMove: evaluation.IntPtr(10),
		},
		isBetter: false,
	},
	{
		eval1: &evaluation.Evaluation{
			Score:         evaluation.Float64Ptr(math.Inf(-1)),
			BestMove:      nil,
			RemainingMove: evaluation.IntPtr(10),
		},
		eval2: &evaluation.Evaluation{
			Score:         evaluation.Float64Ptr(math.Inf(-1)),
			BestMove:      nil,
			RemainingMove: evaluation.IntPtr(5),
		},
		isBetter: true,
	},
}

func TestEvaluation(t *testing.T) {
	for _, test := range compareEvaluations {
		if test.eval1.IsBetterThan(test.eval2) != test.isBetter {
			t.Errorf("Expected \n%v \n\nto be better than \n\n%v", test.eval1, test.eval2)
		}
	}
}
