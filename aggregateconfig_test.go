package cqrs_test

import (
	"errors"
	"testing"

	"github.com/bounoable/cqrs"
	mock_cqrs "github.com/bounoable/cqrs/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewAggregateConfig(t *testing.T) {
	cfg := cqrs.NewAggregateConfig()
	assert.Len(t, cfg.Factories(), 0)
}

func TestRegisterAggregate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := cqrs.NewAggregateConfig()
	typ := cqrs.AggregateType("test")
	id := uuid.New()

	a, err := cfg.New(typ, id)
	assert.True(t, errors.Is(cqrs.UnregisteredAggregateError{
		AggregateType: typ,
	}, err))

	aggregate := mock_cqrs.NewMockAggregate(ctrl)
	factory := cqrs.AggregateFactory(func(id uuid.UUID) cqrs.Aggregate {
		return aggregate
	})

	cfg.Register(typ, factory)

	a, err = cfg.New(typ, id)
	assert.Nil(t, err)
	assert.Equal(t, aggregate, a)
}
