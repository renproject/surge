# `surge`

A library for fast binary (un)marshaling.

```go
package main

import (
    "bytes"
    "github.com/renproject/surge"
)

type Person struct {
    Name    string   `surge:"0"`
    Age     uint64   `surge:"1"`
    Friends []Person `surge:"2"`
}

func main() {
    alice := Person{
        Name: "Alice",
        Age: 25,
        Friends: []Person{
            Person{
                Name: "Bob",
                Age: 26,
                Friends: []Person{},
            },
        },
    }

    buf := new(bytes.Buffer)
    buf.Grow(surge.SizeHint(person))
    if err := surge.Marshal(buf, person); err != nil {
        panic(err)
    }

    aliceAgain := Person{}
    if err := surge.Unmarshal(buf, &aliceAgain); err != nil {
        panic(err)
    }

    if alice != aliceAgain {
        panic("bad marshal or unmarshal")
    }
}
```