# uuid

[English](./README.md) | [日本語](./README_ja.md)

`github.com/WhiteRaven777/uuid` は、Go 1.27 の標準ライブラリに含まれる [`uuid`](https://pkg.go.dev/uuid) パッケージを拡張する小さなパッケージです。

## 標準ライブラリのラッパー

標準 API を薄いラッパーとして公開しています。

- `Parse`、`MustParse`、`New`、`Nil`、`Max`、`NewV4`、`NewV7`
- `UUID.String`、`MarshalText`、`AppendText`、`UnmarshalText`、`Compare`

このパッケージの `UUID` 型は、標準ライブラリの型と同じ `[16]byte` 表現を持ちます。
UUID の解析、version 4 の生成、version 7 の生成は標準ライブラリへ委譲しています。

## 独自拡張

### 追加の出力形式

このパッケージの `UUID` 型には、標準ライブラリの型にはない 3 つの出力メソッドがあります。
次の UUID を例にします。

```text
01234567-89ab-7def-8123-456789abcdef
```

各メソッドは次の値を返します。

| メソッド | 形式 | 出力 |
| --- | --- | --- |
| `Short()` | ハイフンを含まない 32 桁の小文字 16 進数 | `0123456789ab7def8123456789abcdef` |
| `URN()` | UUID URN | `urn:uuid:01234567-89ab-7def-8123-456789abcdef` |
| `Base64()` | パディングを含まない標準 Base64 | `ASNFZ4mrfe+BI0VniavN7w` |

`Short` と `URN` の戻り値は `Parse` へ渡せます。
Base64 は出力専用形式であり、`Parse` では解析できません。

### 追加の UUID バージョン

#### UUID バージョン 1

`NewV1` は、標準ライブラリが提供する version 4 および version 7 の生成機能を補完します。
UUIDv1 は、60 bit の Gregorian timestamp、14 bit の clock sequence、48 bit の Node ID で構成されます。

`NewV1` は既定で、ハードウェアの MAC アドレスを公開しないよう、暗号学的に安全な乱数から生成した Node ID を使用します。

```go
u := uuid.NewV1()
```

Node ID に MAC アドレスを使う場合は、`true` を渡します。

```go
u := uuid.NewV1(true)
```

適切な MAC アドレスを取得できない場合は、ランダムな Node ID へフォールバックします。
version 1 は、Gregorian timestamp 形式を必要とするシステムとの互換性を目的としています。
新しい時系列 UUID には、通常 version 7 のほうが適しています。

#### UUID バージョン 6

`NewV6` は UUIDv1 と同じ Gregorian timestamp、clock sequence、Node ID を使用し、timestamp を最上位側から順番に格納します。
これにより、UUID が生成順に並び、データベースインデックスの局所性が改善します。

`NewV1` と同様、Node ID は既定でランダムです。
MAC アドレスを使う場合は `true` を渡します。
取得できない場合はランダムな Node ID へフォールバックします。

```go
ordered := uuid.NewV6()
hardwareBased := uuid.NewV6(true)
```

UUIDv6 は、UUIDv1 互換のフィールドと生成時刻順のソートを必要とするシステムに適しています。
UUIDv1 との互換性を必要としない新しいシステムには、通常 UUIDv7 が適しています。

## 動作要件

Go 1.27 以降が必要です。

## 使用例

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

## ライセンス

Copyright (c) 2021-present [WhiteRaven777](https://github.com/WhiteRaven777)

[MIT License](./LICENSE) の下で公開されています。
