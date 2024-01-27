package bots

import (
	"CongratulatorBot/internal/model"
	"CongratulatorBot/internal/storage"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"os"
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
			// DefaultMemberPermissions: &adminMemberPermissions,
		},
		{
			Name:        "gz",
			Description: "Поздравление с праздником",
			// DefaultMemberPermissions: &adminMemberPermissions,
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
			pathToPostcard, err := storage.GetRandomPostcardPath()
			if err != nil {
				panic(err)
			}

			file, err := os.Open(pathToPostcard)
			if err != nil {
				panic(err)
			}

			discordFile := discordgo.File{
				ContentType: "text/plain",
				Name:        pathToPostcard,
				Reader:      file,
			}

			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Files: []*discordgo.File{
						&discordFile,
					},
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

func StartDiscordBot(token string, errorChan chan<- error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	initCommands(s, errorChan)
	startInfinitePolling(s, errorChan)
}

func startInfinitePolling(s *discordgo.Session, errorChan chan<- error) {
	if err := s.Open(); err != nil {
		errorChan <- err
	}
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
