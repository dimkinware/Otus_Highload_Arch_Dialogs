package storage

import (
	"HighArch-dialogs/entity"
	"github.com/gocql/gocql"
)

type DialogStore interface {
	GetDialog(from, to string) ([]entity.DialogMessage, error)
	AddMessage(from, to, text string) error
}

type dialogStore struct {
	session *gocql.Session
}

func NewDialogStore(session *gocql.Session) DialogStore {
	return &dialogStore{
		session: session,
	}
}

func (d dialogStore) GetDialog(from, to string) ([]entity.DialogMessage, error) {
	query := "SELECT id, author_uid, receiver_uid, message, message_time FROM dialogs_space.messages WHERE author_uid = ? AND receiver_uid =?;"
	iter := d.session.Query(query, from, to).Iter()
	var messages []entity.DialogMessage
	var msg entity.DialogMessage
	for iter.Scan(&msg.Id, &msg.From, &msg.To, &msg.Text, &msg.Time) {
		messages = append(messages, msg)
	}
	err := iter.Close()
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (d dialogStore) AddMessage(from, to, text string) error {
	query := "INSERT INTO dialogs_space.messages(id, author_uid, receiver_uid, message, message_time) VALUES(uuid(),?,?,?,toTimestamp(now()))"
	if err := d.session.Query(query, from, to, text).Exec(); err != nil {
		return err
	}
	return nil
}
