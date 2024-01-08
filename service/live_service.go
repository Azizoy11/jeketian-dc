package service

import (
	"context"
	"jkt48lab/model"
	"jkt48lab/repository"
)

type LiveService interface {
	FindAll(ctx context.Context) ([]model.Live, error)
	Find(ctx context.Context, username string) (model.Live, error)
	IsRecording(ctx context.Context, onLives *model.OnLives, username string) bool
}

type LiveServiceImpl struct {
	LiveRepository repository.LiveRepository
}

func (service *LiveServiceImpl) IsRecording(ctx context.Context, onLives *model.OnLives, username string) bool {
	return service.LiveRepository.IsRecording(ctx, onLives, username)
}

func (service *LiveServiceImpl) FindAll(ctx context.Context) ([]model.Live, error) {
	lives, err := service.LiveRepository.FindAll(ctx)
	return lives, err
}

func (service *LiveServiceImpl) Find(ctx context.Context, username string) (model.Live, error) {
	live, err := service.LiveRepository.Find(ctx, username)
	return live, err
}

func NewLiveService(liveRepository repository.LiveRepository) LiveService {
	return &LiveServiceImpl{
		LiveRepository: liveRepository,
	}
}
