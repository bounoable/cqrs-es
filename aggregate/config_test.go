package aggregate_test

import (
	"errors"
	"testing"

	"github.com/bounoable/cqrs-es"
	"github.com/bounoable/cqrs-es/aggregate"
	mock_cqrs "github.com/bounoable/cqrs-es/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	cfg := aggregate.Config()
	assert.Len(t, cfg.Factories(), 0)
}

func TestRegisterAggregate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := aggregate.Config()
	typ := cqrs.AggregateType("test")
	id := uuid.New()

	a, err := cfg.New(typ, id)
	assert.True(t, errors.Is(aggregate.UnregisteredError{
		AggregateType: typ,
	}, err))

	agg := mock_cqrs.NewMockAggregate(ctrl)
	factory := cqrs.AggregateFactory(func(id uuid.UUID) cqrs.Aggregate {
		return agg
	})

	cfg.Register(typ, factory)

	a, err = cfg.New(typ, id)
	assert.Nil(t, err)
	assert.Equal(t, agg, a)
}
