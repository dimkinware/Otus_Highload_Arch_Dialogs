package entity

type DialogMessage struct {
	Id   string `json:"id"`
	From string `json:"author_id"`
	To   string `json:"receiver_uid"`
	Text string `json:"message"`
	Time int64  `json:"message_time"`
}
