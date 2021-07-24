# uuid

[![GoDoc Widget]][GoDoc]

`uuid` provides asynchronous unique numbering on the server side, based on [a draft update of RFC4122](https://www.ietf.org/archive/id/draft-peabody-dispatch-new-uuid-format-01.html).


## Install
```bash
$ go get -u github.com/WhiteRaven777/uuid
```

## Features
* Generate
  * [ ] UUID v1
    * Method to generate UUIDs using time and MAC address
  * [ ] UUID v2
    * Method to generate UUIDs using POSIX user ID or group ID to replace part of UUID v1
  * [ ] UUID v3
    * Method to generate UUIDs using MD5 in combination with information that is unique in a namespace
  * [ ] UUID v4
    * Method to generate UUIDs using random numbers
  * [ ] UUID v5
    * Method to generate UUIDs using SHA1 based on UUID v3
  * [ ] UUID v6 (draft)
    * Method to generate UUIDs using time based on the Gregorian calendar as well as v1, the clock sequence and random numbers
  * [x] UUID v7 (draft)
    * Method to generate UUIDs using unix time with millisecond or microsecond or nanosecond precision and random numbers
  * [ ] UUID v8 (draft)
    * Method to generate UUIDs using any precision and format of time and node when no other method is available
* Encode
  * [x] Short string 
  * [x] UUID format
  * [x] URN format 
  * [x] base64 format
* Decode (parse)
  * [x] Short string
  * [x] UUID format
  * [x] URN format
  * [x] base64 format

## Example
See [_examples/](https://github.com/WhiteRaven777/uuid/blob/master/_examples/) for a variety of examples.

```go
package main

import (
	"fmt"
	"time"

	"github.com/WhiteRaven777/uuid"
)

func main() {
	var (
		id  uuid.ID
		uid uuid.UUID
		t time.Time
	)

	id = uuid.New(uuid.UUIDv7)

	// automatically generate based on generation timing
	uid = id.Generate()

	// you can convert it to a string
	fmt.Println("UUID(short)", uid.Short())
	fmt.Println("UUID(default)", uid.String())

	// you can decode (parse) the string into the UUID
	if buf, err := uuid.Parse(uid.String()); err == nil {
		t = id.Decode(buf)
	}
	
	// generate based on any time
	uid = id.Convert(t)

	// you can convert it to a string
	fmt.Println("URN", uid.URN())
	fmt.Println("base64", uid.Base64())
	
	// Note.
	// Regenerated UUIDs from the same time information will almost never have the same value.
	// This is because the time information is combined with random number information.
}
```

## License

Copyright (c) 2021-present [WhiteRaven777](https://github.com/WhiteRaven777)

Licensed under [MIT License](./LICENSE)

[GoDoc]: https://pkg.go.dev/github.com/WhiteRaven777/uuid?tab=versions
[GoDoc Widget]: https://godoc.org/github.com/WhiteRaven777/uuid?status.svg