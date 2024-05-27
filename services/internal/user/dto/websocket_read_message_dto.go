package dto

import "encoding/json"

type NewMessageDTO struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func NewWebSocketReadMessageDTO(messageByte []byte) (*NewMessageDTO, error) {

	var message NewMessageDTO
	err := json.Unmarshal(messageByte, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}
