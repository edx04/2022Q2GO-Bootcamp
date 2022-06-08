package entity

type Quote struct {
	Id     int64  `json:"id"`
	Author string `json:"author"`
	Text   string `json:"text"`
}
