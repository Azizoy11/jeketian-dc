package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"jkt48lab/helper"
	"jkt48lab/model"
	"log"
	"net/http"
	"os"
	"strings"
)

type LiveRepository interface {
	FindAll(ctx context.Context) ([]model.Live, error)
	Find(ctx context.Context, username string) (model.Live, error)
	IsRecording(ctx context.Context, onLives *model.OnLives, username string) bool
}

type LiveRepositoryImpl struct {
}

func (repository *LiveRepositoryImpl) FindAll(ctx context.Context) ([]model.Live, error) {
	jsonLives, err := os.Open("../data/lives.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonLives.Close()

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
			if !strings.Contains(data.RoomUrlKey, "48_") {
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
			if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
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

func (repository *LiveRepositoryImpl) Find(ctx context.Context, username string) (model.Live, error) {
	jsonLives, err := os.Open("../data/lives.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonLives.Close()
	return model.Live{}, nil
}

func (repository *LiveRepositoryImpl) IsRecording(ctx context.Context, onLives *model.OnLives, username string) bool {
	_, err := repository.Find(ctx, username)
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

func NewLiveRepository() LiveRepository {
	return &LiveRepositoryImpl{}
}
