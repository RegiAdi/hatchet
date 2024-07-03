// Code generated by mockery v2.33.2. DO NOT EDIT.

package servicemocks

import (
	responses "github.com/RegiAdi/venera/responses"
	mock "github.com/stretchr/testify/mock"
)

// UserService is an autogenerated mock type for the UserService type
type UserService struct {
	mock.Mock
}

// GetUserDetail provides a mock function with given fields: APIToken
func (_m *UserService) GetUserDetail(APIToken string) (responses.UserResponse, error) {
	ret := _m.Called(APIToken)

	var r0 responses.UserResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (responses.UserResponse, error)); ok {
		return rf(APIToken)
	}
	if rf, ok := ret.Get(0).(func(string) responses.UserResponse); ok {
		r0 = rf(APIToken)
	} else {
		r0 = ret.Get(0).(responses.UserResponse)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(APIToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserService creates a new instance of UserService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserService(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserService {
	mock := &UserService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
