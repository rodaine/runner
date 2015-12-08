package runner

import "sync"

// A Context object is used as both the input and output of a Command. It allows passing information to later Commands
// in the Sequence. Context also captures the error state of a command and will trigger a rollback of previously completed
// commands.
type Context struct {
	rwm  sync.RWMutex
	ctxs []map[interface{}]interface{}
	err  error
}

// NewContext returns a Context instance that has been pre-allocated with the specified capacity. The capacity is not fixed,
// however predetermining the maximum depth of commands in the sequence reduces the number of allocations required at
// runtime.
func NewContext(capacity int) *Context {
	return &Context{
		ctxs: make([]map[interface{}]interface{}, 0, capacity),
	}
}

// Err returns the current error of the Command Sequence, otherwise nil. If an error is returned, Sequence will trigger
// a rollback of all previously executed commands that implement the Rollbacker interface in reverse order.
func (s *Context) Err() error {
	s.rwm.RLock()
	defer s.rwm.RUnlock()
	return s.err
}

// SetErr allows a command to set the error of the Context. Be aware that this error can be suppressed by passing in nil,
// however if a rollback is already underway, it cannot be stopped. Use commands.FailableCommand to allow a command to
// fail.
func (s *Context) SetErr(err error) {
	s.rwm.Lock()
	s.err = err
	s.rwm.Unlock()
}

// Set allows a Command to assign arbitrary information into the Context to be shared with subsequent Commands in the Sequence.
// Previous Commands cannot access subsequent Command properties. Set will override previous values of a key, but the
// change is not destructive. The key provided must be a valid map key (the key must be "comparable").
func (s *Context) Set(key, value interface{}) {
	s.rwm.Lock()
	s.ctxs[len(s.ctxs)-1][key] = value
	s.rwm.Unlock()
}

// Get returns a value stored by the current or a previous Commands. Get cannot access values set by subsequent Commands.
// If the key is not set, nil will be returned.
func (s *Context) Get(key interface{}) interface{} {
	s.rwm.RLock()
	defer s.rwm.RUnlock()
	for i := len(s.ctxs) - 1; i >= 0; i-- {
		if val, ok := s.ctxs[i][key]; ok {
			return val
		}
	}
	return nil
}

// push adds a new context layer for a subsequent Command.
func (s *Context) push() {
	s.rwm.Lock()
	s.ctxs = append(s.ctxs, map[interface{}]interface{}{})
	s.rwm.Unlock()
}

// pop removes a context layer when bubbling back up from a subsequent Command, used during a rollback.
func (s *Context) pop() {
	s.rwm.Lock()
	s.ctxs = s.ctxs[:len(s.ctxs)-1]
	s.rwm.Unlock()
}
