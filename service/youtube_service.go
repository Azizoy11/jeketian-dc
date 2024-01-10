package service

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"jkt48lab/model"
	"log"
	"os"
	"time"
)

func UploadVideoToYoutube(live model.Live) *youtube.Video {
	godotenv.Load()

	location, _ := time.LoadLocation("Asia/Jakarta")
	startedAt := time.Unix(int64(live.StartedAt), 0).In(location)
	year, month, date := startedAt.Date()
	videoSetting := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title: fmt.Sprintf("%sLive %s [%d-%s-%d]", live.Platform, live.MemberDisplayName, date, month, year),
		},
	}

	ctx := context.Background()

	b, err := os.ReadFile("./client_secret.json")
	if err != nil {
		log.Fatal(err)
	}
	config, err := google.JWTConfigFromJSON(b, youtube.YoutubeUploadScope)
	if err != nil {
		log.Fatal(err)
	}
	client := config.Client(ctx)

	youtubeService, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatal(err)
	}

	call := youtubeService.Videos.Insert([]string{"snipper"}, videoSetting)

	var filename string
	if live.Platform == "IDN" {
		filename = fmt.Sprintf("./download/idn/%s.mp4", live.MemberUsername)
	} else if live.Platform == "Showroom" {
		filename = fmt.Sprintf("./download/sr/%s.mp4", live.MemberUsername)
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	response, err := call.Media(file).Do()
	if err != nil {
		log.Fatal(err)
	}
	return response
}
