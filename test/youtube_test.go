package main

import (
	"github.com/joho/godotenv"
	"jkt48lab/model"
	"jkt48lab/service"
	"testing"
)

func TestUpload(t *testing.T) {
	godotenv.Load()
	service.UploadVideoToYoutube(model.Live{
		MemberUsername:    "jkt48_elin_1704878479",
		MemberDisplayName: "Elin",
		Platform:          "IDN",
		Title:             "Tes Live",
		StreamUrl:         "",
		Views:             0,
		ImageUrl:          "",
		StartedAt:         1704878479,
	})
}
