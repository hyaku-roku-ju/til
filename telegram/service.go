package telegram

import "context"

type TelegramRepository interface {
	Create(ctx context.Context, userId string, telegramId string) error
	GetTelegramIdByUserId(ctx context.Context, userId string) (string, error)
	GetUserIdByTelegramId(ctx context.Context, telegramId string) (string, error)
}

// Maybe better name? Used as internal dependency of handler, 
// But it has only one method, such that we won't need to implement everything
// for telegram handler.
type TelegramServiceEntrypoint interface {
  ProcessDecodedMessage(ctx context.Context, update Update) error
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

func (t *TelegramService) GetTelegramIdByUserId(ctx context.Context, userId string) (string, error) {
	id, err := t.Repo.GetTelegramIdByUserId(ctx, userId)
	return id, err
}

func (t *TelegramService) GetUserIdByTelegramId(ctx context.Context, telegramId string) (string, error) {
	userId, err := t.Repo.GetUserIdByTelegramId(ctx, telegramId)
	return userId, err
}
