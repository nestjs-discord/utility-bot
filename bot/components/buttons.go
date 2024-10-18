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
//
// NOTE:
// When making changes to the existing JSON fields, it's important to be cautious.
// Adding more fields is safe, but altering or removing existing ones
// can cause issues when the bot receives updates from old messages.
//
// Remember, marking the JSON tag with omitempty reduces the payload size when the fields are empty.
type CustomID struct {
	Action    string `json:"a"`
	UserId    string `json:"u,omitempty"`
	ChannelId string `json:"c,omitempty"`
	FormId    string `json:"f,omitempty"`
	// More fields can be added here...
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
