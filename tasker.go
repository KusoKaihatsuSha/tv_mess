package main

import (
	"context"
	"strings"
	"sync"
	"time"
)

type CtxKey string

// Tasker work like workers pool
type Tasker struct {
	Hands  chan struct{} //same worker
	Things chan Thing
	Wg     *sync.WaitGroup
	Branch BranchContext
	M      *sync.RWMutex
}

// BranchContext help work with context inside Tasker
type BranchContext struct {
	Context context.Context
	Cancel  context.CancelFunc
}

// Thing  = task element
type Thing struct {
	Input   any
	Message *Message
	Action  func(*Tasker, Thing, *Message)
	Output  chan any
}

// AddCtx(*Tasker, string, any)
// add key and value to context for massage
func (m *Message) AddCtx(o *Tasker, key string, val any) {
	o.M.Lock()
	var ctx, ctx2 context.Context
	var fn context.CancelFunc
	ctx = context.WithValue(o.Branch.Context, CtxKey(strings.TrimSpace(m.UUID+key)), val)
	ctx2, fn = context.WithCancel(ctx)
	o.Branch = BranchContext{ctx2, fn}
	o.M.Unlock()
}

// AddCtx(*Tasker, string, any)
// add key and value to context
func AddCtx(o *Tasker, key string, val any) {
	o.M.Lock()
	var ctx, ctx2 context.Context
	var fn context.CancelFunc
	ctx = context.WithValue(o.Branch.Context, CtxKey(strings.TrimSpace(key)), val)
	ctx2, fn = context.WithCancel(ctx)
	o.Branch = BranchContext{ctx2, fn}
	o.M.Unlock()
}

// GetCtx[T any](*Tasker, string, *Message) T
// Get context values with waiting
func GetCtx[T any](o *Tasker, key string, message *Message) T {
	var N T
	var notfound = make(chan struct{})
	timeout := false
	go func() {
		for {
			if o.Branch.Context.Value(CtxKey(strings.TrimSpace(message.UUID+key))) != nil || timeout {
				notfound <- struct{}{}
				defer close(notfound)
				return
			}
		}
	}()
	select {
	case <-time.After(10 * time.Minute):
		timeout = true
	case <-notfound:
	}
	o.M.RLock()
	m := o.Branch.Context.Value(CtxKey(strings.TrimSpace(message.UUID + key)))
	o.M.RUnlock()
	if m != nil {
		return m.(T)
	}
	return N
}

// Init(int, int) *Tasker
// initialize Tasker
func (o *Tasker) Init(cpuCapability, taskCapability int) *Tasker {
	o.Wg = &sync.WaitGroup{}
	o.M = &sync.RWMutex{}
	o.Branch = BranchContext{Context: context.Background()}
	o.Hands = make(chan struct{}, cpuCapability)
	o.Things = make(chan Thing, taskCapability)
	for i := 1; i <= cap(o.Hands); i++ {
		o.Hands <- struct{}{}
	}
	go o.Pull()
	return o
}

// Add(any, func(*Tasker, Thing, *Message), *Message) *Thing
// add tasks
func (o *Tasker) Add(val any, fn func(*Tasker, Thing, *Message), mm *Message) *Thing {
	select {
	case <-o.Branch.Context.Done():
	default:
		o.Wg.Add(1)
		n := Thing{}
		n.Input = val
		n.Message = mm
		n.Action = fn
		o.Things <- n
		return &n
	}
	return nil
}

// Pull() Handle tasks from chan
func (o *Tasker) Pull() {
	defer close(o.Hands)
	for range o.Hands {
		select {
		case <-o.Branch.Context.Done():
			break
		default:
			go o.Work()
		}
	}

}

// Work() run tasks fn
func (o *Tasker) Work() {
	defer close(o.Things)
	for c := range o.Things {
		select {
		case <-o.Branch.Context.Done():
			c.Action(o, c, c.Message)
			o.Wg.Done()
			break
		default:
			c.Action(o, c, c.Message)
			o.Wg.Done()
		}
	}
}
