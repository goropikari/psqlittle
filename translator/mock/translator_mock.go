// Code generated by MockGen. DO NOT EDIT.
// Source: translator/translator.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	core "github.com/goropikari/mysqlite2/core"
	translator "github.com/goropikari/mysqlite2/translator"
)

// MockRow is a mock of Row interface.
type MockRow struct {
	ctrl     *gomock.Controller
	recorder *MockRowMockRecorder
}

// MockRowMockRecorder is the mock recorder for MockRow.
type MockRowMockRecorder struct {
	mock *MockRow
}

// NewMockRow creates a new mock instance.
func NewMockRow(ctrl *gomock.Controller) *MockRow {
	mock := &MockRow{ctrl: ctrl}
	mock.recorder = &MockRowMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRow) EXPECT() *MockRowMockRecorder {
	return m.recorder
}

// GetValueByColName mocks base method.
func (m *MockRow) GetValueByColName(arg0 core.ColName) core.Value {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValueByColName", arg0)
	ret0, _ := ret[0].(core.Value)
	return ret0
}

// GetValueByColName indicates an expected call of GetValueByColName.
func (mr *MockRowMockRecorder) GetValueByColName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValueByColName", reflect.TypeOf((*MockRow)(nil).GetValueByColName), arg0)
}

// MockExpr is a mock of Expr interface.
type MockExpr struct {
	ctrl     *gomock.Controller
	recorder *MockExprMockRecorder
}

// MockExprMockRecorder is the mock recorder for MockExpr.
type MockExprMockRecorder struct {
	mock *MockExpr
}

// NewMockExpr creates a new mock instance.
func NewMockExpr(ctrl *gomock.Controller) *MockExpr {
	mock := &MockExpr{ctrl: ctrl}
	mock.recorder = &MockExprMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExpr) EXPECT() *MockExprMockRecorder {
	return m.recorder
}

// Eval mocks base method.
func (m *MockExpr) Eval() func(translator.Row) core.Value {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Eval")
	ret0, _ := ret[0].(func(translator.Row) core.Value)
	return ret0
}

// Eval indicates an expected call of Eval.
func (mr *MockExprMockRecorder) Eval() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Eval", reflect.TypeOf((*MockExpr)(nil).Eval))
}