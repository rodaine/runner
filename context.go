package runner

import "sync"

type Context interface {
	Err() error
	SetErr(err error)

	Get(key interface{}) (val interface{}, found bool)
	Set(key, val interface{})

	push()
	pop()
}

type hash map[interface{}]interface{}

func NewContext() Context {
	return &ctx{
		kvs: []hash{make(hash)},
	}
}

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
