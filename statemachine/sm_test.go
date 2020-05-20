package statemachine

import (
	"context"
	"reflect"
	"testing"
)

type StateMachineBuilderFactory interface {
	Create() StateMachineBuilder
}

type StateMachineBuilder interface {
	Build() StateMachine
}
type StateMachine interface {
	Process(context.Context, interface{}) (interface{}, error)
}

func TestT(t *testing.T) {
	var factory StateMachineBuilderFactory
	builder := factory.Create()
	sm := builder.Build()
	_, err := sm.Process(context.TODO(), nil)
	t.Log(err)
}

type OrderStateMachine struct {
	states map[string]interface{}
}

func (sm *OrderStateMachine) Process(context.Context, interface{}) error {
	return nil
}
func (sm *OrderStateMachine) doProcess(context.Context, interface{}) error {
	return nil
}

type CreateOrderTransition struct {
}

func (sm *CreateOrderTransition) Process(ctx context.Context, params interface{}) error {
	return sm.doProcess(ctx, params.(CreateOrderRequest))
}

type CreateOrderRequest struct {
}

func (sm *CreateOrderTransition) doProcess(ctx context.Context, request CreateOrderRequest) error {

	return nil
}

func TestX(t *testing.T) {
	var factory StateMachineBuilderFactory

	t.Log(reflect.TypeOf(factory))
}
