// Code generated by MockGen. DO NOT EDIT.
// Source: parser.go

// Package mock_parser is a generated GoMock package.
package mock_parser

import (
	reflect "reflect"

	model "github.com/ISDL-dev/MaP1058-socket-client/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockParser is a mock of Parser interface.
type MockParser struct {
	ctrl     *gomock.Controller
	recorder *MockParserMockRecorder
}

// MockParserMockRecorder is the mock recorder for MockParser.
type MockParserMockRecorder struct {
	mock *MockParser
}

// NewMockParser creates a new mock instance.
func NewMockParser(ctrl *gomock.Controller) *MockParser {
	mock := &MockParser{ctrl: ctrl}
	mock.recorder = &MockParserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockParser) EXPECT() *MockParserMockRecorder {
	return m.recorder
}

// ToAnalysis mocks base method.
func (m *MockParser) ToAnalysis(c model.Command) (model.AnalysisType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToAnalysis", c)
	ret0, _ := ret[0].(model.AnalysisType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToAnalysis indicates an expected call of ToAnalysis.
func (mr *MockParserMockRecorder) ToAnalysis(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToAnalysis", reflect.TypeOf((*MockParser)(nil).ToAnalysis), c)
}

// ToChannelCal mocks base method.
func (m *MockParser) ToChannelCal(c model.Command) (model.ChannelCal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToChannelCal", c)
	ret0, _ := ret[0].(model.ChannelCal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToChannelCal indicates an expected call of ToChannelCal.
func (mr *MockParserMockRecorder) ToChannelCal(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToChannelCal", reflect.TypeOf((*MockParser)(nil).ToChannelCal), c)
}

// ToChannelPower mocks base method.
func (m *MockParser) ToChannelPower(c model.Command) (model.ChannelPower, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToChannelPower", c)
	ret0, _ := ret[0].(model.ChannelPower)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToChannelPower indicates an expected call of ToChannelPower.
func (mr *MockParserMockRecorder) ToChannelPower(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToChannelPower", reflect.TypeOf((*MockParser)(nil).ToChannelPower), c)
}

// ToCommand mocks base method.
func (m *MockParser) ToCommand(s string) (model.Command, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToCommand", s)
	ret0, _ := ret[0].(model.Command)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToCommand indicates an expected call of ToCommand.
func (mr *MockParserMockRecorder) ToCommand(s interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToCommand", reflect.TypeOf((*MockParser)(nil).ToCommand), s)
}

// ToSignals mocks base method.
func (m *MockParser) ToSignals(b []byte) (*model.Signals, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToSignals", b)
	ret0, _ := ret[0].(*model.Signals)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToSignals indicates an expected call of ToSignals.
func (mr *MockParserMockRecorder) ToSignals(b interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToSignals", reflect.TypeOf((*MockParser)(nil).ToSignals), b)
}

// ToTrendRange mocks base method.
func (m *MockParser) ToTrendRange(c model.Command) (model.TrendRange, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToTrendRange", c)
	ret0, _ := ret[0].(model.TrendRange)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToTrendRange indicates an expected call of ToTrendRange.
func (mr *MockParserMockRecorder) ToTrendRange(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToTrendRange", reflect.TypeOf((*MockParser)(nil).ToTrendRange), c)
}
