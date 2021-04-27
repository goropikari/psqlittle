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

// CreateTable mocks base method.
func (m *MockDB) CreateTable(arg0 string, arg1 core.Cols) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTable", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTable indicates an expected call of CreateTable.
func (mr *MockDBMockRecorder) CreateTable(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTable", reflect.TypeOf((*MockDB)(nil).CreateTable), arg0, arg1)
}

// DropTable mocks base method.
func (m *MockDB) DropTable(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DropTable", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DropTable indicates an expected call of DropTable.
func (mr *MockDBMockRecorder) DropTable(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DropTable", reflect.TypeOf((*MockDB)(nil).DropTable), arg0)
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

// Delete mocks base method.
func (m *MockTable) Delete(arg0 func(backend.Row) core.Value) (backend.Table, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(backend.Table)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockTableMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTable)(nil).Delete), arg0)
}

// GetColNames mocks base method.
func (m *MockTable) GetColNames() core.ColumnNames {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetColNames")
	ret0, _ := ret[0].(core.ColumnNames)
	return ret0
}

// GetColNames indicates an expected call of GetColNames.
func (mr *MockTableMockRecorder) GetColNames() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetColNames", reflect.TypeOf((*MockTable)(nil).GetColNames))
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

// InsertValues mocks base method.
func (m *MockTable) InsertValues(arg0 core.ColumnNames, arg1 core.ValuesList) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertValues", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertValues indicates an expected call of InsertValues.
func (mr *MockTableMockRecorder) InsertValues(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertValues", reflect.TypeOf((*MockTable)(nil).InsertValues), arg0, arg1)
}

// Project mocks base method.
func (m *MockTable) Project(arg0 core.ColumnNames, arg1 []func(backend.Row) core.Value) (backend.Table, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Project", arg0, arg1)
	ret0, _ := ret[0].(backend.Table)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Project indicates an expected call of Project.
func (mr *MockTableMockRecorder) Project(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Project", reflect.TypeOf((*MockTable)(nil).Project), arg0, arg1)
}

// Update mocks base method.
func (m *MockTable) Update(arg0 core.ColumnNames, arg1 func(backend.Row) core.Value, arg2 []func(backend.Row) core.Value) (backend.Table, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(backend.Table)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockTableMockRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTable)(nil).Update), arg0, arg1, arg2)
}

// UpdateTableName mocks base method.
func (m *MockTable) UpdateTableName(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateTableName", arg0)
}

// UpdateTableName indicates an expected call of UpdateTableName.
func (mr *MockTableMockRecorder) UpdateTableName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTableName", reflect.TypeOf((*MockTable)(nil).UpdateTableName), arg0)
}

// Where mocks base method.
func (m *MockTable) Where(arg0 func(backend.Row) core.Value) (backend.Table, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Where", arg0)
	ret0, _ := ret[0].(backend.Table)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Where indicates an expected call of Where.
func (mr *MockTableMockRecorder) Where(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Where", reflect.TypeOf((*MockTable)(nil).Where), arg0)
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
func (m *MockRow) GetValueByColName(arg0 core.ColumnName) core.Value {
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

// SetColNames mocks base method.
func (m *MockRow) SetColNames(arg0 core.ColumnNames) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetColNames", arg0)
}

// SetColNames indicates an expected call of SetColNames.
func (mr *MockRowMockRecorder) SetColNames(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetColNames", reflect.TypeOf((*MockRow)(nil).SetColNames), arg0)
}

// SetValues mocks base method.
func (m *MockRow) SetValues(arg0 core.Values) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetValues", arg0)
}

// SetValues indicates an expected call of SetValues.
func (mr *MockRowMockRecorder) SetValues(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetValues", reflect.TypeOf((*MockRow)(nil).SetValues), arg0)
}

// UpdateValue mocks base method.
func (m *MockRow) UpdateValue(arg0 core.ColumnName, arg1 core.Value) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateValue", arg0, arg1)
}

// UpdateValue indicates an expected call of UpdateValue.
func (mr *MockRowMockRecorder) UpdateValue(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateValue", reflect.TypeOf((*MockRow)(nil).UpdateValue), arg0, arg1)
}
