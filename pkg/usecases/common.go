package usecases

//GomsRepository interface that represents all the methods available to
// consume goms microservice
type GomsRepository interface {
	GetHealthcheck() (string, error)
}
