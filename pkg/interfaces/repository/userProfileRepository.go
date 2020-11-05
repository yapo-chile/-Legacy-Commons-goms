package repository

import (
	"context"
	"crypto/sha1" // nolint: gosec
	"encoding/json"
	"fmt"

	"github.mpi-internal.com/Yapo/goms/pkg/usecases"
)

const (
	errorNoUserDataFound string = "there was no user data found using the email: %s"
	errorUnmarshal       string = "there was an error parsing the user data %s"
)

// UserProfileRepository wrapper struct for the RedisHandler
type UserProfileRepository struct {
	Handler HTTPHandler
	Path    string
}

// NewUserProfileRepository constructor
func NewUserProfileRepository(handler HTTPHandler, path string) usecases.UserProfileRepository {
	return &UserProfileRepository{
		Handler: handler,
		Path:    path,
	}
}

// GetUserProfileData makes a http request to profile service
// to get the user profile data
// it sends the sha1 representation of the provided email
func (repo *UserProfileRepository) GetUserProfileData(
	ctx context.Context,
	email string,
) (usecases.UserBasicData, error) {
	h := sha1.New()        // nolint: gosec
	h.Write([]byte(email)) // nolint: gosec, errcheck
	sha1Email := fmt.Sprintf("%x", h.Sum(nil))
	request := repo.Handler.NewRequest(ctx).SetMethod("GET").SetPath(fmt.Sprintf(repo.Path, sha1Email))

	response, err := repo.Handler.Send(request)
	var JSONResp = response.GetBodyString()

	if err == nil && JSONResp != "" {
		resp := fmt.Sprintf("%s", JSONResp)
		var userData map[string]usecases.UserBasicData

		err := json.Unmarshal([]byte(resp), &userData)
		if err != nil {
			return usecases.UserBasicData{}, fmt.Errorf(errorUnmarshal, email)
		}

		val, ok := userData[sha1Email]
		if !ok {
			return usecases.UserBasicData{}, fmt.Errorf(errorNoUserDataFound, email)
		}

		return val, err
	}

	return usecases.UserBasicData{}, fmt.Errorf(errorNoUserDataFound, email)
}
