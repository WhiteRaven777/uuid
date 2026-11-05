# uuid

[English](./README.md) | [日本語](./README_ja.md)

`github.com/WhiteRaven777/uuid` is a small extension of the Go 1.27 standard library's [`uuid`](https://pkg.go.dev/uuid) package.

## Standard library wrappers

The standard API is exposed through thin wrappers:

- `Parse`, `MustParse`, `New`, `Nil`, `Max`, `NewV4`, and `NewV7`
- `UUID.String`, `MarshalText`, `AppendText`, `UnmarshalText`, and `Compare`

The local `UUID` type has the same `[16]byte` representation as the standard library type.
UUID parsing, version 4 generation, and version 7 generation are delegated to the standard library.

## Extensions

### Additional output formats

The local `UUID` type provides three output methods that are not present on the standard library type.
Given this UUID:

```text
01234567-89ab-7def-8123-456789abcdef
```

The methods return:

| Method | Format | Output |
| --- | --- | --- |
| `Short()` | 32 lowercase hexadecimal digits without hyphens | `0123456789ab7def8123456789abcdef` |
| `URN()` | UUID URN | `urn:uuid:01234567-89ab-7def-8123-456789abcdef` |
| `Base64()` | Unpadded standard Base64 | `ASNFZ4mrfe+BI0VniavN7w` |

Values returned by `Short` and `URN` can be passed to `Parse`.
Base64 is an output format only and is not accepted by `Parse`.

### Additional UUID versions

#### UUID version 1

`NewV1` complements the version 4 and version 7 generators provided by the standard library.
UUIDv1 contains a 60-bit Gregorian timestamp, a 14-bit clock sequence, and a 48-bit node identifier.

By default, `NewV1` uses a cryptographically random node identifier rather than exposing a hardware MAC address:

```go
u := uuid.NewV1()
```

Pass `true` to request a MAC address as the node identifier:

```go
u := uuid.NewV1(true)
```

If no suitable MAC address is available, generation falls back to a random node identifier.
Version 1 is intended for compatibility with systems that require its Gregorian timestamp layout.
For new time-ordered identifiers, version 7 is generally the better choice.

#### UUID version 2

`NewV2` creates a DCE Security UUID by replacing parts of the UUIDv1 layout with a 32-bit local identifier and an 8-bit local domain.
Use `Person` for a POSIX user ID, `Group` for a POSIX group ID, or `Organisation` for an organisation-defined identifier.

```go
user := uuid.NewV2(uuid.Person, 1000)
group := uuid.NewV2(uuid.Group, 1000)
```

The node identifier is random by default.
Pass `true` as the final argument to request a MAC address, with a random node identifier used as a fallback.
An unsupported `LocalDomain` value returns the Nil UUID.

UUIDv2 replaces the low 32 timestamp bits and the low 8 clock-sequence bits used by UUIDv1.
Its timestamp therefore has a resolution of approximately seven minutes, and only 6 clock-sequence bits remain.
At most 64 distinct clock-sequence values are available for the same node, local domain, local ID, and timestamp interval.
Embedding a UID, GID, or MAC address can also expose host identity information, so UUIDv2 should be used only for DCE compatibility.

#### UUID version 3

`NewV3` creates a deterministic, name-based UUID by hashing a namespace UUID followed by a name with MD5.
The same namespace and name produce the same UUID, while the namespace keeps identical names in different naming systems separate.

RFC 9562 defines namespace UUIDs for DNS names, URLs, OIDs, and X.500 Distinguished Names.
Select one with the `DNS`, `URL`, `OID`, or `X500` constant.

```go
u := uuid.NewV3(uuid.DNS, []byte("www.example.com"))
```

`NewV3` accepts any byte sequence and does not canonicalise it.
Callers are responsible for converting names to the canonical representation required by their namespace.
An unsupported `Namespace` value returns the Nil UUID.
Because UUIDv3 uses MD5, it is intended for compatibility rather than security-sensitive applications.

#### UUID version 5

`NewV5` creates a deterministic, name-based UUID in the same way as `NewV3`, but uses SHA-1 instead of MD5.
It accepts the same `DNS`, `URL`, `OID`, and `X500` namespace constants.

```go
u := uuid.NewV5(uuid.DNS, []byte("www.example.com"))
```

SHA-1 produces 160 bits, of which `NewV5` uses the first 128 bits before setting the version and variant fields.
As with `NewV3`, callers are responsible for canonicalising names, and an unsupported `Namespace` value returns the Nil UUID.
UUIDv5 is generally preferable to UUIDv3 for deterministic name-based identifiers, but SHA-1 should not be used for security-sensitive hashing.

#### UUID version 6

`NewV6` uses the same Gregorian timestamp, clock sequence, and node identifier as UUIDv1, but stores the timestamp from most to least significant.
This makes the UUIDs sort in generation order and improves locality in database indexes.

As with `NewV1`, the node identifier is random by default.
Pass `true` to request a MAC address, with a random node identifier used as a fallback:

```go
ordered := uuid.NewV6()
hardwareBased := uuid.NewV6(true)
```

UUIDv6 is primarily useful for systems that need UUIDv1-compatible fields with better sorting behaviour.
New systems that do not require UUIDv1 compatibility should generally use UUIDv7.

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
