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
	liveRepository := repository.NewLiveRepository()
	liveService := service.NewLiveService(liveRepository)
	ctx := context.Background()

	log.Println("Running...")

	var onLives model.OnLives

	for {
		lives, _ := liveService.FindAll(ctx)
		for _, live := range lives {
			live := live
			go func() {
				IsRecording := liveService.IsRecording(ctx, &onLives, live.MemberUsername)
				if !IsRecording {
					os.Mkdir(fmt.Sprintf("download/%s", live.MemberUsername), os.ModePerm)
					DL := hlsdl.NewRecorder(live.StreamUrl, fmt.Sprintf("download/%s", live.MemberUsername))
					filepath, err := DL.Start()
					if err != nil {
						log.Fatal(err)
					}
					log.Println(fmt.Sprintf("%s | %s", live.MemberUsername, filepath))
				}
			}()
		}
		time.Sleep(10 * time.Second)
	}
}
