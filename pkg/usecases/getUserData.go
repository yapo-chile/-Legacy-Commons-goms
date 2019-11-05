package usecases

import (
	"fmt"
)

// GetUserDataInteractor defines the interactor
type GetUserDataInteractor struct {
	UserProfileRepository UserProfileRepository
	Logger                GetUserPrometheusDefaultLogger
}

// UserProfilerometheusLogger logs UserProfilerometheusLogger events
type GetUserPrometheusDefaultLogger interface {
	LogBadInput(mail string)
}

// GetUser retrieves the basic data of a user given a mail

func (interactor *GetUserDataInteractor) GetUser(mail string) (UserBasicData, error) {

	userProfile, err := interactor.UserProfileRepository.GetUserProfileData(mail)
	if err != nil {
		interactor.Logger.LogBadInput(mail)
		return userProfile, fmt.Errorf("cannot retrieve the user's profile error %+v", err)
	}
	return userProfile, err
}
