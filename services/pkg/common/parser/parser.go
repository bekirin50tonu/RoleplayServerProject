package parser

import "encoding/json"

func ParseDataWithInterface[T any](data []byte) (*T, error) {
	var x T
	err := json.Unmarshal(data, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}
