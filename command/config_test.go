package command_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	cqrs "github.com/bounoable/cqrs-es"
	"github.com/bounoable/cqrs-es/command"
	mock_cqrs "github.com/bounoable/cqrs-es/mocks"
	"github.com/golang/mock/gomock"
)

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := command.Config()

	typ := cqrs.CommandType("test")
	handler := mock_cqrs.NewMockCommandHandler(ctrl)

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
