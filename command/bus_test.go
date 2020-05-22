package command_test

import (
	"context"
	"testing"

	cqrs "github.com/bounoable/cqrs-es"
	"github.com/bounoable/cqrs-es/command"
	mock_cqrs "github.com/bounoable/cqrs-es/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDispatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := mock_cqrs.NewMockCommandConfig(ctrl)
	bus := command.BusWithConfig(cfg)

	ctx := context.Background()
	cmd := mock_cqrs.NewMockCommand(ctrl)
	cmdType := cqrs.CommandType("test")
	cmdHandler := mock_cqrs.NewMockCommandHandler(ctrl)

	cmd.EXPECT().Type().Return(cmdType)
	cfg.EXPECT().Handler(cmdType).Return(cmdHandler, nil)
	cmdHandler.EXPECT().Handle(ctx, cmd).Return(nil)

	err := bus.Dispatch(ctx, cmd)
	assert.Nil(t, err)
}
