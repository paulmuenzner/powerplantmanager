package statistic

import (
	"sort"

	"gonum.org/v1/gonum/stat"
)

func Quantile(data []float64) (q25 float64, q75 float64, iqr float64, lowerBound float64, upperBound float64, outliers []float64, q90 float64, q95 float64) {

	// Sort the data
	sort.Float64s(data)

	// Calculate the first and third quartiles
	q25 = stat.Quantile(0.25, stat.Empirical, data, nil)
	q75 = stat.Quantile(0.75, stat.Empirical, data, nil)
	q90 = stat.Quantile(0.9, stat.Empirical, data, nil)
	q95 = stat.Quantile(0.95, stat.Empirical, data, nil)

	// Calculate the interquartile range (IQR)
	iqr = q75 - q25

	// Define the lower and upper bounds for outliers
	lowerBound = q25 - 1.5*iqr
	upperBound = q75 + 1.5*iqr

	// Identify outliers
	outliers = make([]float64, 0)
	for _, value := range data {
		if value < lowerBound || value > upperBound {
			outliers = append(outliers, value)
		}
	}

	return q25, q75, iqr, lowerBound, upperBound, outliers, q90, q95
}
