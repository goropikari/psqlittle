// Code generated by MockGen. DO NOT EDIT.
// Source: backend.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	backend "github.com/goropikari/mysqlite2/backend"
	core "github.com/goropikari/mysqlite2/core"
)

// MockDB is a mock of DB interface.
type MockDB struct {
	ctrl     *gomock.Controller
	recorder *MockDBMockRecorder
}

// MockDBMockRecorder is the mock recorder for MockDB.
type MockDBMockRecorder struct {
	mock *MockDB
}

// NewMockDB creates a new mock instance.
func NewMockDB(ctrl *gomock.Controller) *MockDB {
	mock := &MockDB{ctrl: ctrl}
	mock.recorder = &MockDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDB) EXPECT() *MockDBMockRecorder {
	return m.recorder
}

// GetTable mocks base method.
func (m *MockDB) GetTable(arg0 string) (backend.Table, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTable", arg0)
	ret0, _ := ret[0].(backend.Table)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTable indicates an expected call of GetTable.
func (mr *MockDBMockRecorder) GetTable(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTable", reflect.TypeOf((*MockDB)(nil).GetTable), arg0)
}

// MockTable is a mock of Table interface.
type MockTable struct {
	ctrl     *gomock.Controller
	recorder *MockTableMockRecorder
}

// MockTableMockRecorder is the mock recorder for MockTable.
type MockTableMockRecorder struct {
	mock *MockTable
}

// NewMockTable creates a new mock instance.
func NewMockTable(ctrl *gomock.Controller) *MockTable {
	mock := &MockTable{ctrl: ctrl}
	mock.recorder = &MockTableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTable) EXPECT() *MockTableMockRecorder {
	return m.recorder
}

// Copy mocks base method.
func (m *MockTable) Copy() backend.Table {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Copy")
	ret0, _ := ret[0].(backend.Table)
	return ret0
}

// Copy indicates an expected call of Copy.
func (mr *MockTableMockRecorder) Copy() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Copy", reflect.TypeOf((*MockTable)(nil).Copy))
}

// GetRows mocks base method.
func (m *MockTable) GetRows() []backend.Row {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRows")
	ret0, _ := ret[0].([]backend.Row)
	return ret0
}

// GetRows indicates an expected call of GetRows.
func (mr *MockTableMockRecorder) GetRows() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRows", reflect.TypeOf((*MockTable)(nil).GetRows))
}

// SetRows mocks base method.
func (m *MockTable) SetRows(arg0 []backend.Row) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetRows", arg0)
}

// SetRows indicates an expected call of SetRows.
func (mr *MockTableMockRecorder) SetRows(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRows", reflect.TypeOf((*MockTable)(nil).SetRows), arg0)
}

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

// GetValues mocks base method.
func (m *MockRow) GetValues() core.Values {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValues")
	ret0, _ := ret[0].(core.Values)
	return ret0
}

// GetValues indicates an expected call of GetValues.
func (mr *MockRowMockRecorder) GetValues() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValues", reflect.TypeOf((*MockRow)(nil).GetValues))
}