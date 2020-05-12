# `🔌 surge`

[![GitHub](https://github.com/renproject/surge/workflows/test/badge.svg)](https://github.com/renproject/surge/workflows/test/badge.svg)
[![Coverage](https://coveralls.io/repos/github/renproject/surge/badge.svg?branch=master)](https://coveralls.io/github/renproject/surge?branch=master)
[![Report](https://goreportcard.com/badge/github.com/renproject/surge)](https://goreportcard.com/badge/github.com/renproject/surge)

[Documentation](https://godoc.org/github.com/renproject/surge)

A library for fast binary (un)marshaling. Designed to be used in Byzantine networks, `🔌 surge` never explicitly panics and protects against malicious inputs. It supports the (un)marshaling of:

- [x] scalars,
- [x] arrays,
- [x] slices,
- [x] maps, and
- [x] custom implementations (using the `Marshaler` and `Unmarshaler` interfaces).

Built with ❤ by Ren. 
