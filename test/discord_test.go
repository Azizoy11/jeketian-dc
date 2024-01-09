package main

import (
	"jkt48lab/model"
	"jkt48lab/service"
	"testing"
)

func TestWebhook(t *testing.T) {
	service.SendDiscordNotification(model.Live{
		MemberUsername:    "agil",
		MemberDisplayName: "",
		Platform:          "IDN",
		Title:             "Asasa",
		StreamUrl:         "BBB",
		Views:             0,
		ImageUrl:          "https://cdn.idntimes.com/content-images/post/20240109/img-20230701-083458-434-66ed4bb0550b59854de64afcf32ee7ac.jpg",
		StartedAt:         0,
	})
}
