// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/observer/observer.go
//
// Generated by this command:
//
//	mockgen -source=./internal/observer/observer.go -destination=./internal/observer/mock_observer.go -package=observer
//

// Package observer is a generated GoMock package.
package observer

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockObserver is a mock of Observer interface.
type MockObserver struct {
	ctrl     *gomock.Controller
	recorder *MockObserverMockRecorder
}

// MockObserverMockRecorder is the mock recorder for MockObserver.
type MockObserverMockRecorder struct {
	mock *MockObserver
}

// NewMockObserver creates a new mock instance.
func NewMockObserver(ctrl *gomock.Controller) *MockObserver {
	mock := &MockObserver{ctrl: ctrl}
	mock.recorder = &MockObserverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockObserver) EXPECT() *MockObserverMockRecorder {
	return m.recorder
}

// Update mocks base method.
func (m *MockObserver) Update(data map[string]any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Update", data)
}

// Update indicates an expected call of Update.
func (mr *MockObserverMockRecorder) Update(data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockObserver)(nil).Update), data)
}
