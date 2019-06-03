package usecases

// GomsRepository interface that represents all the methods available to
// interact with goms microservice
type GomsRepository interface {
	GetHealthcheck() (string, error)
}
