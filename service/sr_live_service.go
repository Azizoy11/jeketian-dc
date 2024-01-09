package service

import (
	"context"
	"jkt48lab/model"
	"jkt48lab/repository"
)

type SRLiveService interface {
	FindAllSR(ctx context.Context) ([]model.Live, error)
	FindSR(ctx context.Context, username string) (model.Live, error)
	IsRecordingSR(ctx context.Context, onLives *model.OnLives, username string) bool
}

type SRLiveServiceImpl struct {
	SRLiveRepository repository.SRLiveRepository
}

func (service *SRLiveServiceImpl) IsRecordingSR(ctx context.Context, onLives *model.OnLives, username string) bool {
	return service.SRLiveRepository.IsRecordingSR(ctx, onLives, username)
}

func (service *SRLiveServiceImpl) FindAllSR(ctx context.Context) ([]model.Live, error) {
	lives, err := service.SRLiveRepository.FindAllSR(ctx)
	return lives, err
}

func (service *SRLiveServiceImpl) FindSR(ctx context.Context, username string) (model.Live, error) {
	live, err := service.SRLiveRepository.FindSR(ctx, username)
	return live, err
}

func NewSRLiveService(liveRepository repository.SRLiveRepository) SRLiveService {
	return &SRLiveServiceImpl{
		SRLiveRepository: liveRepository,
	}
}
