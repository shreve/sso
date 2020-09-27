package sso

import (
	"errors"
)

var EmptyCredsError = errors.New("A username and password must be supplied")

var AlreadyRegisteredError = errors.New("That username has already been registered")
var IncorrectPasswordError = errors.New("That password is incorrect")

var NotSignedInError = errors.New("There is not a signed in user.")

var UserNotFoundError = errors.New("That user does not exist.")
var UserNotCreatedError = errors.New("There was a problem creating that user.")
