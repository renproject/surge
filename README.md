# `🔌 surge`

[Documentation](https://godoc.org/github.com/renproject/surge)

A library for fast binary (un)marshaling. Designed to be used in Byzantine networks, `🔌 surge` never explicitly panics. It supports the (un)marshaling of:

- [x] scalars,
- [x] arrays,
- [x] slices,
- [x] maps,
- [x] custom structs (using the `surge:` tag), and
- [x] custom implementations (using the `Marshaler` and `Unmarshaler` interfaces).

Example:

```go
package main

import (
    "bytes"
    "reflect"
    "github.com/renproject/surge"
)

type Person struct {
    Name    string            `surge:"0"`
    Age     uint64            `surge:"1"`
    Friends map[string]Person `surge:"2"`
}

func main() {
    alice := Person{
        Name:    "Alice",
        Age:     25,
        Friends: map[string]Person{
            "Bob": Person{
                Name:   "Bob",
                Age:     26,
                Friends: map[string]Person{},
            },
        },
    }

    data, err := surge.ToBinary(alice)
    if err != nil {
        panic(err)
    }

    alice2 := Person{}
    if err := surge.FromBinary(&alice2, data); err != nil {
        panic(err)
    }

    if !reflect.DeepEqual(alice, alice2) {
        panic("bad (un)marshal")
    }
}
```

Built with ❤ by Ren. 
