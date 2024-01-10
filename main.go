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
	"os/signal"
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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		fmt.Printf("You pressed ctrl + C. User interrupted infinite loop.")
		os.Exit(0)
	}()

	for {
		srLives, _ := srLiveService.FindAllSR(ctx)
		for _, live := range srLives {
			live := live
			go func() {
				log.Println("[Start Recording]", live.MemberUsername)
				IsRecording := srLiveService.IsRecordingSR(ctx, &onLives, live.MemberUsername)
				if !IsRecording {
					service.SendDiscordNotification(live)
					os.Mkdir("download/sr", os.ModePerm)
					DL := hlsdl.NewRecorder(live.StreamUrl, "download/sr")
					filepath, err := DL.Start(fmt.Sprintf("%s_%d.mp4", live.MemberUsername, time.Now().Unix()))
					if err != nil {
						log.Fatal(err)
					}
					log.Println(fmt.Sprintf("%s | %s", live.MemberUsername, filepath))
				}
			}()
		}

		idnLives, _ := idnLiveService.FindAllIDN(ctx)
		for _, live := range idnLives {
			live := live
			go func() {
				IsRecording, isEnd := idnLiveService.IsRecordingIDN(ctx, &onLives, live.MemberUsername)
				if !IsRecording {
					service.SendDiscordNotification(live)
					log.Println("[Start Recording]", live.MemberUsername)
					os.Mkdir("download/idn", os.ModePerm)
					DL := hlsdl.NewRecorder(live.StreamUrl, "download/idn")
					filepath, err := DL.Start(fmt.Sprintf("%s_%d.mp4", live.MemberUsername, time.Now().Unix()))
					if err != nil {
						log.Println(err)
					}
					log.Println(fmt.Sprintf("%s | %s", live.MemberUsername, filepath))
				}
				if isEnd {
					service.SendDiscordEndNotification(live)
				}
			}()
		}
		time.Sleep(10 * time.Second)
	}
}
