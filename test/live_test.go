package main

import (
	"context"
	"fmt"
	"github.com/canhlinh/hlsdl"
	"io"
	"jkt48lab/model"
	"jkt48lab/repository"
	"jkt48lab/service"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestRequest(t *testing.T) {
	resp, _ := http.Get("https://www.showroom-live.com/api/live/onlives")
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	log.Println(string(body))
}

func TestIDN(t *testing.T) {
	idnLiveRepository := repository.NewIDNLiveRepository()
	idnLiveService := service.NewIDNLiveService(idnLiveRepository)
	ctx := context.Background()

	var onLives model.OnLives

	for {
		idnLives, _ := idnLiveService.FindAllIDN(ctx)
		go func() {
			for _, live := range idnLives {
				log.Println(live)
				IsRecording := idnLiveService.IsRecordingIDN(ctx, &onLives, live.MemberUsername)
				if !IsRecording {
					os.Mkdir(fmt.Sprintf("download/%s", live.MemberUsername), os.ModePerm)
					DL := hlsdl.NewRecorder(live.StreamUrl, fmt.Sprintf("download/%s", live.MemberUsername))
					filepath, err := DL.Start()
					if err != nil {
						log.Println(err)
					}
					log.Println(fmt.Sprintf("%s | %s", live.MemberUsername, filepath))
				}
				time.Sleep(10 * time.Second)
			}
		}()
	}
}
