# uuid

[English](./README.md) | [日本語](./README_ja.md)

[v0.0.1 から v1.0.0 への移行](./MIGRATION_ja.md)

`github.com/WhiteRaven777/uuid` は、Go 1.27 の標準ライブラリに含まれる [`uuid`](https://pkg.go.dev/uuid) パッケージを拡張する小さなパッケージです。

## 標準ライブラリのラッパー

標準 API を薄いラッパーとして公開しています。

- `Parse`、`MustParse`、`New`、`Nil`、`Max`、`NewV4`、`NewV7`
- `UUID.String`、`MarshalText`、`AppendText`、`UnmarshalText`、`Compare`

このパッケージの `UUID` 型は、標準ライブラリの型と同じ `[16]byte` 表現を持ちます。
UUID の解析、version 4 の生成、version 7 の生成は標準ライブラリへ委譲しています。

## UUID の各バージョン

このパッケージは UUIDv1 から UUIDv8 までの全バージョンを提供します。
UUIDv4 と UUIDv7 の生成は標準ライブラリへ委譲し、それ以外のバージョンはこのパッケージで実装しています。

| バージョン | 適した用途 | 主な特徴 | 取り出せる時間情報 | 注意点 |
| --- | --- | --- | --- | --- |
| v1 | Gregorian time 形式を要求する既存システム | 時刻、clock sequence、Node ID を格納 | Gregorian timestamp | UUID 自体は時系列順に並ばず、時刻や Node ID が露出する可能性があるため、新規システムでは v7 を推奨 |
| v2 | UID、GID、local ID を埋め込む DCE Security 互換用途 | v1 の一部を local ID と domain で置換 | 約 7 分の分解能を持つ概算の Gregorian timestamp | clock sequence は 6 bit のみ |
| v3 | MD5 による決定的な名前ベース UUID を要求する既存システム | 同じ namespace と name から常に同じ UUID を生成 | なし | MD5 は互換性維持のために使用し、可能なら v5 を推奨 |
| v4 | 内容を推測させない汎用 ID | 122 bit の乱数 | なし | 生成時刻順には並ばない |
| v5 | 決定的な名前ベース ID | namespace と name の SHA-1 hash | なし | name を一貫した方法で正規化する必要がある |
| v6 | UUIDv1 とフィールド互換で時刻順に並ぶ ID | 並べ替えた Gregorian timestamp により時系列順にソート可能 | Gregorian timestamp | v1 互換性が不要なら v7 を推奨 |
| v7 | 汎用的な時系列 ID | Unix millisecond timestamp により時系列順にソート可能 | millisecond 精度の Unix timestamp | 新しいデータベース用 ID の既定候補 |
| v8 | 文書化された独自 UUID 形式 | 122 bit を用途ごとに定義 | 独自配置による | 一意性と相互運用性はアプリケーション側の責任 |

各バージョンの説明には、128 bit の配置を RFC 形式の図で示しています。
`ver` は version、`var` は RFC variant です。

#### UUID バージョン 1

`NewV1` は、UUIDv1 のフィールド配置を要求するシステムとの互換性を目的とした Gregorian time-based UUID を生成します。
UUIDv1 は、60 bit の Gregorian timestamp、14 bit の clock sequence、48 bit の Node ID で構成されます。

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

`NewV1` は既定で、ハードウェアの MAC アドレスを公開しないよう、暗号学的に安全な乱数から生成した Node ID を使用します。

```go
u := uuid.NewV1()
createdAt, ok := u.Time()
```

Node ID に MAC アドレスを使う場合は、`true` を渡します。

```go
u := uuid.NewV1(true)
```

適切な MAC アドレスを取得できない場合は、ランダムな Node ID へフォールバックします。
version 1 は、Gregorian timestamp 形式を必要とするシステムとの互換性を目的としています。
新しい時系列 UUID には、通常 version 7 のほうが適しています。

#### UUID バージョン 2

`NewV2` は UUIDv1 の一部を 32 bit の local ID と 8 bit の local domain で置き換え、DCE Security UUID を生成します。
POSIX user ID には `Person`、POSIX group ID には `Group`、組織独自の ID には `Organisation` を指定します。

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

UUIDv2 の `Time` は、保持されている timestamp の bit が表す約 7 分区間の開始時刻を返します。

Node ID は既定でランダムです。
最後の引数に `true` を渡すと MAC アドレスを使い、取得できない場合はランダムな Node ID へフォールバックします。
未対応の `LocalDomain` 値を渡すと Nil UUID を返します。

UUIDv2 は、UUIDv1 の timestamp 下位 32 bit と clock sequence 下位 8 bit を置き換えます。
そのため timestamp の分解能は約 7 分となり、clock sequence は 6 bit しか残りません。
同じ Node ID、local domain、local ID、timestamp 区間で利用できる clock sequence は最大 64 通りです。
UID、GID、MAC アドレスを埋め込むとホストの識別情報が露出する可能性もあるため、UUIDv2 は DCE 互換性が必要な場合に限って使用してください。

#### UUID バージョン 3

`NewV3` は、namespace UUID と name を連結した値を MD5 でハッシュし、決定的な名前ベース UUID を生成します。
同じ namespace と name からは同じ UUID が生成され、namespace が異なれば同じ name でも別の UUID になります。

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

RFC 9562 では、DNS name、URL、OID、X.500 Distinguished Name の namespace UUID が定義されています。
`DNS`、`URL`、`OID`、`X500` のいずれかの定数を選択します。

```go
u := uuid.NewV3(uuid.DNS, []byte("www.example.com"))
```

`NewV3` は任意のバイト列を受け取り、正規化は行いません。
呼び出し側で、namespace に応じた正規のバイト表現へ name を変換する必要があります。
未対応の `Namespace` 値を渡すと Nil UUID を返します。
UUIDv3 は MD5 を使用するため、セキュリティ用途ではなく互換性維持を目的としています。

#### UUID バージョン 4

`NewV4` は Go 標準ライブラリへ処理を委譲し、暗号学的に安全な乱数から内容を推測しにくい UUID を生成します。

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

UUIDv4 は、生成時刻、順序、アプリケーション固有情報を ID から読み取れないようにしたい場合に適しています。
version と variant を除く 122 bit をすべて乱数にする場合は UUIDv4 を使用します。
時刻情報を持たないため、文字列表現やバイナリ表現の並び順は生成順を表しません。

#### UUID バージョン 5

`NewV5` は `NewV3` と同じ方法で決定的な名前ベース UUID を生成しますが、MD5 の代わりに SHA-1 を使用します。
namespace には、同じ `DNS`、`URL`、`OID`、`X500` 定数を指定できます。

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

SHA-1 が生成する 160 bit のうち、`NewV5` は先頭 128 bit を使用してから version と variant のフィールドを設定します。
`NewV3` と同様、name の正規化は呼び出し側で行い、未対応の `Namespace` 値を渡すと Nil UUID を返します。
決定的な名前ベース UUID には通常 UUIDv3 より UUIDv5 が適していますが、SHA-1 をセキュリティ用途のハッシュとして使うべきではありません。

#### UUID バージョン 6

`NewV6` は UUIDv1 と同じ Gregorian timestamp、clock sequence、Node ID を使用し、timestamp を最上位側から順番に格納します。
これにより、UUID が生成順に並び、データベースインデックスの局所性が改善します。

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

`NewV1` と同様、Node ID は既定でランダムです。
MAC アドレスを使う場合は `true` を渡します。
取得できない場合はランダムな Node ID へフォールバックします。

```go
ordered := uuid.NewV6()
hardwareBased := uuid.NewV6(true)
createdAt, ok := ordered.Time()
```

UUIDv6 は、UUIDv1 互換のフィールドと生成時刻順のソートを必要とするシステムに適しています。
UUIDv1 との互換性を必要としない新しいシステムには、通常 UUIDv7 が適しています。

#### UUID バージョン 7

`NewV7` は Go 標準ライブラリへ処理を委譲し、時刻順に並ぶ UUID を生成します。
48 bit の Unix millisecond timestamp に乱数を続けて格納します。

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

UUIDv7 は、MAC アドレスや UUIDv1、UUIDv6 の Gregorian timestamp 形式を使用せず、データベースインデックスの局所性を得られます。
生成時刻順に並ぶ ID が必要な新規システムでは、汎用的な選択肢として UUIDv7 を推奨します。

#### UUID バージョン 8

`NewV8` は RFC 9562 で定義された 3 つのカスタム領域から、実験的またはベンダー固有の UUID を構築します。

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

有効範囲は、`customA` が `0x000000000000` から `0xffffffffffff`、`customB` が `0x000` から `0xfff`、`customC` が `0x0000000000000000` から `0x3fffffffffffffff` です。
範囲を超えた値には各領域の bit mask を適用し、上位 bit を切り捨てます。
version を 8、variant を `10` に設定します。

UUIDv8 ではカスタム領域の生成方法が規定されておらず、一意性も保証されません。
用途ごとに定義した構造から生成する方法や、hash digest を 3 つの領域に分割する方法などを利用できます。
利用可能な 122 bit をすべて乱数で埋める場合は UUIDv4 を使用してください。
通常の時系列 ID には UUIDv7 が適しています。

## 独自拡張

### 埋め込まれた時刻

`UUID.Time()` は、UUIDv1、UUIDv2、UUIDv6、UUIDv7 に埋め込まれた時刻を返します。

```go
createdAt, ok := u.Time()
```

| UUID version | 返す時刻 | 精度と挙動 |
| --- | --- | --- |
| v1 | Gregorian timestamp | UUID の 100 nanosecond timestamp 単位で格納 |
| v2 | 概算の Gregorian timestamp | 保持されている timestamp の bit が表す約 7 分区間の開始時刻 |
| v6 | Gregorian timestamp | UUID の 100 nanosecond timestamp 単位で格納 |
| v7 | Unix timestamp | millisecond 精度 |
| v3、v4、v5、v8 | `time.Time` のゼロ値 | 標準の時刻フィールドがないため `false` を返す |

取得できた時刻は UTC で返します。

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
