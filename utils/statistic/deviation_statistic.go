package statistic

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/stat"
)

func StandardDeviation(data []float64, weights []float64) (float64, error) {

	if weights != nil {
		if len(weights) != len(data) {
			return 0, fmt.Errorf("error in 'StandardDeviation()'. slice of weights and data must have same lengths. slice length weights: %d. slice length data: %d", len(weights), len(data))
		}
	}

	// Calculate variance
	variance := stat.StdDev(data, weights)

	// Calculate standard deviation (square root of variance)
	stdDev := math.Sqrt(variance)

	return stdDev, nil
}
