# uuid

[English](./README.md) | [日本語](./README_ja.md)

[Migration from v0.0.1 to v1.0.0](./MIGRATION.md)

`github.com/WhiteRaven777/uuid` is a small extension of the Go 1.27 standard library's [`uuid`](https://pkg.go.dev/uuid) package.

## Standard library wrappers

The standard API is exposed through thin wrappers:

- `Parse`, `MustParse`, `New`, `Nil`, `Max`, `NewV4`, and `NewV7`
- `UUID.String`, `MarshalText`, `AppendText`, `UnmarshalText`, and `Compare`

The local `UUID` type has the same `[16]byte` representation as the standard library type.
UUID parsing, version 4 generation, and version 7 generation are delegated to the standard library.

## UUID versions

This package exposes every UUID version from 1 through 8.
UUIDv4 and UUIDv7 delegate generation to the standard library, while the other versions are implemented by this package.

| Version | Use when | Main characteristics | Recoverable time information | Important considerations |
| --- | --- | --- | --- | --- |
| v1 | A legacy system requires the Gregorian time layout | Time-based; includes a clock sequence and node ID | Gregorian timestamp | Not naturally sortable by time; can expose time and node information; prefer v7 for new systems |
| v2 | DCE Security compatibility requires an embedded UID, GID, or local ID | Derived from v1 with a local ID and domain | Approximate Gregorian timestamp with about seven-minute resolution | Only 6 clock-sequence bits remain |
| v3 | An existing system requires deterministic MD5 name-based UUIDs | Same namespace and name always produce the same UUID | None | MD5 is retained for compatibility; prefer v5 where possible |
| v4 | A general-purpose opaque identifier is required | 122 random bits | None | Not naturally ordered |
| v5 | A deterministic name-based identifier is required | SHA-1 hash of a namespace and name | None | Names must be canonicalised consistently |
| v6 | A sortable UUIDv1-compatible representation is required | Sortable by its reordered Gregorian timestamp | Gregorian timestamp | Prefer v7 unless v1 field compatibility matters |
| v7 | A general-purpose, time-ordered identifier is required | Sortable by its Unix millisecond timestamp | Unix timestamp with millisecond resolution | Recommended default for new database identifiers |
| v8 | An application needs a documented custom UUID layout | 122 application-defined bits | Depends on the custom layout | Uniqueness and interoperability are the application's responsibility |

Each section below includes an RFC-style diagram of its 128-bit layout.
`ver` is the version field and `var` is the RFC variant field.

#### UUID version 1

`NewV1` creates a Gregorian time-based UUID for compatibility with systems that require the UUIDv1 field layout.
UUIDv1 contains a 60-bit Gregorian timestamp, a 14-bit clock sequence, and a 48-bit node identifier.

```text
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                           time_low                            |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|          time_mid             |       time_hi_and_version     |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
| clk_seq_hi_res |  clk_seq_low |          node (0-1)           |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                         node (2-5)                            |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

By default, `NewV1` uses a cryptographically random node identifier rather than exposing a hardware MAC address:

```go
u := uuid.NewV1()
createdAt, ok := u.Time()
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

```text
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                           local_id                            |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|          time_mid             |       time_hi_and_version     |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
| clk_seq_hi_res | local_domain |          node (0-1)           |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                         node (2-5)                            |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

```go
user := uuid.NewV2(uuid.Person, 1000)
group := uuid.NewV2(uuid.Group, 1000)
intervalStart, ok := user.Time()
```

For UUIDv2, `Time` returns the beginning of the approximately seven-minute interval represented by the retained timestamp bits.

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

```text
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                            hash_a                             |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|            hash_a             |  ver  |        hash_b          |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|var|                         hash_c                            |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                            hash_c                             |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

RFC 9562 defines namespace UUIDs for DNS names, URLs, OIDs, and X.500 Distinguished Names.
Select one with the `DNS`, `URL`, `OID`, or `X500` constant.

```go
u := uuid.NewV3(uuid.DNS, []byte("www.example.com"))
```

`NewV3` accepts any byte sequence and does not canonicalise it.
Callers are responsible for converting names to the canonical representation required by their namespace.
An unsupported `Namespace` value returns the Nil UUID.
Because UUIDv3 uses MD5, it is intended for compatibility rather than security-sensitive applications.

#### UUID version 4

`NewV4` creates an opaque UUID from cryptographically secure random data by delegating to the Go standard library.

