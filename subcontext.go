package runner

type subCtx struct {
	parent Context
	ctx    Context
}

func newSubContext(parent Context) Context {
	return &subCtx{
		parent: parent,
		ctx:    NewContext(),
	}
}

func (sc *subCtx) Err() error {
	err := sc.ctx.Err()
	if sc.ctx.Err() == nil {
		err = sc.parent.Err()
	}
	return err
}

func (sc *subCtx) SetErr(err error) {
	sc.ctx.SetErr(err)
}

func (sc *subCtx) Get(key interface{}) (val interface{}, found bool) {
	if val, found = sc.ctx.Get(key); !found {
		val, found = sc.parent.Get(key)
	}
	return
}

func (sc *subCtx) Set(key, val interface{}) {
	sc.ctx.Set(key, val)
}

func (sc *subCtx) push() {
	sc.ctx.push()
}

func (sc *subCtx) pop() {
	sc.ctx.pop()
}
