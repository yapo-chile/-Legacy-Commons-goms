package usecases

import "time"

// GomsRepository interface that represents all the methods available to
// interact with goms microservice
type GomsRepository interface {
	GetHealthcheck() (string, error)
}

// UserData contains the user information needed to send a notification
// through knocker api
type UserData struct {
	RoleID       int       `json:"roleId"`
	CountryID    int       `json:"countryId"`
	RegionID     int       `json:"regionId"`
	CommuneID    int       `json:"communeId"`
	UUID         string    `json:"uuid"`
	CreationDate time.Time `json:"creationDate"`
}

// UserRepository defines the methods that a User repository should have
type UserRepository interface {
	// GetUserData gets the user data based on his email
	GetUserData(email string) (UserData, error)
}
type UserProfileRepository interface {
	// GetUserData gets the user data based on his email
	GetUserProfileData(email string) (UserProfileData, error)
	GetUserPrivateData(accSession string) (UserPrivateData, error)
}

// UserData contains the user information retrieved from the user profile
type UserProfileData struct {
	RoleID       int       `json:"roleId"`
	CountryID    int       `json:"countryId"`
	RegionID     int       `json:"regionId"`
	CommuneID    int       `json:"communeId"`
	UUID         string    `json:"uuid"`
	CreationDate time.Time `json:"creationDate"`
}

// UserPrivateData contains the private user information retrieved from the user profile
type UserPrivateData struct {
	Email string `json:"email"`
}
