package mocks

import "gopkg.in/inconshreveable/log15.v2"
import "github.com/stretchr/testify/mock"

type Logger struct {
	mock.Mock
}

func (m *Logger) New(ctx ...interface{}) log15.Logger {
	ret := m.Called(ctx)

	r0 := ret.Get(0).(log15.Logger)

	return r0
}
func (m *Logger) SetHandler(h log15.Handler) {
	m.Called(h)
}
func (m *Logger) Debug(msg string, ctx ...interface{}) {
	m.Called(msg, ctx)
}
func (m *Logger) Info(msg string, ctx ...interface{}) {
	m.Called(msg, ctx)
}
func (m *Logger) Warn(msg string, ctx ...interface{}) {
	m.Called(msg, ctx)
}
func (m *Logger) Error(msg string, ctx ...interface{}) {
	m.Called(msg, ctx)
}
func (m *Logger) Crit(msg string, ctx ...interface{}) {
	m.Called(msg, ctx)
}
