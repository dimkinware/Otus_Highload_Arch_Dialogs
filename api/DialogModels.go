package api

type DialogMessage struct {
	Id   string `json:"id"`
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
	Time int64  `json:"message_time"`
}

type DialogMessageSendApiModel struct {
	Text string `json:"text"`
}
