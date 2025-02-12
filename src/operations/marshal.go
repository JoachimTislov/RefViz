package op

import (
	"encoding/json"
	"fmt"
	"os"
)

func marshal(v any) ([]byte, error) {
	bytes, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("error marshalling: %s", err)
	}
	return bytes, nil
}

func marshalAndWriteToFile(v any, filePath string) error {
	bytes, err := marshal(v)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filePath, bytes, 0o644); err != nil {
		return fmt.Errorf("error when writing to file: %s, err: %s", filePath, err)
	}
	return nil
}
