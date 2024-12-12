package railway

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
)

func ReconstructLogLines(logs []railwayLog) ([]byte, error) {
	var jsonObject []byte

	jsonObject = append(jsonObject, []byte(`[`)...)

	for i := range logs {
		logObject, err := ReconstructLogLine(logs[i])
		if err != nil {
			return nil, err
		}

		jsonObject = append(jsonObject, logObject...)

		if (i + 1) < len(logs) {
			jsonObject = append(jsonObject, []byte(`,`)...)
		}
	}

	jsonObject = append(jsonObject, []byte(`[`)...)

	return jsonObject, nil
}

func ReconstructLogLine(log railwayLog) ([]byte, error) {
	jsonObject := []byte("{}")

	jsonObject, err := jsonparser.Set(jsonObject, []byte(strconv.Quote(log.Message)), "message")
	if err != nil {
		return nil, fmt.Errorf("failed to append message attribute to object: %w", err)
	}

	metadata, err := json.Marshal(log.Tags)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata object: %w", err)
	}

	jsonObject, err = jsonparser.Set(jsonObject, metadata, "_metadata")
	if err != nil {
		return nil, fmt.Errorf("failed to append metadata attribute to object: %w", err)
	}

	for i := range log.Attributes {
		jsonObject, err = jsonparser.Set(jsonObject, []byte(log.Attributes[i].Value), log.Attributes[i].Key)
		if err != nil {
			return nil, fmt.Errorf("failed to append json attribute to object: %w", err)
		}
	}

	timeStamp := []byte(strconv.Quote(log.Timestamp.Format(time.RFC3339Nano)))

	jsonObject, err = jsonparser.Set(jsonObject, timeStamp, "timestamp")
	if err != nil {
		return nil, fmt.Errorf("failed to append timestamp attribute to object: %w", err)
	}

	// set severity in all situations for backwards compatibility
	// railway already normilizes the level attribute into the severity field, or vice versa
	jsonObject, err = jsonparser.Set(jsonObject, []byte(strconv.Quote(log.Severity)), "severity")
	if err != nil {
		return nil, fmt.Errorf("failed to append severity attribute to object: %w", err)
	}

	return jsonObject, nil
}
