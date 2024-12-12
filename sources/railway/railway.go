package railway

import (
	"github.com/ferretcode/pricetag/sink"
	"github.com/ferretcode/pricetag/types"
)

// REQUIRED ENVIRONMENT VARIABLES:
// RAILWAY_API_KEY=
// RAILWAY_ENVIRONMENT_ID=

func CreateSink() sink.Sink {
	newLogChan := make(chan []types.Log)

	return sink.Sink{
		NewLog: newLogChan,
	}
}
