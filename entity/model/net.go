package model

type DataPacket struct {
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}

type PlayerIDPacket struct {
	PlayerID string `json:"playerID"`
}

type ClientInputPacket struct {
}
