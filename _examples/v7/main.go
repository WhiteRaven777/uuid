package main

import (
	"fmt"
	"time"

	"github.com/WhiteRaven777/uuid"
)

func main() {
	var (
		precisions = []uuid.Precision{
			uuid.PrecisionDefault,
			uuid.PrecisionLow,
			uuid.PrecisionMid,
			uuid.PrecisionHigh,
		}

		uid, reUid uuid.UUID
		t          time.Time
		us         string
	)

	id := uuid.New(uuid.UUIDv7)

	for _, p := range precisions {
		id.(*uuid.V7).SetPrecision(p)

		// automatically generate based on generation timing
		uid = id.Generate()

		us = uid.String()

		// you can decode (parse) the string into the UUID
		if buf, err := uuid.Parse(us); err == nil {
			t = id.Decode(buf)
		}

		// generate based on any time
		reUid = id.Convert(t)

		// you can convert it to a string

		fmt.Printf(`precision: %s
  - time resolution:       %s
  - [org]   UUID(default): %s
  - [org]   UUID(short):   %s
  - [renew] UUID(URN):     %s
  - [renew] UUID(base64):  %s
`,
			p.String(),
			t.UTC().Format(time.RFC3339Nano),
			us,
			uid.Short(),
			reUid.URN(),
			reUid.Base64(),
		)
	}

}
