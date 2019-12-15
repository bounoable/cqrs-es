package cqrs_test

import (
	"errors"
	"testing"

	"github.com/bounoable/cqrs"
	"github.com/stretchr/testify/assert"

	mock_cqrs "github.com/bounoable/cqrs/mocks"
	"github.com/golang/mock/gomock"
)

func TestRegisterCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := cqrs.NewCommandConfig()

	typ := cqrs.CommandType("test")
	handler := mock_cqrs.NewMockCommandHandler(ctrl)

	h, err := cfg.Handler(typ)
	assert.True(t, errors.Is(cqrs.UnregisteredCommandError{
		CommandType: typ,
	}, err))
	assert.Nil(t, h)

	cfg.Register(typ, handler)

	h, err = cfg.Handler(typ)
	assert.Nil(t, err)
	assert.Equal(t, handler, h)
}
