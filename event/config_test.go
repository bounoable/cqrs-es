package event_test

import (
	"errors"
	"testing"

	"github.com/bounoable/cqrs-es"
	"github.com/bounoable/cqrs-es/event"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	type testEventData struct{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := event.Config()
	typ := cqrs.EventType("test")

	ed, err := cfg.NewData(typ)
	assert.True(t, errors.Is(event.UnregisteredError{
		EventType: typ,
	}, err))
	assert.Nil(t, ed)

	cfg.Register(typ, testEventData{})

	ed, err = cfg.NewData(typ)
	assert.Nil(t, err)
	assert.Equal(t, testEventData{}, ed)
}
