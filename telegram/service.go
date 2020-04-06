package telegram

import "context"

type TelegramRepository interface {
	Create(ctx context.Context, userId string, telegramId string) error
	GetIdByUserId(ctx context.Context, userId string) (string, error)
}