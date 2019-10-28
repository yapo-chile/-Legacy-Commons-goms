package usecases

import (
	"fmt"
)

// UserProfileInteractor defines the interactor
type UserProfileInteractor struct {
	UserProfileRepository UserProfileRepository
}

// UserProfilerometheusLogger logs UserProfilerometheusLogger events
type UserProfilerometheusLogger interface {
	LogBadInput(string)
	LogRepositoryError(string, UserBasicData, error)
}

// GetUser retrieves the basic data of a user given a mail

func (interactor *UserProfileInteractor) GetUser(mail string) (UserBasicData, error) {

	userProfile, err := interactor.UserProfileRepository.GetUserProfileData(mail)

	if err != nil {
		return userProfile, fmt.Errorf("Cannot retrieve the user's profile error %+v", err)
	}
	return userProfile, err

}
