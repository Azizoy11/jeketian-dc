package service

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"io"
	"jkt48lab/model"
	"log"
	"os"
)

type Webhook struct {
	Id    string `json:"id"`
	Token string `json:"token"`
}

func SendDiscordNotification(live model.Live) {
	godotenv.Load()
	bot, err := discordgo.New(os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	jsonFile, _ := os.Open("data/notification_webhooks.json")
	defer jsonFile.Close()

	webhooksJson, _ := io.ReadAll(jsonFile)

	var webhooks []Webhook
	json.Unmarshal(webhooksJson, &webhooks)

	var liveUrl string
	if live.Platform == "IDN" {
		liveUrl = fmt.Sprintf("https://jkt48.safatanc.com/live/idn/%s", live.MemberUsername)
	} else if live.Platform == "Showroom" {
		liveUrl = fmt.Sprintf("https://jkt48.safatanc.com/live/sr/%s", live.MemberUsername)
	}

	var liveUrlOriginal string
	if live.Platform == "IDN" {
		liveUrlOriginal = fmt.Sprintf("https://idn.app/%s", live.MemberUsername)
	} else if live.Platform == "Showroom" {
		liveUrlOriginal = fmt.Sprintf("https://showroom-live.com/r/%s", live.MemberUsername)
	}

	for _, webhook := range webhooks {
		embed := &discordgo.MessageEmbed{
			URL:         "https://jkt48.safatanc.com",
			Title:       "Notifikasi Live",
			Description: fmt.Sprintf("**%s** sedang live di `%s`", live.MemberUsername, live.Platform),
			Color:       0xccec1c,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "JKT48Lab by safatanc.com",
			},
			Image: nil,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: live.ImageUrl,
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Tonton di JKT48Lab",
					Value:  fmt.Sprintf("[Klik disini](%s)", liveUrl),
					Inline: true,
				},
				{
					Name:   fmt.Sprintf("Tonton di %s", live.Platform),
					Value:  fmt.Sprintf("[Klik disini](%s)", liveUrlOriginal),
					Inline: true,
				},
			},
		}
		_, err = bot.WebhookExecute(webhook.Id, webhook.Token, false, &discordgo.WebhookParams{
			Username: "JKT48Lab",
			Embeds: []*discordgo.MessageEmbed{
				embed,
			},
		})
		if err != nil {
			log.Fatal(err)
		}
	}

}
