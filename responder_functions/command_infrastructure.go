package responder_functions

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/medzernik/SlovakiaDiscordBotGo/config"
	"log"
	"strings"
	"time"
)

func RegisterPlugin(s *discordgo.Session) {
	s.AddHandler(Ready)

}

// Bot parameters

var (
	BotToken = flag.String("token", config.Cfg.ServerInfo.ServerToken, "Bot access token")
	Cleanup  = flag.Bool("cleanup", true, "Cleanup of commands")
)

func init() {
	flag.Parse()
}

/*
func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

*/

var (
	BotCommands = []*discordgo.ApplicationCommand{
		{
			Name: "slovakia",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: ":)",
		},
		{
			Name:        "age",
			Description: "lists the user account age",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user-mention",
					Description: "user-mention",
					Required:    true,
				},
			},
		},
		{
			Name:        "mute",
			Description: "mutes a user (adds a mute role)",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user-mention",
					Description: "user-mention",
					Required:    true,
				},
			},
		},
		{
			Name: "Mute User",
			Type: discordgo.UserApplicationCommand,
		},
		{
			Name: "Purge To Here",
			Type: discordgo.MessageApplicationCommand,
		},
		{
			Name: "Purge To Here User",
			Type: discordgo.MessageApplicationCommand,
		},
		{
			Name:        "unmute",
			Description: "umutes a user (removes a mute role)",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user-mention",
					Description: "user-mention",
					Required:    true,
				},
			},
		},
		{
			Name:        "kick",
			Description: "kicks a user (optional reason)",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user-mention",
					Description: "user-mention",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "dovod",
					Description: "Reason for kick",
					Required:    false,
				},
			},
		},
		{
			Name:        "ban",
			Description: "ban a user (optional reason)",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user-mention",
					Description: "user-mention",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "dovod",
					Description: "Reason for ban",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "delete",
					Description: "How many days back should we delete messages? (default 7)",
					Required:    false,
				},
			},
		},
		{
			Name:        "check-users",
			Description: "checks the users who joined less than 24 hours ago (default, optional custom time frame in hours)",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "timeframe",
					Description: "How far back should we look for new users?",
					Required:    false,
				},
			},
		},
		{
			Name:        "planned",
			Description: "Outputs the currently planned games",
		},
		{
			Name:        "plan",
			Description: "Plans a game to play",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "time",
					Description: "Time to plan the new game at (HH:MM)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "gamename",
					Description: "Name of the game to play",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "Who to ping when the time to play is nigh",
					Required:    true,
				},
			},
		},
		{
			Name:        "topic",
			Description: "A random topic for discussion",
		},
		{
			Name:        "basic-command-with-files",
			Description: "Basic command with files",
		},
		{
			Name:        "weather",
			Description: "Check the weather in your city",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "city",
					Description: "City name",
					Required:    true,
				},
			},
		},
		{
			Name:        "purge",
			Description: "Purges 1-100 messages in the current channel",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "amount",
					Description: "Messages to purge 1-100",
					Required:    true,
				},
			},
		},
		{
			Name:        "members",
			Description: "Shows how many members there are on the server",
		},
		{
			Name:        "covid-vaccines-available",
			Description: "Get the current covid vaccines in slovakia",
		},
		{
			Name:        "covid-number-vaccinated",
			Description: "Checks how many people were vaccinated in the SVK to date.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "days",
					Description: "Days to check for",
					Required:    true,
				},
			},
		},
		{
			Name:        "prune-count",
			Description: "Checks how many members to prune with 7-30 days of inactivity",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "days",
					Description: "Days of inactivity (7-30)",
					Required:    true,
				},
			},
		},
		{
			Name:        "prune-members",
			Description: "Prunes members with 7-30 days of inactivity",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "days",
					Description: "Days of inactivity (7-30)",
					Required:    true,
				},
			},
		},
		{
			Name:        "setroleperm",
			Description: "Sets the current channel permissions by UID bits for a role",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "allow",
					Description: "Allow or deny the permissions",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role",
					Description: "Role to set the permissions to",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "permissions",
					Description: "Permission ID",
					Required:    true,
				},
			},
		},
		{
			Name:        "setuserperm",
			Description: "Sets the current channel permissions by UID bits for a user",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "allow",
					Description: "Allow or deny the permissions",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User to set the permissions to",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "permissions",
					Description: "Permission ID",
					Required:    true,
				},
			},
		},
		{
			Name:        "redirect",
			Description: "Redirects the discussion by setting a slowmode in current channel",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "Channel to redirect to",
					Required:    true,
				},
			},
		},
		{
			Name:        "slow",
			Description: "Sets a slowmode in the current channel",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "duration",
					Description: "Duration in seconds",
					Required:    true,
				},
			},
		},
		{
			Name:        "voicechannelmodify",
			Description: "lists the user account age",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "channel-name",
					Description: "Voice channel name",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "bitrate",
					Description: "Specify bitrate (in kbps)",
					Required:    false,
				},
			},
		},

		//EXAMPLES BELOW
		{
			Name:        "options",
			Description: "Command for demonstrating options",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "string-option",
					Description: "String option",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "integer-option",
					Description: "Integer option",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "bool-option",
					Description: "Boolean option",
					Required:    true,
				},

				// Required options must be listed first since optional parameters
				// always come after when they're used.
				// The same concept applies to Discord's Slash-commands API

				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel-option",
					Description: "Channel option",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user-option",
					Description: "User option",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role-option",
					Description: "Role option",
					Required:    false,
				},
			},
		},
		{
			Name:        "subcommands",
			Description: "Subcommands and command groups example",
			Options: []*discordgo.ApplicationCommandOption{
				// When a command has subcommands/subcommand groups
				// It must not have top-level options, they aren't accesible in the UI
				// in this case (at least not yet), so if a command has
				// subcommands/subcommand any groups registering top-level options
				// will cause the registration of the command to fail

				{
					Name:        "scmd-grp",
					Description: "Subcommands group",
					Options: []*discordgo.ApplicationCommandOption{
						// Also, subcommand groups aren't capable of
						// containing options, by the name of them, you can see
						// they can only contain subcommands
						{
							Name:        "nst-subcmd",
							Description: "Nested subcommand",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
						},
					},
					Type: discordgo.ApplicationCommandOptionSubCommandGroup,
				},
				// Also, you can create both subcommand groups and subcommands
				// in the command at the same time. But, there's some limits to
				// nesting, count of subcommands (top level and nested) and options.
				// Read the intro of slash-commands docs on Discord dev portal
				// to get more information
				{
					Name:        "subcmd",
					Description: "Top-level subcommand",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
		{
			Name:        "responses",
			Description: "Interaction responses testing initiative",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "resp-type",
					Description: "VaccinatedSlovakiaResponse type",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Channel message with source",
							Value: 4,
						},
						{
							Name:  "Deferred response With Source",
							Value: 5,
						},
					},
					Required: true,
				},
			},
		},
		{
			Name:        "followups",
			Description: "Followup messages",
		},
	}
	//Engaging the command handlers.
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		//This command just runs a basic test
		"slovakia": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			FoxTest(s, i)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "SPUTNIK V 150% UCINNOST GDE KOTLEBA BGATIA OMG GDE HRAZDOVE RUKY",
				},
			})
			return
		},
		//This command runs the AgeJoinedCMD function
		"age": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⠀",
				},
			})
			argumentArray := []interface{}{
				i.ApplicationCommandData().Options[0].UserValue(s).ID,
			}
			go AgeJoinedCMD(s, i, argumentArray)
			return
		},
		//This command runs the AgeJoinedCMD function
		"mute": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⠀",
				},
			})
			argumentArray := []interface{}{
				i.ApplicationCommandData().Options[0].UserValue(s).ID,
			}
			go MuteCMD(s, i, argumentArray)
			return
		},
		//This command runs the AgeJoinedCMD function
		"Mute User": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⠀",
				},
			})

			argumentArray := []interface{}{
				i.ApplicationCommandData().TargetID,
			}

			go MuteCMD(s, i, argumentArray)
			return
		},

		"unmute": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⠀",
				},
			})
			argumentArray := []interface{}{
				i.ApplicationCommandData().Options[0].UserValue(s).ID,
			}
			go UnmuteCMD(s, i, argumentArray)
			return
		},
		"kick": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⠀",
				},
			})
			argumentArray := []interface{}{
				i.ApplicationCommandData().Options[0].UserValue(s).ID,
			}
			if len(i.ApplicationCommandData().Options) > 1 {
				argumentArray = append(argumentArray, i.ApplicationCommandData().Options[1].StringValue())
			}
			go KickUserCMD(s, i, argumentArray)
			return
		},
		"ban": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⠀",
				},
			})
			argumentArray := []interface{}{
				i.ApplicationCommandData().Options[0].UserValue(s).ID,
			}
			if len(i.ApplicationCommandData().Options) > 1 {
				argumentArray = append(argumentArray, i.ApplicationCommandData().Options[1].StringValue())

			}
			if len(i.ApplicationCommandData().Options) > 2 {
				argumentArray = append(argumentArray, i.ApplicationCommandData().Options[2].UintValue())
			}

			go BanUserCMD(s, i, argumentArray)
			return
		},
		"check-users": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⠀",
				},
			})
			var argumentArray []interface{}
			if len(i.ApplicationCommandData().Options) > 0 {
				argumentArray = append(argumentArray, i.ApplicationCommandData().Options[0].IntValue())
			}
			go CheckUsersCMD(s, i, argumentArray)
			return
		},
		"planned": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⠀",
				},
			})
			var argumentArray []interface{}

			go PlannedGamesCMD(s, i, argumentArray)
			return
		},
		"plan": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⠀",
				},
			})
			var argumentArray []interface{}

			argumentArray = []interface{}{
				i.ApplicationCommandData().Options[0].StringValue(),
				i.ApplicationCommandData().Options[1].StringValue(),
				i.ApplicationCommandData().Options[2].UserValue(s).Mention(),
			}

			go PlanGameCMD(s, i, argumentArray)
			return
		},
		"topic": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⠀",
				},
			})
			var argumentArray []interface{}

			argumentArray = []interface{}{}

			go TopicCMD(s, i, argumentArray)
			return
		},
		"weather": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⠀",
				},
			})
			var argumentArray []interface{}

			argumentArray = []interface{}{
				i.ApplicationCommandData().Options[0].StringValue(),
			}

			go GetWeatherCMD(s, i, argumentArray)
			return
		},
		"purge": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⠀",
				},
			})
			var argumentArray []interface{}

			argumentArray = []interface{}{
				i.ApplicationCommandData().Options[0].IntValue(),
			}

			go PurgeMessagesCMD(s, i, argumentArray)
			return
		},
		"Purge To Here": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⠀",
				},
			})
			argumentArray := []interface{}{
				i.ApplicationCommandData().TargetID,
			}

			go PurgeMessagesCMDMessage(s, i, argumentArray)
			return
		},
		"Purge To Here User": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⠀",
				},
			})
			argumentArray := []interface{}{
				i.ApplicationCommandData().TargetID,
			}

			go PurgeMessagesCMDMessageOnlyAuthor(s, i, argumentArray)
			return
		},
		"members": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⠀",
				},
			})
			var argumentArray []interface{}

			argumentArray = []interface{}{}

			go MembersCMD(s, i, argumentArray)
			return
		},
		"covid-vaccines-available": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⠀",
				},
			})

			go COVIDVaccinesAvailable(s, i)
			return
		},
		"covid-number-vaccinated": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⠀",
				},
			})
			var argumentArray []interface{}

			/*
				argumentArray = []interface{}{
					i.ApplicationCommandData().Options[0].IntValue(),
				}

			*/

			go COVIDNumberOfVaccinated(s, i, argumentArray)
			return
		},
		"prune-count": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⠀",
				},
			})
			var argumentArray []interface{}

			argumentArray = []interface{}{
				i.ApplicationCommandData().Options[0].IntValue(),
			}

			go PruneCountCMD(s, i, argumentArray)
			return
		},
		"prune-members": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⠀",
				},
			})
			var argumentArray []interface{}

			argumentArray = []interface{}{
				i.ApplicationCommandData().Options[0].IntValue(),
			}

			go PruneMembersCMD(s, i, argumentArray)
			return

		},
		"setroleperm": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⠀",
				},
			})
			var argumentArray []interface{}

			argumentArray = []interface{}{
				i.ApplicationCommandData().Options[0].BoolValue(),
				i.ApplicationCommandData().Options[1].RoleValue(s, i.GuildID).Mention(),
				i.ApplicationCommandData().Options[2].IntValue(),
			}

			go SetRoleChannelPermCMD(s, i, argumentArray)
			return

		},
		"setuserperm": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⠀",
				},
			})
			var argumentArray []interface{}

			argumentArray = []interface{}{
				i.ApplicationCommandData().Options[0].BoolValue(),
				i.ApplicationCommandData().Options[1].UserValue(s).Mention(),
				i.ApplicationCommandData().Options[2].IntValue(),
			}

			go SetUserChannelPermCMD(s, i, argumentArray)
			return
		},
		"redirect": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⠀",
				},
			})
			var argumentArray []interface{}

			argumentArray = []interface{}{
				i.ApplicationCommandData().Options[0].ChannelValue(s).Mention(),
			}

			go RedirectDiscussionCMD(s, i, argumentArray)
			return
		},
		"slow": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⠀",
				},
			})
			var argumentArray []interface{}

			argumentArray = []interface{}{
				i.ApplicationCommandData().Options[0].UintValue(),
			}

			go SlowModeChannelCMD(s, i, argumentArray)
			return
		},
		"voicechannelmodify": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⠀",
				},
			})
			var argumentArray []interface{}

			argumentArray = []interface{}{
				i.ApplicationCommandData().Options[0].StringValue(),
			}
			if len(i.ApplicationCommandData().Options) > 1 {
				argumentArray = append(argumentArray, i.ApplicationCommandData().Options[1].UintValue())
			}

			go ChangeVoiceChannelCurrentCMD(s, i, argumentArray)
			return
		},

		//BELOW THIS STARTS THE EXAMPLE FILE
		"basic-command-with-files": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey there! Congratulations, you just executed your first slash command with a file in the response",
					Files: []*discordgo.File{
						{
							ContentType: "text/plain",
							Name:        "test.txt",
							Reader:      strings.NewReader("Hello Discord!!"),
						},
					},
				},
			})
		},

		"options": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			margs := []interface{}{
				// Here we need to convert raw interface{} value to wanted type.
				// Also, as you can see, here is used utility functions to convert the value
				// to particular type. Yeah, you can use just switch type,
				// but this is much simpler
				i.ApplicationCommandData().Options[0].StringValue(),
				i.ApplicationCommandData().Options[1].IntValue(),
				i.ApplicationCommandData().Options[2].BoolValue(),
			}
			msgformat :=
				` Now you just learned how to use command options. Take a look to the value of which you've just entered:
				> string_option: %s
				> integer_option: %d
				> bool_option: %v
`
			if len(i.ApplicationCommandData().Options) >= 4 {
				margs = append(margs, i.ApplicationCommandData().Options[3].ChannelValue(nil).ID)
				msgformat += "> channel-option: <#%s>\n"
			}
			if len(i.ApplicationCommandData().Options) >= 5 {
				margs = append(margs, i.ApplicationCommandData().Options[4].UserValue(nil).ID)
				msgformat += "> user-option: <@%s>\n"
			}
			if len(i.ApplicationCommandData().Options) >= 6 {
				margs = append(margs, i.ApplicationCommandData().Options[5].RoleValue(nil, "").ID)
				msgformat += "> role-option: <@&%s>\n"
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Ignore type for now, we'll discuss them in "responses" part
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						msgformat,
						margs...,
					),
				},
			})
		},
		"subcommands": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			content := ""

			// As you can see, the name of subcommand (nested, top-level) or subcommand group
			// is provided through arguments.
			switch i.ApplicationCommandData().Options[0].Name {
			case "subcmd":
				content =
					"The top-level subcommand is executed. Now try to execute the nested one."
			default:
				if i.ApplicationCommandData().Options[0].Name != "scmd-grp" {
					return
				}
				switch i.ApplicationCommandData().Options[0].Options[0].Name {
				case "nst-subcmd":
					content = "Nice, now you know how to execute nested commands too"
				default:
					// I added this in the case something might go wrong
					content = "Oops, something gone wrong.\n" +
						"Hol' up, you aren't supposed to see this message."
				}
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
		},
		"responses": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Responses to a command are very important.
			// First of all, because you need to react to the interaction
			// by sending the response in 3 seconds after receiving, otherwise
			// interaction will be considered invalid and you can no longer
			// use the interaction token and ID for responding to the user's request

			content := ""
			// As you can see, the response type names used here are pretty self-explanatory,
			// but for those who want more information see the official documentation
			switch i.ApplicationCommandData().Options[0].IntValue() {
			case int64(discordgo.InteractionResponseChannelMessageWithSource):
				content =
					"You just responded to an interaction, sent a message and showed the original one. " +
						"Congratulations!"
				content +=
					"\nAlso... you can edit your response, wait 5 seconds and this message will be changed"
			default:
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
				})
				if err != nil {
					s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
				}
				return
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
			if err != nil {
				s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something went wrong",
				})
				return
			}
			time.AfterFunc(time.Second*5, func() {
				_, err = s.InteractionResponseEdit(s.State.User.ID, i.Interaction, &discordgo.WebhookEdit{
					Content: content + "\n\nWell, now you know how to create and edit responses. " +
						"But you still don't know how to delete them... so... wait 10 seconds and this " +
						"message will be deleted.",
				})
				if err != nil {
					s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
					return
				}
				time.Sleep(time.Second * 10)
				s.InteractionResponseDelete(s.State.User.ID, i.Interaction)
			})
		},
		"followups": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Followup messages are basically regular messages (you can create as many of them as you wish)
			// but work as they are created by webhooks and their functionality
			// is for handling additional messages after sending a response.

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					// Note: this isn't documented, but you can use that if you want to.
					// This flag just allows you to create messages visible only for the caller of the command
					// (user who triggered the command)
					Flags:   1 << 6,
					Content: "Surprise!",
				},
			})
			msg, err := s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
				Content: "Followup message has been created, after 5 seconds it will be edited",
			})
			if err != nil {
				s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something went wrong",
				})
				return
			}
			time.Sleep(time.Second * 5)

			s.FollowupMessageEdit(s.State.User.ID, i.Interaction, msg.ID, &discordgo.WebhookEdit{
				Content: "Now the original message is gone and after 10 seconds this message will ~~self-destruct~~ be deleted.",
			})

			time.Sleep(time.Second * 10)

			s.FollowupMessageDelete(s.State.User.ID, i.Interaction, msg.ID)

			s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
				Content: "For those, who didn't skip anything and followed tutorial along fairly, " +
					"take a unicorn :unicorn: as reward!\n" +
					"Also, as bonus... look at the original interaction response :D",
			})
		},
	}
)

func initialization(s *discordgo.Session) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

var ReadyInfoPublic *discordgo.Ready

func Ready(s *discordgo.Session, readyInfo *discordgo.Ready) {
	initialization(s)
	ReadyInfoPublic = readyInfo

	for i := range readyInfo.Guilds {
		for _, v := range BotCommands {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, readyInfo.Guilds[i].ID, v)
			if err != nil {
				log.Panicf("Cannot create '%v' command: %v", v.Name, err)
			}
		}
	}
}
