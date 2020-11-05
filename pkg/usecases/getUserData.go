package usecases

import (
	"context"
	"fmt"
)

// GetUserDataInteractor defines the interactor
type GetUserDataInteractor struct {
	UserProfileRepository UserProfileRepository
}

// GetUser retrieves the basic data of a user given a mail
func (interactor *GetUserDataInteractor) GetUser(ctx context.Context, mail string) (UserBasicData, error) {
	userProfile, err := interactor.UserProfileRepository.GetUserProfileData(ctx, mail)
	if err != nil {
		return userProfile, fmt.Errorf("cannot retrieve the user's profile")
	}

	return userProfile, nil
}
