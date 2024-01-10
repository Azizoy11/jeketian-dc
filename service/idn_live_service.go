package service

import (
	"context"
	"jkt48lab/model"
	"jkt48lab/repository"
)

type IDNLiveService interface {
	FindAllIDN(ctx context.Context) ([]model.Live, error)
	FindIDN(ctx context.Context, username string) (model.Live, error)
	IsRecordingIDN(ctx context.Context, onLives *model.OnLives, username string) (bool, bool)
}

type IDNLiveServiceImpl struct {
	IDNLiveRepository repository.IDNLiveRepository
}

func (service *IDNLiveServiceImpl) IsRecordingIDN(ctx context.Context, onLives *model.OnLives, username string) (bool, bool) {
	return service.IDNLiveRepository.IsRecordingIDN(ctx, onLives, username)
}

func (service *IDNLiveServiceImpl) FindAllIDN(ctx context.Context) ([]model.Live, error) {
	lives, err := service.IDNLiveRepository.FindAllIDN(ctx)
	return lives, err
}

func (service *IDNLiveServiceImpl) FindIDN(ctx context.Context, username string) (model.Live, error) {
	live, err := service.IDNLiveRepository.FindIDN(ctx, username)
	return live, err
}

func NewIDNLiveService(liveRepository repository.IDNLiveRepository) IDNLiveService {
	return &IDNLiveServiceImpl{
		IDNLiveRepository: liveRepository,
	}
}
