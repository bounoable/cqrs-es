package cqrs_test

import (
	"context"
	"testing"

	"github.com/bounoable/cqrs"
	mock_cqrs "github.com/bounoable/cqrs/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDispatchCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := mock_cqrs.NewMockCommandConfig(ctrl)
	bus := cqrs.NewCommandBusWithConfig(cfg, nil)

	ctx := context.Background()
	cmd := mock_cqrs.NewMockCommand(ctrl)
	cmdType := cqrs.CommandType("test")
	cmdHandler := mock_cqrs.NewMockCommandHandler(ctrl)

	cmd.EXPECT().CommandType().Return(cmdType)
	cfg.EXPECT().Handler(cmdType).Return(cmdHandler, nil)
	cmdHandler.EXPECT().HandleCommand(ctx, cmd).Return(nil)

	err := bus.Dispatch(ctx, cmd)
	assert.Nil(t, err)
}
