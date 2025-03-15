// Code generated by mockery. DO NOT EDIT.

package repositories

import (
	context "context"

	db "github.com/raythx98/url-shortener/sqlc/db"
	mock "github.com/stretchr/testify/mock"
)

// MockIRepository is an autogenerated mock type for the IRepository type
type MockIRepository struct {
	mock.Mock
}

type MockIRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIRepository) EXPECT() *MockIRepository_Expecter {
	return &MockIRepository_Expecter{mock: &_m.Mock}
}

// CreateRedirect provides a mock function with given fields: ctx, arg
func (_m *MockIRepository) CreateRedirect(ctx context.Context, arg db.CreateRedirectParams) error {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for CreateRedirect")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, db.CreateRedirectParams) error); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockIRepository_CreateRedirect_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateRedirect'
type MockIRepository_CreateRedirect_Call struct {
	*mock.Call
}

// CreateRedirect is a helper method to define mock.On call
//   - ctx context.Context
//   - arg db.CreateRedirectParams
func (_e *MockIRepository_Expecter) CreateRedirect(ctx interface{}, arg interface{}) *MockIRepository_CreateRedirect_Call {
	return &MockIRepository_CreateRedirect_Call{Call: _e.mock.On("CreateRedirect", ctx, arg)}
}

func (_c *MockIRepository_CreateRedirect_Call) Run(run func(ctx context.Context, arg db.CreateRedirectParams)) *MockIRepository_CreateRedirect_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(db.CreateRedirectParams))
	})
	return _c
}

func (_c *MockIRepository_CreateRedirect_Call) Return(_a0 error) *MockIRepository_CreateRedirect_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIRepository_CreateRedirect_Call) RunAndReturn(run func(context.Context, db.CreateRedirectParams) error) *MockIRepository_CreateRedirect_Call {
	_c.Call.Return(run)
	return _c
}

// CreateUrl provides a mock function with given fields: ctx, arg
func (_m *MockIRepository) CreateUrl(ctx context.Context, arg db.CreateUrlParams) (db.Url, error) {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for CreateUrl")
	}

	var r0 db.Url
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, db.CreateUrlParams) (db.Url, error)); ok {
		return rf(ctx, arg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, db.CreateUrlParams) db.Url); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Get(0).(db.Url)
	}

	if rf, ok := ret.Get(1).(func(context.Context, db.CreateUrlParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIRepository_CreateUrl_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateUrl'
type MockIRepository_CreateUrl_Call struct {
	*mock.Call
}

// CreateUrl is a helper method to define mock.On call
//   - ctx context.Context
//   - arg db.CreateUrlParams
func (_e *MockIRepository_Expecter) CreateUrl(ctx interface{}, arg interface{}) *MockIRepository_CreateUrl_Call {
	return &MockIRepository_CreateUrl_Call{Call: _e.mock.On("CreateUrl", ctx, arg)}
}

func (_c *MockIRepository_CreateUrl_Call) Run(run func(ctx context.Context, arg db.CreateUrlParams)) *MockIRepository_CreateUrl_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(db.CreateUrlParams))
	})
	return _c
}

func (_c *MockIRepository_CreateUrl_Call) Return(_a0 db.Url, _a1 error) *MockIRepository_CreateUrl_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIRepository_CreateUrl_Call) RunAndReturn(run func(context.Context, db.CreateUrlParams) (db.Url, error)) *MockIRepository_CreateUrl_Call {
	_c.Call.Return(run)
	return _c
}

// CreateUser provides a mock function with given fields: ctx, arg
func (_m *MockIRepository) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 db.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, db.CreateUserParams) (db.User, error)); ok {
		return rf(ctx, arg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, db.CreateUserParams) db.User); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Get(0).(db.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, db.CreateUserParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIRepository_CreateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateUser'
type MockIRepository_CreateUser_Call struct {
	*mock.Call
}

// CreateUser is a helper method to define mock.On call
//   - ctx context.Context
//   - arg db.CreateUserParams
func (_e *MockIRepository_Expecter) CreateUser(ctx interface{}, arg interface{}) *MockIRepository_CreateUser_Call {
	return &MockIRepository_CreateUser_Call{Call: _e.mock.On("CreateUser", ctx, arg)}
}

func (_c *MockIRepository_CreateUser_Call) Run(run func(ctx context.Context, arg db.CreateUserParams)) *MockIRepository_CreateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(db.CreateUserParams))
	})
	return _c
}

