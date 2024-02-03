package statistic

import (
	"fmt"

	"gonum.org/v1/gonum/stat"
)

func Variance(data []float64, weights []float64) (float64, error) {

	if weights != nil {
		if len(weights) != len(data) {
			return 0, fmt.Errorf("error in 'Variance()'. slice of weights and data must have same lengths. slice length weights: %d. slice length data: %d", len(weights), len(data))
		}
	}

	variance := stat.Variance(data, weights)
	return variance, nil
}
