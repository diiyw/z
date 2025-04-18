# Object Types

## Table of Contents

- [Z Objects](#z-objects)
- [Runtime Object Types](#runtime-object-types)
- [User Object Types](#user-object-types)

## Z Objects

In Z, all object types _(both [runtime types](#runtime-object-types) and
[user types](#user-object-types))_ must implement
[Object](https://godoc.org/github.com/diiyw/z#Object) interface.

### Object Interface

```golang
TypeName() string
```

TypeName method should return the name of the type. Type names are not directly
used by the runtime _(except when it reports a run-time error)_, but, it is
generally a good idea to keep it short but distinguishable from other types.

```golang
String() string
```

String method should return a string representation of the underlying value.
The value returned by String method will be used whenever string formatting for
the value is required, most commonly when being converted into String value.

```golang
BinaryOp(op token.Token, rhs Object) (res Object, err error)
```

In Z, a type can overload binary operators
(`+`, `-`, `*`, `/`, `%`, `&`, `|`, `^`, `&^`, `>>`, `<<`, `>`, `>=`; _note
that `<` and `<=` operators are not overloadable as they're simply implemented
by switching left-hand side and right-hand side of `>`/`>=` operator_) by
implementing BinaryOp method. BinaryOp method takes the operator `op` and the
right-hand side object `rhs`, and, should return a resulting value `res`. 

**Error value vs runtime error**

If BinaryOp method returns an error `err` (the second return value), it will be
treated as a run-time error, which will halt the execution (`VM.Run() error`)
and will return the error to the user. All runtime type implementations, for
example, will return an `ErrInvalidOperator` error when the given operator is
not supported by the type.

Alternatively the method can return an `Error` value as its result `res`
(the first return value), which will not halt the runtime and will be treated
like any other values. As a dynamically typed language, the receiver (another
expression or statement) can determine how to translate `Error` value returned
from binary operator expression.

```golang
IsFalsy() bool
```

IsFalsy method should return true if the underlying value is considered to be
[falsy](https://github.com/diiyw/z/blob/master/docs/runtime-types.md#objectisfalsy).

```golang
Equals(o Object) bool
```

Equals method should return true if the underlying value is considered to be
equal to the underlying value of another object `o`. When comparing values of
different types, the runtime does not guarantee or force anything, but, it's
generally a good idea to make the result consistent. For example, a custom
integer type may return true when comparing against String value, but, it
should return the same result for the same inputs.

```golang
Copy() Object
```

Copy method should return a _new_ copy of the object. Builtin function `copy`
uses this method to copy values. Default implementation of all runtime types
return a deep-copy values, but, it's not a requirement by the runtime.

```golang
IndexGet(index Object) (value Object, err error)
```

IndexGet should take an index Object and return a result Object or an error for
indexable objects. Indexable is an object that can take an index and return an
object. If a type is indexable, its values support dot selector
(value = object.index) and indexer (value = object[index]) syntax.

If Object is not indexable, ErrNotIndexable should be returned as error. If nil
is returned as value, it will be converted to Undefined value by the runtime.

If `IndexGet` returns an error (`err`), the VM will treat it as a run-time
error and ignore the returned value.

Array and Map implementation forces the type of index Object to be Int and
String respectively, but, it's not a required behavior of the VM. It is
completely okay to take various index types as long as it is consistent.

By convention, Array or Array-like types and Map or Map-like types return
`Undefined` value when the key does not exist. But, again, this is not a
required behavior.

```golang
IndexSet(index, value Object) error
```

IndexSet should take an index Object and a value Object for index assignable
objects. Index assignable is an object that can take an index and a value on
the left-hand side of the assignment statement. If a type is index assignable,
its values support assignment using dot selector (`object.index = value`) and
indexer (`object[index] = value`) in the assignment statements.

If Object is not index assignable, ErrNotIndexAssignable should be returned as
error. If an error is returned, it will be treated as a run-time error.

Array and Map implementation forces the type of index Object to be Int and
String respectively, but, it's not a required behavior of the VM. It is
completely okay to take various index types as long as it is consistent.

#### Callable Objects

If the type is Callable, its values can be invoked as if they were functions.
Two functions need to be implemented for Callable objects.

```golang
CanCall() bool
```

CanCall should return whether the Object can be called. When this function
returns true, the Object is considered Callable.

```golang
Call(args ...Object) (ret Object, err error)
```

Call should take an arbitrary number of arguments and return a return value
and/or an error, which the VM will consider as a run-time error.

#### Iterable Objects

If a type is iterable, its values can be used in `for-in` statements
(`for key, value in object { ... }`). Two functions need to be implemented
for Iterable Objects

```golang
CanIterate() bool
```

CanIterate should return whether the Object can be Iterated.

```golang
Iterate() Iterator
```

The Iterate method should return another object that implements
[Iterator](https://godoc.org/github.com/diiyw/z#Iterator) interface.

### Iterator Interface

```golang
Next() bool
```

Next method should return true if there are more elements to iterate. When used
with `for-in` statements, the compiler uses Key and Value methods to populate
the current element's key (or index) and value from the object that this
iterator represents. The runtime will stop iterating in `for-in` statement
when this method returns false.

```golang
Key() Object
```

Key method should return a key (or an index) Object for the current element of
the underlying object. It should return the same value until Next method is
called again. By convention, iterators for the map or map-like objects returns
the String key, and, iterators for array or array-like objects returns the Int
ndex. But, it's not a requirement by the VM.

```golang
Value() Object
```

Value method should return a value Object for the current element of the
underlying object. It should return the same value until Next method is called
again.

## Runtime Object Types

These are the basic types Z runtime supports out of the box:

- Primitive value types: [Int](https://godoc.org/github.com/diiyw/z#Int),
  [String](https://godoc.org/github.com/diiyw/z#String),
  [Float](https://godoc.org/github.com/diiyw/z#Float),
  [Bool](https://godoc.org/github.com/diiyw/z#ArrayIterator),
  [Char](https://godoc.org/github.com/diiyw/z#Char),
  [Bytes](https://godoc.org/github.com/diiyw/z#Bytes),
  [Time](https://godoc.org/github.com/diiyw/z#Time)
- Composite value types: [Array](https://godoc.org/github.com/diiyw/z#Array),
  [ImmutableArray](https://godoc.org/github.com/diiyw/z#ImmutableArray),
  [Map](https://godoc.org/github.com/diiyw/z#Map),
  [ImmutableMap](https://godoc.org/github.com/diiyw/z#ImmutableMap)
- Functions:
  [CompiledFunction](https://godoc.org/github.com/diiyw/z#CompiledFunction),
  [BuiltinFunction](https://godoc.org/github.com/diiyw/z#BuiltinFunction),
  [UserFunction](https://godoc.org/github.com/diiyw/z#UserFunction)
- [Iterators](https://godoc.org/github.com/diiyw/z#Iterator):
  [StringIterator](https://godoc.org/github.com/diiyw/z#StringIterator),
  [ArrayIterator](https://godoc.org/github.com/diiyw/z#ArrayIterator),
  [MapIterator](https://godoc.org/github.com/diiyw/z#MapIterator),
  [ImmutableMapIterator](https://godoc.org/github.com/diiyw/z#ImmutableMapIterator)
- [Error](https://godoc.org/github.com/diiyw/z#Error)
- [Undefined](https://godoc.org/github.com/diiyw/z#Undefined)
- Other internal objects: [Break](https://godoc.org/github.com/diiyw/z#Break),
  [Continue](https://godoc.org/github.com/diiyw/z#Continue),
  [ReturnValue](https://godoc.org/github.com/diiyw/z#ReturnValue)

See
[Runtime Types](https://github.com/diiyw/z/blob/master/docs/runtime-types.md)
for more details on these runtime types.

## User Object Types

Users can easily extend and add their own types by implementing the same
[Object](https://godoc.org/github.com/diiyw/z#Object) interface and the
default `ObjectImpl` implementation. Z runtime will treat them in the
same way as its runtime types with no performance overhead.

Here's an example user type implementation, `StringArray`:

```golang
type StringArray struct {
    z.ObjectImpl
    Value []string
}

func (o *StringArray) String() string {
    return strings.Join(o.Value, ", ")
}

func (o *StringArray) BinaryOp(op token.Token, rhs z.Object) (z.Object, error) {
    if rhs, ok := rhs.(*StringArray); ok {
        switch op {
        case token.Add:
            if len(rhs.Value) == 0 {
                return o, nil
            }
            return &StringArray{Value: append(o.Value, rhs.Value...)}, nil
        }
    }

    return nil, z.ErrInvalidOperator
}

func (o *StringArray) IsFalsy() bool {
    return len(o.Value) == 0
}

func (o *StringArray) Equals(x z.Object) bool {
    if x, ok := x.(*StringArray); ok {
        if len(o.Value) != len(x.Value) {
            return false
        }

        for i, v := range o.Value {
            if v != x.Value[i] {
                return false
            }
        }

        return true
    }

    return false
}

func (o *StringArray) Copy() z.Object {
    return &StringArray{
        Value: append([]string{}, o.Value...),
    }
}

func (o *StringArray) TypeName() string {
    return "string-array"
}
```

You can use a user type via either
[Script.Add](https://godoc.org/github.com/diiyw/z#Script.Add) or by directly
manipulating the symbol table and the global variables. Here's an example code
to add `StringArray` to the script:

```golang
// script that uses 'my_list'
s := z.NewScript([]byte(`
    fmt := import("fmt")
    fmt.println(my_list, ", three")
`))

s.SetImports(stdlib.GetModuleMap("fmt"))
myList := &StringArray{Value: []string{"one", "two"}}
s.Add("my_list", myList)  // add StringArray value 'my_list'
s.Run()                   // prints "one, two, three"
```

It can also implement `IndexGet` and `IndexSet`:

```golang
func (o *StringArray) IndexGet(index z.Object) (z.Object, error) {
    intIdx, ok := index.(*z.Int)
    if ok {
        if intIdx.Value >= 0 && intIdx.Value < int64(len(o.Value)) {
            return &z.String{Value: o.Value[intIdx.Value]}, nil
        }

        return nil, z.ErrIndexOutOfBounds
    }

    strIdx, ok := index.(*z.String)
    if ok {
        for vidx, str := range o.Value {
            if strIdx.Value == str {
                return &z.Int{Value: int64(vidx)}, nil
            }
        }

        return z.UndefinedValue, nil
    }

    return nil, z.ErrInvalidIndexType
}

func (o *StringArray) IndexSet(index, value z.Object) error {
    strVal, ok := z.ToString(value)
    if !ok {
        return z.ErrInvalidIndexValueType
    }

    intIdx, ok := index.(*z.Int)
    if ok {
        if intIdx.Value >= 0 && intIdx.Value < int64(len(o.Value)) {
            o.Value[intIdx.Value] = strVal
            return nil
        }

        return z.ErrIndexOutOfBounds
    }

    return z.ErrInvalidIndexType
}
```

If we implement `CanCall` and `Call`:

```golang
func (o *StringArray) CanCall() bool {
    return true
}

func (o *StringArray) Call(args ...z.Object) (ret z.Object, err error) {
    if len(args) != 1 {
        return nil, z.ErrWrongNumArguments
    }

    s1, ok := z.ToString(args[0])
    if !ok {
        return nil, z.ErrInvalidArgumentType{
            Name:     "first",
            Expected: "string",
            Found:    args[0].TypeName(),
        }
    }

    for i, v := range o.Value {
        if v == s1 {
            return &z.Int{Value: int64(i)}, nil
        }
    }

    return z.UndefinedValue, nil
}
```

Then it can be "invoked":

```golang
s := z.NewScript([]byte(`
    print(my_list("two"))
`))

myList := &StringArray{Value: []string{"one", "two", "three"}}
s.Add("my_list", myList)  // add StringArray value 'my_list'
s.Run()                   // prints "1" (index of "two")
```

We can also make `StringArray` iterable:

```golang
func (o *StringArray) CanIterate() bool {
    return true
}

func (o *StringArray) Iterate() z.Iterator {
    return &StringArrayIterator{
        strArr: o,
    }
}

type StringArrayIterator struct {
    z.ObjectImpl
    strArr *StringArray
    idx    int
}

func (i *StringArrayIterator) TypeName() string {
    return "string-array-iterator"
}

func (i *StringArrayIterator) Next() bool {
    i.idx++
    return i.idx <= len(i.strArr.Value)
}

func (i *StringArrayIterator) Key() z.Object {
    return &z.Int{Value: int64(i.idx - 1)}
}

func (i *StringArrayIterator) Value() z.Object {
    return &z.String{Value: i.strArr.Value[i.idx-1]}
}
```

### ObjectImpl

ObjectImpl represents a default Object Implementation. To defined a new value
type, one can embed ObjectImpl in their type declarations to avoid implementing
all non-significant methods. TypeName() and String() methods still need to be
implemented.
