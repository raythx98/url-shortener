// Code generated by mockery. DO NOT EDIT.

package random

import mock "github.com/stretchr/testify/mock"

// MockIRandom is an autogenerated mock type for the IRandom type
type MockIRandom struct {
	mock.Mock
}

type MockIRandom_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIRandom) EXPECT() *MockIRandom_Expecter {
	return &MockIRandom_Expecter{mock: &_m.Mock}
}

// GenerateAlphaNum provides a mock function with given fields: length
func (_m *MockIRandom) GenerateAlphaNum(length int) string {
	ret := _m.Called(length)

	if len(ret) == 0 {
		panic("no return value specified for GenerateAlphaNum")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(int) string); ok {
		r0 = rf(length)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockIRandom_GenerateAlphaNum_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateAlphaNum'
type MockIRandom_GenerateAlphaNum_Call struct {
	*mock.Call
}

// GenerateAlphaNum is a helper method to define mock.On call
//   - length int
func (_e *MockIRandom_Expecter) GenerateAlphaNum(length interface{}) *MockIRandom_GenerateAlphaNum_Call {
	return &MockIRandom_GenerateAlphaNum_Call{Call: _e.mock.On("GenerateAlphaNum", length)}
}

func (_c *MockIRandom_GenerateAlphaNum_Call) Run(run func(length int)) *MockIRandom_GenerateAlphaNum_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *MockIRandom_GenerateAlphaNum_Call) Return(_a0 string) *MockIRandom_GenerateAlphaNum_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIRandom_GenerateAlphaNum_Call) RunAndReturn(run func(int) string) *MockIRandom_GenerateAlphaNum_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIRandom creates a new instance of MockIRandom. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIRandom(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIRandom {
	mock := &MockIRandom{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
