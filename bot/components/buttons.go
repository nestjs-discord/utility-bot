package components

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrInvalidAction   = errors.New("invalid action")
	ErrCustomIdTooLong = errors.New("custom id too long")
)

// CustomID
// See: https://discord.com/developers/docs/interactions/message-components#custom-id
type CustomID struct {
	Action    string `json:"a"`
	UserId    string `json:"u,omitempty"`
	ChannelId string `json:"c,omitempty"`
	FormId    string `json:"f,omitempty"`
	// TODO: do we need other fields?
}

func EncodeCustomId(id *CustomID) (string, error) {
	if id.Action == "" {
		return "", ErrInvalidAction
	}

	encoded, err := json.Marshal(id)
	if err != nil {
		return "", fmt.Errorf("encode failed: %v", err)
	}

	result := string(encoded)
	if len(result) > 99 {
		return "", ErrCustomIdTooLong
	}

	return result, nil
}

func DecodeCustomId(encoded string) (*CustomID, error) {
	var id CustomID
	err := json.Unmarshal([]byte(encoded), &id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
