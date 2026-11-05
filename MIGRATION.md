# Migrating from v0.0.1 to v1.0.0

[English](./MIGRATION.md) | [日本語](./MIGRATION_ja.md)

Version 1.0.0 replaces the draft-based UUIDv7 implementation in v0.0.1 with an extension of the Go 1.27 standard library's `uuid` package.
This is a breaking release and requires Go 1.27 or newer.

## API changes

| v0.0.1 | v1.0.0 | Migration |
| --- | --- | --- |
| `uuid.New(uuid.UUIDv7).Generate()` | `uuid.NewV7()` | Use the version-specific constructor directly. |
| `uuid.New(version)` | `uuid.New()` | `New` no longer accepts a version and delegates to the standard library's default generator. Use `NewV1` through `NewV8` when the version must be explicit. |
| `id.Decode(u)` | `u.Time()` | Check the returned Boolean because only UUIDv1, UUIDv2, UUIDv6, and UUIDv7 contain standard time fields. |
| `id.Convert(t)` | No direct replacement | Generating UUIDv7 for an arbitrary supplied time is no longer supported. |
| `id.SetPrecision(p)` | No replacement | UUIDv7 now uses the standard RFC 9562 millisecond timestamp layout. |
| `uuid.Parse(value)` | `uuid.Parse(text)` | `Parse` now accepts only strings. |
| `u.Version()` | No direct replacement | Read the version with `u[6] >> 4` if needed. |
| `Version`, `Variant`, `ID`, and `V7` | Removed | Use the version-specific constructors and the `UUID` methods. |

The `Short`, `String`, `URN`, and `Base64` output methods retain their previous purposes.

## UUIDv7 generation and time extraction

Replace the generator interface:

```go
id := uuid.New(uuid.UUIDv7)
u := id.Generate()
createdAt := id.Decode(u)
```

with the version-specific constructor and `Time` method:

```go
u := uuid.NewV7()
createdAt, ok := u.Time()
if !ok {
	// The UUID has no standard time field.
}
```

UUIDv7 generation is delegated to the Go standard library and follows RFC 9562.
The custom millisecond, microsecond, and nanosecond precision modes from v0.0.1 have been removed.

## Existing v0.0.1 UUIDv7 values

UUIDv7 in v0.0.1 followed an earlier draft layout rather than the final RFC 9562 layout.
Its textual values can still be parsed as UUIDs, but `UUID.Time()` in v1.0.0 must not be used to recover their original timestamps because it interprets UUIDv7 using the final 48-bit Unix millisecond field.

Decode or export timestamps from persisted v0.0.1 UUIDv7 values with the v0.0.1 implementation before migrating.
Applications that must retain old values should record which layout produced each value or preserve a legacy decoder.

## Parsing

`Parse` now wraps the standard library and accepts a string:

```go
u, err := uuid.Parse("01234567-89ab-7def-8123-456789abcdef")
```

Passing `[]byte` directly is no longer supported.
Copy a 16-byte binary value into `UUID` explicitly:

```go
var u uuid.UUID
if len(data) != len(u) {
	return errors.New("invalid UUID length")
}
copy(u[:], data)
```

Base64 remains available as an output format through `UUID.Base64()`, but `Parse` no longer accepts Base64 input.
Decode Base64 explicitly with `encoding/base64` before copying the resulting 16 bytes into a `UUID`.

## UUID representation

The local `UUID` type remains a 16-byte array and has the same representation as the standard library's `uuid.UUID` type.
Conversion between the types is explicit:

```go
standard := stduuid.UUID(local)
local = uuid.UUID(standard)
```

## Newly available functionality

Version 1.0.0 adds:

- UUIDv1 and UUIDv6 with random node identifiers by default and optional MAC address use
- DCE Security UUIDv2
- name-based UUIDv3 and UUIDv5
- custom UUIDv8 fields
- standard library wrappers for UUIDv4 and UUIDv7
- `MustParse`, `Nil`, `Max`, text encoding methods, and `Compare`
- embedded time extraction through `UUID.Time()`

