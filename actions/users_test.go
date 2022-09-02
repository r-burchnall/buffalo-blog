package actions

import (
	"social_api/models"
)

func (as *ActionSuite) Test_Users_Create() {
	count, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(0, count)

	u := &models.User{
		Email:                "mark@example.com",
		Password:             "password",
		PasswordConfirmation: "password",
	}

	res := as.JSON("/users").Post(u)
	as.Equal(200, res.Code)

	count, err = as.DB.Count("users")
	as.NoError(err)
	as.Equal(1, count)
}
