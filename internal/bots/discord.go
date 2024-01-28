package bots

import (
	"CongratulatorBot/internal/model"
	"CongratulatorBot/internal/storage"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"strconv"
)

var (
	pathToConfig                 = "internal/config/discordConfig.json"
	config                       = parseConfig(pathToConfig) // TODO multilingual
	appID                        = "1170807660279505076"
	guildID                      = "1200170010111393792"
	adminMemberPermissions int64 = discordgo.PermissionAdministrator

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "help",
			Description: "Вывод списка команд и их описания",
		},
		{
			Name:        "gz",
			Description: "Поздравление с праздником",
		},
		{
			Name:                     "settime",
			Description:              "Установить время рассылки",
			DefaultMemberPermissions: &adminMemberPermissions,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "time",
					Description: "Время рассылки",
					Required:    true,
				},
			},
		},
		{
			Name:                     "removemailing",
			Description:              "Убрать рассылку",
			DefaultMemberPermissions: &adminMemberPermissions,
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"help": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: config.Localization["ru"].HelpMessage,
				},
			}); err != nil {
				panic(err)
			}
		},
		"gz": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: config.Localization["ru"].LoadingMessage,
				},
			}); err != nil {
				panic(err)
			}

			pathToPostcard, err := storage.GetRandomPostcardPath()
			if err != nil {
				panic(err)
			}

			if pathToPostcard == "" {
				noHolidaysMessage := config.Localization["ru"].NoHolidaysMessage
				if _, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &noHolidaysMessage,
				}); err != nil {
					panic(err)
				}
				return
			}

			file, err := os.Open(pathToPostcard)
			if err != nil {
				panic(err)
			}

			var clearMessage string

			discordFile := discordgo.File{
				ContentType: "text/plain",
				Name:        pathToPostcard,
				Reader:      file,
			}
			if _, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &clearMessage,
				Files: []*discordgo.File{
					&discordFile,
				},
			}); err != nil {
				panic(err)
			}
		},
		"settime": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			time := i.ApplicationCommandData().Options[0].StringValue()

			offsetFromUTC, err := strconv.Atoi(os.Getenv("TIME_OFFSET"))
			if err != nil {
				panic(err)
			}

			intChannelID, err := strconv.Atoi(i.ChannelID)
			int64ChannelID := int64(intChannelID)
			if err != nil {
				panic(err)
			}

			serverIsMailing, err := storage.GetIfUserIsMailing(int64ChannelID, true)
			if err != nil {
				panic(err)
			}

			if valid, _, secondsSinceMidnight := timeIsValid(time); valid {
				time_ := int64(secondsSinceMidnight - offsetFromUTC*3600)
				if serverIsMailing {
					if err := storage.RemoveUserFromMailing(int64ChannelID, true); err != nil {
						panic(err)
					}
				}
				if err := storage.AddUserToMailing(int64ChannelID, time_, true); err != nil {
					panic(err)
				}

				if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: config.Localization["ru"].MailingEnableMessage,
					},
				}); err != nil {
					panic(err)
				}
			} else {
				if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: config.Localization["ru"].BadTimeMessage,
					},
				}); err != nil {
					panic(err)
				}

			}
		},
		"removemailing": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			intGuildID, err := strconv.Atoi(i.GuildID)
			if err != nil {
				panic(err)
			}

			int64GuildID := int64(intGuildID)

			serverIsMailing, err := storage.GetIfUserIsMailing(int64GuildID, true)
			if err != nil {
				panic(err)
			}

			if serverIsMailing {
				if err := storage.RemoveUserFromMailing(int64GuildID, true); err != nil {
					panic(err)
				}
			}

			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: config.Localization["ru"].MailingDisableMessage,
				},
			}); err != nil {
				panic(err)
			}
		},
	}
)

func parseConfig(path string) *model.DiscordConfig {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	var cfg model.DiscordConfig
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&cfg); err != nil {
		panic(err)
	}
	return &cfg
}

func sendDiscordPostcard(s *discordgo.Session, channelID string) {
	pathToPostcard, err := storage.GetRandomPostcardPath()
	if err != nil {
		panic(err)
	}

	file, err := os.Open(pathToPostcard)
	if err != nil {
		panic(err)
	}

	if _, err = s.ChannelFileSend(channelID, pathToPostcard, file); err != nil {
		panic(err)
	}
}

func NewDiscordBot(token string) (*discordgo.Session, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	return s, nil
}

func StartDiscordBot(s *discordgo.Session, errorChan chan<- error) {
	initCommands(s, errorChan)
	startInfinitePolling(s, errorChan)
}

func startInfinitePolling(s *discordgo.Session, errorChan chan<- error) {
	if err := s.Open(); err != nil {
		errorChan <- err
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	fmt.Println("Press Ctrl+C to exit")
	<-stop
	removeCommands(s)
}

func removeCommands(s *discordgo.Session) {
	fmt.Println("Removing commands...")
	registeredCommands, err := s.ApplicationCommands(s.State.User.ID, guildID)
	if err != nil {
		fmt.Printf("Could not fetch registered commands: %v", err)
	}
	for _, v := range registeredCommands {
		err = s.ApplicationCommandDelete(s.State.User.ID, guildID, v.ID)
		if err != nil {
			fmt.Printf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	fmt.Println("Gracefully shutting down.")
}

func initCommands(s *discordgo.Session, errorChan chan<- error) {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(appID, guildID, v)
		if err != nil {
			errorChan <- err
		}
		registeredCommands[i] = cmd
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

}
