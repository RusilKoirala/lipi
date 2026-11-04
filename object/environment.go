package object

// Environment stores variables: rakha x = 5;  ->  store["x"] = 5
type Environment struct {
	store map[string]Object
	outer *Environment // the scope around this one (this is how closures work)
}

// the top level scope
func NewEnvironment() *Environment {
	return &Environment{store: map[string]Object{}}
}

// new scope inside another one (every function call makes one of these)
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// look here first, then in the outer scope
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return obj, ok
}

// always writes to THIS scope, so rakha inside a function dont change the outer variable
func (e *Environment) Set(name string, val Object) {
	e.store[name] = val
}
