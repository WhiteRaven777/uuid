# uuid

`github.com/WhiteRaven777/uuid` is a small extension of the Go 1.27 standard
library's [`uuid`](https://pkg.go.dev/uuid) package.

The standard API is exposed through thin wrappers:

- `Parse`, `MustParse`, `New`, `Nil`, `Max`, `NewV4`, and `NewV7`
- `UUID.String`, `MarshalText`, `AppendText`, `UnmarshalText`, and `Compare`

The local `UUID` type has the same `[16]byte` representation and adds `Short`,
`URN`, and `Base64` convenience encodings.

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