func (_c *MockIRepository_CreateUser_Call) Return(_a0 db.User, _a1 error) *MockIRepository_CreateUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIRepository_CreateUser_Call) RunAndReturn(run func(context.Context, db.CreateUserParams) (db.User, error)) *MockIRepository_CreateUser_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteUrl provides a mock function with given fields: ctx, id
func (_m *MockIRepository) DeleteUrl(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteUrl")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockIRepository_DeleteUrl_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteUrl'
type MockIRepository_DeleteUrl_Call struct {
	*mock.Call
}

// DeleteUrl is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
func (_e *MockIRepository_Expecter) DeleteUrl(ctx interface{}, id interface{}) *MockIRepository_DeleteUrl_Call {
	return &MockIRepository_DeleteUrl_Call{Call: _e.mock.On("DeleteUrl", ctx, id)}
}

func (_c *MockIRepository_DeleteUrl_Call) Run(run func(ctx context.Context, id int64)) *MockIRepository_DeleteUrl_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *MockIRepository_DeleteUrl_Call) Return(_a0 error) *MockIRepository_DeleteUrl_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIRepository_DeleteUrl_Call) RunAndReturn(run func(context.Context, int64) error) *MockIRepository_DeleteUrl_Call {
	_c.Call.Return(run)
	return _c
}

// GetRedirectsByUrlId provides a mock function with given fields: ctx, urlId
func (_m *MockIRepository) GetRedirectsByUrlId(ctx context.Context, urlId *int64) ([]db.Redirect, error) {
	ret := _m.Called(ctx, urlId)

	if len(ret) == 0 {
		panic("no return value specified for GetRedirectsByUrlId")
	}

	var r0 []db.Redirect
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *int64) ([]db.Redirect, error)); ok {
		return rf(ctx, urlId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *int64) []db.Redirect); ok {
		r0 = rf(ctx, urlId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]db.Redirect)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *int64) error); ok {
		r1 = rf(ctx, urlId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIRepository_GetRedirectsByUrlId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRedirectsByUrlId'
type MockIRepository_GetRedirectsByUrlId_Call struct {
	*mock.Call
}

// GetRedirectsByUrlId is a helper method to define mock.On call
//   - ctx context.Context
//   - urlId *int64
func (_e *MockIRepository_Expecter) GetRedirectsByUrlId(ctx interface{}, urlId interface{}) *MockIRepository_GetRedirectsByUrlId_Call {
	return &MockIRepository_GetRedirectsByUrlId_Call{Call: _e.mock.On("GetRedirectsByUrlId", ctx, urlId)}
}

func (_c *MockIRepository_GetRedirectsByUrlId_Call) Run(run func(ctx context.Context, urlId *int64)) *MockIRepository_GetRedirectsByUrlId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*int64))
	})
	return _c
}

func (_c *MockIRepository_GetRedirectsByUrlId_Call) Return(_a0 []db.Redirect, _a1 error) *MockIRepository_GetRedirectsByUrlId_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIRepository_GetRedirectsByUrlId_Call) RunAndReturn(run func(context.Context, *int64) ([]db.Redirect, error)) *MockIRepository_GetRedirectsByUrlId_Call {
	_c.Call.Return(run)
	return _c
}

// GetUrl provides a mock function with given fields: ctx, id
func (_m *MockIRepository) GetUrl(ctx context.Context, id int64) (db.Url, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetUrl")
	}

	var r0 db.Url
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (db.Url, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) db.Url); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(db.Url)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIRepository_GetUrl_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUrl'
type MockIRepository_GetUrl_Call struct {
	*mock.Call
}

