package cqrs_test

import (
	"errors"
	"testing"

	"github.com/bounoable/cqrs"
	mock_cqrs "github.com/bounoable/cqrs/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRegisterEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := cqrs.NewEventConfig()
	typ := cqrs.EventType("test")

	ed, err := cfg.NewData(typ)
	assert.True(t, errors.Is(cqrs.UnregisteredEventError{
		EventType: typ,
	}, err))
	assert.Nil(t, ed)

	eventData := mock_cqrs.NewMockEventData(ctrl)
	factory := cqrs.EventDataFactory(func() cqrs.EventData {
		return eventData
	})

	cfg.Register(typ, factory)

	ed, err = cfg.NewData(typ)
	assert.Nil(t, err)
	assert.Equal(t, eventData, ed)
}
