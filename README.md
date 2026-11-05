# uuid

[English](./README.md) | [日本語](./README_ja.md)

`github.com/WhiteRaven777/uuid` is a small extension of the Go 1.27 standard
library's [`uuid`](https://pkg.go.dev/uuid) package.

## Standard library wrappers

The standard API is exposed through thin wrappers:

- `Parse`, `MustParse`, `New`, `Nil`, `Max`, `NewV4`, and `NewV7`
- `UUID.String`, `MarshalText`, `AppendText`, `UnmarshalText`, and `Compare`

The local `UUID` type has the same `[16]byte` representation as the standard
library type. UUID parsing, version 4 generation, and version 7 generation are
delegated to the standard library.

## Extensions

### Additional output formats

The local `UUID` type provides three output methods that are not present on the
standard library type. Given this UUID:

```text
01234567-89ab-7def-8123-456789abcdef
```

the methods return:

| Method | Format | Output |
| --- | --- | --- |
| `Short()` | 32 lowercase hexadecimal digits without hyphens | `0123456789ab7def8123456789abcdef` |
| `URN()` | UUID URN | `urn:uuid:01234567-89ab-7def-8123-456789abcdef` |
| `Base64()` | Unpadded standard Base64 | `ASNFZ4mrfe+BI0VniavN7w` |

Values returned by `Short` and `URN` can be passed to `Parse`. Base64 is an
output format only and is not accepted by `Parse`.

### Additional UUID versions

#### UUID version 1

`NewV1` complements the version 4 and version 7 generators provided by the
standard library. UUIDv1 contains a 60-bit Gregorian timestamp, a 14-bit clock
sequence, and a 48-bit node identifier.

By default, `NewV1` uses a cryptographically random node identifier rather than
exposing a hardware MAC address:

```go
u := uuid.NewV1()
```

Pass `true` to request a MAC address as the node identifier:

```go
u := uuid.NewV1(true)
```

If no suitable MAC address is available, generation falls back to a random
node identifier. Version 1 is intended for compatibility with systems that
require its Gregorian timestamp layout. For new time-ordered identifiers,
version 7 is generally the better choice.

#### UUID version 6

`NewV6` uses the same Gregorian timestamp, clock sequence, and node identifier
as UUIDv1, but stores the timestamp from most to least significant. This makes
the UUIDs sort in generation order and improves locality in database indexes.

As with `NewV1`, the node identifier is random by default. Pass `true` to
request a MAC address, with a random node identifier used as a fallback:

```go
ordered := uuid.NewV6()
hardwareBased := uuid.NewV6(true)
```

UUIDv6 is primarily useful for systems that need UUIDv1-compatible fields with
better sorting behaviour. New systems that do not require UUIDv1 compatibility
should generally use UUIDv7.

## Requirements

Go 1.27 or newer is required.

## Example

```go
package main

import (
	"fmt"

	"github.com/WhiteRaven777/uuid"
)

func main() {
	u := uuid.NewV7()
	fmt.Println(u.String())
	fmt.Println(u.Short())
	fmt.Println(u.URN())
	fmt.Println(u.Base64())

	parsed := uuid.MustParse(u.String())
	fmt.Println(parsed == u)
}
```

## License

Copyright (c) 2021-present [WhiteRaven777](https://github.com/WhiteRaven777)

Licensed under the [MIT License](./LICENSE).
