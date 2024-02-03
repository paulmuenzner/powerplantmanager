package statistic

import (
	"fmt"
	typepackage "github.com/paulmuenzner/powerplantmanager/utils/type"
	"sort"

	"gonum.org/v1/gonum/stat"
)

// Median calculates the median of a slice of numbers.
func Median(data []float64, p float64, weights []float64) (float64, error) {
	if p < 0 || p > 1 {
		return 0, fmt.Errorf("error in 'Median()'. p must be between 0 and 1. p: %f", p)
	}

	// Sort the data
	sort.Float64s(data)

	// Check if the slice is not empty
	isSliceEmpty := typepackage.IsSliceEmpty[float64](data)
	if isSliceEmpty {
		return 0, fmt.Errorf("error in 'Median()'. slice must not be empty. slice length: %d", len(data))
	}
	if weights != nil {
		if len(weights) != len(data) {
			return 0, fmt.Errorf("error in 'Median()'. slice of weights and data must have same lengths. slice length weights: %d. slice length data: %d", len(weights), len(data))
		}
	}

	// Calculate the median using Quantile
	median := stat.Quantile(p, stat.Empirical, data, nil)

	return median, nil
}
