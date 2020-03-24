package command_test

import (
	"context"
	"testing"

	"github.com/bounoable/cqrs-es"
	"github.com/bounoable/cqrs-es/command"
	mock_cqrs "github.com/bounoable/cqrs-es/mocks"
	mock_command "github.com/bounoable/cqrs-es/mocks/command"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDispatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := mock_command.NewMockConfig(ctrl)
	bus := command.NewBusWithConfig(cfg)

	ctx := context.Background()
	cmd := mock_cqrs.NewMockCommand(ctrl)
	cmdType := cqrs.CommandType("test")
	cmdHandler := mock_command.NewMockHandler(ctrl)

	cmd.EXPECT().CommandType().Return(cmdType)
	cfg.EXPECT().Handler(cmdType).Return(cmdHandler, nil)
	cmdHandler.EXPECT().HandleCommand(ctx, cmd).Return(nil)

	err := bus.Dispatch(ctx, cmd)
	assert.Nil(t, err)
}
