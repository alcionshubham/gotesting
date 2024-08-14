//go:generate mockgen -source=service.go -destination=mock_service_test.go -package gotesting
package gotesting

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

type User struct {
	ID   string
	Name string
}

// assume this is not a service implemented by us and we will mock it to test
type UserServiceClient interface {
	GetUser(ID string) (*User, error)
	CreateUser(user *User) (ID string, err error)
	Monitor() (chan *User, error) // returns users getting creates
}

type AuthService interface {
	// Simple call to UserServiceClient.GetUser
	GetUser(ID string) (*User, error)
	// CreateUser: will call monitor and concurrently calls CreateUser on sucess
	// monitor will return user with timeout of 10 secs on create success
	// also fails on cancelled context if CreateUser returns error
	CreateUser(user *User) (*User, error)
}

// This is the service we want to test
type authService struct {
	userClient UserServiceClient
	//... This part is not important for example
}

func NewAuthService(userClient UserServiceClient) AuthService {
	return &authService{
		userClient: userClient,
	}
}

// assume authService implements AuthService interface based on comments
// I will upload a fuller example to github

// CreateUser implements AuthService.
func (a *authService) CreateUser(user *User) (*User, error) {
	if user.Name == "" {
		return nil, errors.New("name is required")
	}

	eg, egCtx := errgroup.WithContext(context.TODO())

	var newUser *User
	eg.Go(func() (err error) {
		userChan, err := a.userClient.Monitor()
		if err != nil {
			return err
		}
		for {
			select {
			case <-egCtx.Done():
			case newUser = <-userChan:
				if newUser != nil {
					return nil
				}
			}
		}
	})

	eg.Go(func() (err error) {
		_, err = a.userClient.CreateUser(user)
		if err != nil {
			return errors.New("name is required")
		}

		return err
	})

	err := eg.Wait()
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

// GetUser implements AuthService.
func (a *authService) GetUser(ID string) (*User, error) {
	if ID == "" {
		return nil, errors.New("id is required")
	}
	return a.userClient.GetUser(ID)
}
