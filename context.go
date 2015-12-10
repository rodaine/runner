package runner

import "sync"

// A Context encapsulates the state for a Command to execute with. Contexts are passed down to subsequent Commands,
// allowing them to access data from previous Commands. Setting a non-nil error on the Context will initiate a rollback
// of the run.
type Context interface {
	// Err returns an error if the Command execution has failed. Otherwise, the value is nil. If this value is non-nil,
	// a rollback is triggered starting with the previous Command and proceeds in reverse order.
	Err() error

	// SetErr allows a command to indicate the Command has failed and should trigger a rollback.
	// TODO: Prevent non-Failable commands from un-setting an existing error.
	SetErr(err error)

	// Get returns a stored value from the Context, as well as indicated whether or not a value was found for that key.
	// Values are evaluated and returned in reverse execution of Commands.
	Get(key interface{}) (val interface{}, found bool)

	// Set allows a command to add data to the execution Context that can be used by the current Command in a rollback or
	// by a subsequent Command in its execution. Values of previous Commands are shadowed and not destroyed so their
	// original values are accessible during a rollback or within parallel Commands.
	Set(key, val interface{})

	push()
	pop()
}

// NewContext returns a new root context. This function is a utility to aid in testing Command implementations.
func NewContext() Context {
	return &ctx{
		kvs: []hash{make(hash)},
	}
}

type hash map[interface{}]interface{}

type ctx struct {
	sync.RWMutex

	kvs []hash
	err error
}

func (ctx *ctx) Err() error {
	ctx.RLock()
	defer ctx.RUnlock()
	return ctx.err
}

func (ctx *ctx) SetErr(err error) {
	ctx.Lock()
	ctx.err = err
	ctx.Unlock()
}

func (ctx *ctx) Get(key interface{}) (val interface{}, found bool) {
	ctx.RLock()
	defer ctx.RUnlock()

	for i := len(ctx.kvs) - 1; i >= 0; i-- {
		if val, found = ctx.kvs[i][key]; found {
			break
		}
	}

	return
}

func (ctx *ctx) Set(key interface{}, val interface{}) {
	ctx.Lock()
	ctx.kvs[len(ctx.kvs)-1][key] = val
	ctx.Unlock()
}

func (ctx *ctx) push() {
	ctx.kvs = append(ctx.kvs, make(hash))
}

func (ctx *ctx) pop() {
	if len(ctx.kvs) == 1 {
		panic("cannot pop root context")
	}
	ctx.kvs = ctx.kvs[:len(ctx.kvs)-1]
}
