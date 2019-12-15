package cqrs_test

import (
	"testing"

	"github.com/bounoable/cqrs-es"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewBaseCommand(t *testing.T) {
	typ := cqrs.CommandType("test")
	aggregateType := cqrs.AggregateType("test")
	aggregateID := uuid.New()
	cmd := cqrs.NewBaseCommand(typ, aggregateType, aggregateID)

	assert.Equal(t, typ, cmd.CommandType())
	assert.Equal(t, aggregateType, cmd.AggregateType())
	assert.Equal(t, aggregateID, cmd.AggregateID())
}
