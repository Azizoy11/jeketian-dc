package main

import (
	"context"
	"fmt"
	"jkt48lab/hlsdl"
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
				log.Println("[Start Recording]", live.MemberUsername)
				IsRecording := srLiveService.IsRecordingSR(ctx, &onLives, live.MemberUsername)
				if !IsRecording {
					os.Mkdir("download/sr", os.ModePerm)
					DL := hlsdl.NewRecorder(live.StreamUrl, "download/sr")
					filepath, err := DL.Start(fmt.Sprintf("%s.mp4", live.MemberUsername))
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
				IsRecording := idnLiveService.IsRecordingIDN(ctx, &onLives, live.MemberUsername)
				if !IsRecording {
					log.Println("[Start Recording]", live.MemberUsername)
					os.Mkdir("download/idn", os.ModePerm)
					DL := hlsdl.NewRecorder(live.StreamUrl, "download/idn")
					filepath, err := DL.Start(fmt.Sprintf("%s.mp4", live.MemberUsername))
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
