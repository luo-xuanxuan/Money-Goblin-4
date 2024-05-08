package hi

import (
	"MoneyGoblin4/disc_util"

	"github.com/bwmarrin/discordgo"
)

type Hi_Command struct {
	command *discordgo.ApplicationCommand
}

func (h *Hi_Command) Reference() *discordgo.ApplicationCommand {
	if h.command == nil {
		h.command = &discordgo.ApplicationCommand{
			Name:        "hi",
			Description: "Hi!",
		}
	}
	return h.command
}

func (*Hi_Command) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	err := disc_util.Respond_Ephemeral(s, i, "Hey There!")
	if err != nil {
		return err
	}

	return nil
}
