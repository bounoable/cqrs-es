package command_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bounoable/cqrs-es/command"
	mock_command "github.com/bounoable/cqrs-es/mocks/command"
	"github.com/golang/mock/gomock"
)

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := command.NewConfig()

	typ := command.Type("test")
	handler := mock_command.NewMockHandler(ctrl)

	h, err := cfg.Handler(typ)
	assert.True(t, errors.Is(command.UnregisteredError{
		CommandType: typ,
	}, err))
	assert.Nil(t, h)

	cfg.Register(typ, handler)

	h, err = cfg.Handler(typ)
	assert.Nil(t, err)
	assert.Equal(t, handler, h)
}