```text
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                           random_a                            |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|          random_a             |  ver  |       random_b        |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|var|                        random_c                           |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                           random_c                            |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

```go
u := uuid.NewV4()
```

UUIDv4 is suitable when identifiers should not reveal creation time, ordering, or application-specific information.
It is the appropriate choice when all 122 non-version and non-variant bits are random.
Because UUIDv4 has no time component, its textual or binary order does not represent creation order.

#### UUID version 5

`NewV5` creates a deterministic, name-based UUID in the same way as `NewV3`, but uses SHA-1 instead of MD5.
It accepts the same `DNS`, `URL`, `OID`, and `X500` namespace constants.

```text
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                            hash_a                             |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|            hash_a             |  ver  |        hash_b          |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|var|                         hash_c                            |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                            hash_c                             |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

```go
u := uuid.NewV5(uuid.DNS, []byte("www.example.com"))
```

SHA-1 produces 160 bits, of which `NewV5` uses the first 128 bits before setting the version and variant fields.
As with `NewV3`, callers are responsible for canonicalising names, and an unsupported `Namespace` value returns the Nil UUID.
UUIDv5 is generally preferable to UUIDv3 for deterministic name-based identifiers, but SHA-1 should not be used for security-sensitive hashing.

#### UUID version 6

`NewV6` uses the same Gregorian timestamp, clock sequence, and node identifier as UUIDv1, but stores the timestamp from most to least significant.
This makes the UUIDs sort in generation order and improves locality in database indexes.

```text
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                           time_high                           |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|           time_mid            |      time_low_and_version     |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
| clk_seq_hi_res |  clk_seq_low |          node (0-1)           |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                         node (2-5)                            |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

As with `NewV1`, the node identifier is random by default.
Pass `true` to request a MAC address, with a random node identifier used as a fallback:

```go
ordered := uuid.NewV6()
hardwareBased := uuid.NewV6(true)
createdAt, ok := ordered.Time()
```

UUIDv6 is primarily useful for systems that need UUIDv1-compatible fields with better sorting behaviour.
New systems that do not require UUIDv1 compatibility should generally use UUIDv7.

#### UUID version 7

`NewV7` creates a time-ordered UUID by delegating to the Go standard library.
It stores a 48-bit Unix timestamp in milliseconds followed by random data.

```text
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                           unix_ts_ms                          |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|          unix_ts_ms           |  ver  |        rand_a         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|var|                         rand_b                            |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                            rand_b                             |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

```go
u := uuid.NewV7()
createdAt, ok := u.Time()
```

UUIDv7 provides useful database index locality without embedding a MAC address or using the Gregorian timestamp layout of UUIDv1 and UUIDv6.
It is the recommended general-purpose choice for new systems that need identifiers to follow creation time.

#### UUID version 8

`NewV8` constructs an experimental or vendor-specific UUID from the three custom fields defined by RFC 9562.

```text
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                           custom_a                            |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|          custom_a             |  ver  |       custom_b        |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|var|                        custom_c                           |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                           custom_c                            |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

```go
u := uuid.NewV8(0x123456789abc, 0xdef, 0x0123456789abcdef)
```

The valid ranges are `0x000000000000` to `0xffffffffffff` for `customA`, `0x000` to `0xfff` for `customB`, and `0x0000000000000000` to `0x3fffffffffffffff` for `customC`.
Values outside these ranges are masked to the field width, discarding their more significant bits.
The version and variant fields are set to version 8 and variant `10` respectively.

UUIDv8 does not prescribe how the custom fields are produced and does not guarantee uniqueness.
Applications may derive them from a documented structure or split a digest across the three fields.
If all 122 available bits would contain only random data, use UUIDv4 instead.
For ordinary time-ordered identifiers, use UUIDv7.

## Extensions

### Embedded time

`UUID.Time()` returns the time embedded in UUIDv1, UUIDv2, UUIDv6, and UUIDv7 values.

```go
createdAt, ok := u.Time()
```

| UUID version | Returned time | Precision and behaviour |
| --- | --- | --- |
| v1 | Gregorian timestamp | Stored in 100-nanosecond UUID timestamp units |
| v2 | Approximate Gregorian timestamp | Beginning of the approximately seven-minute interval represented by the retained timestamp bits |
| v6 | Gregorian timestamp | Stored in 100-nanosecond UUID timestamp units |
| v7 | Unix timestamp | Millisecond precision |
| v3, v4, v5, v8 | Zero time | Returns `false` because these versions have no standard time field |

Successful results are returned in UTC.

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
