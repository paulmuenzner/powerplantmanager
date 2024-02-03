package statistic

import (
	"fmt"

	"gonum.org/v1/gonum/stat"
)

func Correlation(x, y []float64) (float64, error) {
	// Ensure both slices have the same length
	if len(x) != len(y) {
		return 0, fmt.Errorf("error in 'Correlation()'. slice of weights and data must have same lengths. slice length x: %d. slice length y: %d", len(x), len(y))
	}

	// Calculate the correlation coefficient
	corr := stat.Correlation(x, y, nil)

	return corr, nil
}
