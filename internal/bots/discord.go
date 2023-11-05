package bots

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
)

type DiscordBot struct {
	*discordgo.Session
}

func NewDiscordBot(token string) *DiscordBot {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil
	}
	return &DiscordBot{session}
}

func StartDiscordBot() {
	bot := NewDiscordBot(os.Getenv("DISCORD_BOT_TOKEN"))

	bot.AddHandler(onCommandHandler)
	bot.Identify.Intents = discordgo.IntentsGuildMessages
	err := bot.Open()
	if err != nil {
		panic(err)
	}
}

func onCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		commandName := i.Data

		fmt.Println(commandName)

		//switch commandName {
		//case "ping":
		//	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		//		Type: discordgo.InteractionResponseChannelMessageWithSource,
		//		Data: &discordgo.InteractionResponseData{
		//			Content: "Pong!",
		//		},
		//	}); err != nil {
		//		panic(err)
		//	}
		//}
	}
}
