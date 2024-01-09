package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"jkt48lab/helper"
	"jkt48lab/model"
	"log"
	"net/http"
	"strings"
)

type SRLiveRepository interface {
	FindAllSR(ctx context.Context) ([]model.Live, error)
	FindSR(ctx context.Context, username string) (model.Live, error)
	IsRecordingSR(ctx context.Context, onLives *model.OnLives, username string) bool
}

type SRLiveRepositoryImpl struct {
}

func (repository *SRLiveRepositoryImpl) FindAllSR(ctx context.Context) ([]model.Live, error) {
	resp, err := helper.Fetch("https://www.showroom-live.com/api/live/onlives")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var result model.LiveShowroomResponses
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		log.Println(string(body))
		log.Println("Gagal mengubah JSON ke LiveShowroomResponses")
	}

	var lives []model.Live
	if len(result.OnLives) > 0 {
		for _, data := range result.OnLives[0].Lives {
			if data.PremiumRoomType == 1 {
				continue
			}
			if !strings.Contains(data.RoomUrlKey, "JKT48") {
				continue
			}
			if data.RoomId == 0 {
				continue
			}

			resp, err := http.Get(fmt.Sprintf("https://www.showroom-live.com/api/live/streaming_url?abr_available=1&room_id=%d", data.RoomId))
			if err != nil {
				log.Fatal(err)
			}
			body, err := io.ReadAll(resp.Body)

			var result model.LiveShowroomStreamingUrlResponses
			if err := json.Unmarshal(body, &result); err != nil {
				log.Println("Gagal mengubah JSON ke LiveShowroomStreamingUrlResponses")
			}
			resp.Body.Close()

			if len(result.StreamingUrlList) > 0 {
				live := model.Live{
					MemberUsername:    data.RoomUrlKey,
					MemberDisplayName: data.MainName,
					Platform:          "Showroom",
					Title:             fmt.Sprintf("%s Live", data.MainName),
					StreamUrl:         result.StreamingUrlList[1].Url,
					Views:             data.ViewNum,
					StartedAt:         data.StartedAt,
				}
				lives = append(lives, live)
			}

		}
	}
	return lives, nil
}

func (repository *SRLiveRepositoryImpl) FindSR(ctx context.Context, username string) (model.Live, error) {
	var live model.Live
	lives, _ := repository.FindAllSR(ctx)
	for _, l := range lives {
		if l.MemberUsername == username {
			live = l
			return live, nil
		}
	}
	return live, errors.New(fmt.Sprintf("%s sedang tidak live"))
}

func (repository *SRLiveRepositoryImpl) IsRecordingSR(ctx context.Context, onLives *model.OnLives, username string) bool {
	_, err := repository.FindSR(ctx, username)
	if err != nil {
		// Member sedang tidak live
		if helper.Contains(onLives.MemberOnLives, username) {
			onLives.MemberOnLives = helper.RemoveStringFromSlice(onLives.MemberOnLives, username)
			return false
		}
	}
	if !helper.Contains(onLives.MemberOnLives, username) {
		onLives.MemberOnLives = append(onLives.MemberOnLives, username)
		return false
	} else {
		return true
	}
}

func NewSRLiveRepository() SRLiveRepository {
	return &SRLiveRepositoryImpl{}
}
