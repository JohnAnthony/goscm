GoSCM is a Scheme interpreter (and toolkit for creating Scheme interpreters) for
extending programs written in Go.


Data Types
==

All GoSCM types must be able to Eval() in an environment (whether they use that
environment to return a value or not) and convert to a string with String().

The (fairly standard) basic built-in types and their purposes:

Boolean :- A wrapper for Go's bools
PlainInt :- A wrapper for Go's ints
Pair :- A cons cell, the building blocks of lists and other data structures
Symbol :- A string-based case-insensitive identifier. Used for environment
          lookups, amonst other things.
Environ :- A wrapper around Go's hash maps. Chains with parents to create
           expanding and contracting lookup tables for environmental bindings
Contin :- An environment and list bundled up together and ready to evaluate.
Proc :- A scheme procedure that can combine with a list of evaluated arguments
Macro :- As above, but is evaluated before Procs and with unevaluated arguments
Foreign :- A Go function that can be combines with an environment and list of
           arguments to act like a scheme procedure.
Special :- As above, but is evaluated before Foreigns and with unevaluated
           arguments. Think of it as a foreign function macro.

And the interfaces:

SCMT :- Requires only the ability to Eval() and be a String()
Func :- Can be Applied() to a list of arguments
PreFunc :- Can be Expanded(), re-arranging its arguments ready for Eval()


Getting Started
==

Your package should include a Read function (you can use simple.Read) and an
Environment that includes the basics for you (you can use simple.Env or build
your own from a NewEnv with nil as the parent). You are then able to read from
IO streams or strings with Read or ReadStr respectively, and then evaluate those
expressions with their .Eval() methods, passing in the environment you wish to
use.

Sound confusing? An example of how to do this is provided in goscm.go's REPL()
function and the examples/simple_interpreter.go program. It's really very easy.


Extending your program
==

You probably want to do something as simple as this to get started:

```
func (mt *MyType) Eval(goscm.Environ) (SCMT, error) {
    return mt, nil
}
```

and

```
func (mt *MyType) String() string {
    return "#<mytype>"
}
```

And you're ready to start passing your type about using your type in GoSCM. You
might then want to create a constructor, like so:

```
func scm_make_mytype(args *goscm.Pair, env *goscm.Environ) (SCMT, error) {
    if args != goscm.SCM_Nil {
        return goscm.SCM_Nil, errors("Expected zero arguments")
    }

    return &MyType { /* Your struct initialisers here */ }, nil
}
```

Further examples of how to create Scheme procedures to interact with your types
are provided in simple/env.go.
