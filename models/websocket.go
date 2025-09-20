package models

import "encoding/json"

// TODO: improve this so it can accept any kinds of data
type WSMessage struct {
	Event  string         `json:"event"`
	Ticket TicketResponse `json:"ticket"`
}

func (m *WSMessage) ToJsonMarshal() ([]byte, error) {
	msg, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
