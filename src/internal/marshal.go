package internal

import (
	"encoding/json"
	"fmt"
	"os"
)

func MarshalAndWriteToFile(v any, filePath string) error {
	bytes, err := marshal(v)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filePath, bytes, 0o644); err != nil {
		return fmt.Errorf("error when writing to file: %s, err: %s", filePath, err)
	}
	return nil
}

func GetAndUnmarshalFile(filePath string, v any) error {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("get content from cache error: %s", err)
	}
	if err := json.Unmarshal(bytes, &v); err != nil {
		return fmt.Errorf("unmarshaling error: %s", err)
	}
	return nil
}

func marshal(v any) ([]byte, error) {
	bytes, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("error marshalling: %s", err)
	}
	return bytes, nil
}