// GetUrl is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
func (_e *MockIRepository_Expecter) GetUrl(ctx interface{}, id interface{}) *MockIRepository_GetUrl_Call {
	return &MockIRepository_GetUrl_Call{Call: _e.mock.On("GetUrl", ctx, id)}
}

func (_c *MockIRepository_GetUrl_Call) Run(run func(ctx context.Context, id int64)) *MockIRepository_GetUrl_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *MockIRepository_GetUrl_Call) Return(_a0 db.Url, _a1 error) *MockIRepository_GetUrl_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIRepository_GetUrl_Call) RunAndReturn(run func(context.Context, int64) (db.Url, error)) *MockIRepository_GetUrl_Call {
	_c.Call.Return(run)
	return _c
}

// GetUrlByShortUrl provides a mock function with given fields: ctx, shortUrl
func (_m *MockIRepository) GetUrlByShortUrl(ctx context.Context, shortUrl string) (*db.Url, error) {
	ret := _m.Called(ctx, shortUrl)

	if len(ret) == 0 {
		panic("no return value specified for GetUrlByShortUrl")
	}

	var r0 *db.Url
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*db.Url, error)); ok {
		return rf(ctx, shortUrl)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *db.Url); ok {
		r0 = rf(ctx, shortUrl)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*db.Url)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, shortUrl)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIRepository_GetUrlByShortUrl_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUrlByShortUrl'
type MockIRepository_GetUrlByShortUrl_Call struct {
	*mock.Call
}

// GetUrlByShortUrl is a helper method to define mock.On call
//   - ctx context.Context
//   - shortUrl string
func (_e *MockIRepository_Expecter) GetUrlByShortUrl(ctx interface{}, shortUrl interface{}) *MockIRepository_GetUrlByShortUrl_Call {
	return &MockIRepository_GetUrlByShortUrl_Call{Call: _e.mock.On("GetUrlByShortUrl", ctx, shortUrl)}
}

func (_c *MockIRepository_GetUrlByShortUrl_Call) Run(run func(ctx context.Context, shortUrl string)) *MockIRepository_GetUrlByShortUrl_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockIRepository_GetUrlByShortUrl_Call) Return(_a0 *db.Url, _a1 error) *MockIRepository_GetUrlByShortUrl_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIRepository_GetUrlByShortUrl_Call) RunAndReturn(run func(context.Context, string) (*db.Url, error)) *MockIRepository_GetUrlByShortUrl_Call {
	_c.Call.Return(run)
	return _c
}

// GetUrlsByUserId provides a mock function with given fields: ctx, userId
func (_m *MockIRepository) GetUrlsByUserId(ctx context.Context, userId *int64) ([]db.Url, error) {
	ret := _m.Called(ctx, userId)

	if len(ret) == 0 {
		panic("no return value specified for GetUrlsByUserId")
	}

	var r0 []db.Url
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *int64) ([]db.Url, error)); ok {
		return rf(ctx, userId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *int64) []db.Url); ok {
		r0 = rf(ctx, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]db.Url)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *int64) error); ok {
		r1 = rf(ctx, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIRepository_GetUrlsByUserId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUrlsByUserId'
type MockIRepository_GetUrlsByUserId_Call struct {
	*mock.Call
}

// GetUrlsByUserId is a helper method to define mock.On call
//   - ctx context.Context
//   - userId *int64
func (_e *MockIRepository_Expecter) GetUrlsByUserId(ctx interface{}, userId interface{}) *MockIRepository_GetUrlsByUserId_Call {
	return &MockIRepository_GetUrlsByUserId_Call{Call: _e.mock.On("GetUrlsByUserId", ctx, userId)}
}

func (_c *MockIRepository_GetUrlsByUserId_Call) Run(run func(ctx context.Context, userId *int64)) *MockIRepository_GetUrlsByUserId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*int64))
	})
	return _c
}

func (_c *MockIRepository_GetUrlsByUserId_Call) Return(_a0 []db.Url, _a1 error) *MockIRepository_GetUrlsByUserId_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIRepository_GetUrlsByUserId_Call) RunAndReturn(run func(context.Context, *int64) ([]db.Url, error)) *MockIRepository_GetUrlsByUserId_Call {
	_c.Call.Return(run)
	return _c
}

