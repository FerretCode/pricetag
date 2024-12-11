package sink

import (
	"github.com/ferretcode/pricetag/types"
)

type Sink struct {
	NewLog chan types.Log
}
