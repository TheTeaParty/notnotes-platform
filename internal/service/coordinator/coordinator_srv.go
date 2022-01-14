package coordinator

import (
	"context"
	"encoding/json"
	"reflect"
	"sync"

	"github.com/TheTeaParty/notnotes-platform/internal/config"
	"github.com/TheTeaParty/notnotes-platform/internal/service"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type coordinatorNATS struct {
	handlers        map[service.CoordinatorTrigger][]service.GenericHandler
	events          []service.CoordinatorTrigger
	lock            *sync.RWMutex
	natsClient      *nats.Conn
	pubSubReceivers map[service.CoordinatorTrigger]*nats.Subscription
	l               *zap.Logger
	c               *config.Config

	cons *nats.ConsumerInfo
}

func (c *coordinatorNATS) UnsubscribeAll() error {
	return nil
}

func (c *coordinatorNATS) argInfo(handler service.GenericHandler) (reflect.Type, int, error) {
	hType := reflect.TypeOf(handler)
	if hType.Kind() != reflect.Func {
		return nil, 0, service.ErrInvalidHandler
	}
	numArgs := hType.NumIn()
	if numArgs != 2 {
		return nil, 0, service.ErrInvalidNumberOfArguments
	}

	numOuts := hType.NumOut()
	if numOuts != 1 {
		return nil, 0, service.ErrInvalidHandlerReturn
	}

	if hType.Out(0).String() != "error" {
		return nil, 0, service.ErrInvalidHandlerReturn
	}

	if hType.In(0).String() != "context.Context" {
		return nil, 0, service.ErrInvalidHandlerFirstArg
	}

	return hType.In(1), numArgs, nil
}

func (c *coordinatorNATS) Events() []service.CoordinatorTrigger {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.events
}

func (c *coordinatorNATS) Handle(name service.CoordinatorTrigger, handler service.GenericHandler) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	l := c.l.With(zap.String("trigger", string(name)))

	argType, _, err := c.argInfo(handler)
	if err != nil {
		c.l.With(zap.Error(err)).Error("Invalid argument type")
	}

	_, ok := c.handlers[name]
	if !ok {
		c.handlers[name] = make([]service.GenericHandler, 0)
		c.events = append(c.events, name)
	}

	c.handlers[name] = append(c.handlers[name], handler)

	_, ok = c.pubSubReceivers[name]
	if ok {
		return
	}

	ctx := context.Background()

	sub, err := c.natsClient.QueueSubscribe(string(name), c.c.AppName, func(msg *nats.Msg) {
		c.lock.RLock()
		handlers := c.handlers[name]
		c.lock.RUnlock()
		msg.InProgress()

		arg := reflect.New(argType.Elem()).Interface()
		if err := json.Unmarshal(msg.Data, &arg); err != nil {
			l.With(zap.Error(err)).Error("Error parsing args")
			msg.Nak()
			return
		}

		vArgs := make([]reflect.Value, 2)
		vArgs[0] = reflect.ValueOf(ctx)
		vArgs[1] = reflect.ValueOf(arg)
		var wg sync.WaitGroup

		for _, h := range handlers {
			wg.Add(1)
			go func(h service.GenericHandler, vArgs []reflect.Value) {

				defer wg.Done()

				hValue := reflect.ValueOf(h)
				vRets := hValue.Call(vArgs)

				if vRets[0].Interface() != nil {
					err := vRets[0].Interface().(error)
					l.With(zap.Error(err)).Error("Error running handler for trigger")
					msg.Nak()
					return
				}

			}(h, vArgs)
		}

		wg.Wait()
		msg.Ack()
	})
	if err != nil {
		l.With(zap.Error(err)).Error("Can't subscribe to ")
		return
	}

	c.pubSubReceivers[name] = sub

	return
}

func (c *coordinatorNATS) Trigger(ctx context.Context, name service.CoordinatorTrigger, arg interface{}) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	l := c.l.With(zap.String("method", "trigger"))

	handlers, ok := c.handlers[name]
	if !ok {
		handlers = make([]service.GenericHandler, 0)
		c.handlers[name] = handlers
	}

	data, err := json.Marshal(arg)
	if err != nil {
		data = []byte("")
	}

	err = c.natsClient.Publish(string(name), data)
	if err != nil {
		l.With(zap.Error(err)).Error("Can't publish message")
		return
	}
}

func NewNATS(l *zap.Logger,
	c *config.Config,
	natsClient *nats.Conn) service.Coordinator {

	return &coordinatorNATS{
		l:               l.With(zap.String("service", "coordinator")),
		c:               c,
		natsClient:      natsClient,
		handlers:        make(map[service.CoordinatorTrigger][]service.GenericHandler),
		lock:            &sync.RWMutex{},
		pubSubReceivers: make(map[service.CoordinatorTrigger]*nats.Subscription),
	}
}
