package repository

import (
	"crypto/sha1" // nolint: gosec
	"encoding/json"
	"fmt"

	"github.mpi-internal.com/Yapo/goms/pkg/usecases"
)

const (
	errorNoUserDataFound  string = "There was no user data found using the email: %s"
	errorNoDataAccSession string = "There was an error getting data from profile: %s"
)

// UserRepository wrapper struct for the RedisHandler
type UserProfileRepository struct {
	Handler HTTPHandler
	Path    string
}

// NewUserRepository constructor for a NotificationRepository
func NewUserProfileRepository(handler HTTPHandler, path string) usecases.UserProfileRepository {
	return &UserProfileRepository{
		Handler: handler,
		Path:    path,
	}
}

// GetUserProfileData makes a http request to profile service
// to get the user profile data
// it sends the sha1 representation of the provided email
func (repo *UserProfileRepository) GetUserProfileData(email string) (usecases.UserProfileData, error) {
	h := sha1.New()        // nolint: gosec
	h.Write([]byte(email)) // nolint: gosec
	sha1Email := fmt.Sprintf("%x", h.Sum(nil))
	request := repo.Handler.NewRequest().SetMethod("GET").SetPath(fmt.Sprintf(repo.Path, sha1Email))
	JSONResp, err := repo.Handler.Send(request)
	if err == nil && JSONResp != "" {
		resp := fmt.Sprintf("%s", JSONResp)
		var userData map[string]usecases.UserProfileData
		err := json.Unmarshal([]byte(resp), &userData)
		val, ok := userData[sha1Email]
		if !ok {
			return usecases.UserProfileData{}, fmt.Errorf(errorNoUserDataFound, email)
		}
		return val, err
	}
	return usecases.UserProfileData{}, fmt.Errorf(errorNoUserDataFound, email)
}

// GetUserPrivateData makes a http request to profile service
// to get the user private profile data
// it sends the Authorization header
func (repo *UserProfileRepository) GetUserPrivateData(accSession string) (usecases.UserPrivateData, error) {
	request := repo.Handler.NewRequest().
		SetMethod("GET").
		SetHeaders(map[string]string{"Authorization": accSession}).
		SetPath(repo.Path)
	JSONResp, err := repo.Handler.Send(request)
	if err == nil && JSONResp != "" {
		resp := fmt.Sprintf("%s", JSONResp)
		var userData usecases.UserPrivateData
		err := json.Unmarshal([]byte(resp), &userData)
		return userData, err
	}
	return usecases.UserPrivateData{}, fmt.Errorf(errorNoDataAccSession, accSession)
}
