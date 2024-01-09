package main

import (
	"context"
	"fmt"
	"github.com/canhlinh/hlsdl"
	"jkt48lab/model"
	"jkt48lab/repository"
	"jkt48lab/service"
	"log"
	"os"
	"time"
)

func main() {
	srLiveRepository := repository.NewSRLiveRepository()
	srLiveService := service.NewSRLiveService(srLiveRepository)

	idnLiveRepository := repository.NewIDNLiveRepository()
	idnLiveService := service.NewIDNLiveService(idnLiveRepository)

	ctx := context.Background()

	log.Println("Running...")

	var onLives model.OnLives

	for {
		srLives, _ := srLiveService.FindAllSR(ctx)
		go func() {
			for _, live := range srLives {
				IsRecording := srLiveService.IsRecordingSR(ctx, &onLives, live.MemberUsername)
				if !IsRecording {
					os.Mkdir(fmt.Sprintf("download/%s", live.MemberUsername), os.ModePerm)
					DL := hlsdl.NewRecorder(live.StreamUrl, fmt.Sprintf("download/%s", live.MemberUsername))
					filepath, err := DL.Start()
					if err != nil {
						log.Fatal(err)
					}
					log.Println(fmt.Sprintf("%s | %s", live.MemberUsername, filepath))
				}
				time.Sleep(10 * time.Second)
			}
		}()

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
