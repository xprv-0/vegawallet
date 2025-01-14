// Code generated by MockGen. DO NOT EDIT.
// Source: code.vegaprotocol.io/vegawallet/service (interfaces: NodeForward)

// Package mocks is a generated GoMock package.
package mocks

import (
	v1 "code.vegaprotocol.io/protos/vega/api/v1"
	v10 "code.vegaprotocol.io/protos/vega/commands/v1"
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockNodeForward is a mock of NodeForward interface
type MockNodeForward struct {
	ctrl     *gomock.Controller
	recorder *MockNodeForwardMockRecorder
}

// MockNodeForwardMockRecorder is the mock recorder for MockNodeForward
type MockNodeForwardMockRecorder struct {
	mock *MockNodeForward
}

// NewMockNodeForward creates a new mock instance
func NewMockNodeForward(ctrl *gomock.Controller) *MockNodeForward {
	mock := &MockNodeForward{ctrl: ctrl}
	mock.recorder = &MockNodeForwardMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNodeForward) EXPECT() *MockNodeForwardMockRecorder {
	return m.recorder
}

// HealthCheck mocks base method
func (m *MockNodeForward) HealthCheck(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HealthCheck", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// HealthCheck indicates an expected call of HealthCheck
func (mr *MockNodeForwardMockRecorder) HealthCheck(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HealthCheck", reflect.TypeOf((*MockNodeForward)(nil).HealthCheck), arg0)
}

// LastBlockHeight mocks base method
func (m *MockNodeForward) LastBlockHeight(arg0 context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastBlockHeight", arg0)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastBlockHeight indicates an expected call of LastBlockHeight
func (mr *MockNodeForwardMockRecorder) LastBlockHeight(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastBlockHeight", reflect.TypeOf((*MockNodeForward)(nil).LastBlockHeight), arg0)
}

// SendTx mocks base method
func (m *MockNodeForward) SendTx(arg0 context.Context, arg1 *v10.Transaction, arg2 v1.SubmitTransactionRequest_Type) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendTx", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendTx indicates an expected call of SendTx
func (mr *MockNodeForwardMockRecorder) SendTx(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendTx", reflect.TypeOf((*MockNodeForward)(nil).SendTx), arg0, arg1, arg2)
}
