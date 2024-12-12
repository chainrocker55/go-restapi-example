// Code generated by mockery v2.49.1. DO NOT EDIT.

package mocks

import (
	entity "creditlimit-connector/app/entity"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// RunningNoRepo is an autogenerated mock type for the RunningNoRepo type
type RunningNoRepo struct {
	mock.Mock
}

// FindByNameAndUpdatedAt provides a mock function with given fields: name, datetime
func (_m *RunningNoRepo) FindByNameAndUpdatedAt(name string, datetime time.Time) (*entity.RunningNoEntity, error) {
	ret := _m.Called(name, datetime)

	if len(ret) == 0 {
		panic("no return value specified for FindByNameAndUpdatedAt")
	}

	var r0 *entity.RunningNoEntity
	var r1 error
	if rf, ok := ret.Get(0).(func(string, time.Time) (*entity.RunningNoEntity, error)); ok {
		return rf(name, datetime)
	}
	if rf, ok := ret.Get(0).(func(string, time.Time) *entity.RunningNoEntity); ok {
		r0 = rf(name, datetime)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.RunningNoEntity)
		}
	}

	if rf, ok := ret.Get(1).(func(string, time.Time) error); ok {
		r1 = rf(name, datetime)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: data
func (_m *RunningNoRepo) Save(data entity.RunningNoEntity) error {
	ret := _m.Called(data)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(entity.RunningNoEntity) error); ok {
		r0 = rf(data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRunningNoRepo creates a new instance of RunningNoRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRunningNoRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *RunningNoRepo {
	mock := &RunningNoRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
