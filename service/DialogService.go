package service

import (
	"HighArch-dialogs/api"
	"HighArch-dialogs/storage"
)

type DialogService struct {
	dialogStore storage.DialogStore
}

func NewDialogService(dialogStore storage.DialogStore) *DialogService {
	return &DialogService{
		dialogStore: dialogStore,
	}
}

func (s *DialogService) GetDialog(from, to string) ([]api.DialogMessage, error) {
	if len(from) <= 0 || len(to) <= 0 {
		return nil, ErrorValidation
	}
	messages, err := s.dialogStore.GetDialog(from, to)
	if err != nil {
		return nil, ErrorStoreError
	}
	var result []api.DialogMessage
	for _, message := range messages {
		apiMsg := api.DialogMessage{
			Id:   message.Id,
			From: message.From,
			To:   message.To,
			Text: message.Text,
			Time: message.Time,
		}
		result = append(result, apiMsg)
	}
	return result, nil
}

func (s *DialogService) AddDialogMessage(from string, to string, text string) error {
	if len(from) <= 0 || len(to) <= 0 || len(text) <= 0 {
		return ErrorValidation
	}
	err := s.dialogStore.AddMessage(from, to, text)
	if err != nil {
		return ErrorStoreError
	}
	return nil
}
