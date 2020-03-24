package command_test

import (
	"testing"

	"github.com/bounoable/cqrs-es"
	"github.com/bounoable/cqrs-es/command"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewBase(t *testing.T) {
	typ := cqrs.CommandType("test")
	aggregateType := cqrs.AggregateType("test")
	aggregateID := uuid.New()
	cmd := command.NewBase(typ, aggregateType, aggregateID)

	assert.Equal(t, typ, cmd.CommandType())
	assert.Equal(t, aggregateType, cmd.AggregateType())
	assert.Equal(t, aggregateID, cmd.AggregateID())
}
