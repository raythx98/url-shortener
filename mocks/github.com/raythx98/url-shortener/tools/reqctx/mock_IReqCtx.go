// Code generated by mockery. DO NOT EDIT.

package reqctx

import (
	context "context"

	reqctx "github.com/raythx98/gohelpme/tool/reqctx"
	mock "github.com/stretchr/testify/mock"
)

// MockIReqCtx is an autogenerated mock type for the IReqCtx type
type MockIReqCtx struct {
	mock.Mock
}

type MockIReqCtx_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIReqCtx) EXPECT() *MockIReqCtx_Expecter {
	return &MockIReqCtx_Expecter{mock: &_m.Mock}
}

// GetValue provides a mock function with given fields: ctx
func (_m *MockIReqCtx) GetValue(ctx context.Context) *reqctx.Value {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetValue")
	}

	var r0 *reqctx.Value
	if rf, ok := ret.Get(0).(func(context.Context) *reqctx.Value); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*reqctx.Value)
		}
	}

	return r0
}

// MockIReqCtx_GetValue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetValue'
type MockIReqCtx_GetValue_Call struct {
	*mock.Call
}

// GetValue is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockIReqCtx_Expecter) GetValue(ctx interface{}) *MockIReqCtx_GetValue_Call {
	return &MockIReqCtx_GetValue_Call{Call: _e.mock.On("GetValue", ctx)}
}

func (_c *MockIReqCtx_GetValue_Call) Run(run func(ctx context.Context)) *MockIReqCtx_GetValue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockIReqCtx_GetValue_Call) Return(_a0 *reqctx.Value) *MockIReqCtx_GetValue_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIReqCtx_GetValue_Call) RunAndReturn(run func(context.Context) *reqctx.Value) *MockIReqCtx_GetValue_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIReqCtx creates a new instance of MockIReqCtx. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIReqCtx(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIReqCtx {
	mock := &MockIReqCtx{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
