package usecases

import (
	"fmt"
)

type UserProfileInteractor struct {
	UserProfileRepository UserProfileRepository
}

type UserBasicData struct {
	Nombre string
	Sexo   string
	Rut    string
}
type UserProfilerometheusLogger interface {
	LogBadInput(string)
	LogRepositoryError(string, UserBasicData, error)
}

// GetInput returns NotificationInput data kind

func (interactor *UserProfileInteractor) GetUser(mail string) (interface{}, error) {

	userProfile, err := interactor.UserProfileRepository.GetUserProfileData(mail)
	if err != nil {
		return nil, fmt.Errorf("Cannot retrieve the user's profile error %+v", err)
	}
	return userProfile, err

}
