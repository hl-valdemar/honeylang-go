# More ideas on the language

## Equivalent of Zig's `comptime_int` and `comptime_float`

Zig has a primitive type called `comptime_int` which is basically any int that can resolve at comptime.

It would be nice to have such a type in Honey as well, but maybe we can just call it `int` (instead of the rather cumbersome `comptime_int`)
Likewise, a `float` type equivalent to Zig's `comptime_float` would also be nice.

## Return type of `main`

The return type of the `main` function should be restricted to a `u8`, an `int` (that resolves to `u8`), and `void` (which implicitly returns `0`).

## Type declaration of constants and variables

What about:

```hon
# constants (`::`)
x :: 32
x int :: 32
x float :: 3.14

# variables
a = 0xf0f1
a int = 0xf0f1

# arrays
arr [3]u8 :: [1, 2, 3]      # const binding, const elements
arr [3]mut u8 :: [1, 2, 3]  # const binding, mutable elements
arr [3]u8 = [1, 2, 3]       # mutable binding, const elements
arr [3]mut u8 = [1, 2, 3]   # mutable binding, mutable elements

arr [_]u8 :: [1, 2, 3]  # size inferred from literal

# single-item pointers (only points to *one* element, no arithmetic)
ptr @T :: &t      # const binding, const pointer (i.e., no write through pointer)
ptr @mut T :: &t  # const binding, mutable pointer
ptr @T = &t       # mutable binding, const pointer
ptr @mut T = &t   # mutable binding, mutable pointer

# c-like many-item pointers (arithmetic supported - primarily used for interop)
ptr *T :: &t      # const binding, const pointer (i.e., no write through pointer)
ptr *mut T :: &t  # const binding, mutable pointer
ptr *T = &t       # mutable binding, const pointer
ptr *mut T = &t   # mutable binding, mutable pointer

# slices (includes both a len field for native processing, and guarantees null termination for c interop)
s []u8 :: arr[..]    # slice including all elements
s []u8 :: arr[..2]   # slice including all elements up to index 2
s []u8 :: arr[2..]   # slice including all elements from index 2
s []u8 :: arr[2..4]  # slice including all elements between index 2 (incl) and 4 (excl)

# arrays with explicit sentinels (must be of same type as array)
s [_..0]u8 :: arr[1, 2, 3]    # null terminator (not default for arrays)
s [3..'A']u8 :: arr[1, 2, 3]  # use 'A' as sentinel

# slices with explicit sentinel (must be of same type as slice)
s [..0]u8 :: arr[..]    # explicit null terminator (default so not necessary)
s [..'A']u8 :: arr[..]  # use 'A' as sentinel

# functions
sum func(a, b int) int { ... }

# to comptime or not to comptime
$x :: 32                                        # comptime evaluable constant
$sum func(a, b int) int { ... }                 # comptime evaluable function
max func($T type, a, b T) T { ... }             # comptime polymorphic type params
make_array func ($size usize) [size]u8 { ... }  # comptime value params

# NB: variables are runtime dependent and as such not comptime compatible

# namespaces carry behavior (and also holds data) - can't be instantiated, no runtime concept of a namespace
person {
    # structs carry only data, no methods, no behavior
    Data struct {
        name []u8  # `[]u8` is basically a string
        age u32
        position Vec2
    }

    # mutable variable scoped to person (sum of distances walked by all people)
    dist_sum float = 0

    walk func(person @mut Data, distance float) boolean {
        person.position.x += distance
        dist_sum += distance
    }
}

# instantiating a struct from a namespace
p :: person.Data{
    name = "Jane",
    age = 32,
    position = Vec2{ x = 0, y = 0 },
}

# using a function from a namespace
person.walk(&p, 2.7)
```

## YOLO publicity rules (who cares)

Everything is public. (This should be easy to change later if need be.)
