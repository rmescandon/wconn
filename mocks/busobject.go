// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/godbus/dbus (interfaces: BusObject)

package mocks

import (
	context "context"
	dbus "github.com/godbus/dbus"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockBusObject is a mock of BusObject interface
type MockBusObject struct {
	ctrl     *gomock.Controller
	recorder *MockBusObjectMockRecorder
}

// MockBusObjectMockRecorder is the mock recorder for MockBusObject
type MockBusObjectMockRecorder struct {
	mock *MockBusObject
}

// NewMockBusObject creates a new mock instance
func NewMockBusObject(ctrl *gomock.Controller) *MockBusObject {
	mock := &MockBusObject{ctrl: ctrl}
	mock.recorder = &MockBusObjectMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockBusObject) EXPECT() *MockBusObjectMockRecorder {
	return _m.recorder
}

// AddMatchSignal mocks base method
func (_m *MockBusObject) AddMatchSignal(_param0 string, _param1 string, _param2 ...dbus.MatchOption) *dbus.Call {
	_s := []interface{}{_param0, _param1}
	for _, _x := range _param2 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "AddMatchSignal", _s...)
	ret0, _ := ret[0].(*dbus.Call)
	return ret0
}

// AddMatchSignal indicates an expected call of AddMatchSignal
func (_mr *MockBusObjectMockRecorder) AddMatchSignal(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "AddMatchSignal", reflect.TypeOf((*MockBusObject)(nil).AddMatchSignal), _s...)
}

// Call mocks base method
func (_m *MockBusObject) Call(_param0 string, _param1 dbus.Flags, _param2 ...interface{}) *dbus.Call {
	_s := []interface{}{_param0, _param1}
	for _, _x := range _param2 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "Call", _s...)
	ret0, _ := ret[0].(*dbus.Call)
	return ret0
}

// Call indicates an expected call of Call
func (_mr *MockBusObjectMockRecorder) Call(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Call", reflect.TypeOf((*MockBusObject)(nil).Call), _s...)
}

// CallWithContext mocks base method
func (_m *MockBusObject) CallWithContext(_param0 context.Context, _param1 string, _param2 dbus.Flags, _param3 ...interface{}) *dbus.Call {
	_s := []interface{}{_param0, _param1, _param2}
	for _, _x := range _param3 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "CallWithContext", _s...)
	ret0, _ := ret[0].(*dbus.Call)
	return ret0
}

// CallWithContext indicates an expected call of CallWithContext
func (_mr *MockBusObjectMockRecorder) CallWithContext(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "CallWithContext", reflect.TypeOf((*MockBusObject)(nil).CallWithContext), _s...)
}

// Destination mocks base method
func (_m *MockBusObject) Destination() string {
	ret := _m.ctrl.Call(_m, "Destination")
	ret0, _ := ret[0].(string)
	return ret0
}

// Destination indicates an expected call of Destination
func (_mr *MockBusObjectMockRecorder) Destination() *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Destination", reflect.TypeOf((*MockBusObject)(nil).Destination))
}

// GetProperty mocks base method
func (_m *MockBusObject) GetProperty(_param0 string) (dbus.Variant, error) {
	ret := _m.ctrl.Call(_m, "GetProperty", _param0)
	ret0, _ := ret[0].(dbus.Variant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProperty indicates an expected call of GetProperty
func (_mr *MockBusObjectMockRecorder) GetProperty(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetProperty", reflect.TypeOf((*MockBusObject)(nil).GetProperty), arg0)
}

// Go mocks base method
func (_m *MockBusObject) Go(_param0 string, _param1 dbus.Flags, _param2 chan *dbus.Call, _param3 ...interface{}) *dbus.Call {
	_s := []interface{}{_param0, _param1, _param2}
	for _, _x := range _param3 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "Go", _s...)
	ret0, _ := ret[0].(*dbus.Call)
	return ret0
}

// Go indicates an expected call of Go
func (_mr *MockBusObjectMockRecorder) Go(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Go", reflect.TypeOf((*MockBusObject)(nil).Go), _s...)
}

// GoWithContext mocks base method
func (_m *MockBusObject) GoWithContext(_param0 context.Context, _param1 string, _param2 dbus.Flags, _param3 chan *dbus.Call, _param4 ...interface{}) *dbus.Call {
	_s := []interface{}{_param0, _param1, _param2, _param3}
	for _, _x := range _param4 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "GoWithContext", _s...)
	ret0, _ := ret[0].(*dbus.Call)
	return ret0
}

// GoWithContext indicates an expected call of GoWithContext
func (_mr *MockBusObjectMockRecorder) GoWithContext(arg0, arg1, arg2, arg3 interface{}, arg4 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1, arg2, arg3}, arg4...)
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GoWithContext", reflect.TypeOf((*MockBusObject)(nil).GoWithContext), _s...)
}

// Path mocks base method
func (_m *MockBusObject) Path() dbus.ObjectPath {
	ret := _m.ctrl.Call(_m, "Path")
	ret0, _ := ret[0].(dbus.ObjectPath)
	return ret0
}

// Path indicates an expected call of Path
func (_mr *MockBusObjectMockRecorder) Path() *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Path", reflect.TypeOf((*MockBusObject)(nil).Path))
}

// RemoveMatchSignal mocks base method
func (_m *MockBusObject) RemoveMatchSignal(_param0 string, _param1 string, _param2 ...dbus.MatchOption) *dbus.Call {
	_s := []interface{}{_param0, _param1}
	for _, _x := range _param2 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "RemoveMatchSignal", _s...)
	ret0, _ := ret[0].(*dbus.Call)
	return ret0
}

// RemoveMatchSignal indicates an expected call of RemoveMatchSignal
func (_mr *MockBusObjectMockRecorder) RemoveMatchSignal(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "RemoveMatchSignal", reflect.TypeOf((*MockBusObject)(nil).RemoveMatchSignal), _s...)
}
