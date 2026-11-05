# uuid

[English](./README.md) | [日本語](./README_ja.md)

`github.com/WhiteRaven777/uuid` は、Go 1.27の標準ライブラリに含まれる
[`uuid`](https://pkg.go.dev/uuid) パッケージを拡張する小さなパッケージです。

## 標準ライブラリのwrapper

標準APIを薄いwrapperとして公開しています。

- `Parse`、`MustParse`、`New`、`Nil`、`Max`、`NewV4`、`NewV7`
- `UUID.String`、`MarshalText`、`AppendText`、`UnmarshalText`、`Compare`

このパッケージの `UUID` 型は、標準ライブラリの型と同じ `[16]byte` 表現を
持ちます。UUIDの解析、version 4の生成、version 7の生成は標準ライブラリへ
委譲しています。

## 独自拡張

### 追加の出力形式

このパッケージの `UUID` 型には、標準ライブラリの型にはない3つの出力
メソッドがあります。次のUUIDを例にすると、

```text
01234567-89ab-7def-8123-456789abcdef
```

各メソッドは次の値を返します。

| メソッド | 形式 | 出力 |
| --- | --- | --- |
| `Short()` | ハイフンを含まない32桁の小文字16進数 | `0123456789ab7def8123456789abcdef` |
| `URN()` | UUID URN | `urn:uuid:01234567-89ab-7def-8123-456789abcdef` |
| `Base64()` | パディングを含まない標準Base64 | `ASNFZ4mrfe+BI0VniavN7w` |

`Short` と `URN` の戻り値は `Parse` へ渡せます。Base64は出力専用形式であり、
`Parse` では解析できません。

### 追加のUUID version

#### UUID version 1

`NewV1` は、標準ライブラリが提供するversion 4およびversion 7の生成機能を
補完します。UUIDv1は、60bitのGregorian timestamp、14bitのclock sequence、
48bitのNode IDで構成されます。

`NewV1` は既定で、ハードウェアのMACアドレスを公開しないよう、暗号学的に
安全な乱数から生成したNode IDを使用します。

```go
u := uuid.NewV1()
```

Node IDとしてMACアドレスの使用を要求する場合は、`true` を渡します。

```go
u := uuid.NewV1(true)
```

適切なMACアドレスを取得できない場合は、ランダムなNode IDへフォールバック
します。version 1は、Gregorian timestamp形式を必要とするシステムとの互換性を
目的としています。新しい時系列UUIDには、通常version 7のほうが適しています。

## 動作要件

Go 1.27以降が必要です。

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
