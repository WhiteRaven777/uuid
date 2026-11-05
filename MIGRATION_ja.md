# v0.0.1 から v1.0.0 への移行

[English](./MIGRATION.md) | [日本語](./MIGRATION_ja.md)

v1.0.0 では、v0.0.1 のドラフト仕様に基づく UUIDv7 実装を廃止し、Go 1.27 の標準ライブラリに含まれる `uuid` パッケージの拡張として再実装しました。
このリリースには破壊的変更があり、Go 1.27 以降が必要です。

## API の変更

| v0.0.1 | v1.0.0 | 移行方法 |
| --- | --- | --- |
| `uuid.New(uuid.UUIDv7).Generate()` | `uuid.NewV7()` | version ごとの constructor を直接使用します。 |
| `uuid.New(version)` | `uuid.New()` | `New` は version を受け取らず、標準ライブラリの既定 generator へ委譲します。version を明示する場合は `NewV1` から `NewV8` を使用します。 |
| `id.Decode(u)` | `u.Time()` | 標準の時刻フィールドを持つ UUIDv1、UUIDv2、UUIDv6、UUIDv7 だけが時刻を返すため、Boolean の戻り値を確認します。 |
| `id.Convert(t)` | 直接の代替なし | 指定した任意時刻から UUIDv7 を生成する機能は廃止しました。 |
| `id.SetPrecision(p)` | 代替なし | UUIDv7 は RFC 9562 の millisecond timestamp 形式を使用します。 |
| `uuid.Parse(value)` | `uuid.Parse(text)` | `Parse` は string だけを受け取ります。 |
| `u.Version()` | 直接の代替なし | 必要な場合は `u[6] >> 4` で version を取得します。 |
| `Version`、`Variant`、`ID`、`V7` | 廃止 | version ごとの constructor と `UUID` の method を使用します。 |

`Short`、`String`、`URN`、`Base64` の各出力 method は、以前と同じ用途で使用できます。

## UUIDv7 の生成と時刻の取得

generator interface を使用していたコードを変更します。

```go
id := uuid.New(uuid.UUIDv7)
u := id.Generate()
createdAt := id.Decode(u)
```

version ごとの constructor と `Time` method を使用します。

```go
u := uuid.NewV7()
createdAt, ok := u.Time()
if !ok {
	// 標準の時刻フィールドを持たない UUID です。
}
```

UUIDv7 の生成は Go 標準ライブラリへ委譲し、RFC 9562 に準拠します。
v0.0.1 にあった millisecond、microsecond、nanosecond の独自 precision mode は廃止しました。

## 既存の v0.0.1 UUIDv7

v0.0.1 の UUIDv7 は、最終版の RFC 9562 ではなく、以前のドラフト仕様に基づく配置を使用していました。
文字列表現は UUID として引き続き解析できますが、v1.0.0 の `UUID.Time()` で元の時刻を取得してはいけません。
この method は、UUIDv7 を最終仕様の 48 bit Unix millisecond field として解釈します。

永続化した v0.0.1 UUIDv7 の時刻は、移行前に v0.0.1 の実装で decode または export してください。
既存の値を維持する必要がある場合は、各値を生成した配置を記録するか、旧形式の decoder を残してください。

## Parse

`Parse` は標準ライブラリの処理を wrap し、string を受け取ります。

```go
u, err := uuid.Parse("01234567-89ab-7def-8123-456789abcdef")
```

`[]byte` を直接渡すことはできません。
16 byte のバイナリ値は、明示的に `UUID` へ copy します。

```go
var u uuid.UUID
if len(data) != len(u) {
	return errors.New("invalid UUID length")
}
copy(u[:], data)
```

Base64 は `UUID.Base64()` による出力形式として利用できますが、`Parse` では解析できません。
Base64 は `encoding/base64` で明示的に decode し、得られた 16 byte を `UUID` へ copy してください。

## UUID の表現

このパッケージの `UUID` 型は引き続き 16 byte の配列で、標準ライブラリの `uuid.UUID` 型と同じ表現を持ちます。
型の変換は明示的に行います。

```go
standard := stduuid.UUID(local)
local = uuid.UUID(standard)
```

## 新しく利用できる機能

v1.0.0 では、次の機能を追加しました。

- 既定でランダムな Node ID を使い、必要に応じて MAC アドレスも使用できる UUIDv1 と UUIDv6
- DCE Security UUIDv2
- 名前ベースの UUIDv3 と UUIDv5
- カスタム領域を持つ UUIDv8
- UUIDv4 と UUIDv7 の標準ライブラリ wrapper
- `MustParse`、`Nil`、`Max`、text encoding method、`Compare`
- `UUID.Time()` による埋め込み時刻の取得

