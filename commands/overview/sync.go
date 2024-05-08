package overview

import (
	"MoneyGoblin4/db"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func create_sync_button() discordgo.Button {
	button := discordgo.Button{
		CustomID: "sync",
		Style:    discordgo.SecondaryButton,
		Emoji: &discordgo.ComponentEmoji{
			Name: "🔄",
		},
	}

	return button
}

func sync_handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	component := i.MessageComponentData()

	if component.CustomID != "sync" {
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		fmt.Println("Failed to respond to interaction:", err)
		return
	}

	err = db.Update_World_Status()
	if err != nil {
		fmt.Println("World update error:", err)
		return
	}

	message := &discordgo.MessageEdit{
		ID:      i.Message.ID,
		Channel: i.Message.ChannelID,
		Embeds:  &[]*discordgo.MessageEmbed{build_overview()},
		Components: &[]discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					create_world_select(),
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					create_sync_button(),
				},
			},
		},
	}

	s.ChannelMessageEditComplex(message)

}
