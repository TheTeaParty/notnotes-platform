package service

import (
	"context"
	"errors"
)

var (
	ErrInvalidHandler           = errors.New("handler is not func")
	ErrInvalidHandlerFirstArg   = errors.New("handler first arg must be context.Context type")
	ErrInvalidNumberOfArguments = errors.New("handler must have first argument as context.Context and one other argument")
	ErrInvalidHandlerArguments  = errors.New("invalid handler arguments")
	ErrInvalidHandlerReturn     = errors.New("handler must return only an error")
)

// CoordinatorTrigger is trigger name. Pattern is %domain%.%event%.%other_sub_events%(optional)
type CoordinatorTrigger string

const (
	TriggerNoteCreated CoordinatorTrigger = "note.created"
	TriggerNoteUpdated CoordinatorTrigger = "note.updated"
	TriggerNoteDeleted CoordinatorTrigger = "note.deleted"

	TriggerTagCreated CoordinatorTrigger = "tag.created"
	TriggerTagDeleted CoordinatorTrigger = "tag.deleted"
)

type Coordinator interface {
	Handle(name CoordinatorTrigger, handler GenericHandler)
	Trigger(ctx context.Context, name CoordinatorTrigger, arg interface{})
	Events() []CoordinatorTrigger
}

type GenericHandler interface{}
