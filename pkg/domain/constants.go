package domain

// MetricType is the datatype to represent a prometheus metric type
type MetricType int

// Metrics exposer constants
const (
	// BadInputError represents counter of bad input errors
	BadInputError MetricType = iota
	// RepositoryError represents counter of repository errors
	RepositoryError
)
