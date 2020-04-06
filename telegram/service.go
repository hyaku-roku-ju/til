package telegram

import "context"

type TelegramRepository interface {
	Create(ctx context.Context, userId string, telegramId string) error
	GetIdByUserId(ctx context.Context, userId string) (string, error)
}

type TelegramService struct {
	Repo TelegramRepository
}

func NewService(repo TelegramRepository) TelegramService {
	return TelegramService{repo}
}

func (t *TelegramService) StoreTelegramId(ctx context.Context, userId string, telegramId string) error {
	err := t.Repo.Create(ctx, userId, telegramId)
	return err
}

func (t *TelegramService) GetIdByUserId(ctx context.Context, userId string) (string, error) {
	id, err := t.Repo.GetIdByUserId(ctx, userId)
	return id, err
}