// GetUser provides a mock function with given fields: ctx, id
func (_m *MockIRepository) GetUser(ctx context.Context, id int64) (db.User, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetUser")
	}

	var r0 db.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (db.User, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) db.User); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(db.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIRepository_GetUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUser'
type MockIRepository_GetUser_Call struct {
	*mock.Call
}

// GetUser is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
func (_e *MockIRepository_Expecter) GetUser(ctx interface{}, id interface{}) *MockIRepository_GetUser_Call {
	return &MockIRepository_GetUser_Call{Call: _e.mock.On("GetUser", ctx, id)}
}

func (_c *MockIRepository_GetUser_Call) Run(run func(ctx context.Context, id int64)) *MockIRepository_GetUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *MockIRepository_GetUser_Call) Return(_a0 db.User, _a1 error) *MockIRepository_GetUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIRepository_GetUser_Call) RunAndReturn(run func(context.Context, int64) (db.User, error)) *MockIRepository_GetUser_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserByEmail provides a mock function with given fields: ctx, email
func (_m *MockIRepository) GetUserByEmail(ctx context.Context, email string) (*db.User, error) {
	ret := _m.Called(ctx, email)

	if len(ret) == 0 {
		panic("no return value specified for GetUserByEmail")
	}

	var r0 *db.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*db.User, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *db.User); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*db.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIRepository_GetUserByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserByEmail'
type MockIRepository_GetUserByEmail_Call struct {
	*mock.Call
}

// GetUserByEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - email string
func (_e *MockIRepository_Expecter) GetUserByEmail(ctx interface{}, email interface{}) *MockIRepository_GetUserByEmail_Call {
	return &MockIRepository_GetUserByEmail_Call{Call: _e.mock.On("GetUserByEmail", ctx, email)}
}

func (_c *MockIRepository_GetUserByEmail_Call) Run(run func(ctx context.Context, email string)) *MockIRepository_GetUserByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockIRepository_GetUserByEmail_Call) Return(_a0 *db.User, _a1 error) *MockIRepository_GetUserByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIRepository_GetUserByEmail_Call) RunAndReturn(run func(context.Context, string) (*db.User, error)) *MockIRepository_GetUserByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserTotalClicks provides a mock function with given fields: ctx, userId
func (_m *MockIRepository) GetUserTotalClicks(ctx context.Context, userId *int64) (int64, error) {
	ret := _m.Called(ctx, userId)

	if len(ret) == 0 {
		panic("no return value specified for GetUserTotalClicks")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *int64) (int64, error)); ok {
		return rf(ctx, userId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *int64) int64); ok {
		r0 = rf(ctx, userId)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *int64) error); ok {
		r1 = rf(ctx, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIRepository_GetUserTotalClicks_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserTotalClicks'
type MockIRepository_GetUserTotalClicks_Call struct {
	*mock.Call
}

// GetUserTotalClicks is a helper method to define mock.On call
//   - ctx context.Context
//   - userId *int64
func (_e *MockIRepository_Expecter) GetUserTotalClicks(ctx interface{}, userId interface{}) *MockIRepository_GetUserTotalClicks_Call {
	return &MockIRepository_GetUserTotalClicks_Call{Call: _e.mock.On("GetUserTotalClicks", ctx, userId)}
}

func (_c *MockIRepository_GetUserTotalClicks_Call) Run(run func(ctx context.Context, userId *int64)) *MockIRepository_GetUserTotalClicks_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*int64))
	})
	return _c
}

func (_c *MockIRepository_GetUserTotalClicks_Call) Return(_a0 int64, _a1 error) *MockIRepository_GetUserTotalClicks_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIRepository_GetUserTotalClicks_Call) RunAndReturn(run func(context.Context, *int64) (int64, error)) *MockIRepository_GetUserTotalClicks_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIRepository creates a new instance of MockIRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIRepository {
	mock := &MockIRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
