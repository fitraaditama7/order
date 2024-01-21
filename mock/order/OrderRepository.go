// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	repository "order-backend/repository"

	time "time"
)

// OrderRepository is an autogenerated mock type for the OrderRepository type
type OrderRepository struct {
	mock.Mock
}

// FindList provides a mock function with given fields: ctx, limit, page, keyword, dateStart, dateEnd
func (_m *OrderRepository) FindList(ctx context.Context, limit int, page int, keyword string, dateStart *time.Time, dateEnd *time.Time) ([]repository.OrderList, error) {
	ret := _m.Called(ctx, limit, page, keyword, dateStart, dateEnd)

	if len(ret) == 0 {
		panic("no return value specified for FindList")
	}

	var r0 []repository.OrderList
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int, string, *time.Time, *time.Time) ([]repository.OrderList, error)); ok {
		return rf(ctx, limit, page, keyword, dateStart, dateEnd)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int, string, *time.Time, *time.Time) []repository.OrderList); ok {
		r0 = rf(ctx, limit, page, keyword, dateStart, dateEnd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]repository.OrderList)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int, string, *time.Time, *time.Time) error); ok {
		r1 = rf(ctx, limit, page, keyword, dateStart, dateEnd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindTotalOrder provides a mock function with given fields: ctx, keyword, dateStart, dateEnd
func (_m *OrderRepository) FindTotalOrder(ctx context.Context, keyword string, dateStart *time.Time, dateEnd *time.Time) (*repository.TotalOrder, error) {
	ret := _m.Called(ctx, keyword, dateStart, dateEnd)

	if len(ret) == 0 {
		panic("no return value specified for FindTotalOrder")
	}

	var r0 *repository.TotalOrder
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *time.Time, *time.Time) (*repository.TotalOrder, error)); ok {
		return rf(ctx, keyword, dateStart, dateEnd)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *time.Time, *time.Time) *repository.TotalOrder); ok {
		r0 = rf(ctx, keyword, dateStart, dateEnd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*repository.TotalOrder)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *time.Time, *time.Time) error); ok {
		r1 = rf(ctx, keyword, dateStart, dateEnd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewOrderRepository creates a new instance of OrderRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOrderRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *OrderRepository {
	mock := &OrderRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
