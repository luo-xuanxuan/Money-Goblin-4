package overview

import (
	"MoneyGoblin4/db"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Overview_Command struct {
	command *discordgo.ApplicationCommand
}

var handled = false

func (oc *Overview_Command) Preload(s *discordgo.Session, guild string) error {
	if handled {
		return nil
	}

	s.AddHandler(world_select_handler)
	s.AddHandler(fc_select_handler)
	s.AddHandler(sync_handler)

	handled = true

	return nil
}

func (oc *Overview_Command) Reference() *discordgo.ApplicationCommand {
	if oc.command == nil {
		oc.command = &discordgo.ApplicationCommand{
			Name:        "overview",
			Description: "Creates an overview",
		}
	}
	return oc.command
}

func (oc *Overview_Command) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{build_overview()},
			Components: []discordgo.MessageComponent{
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
		},
	})
	if err != nil {
		return err
	}

	return err
}

func build_overview() *discordgo.MessageEmbed {

	embed := &discordgo.MessageEmbed{
		Color: 0xF8C8DC,
		Author: &discordgo.MessageEmbedAuthor{
			Name: "Money Goblin 4.0",
		},
	}

	returned_subs := 0
	total_subs := 0

	for _, world := range db.World_Statuses {
		value := ""
		for _, fc := range world.Free_Company_List {
			mark := "✅"
			for _, sub := range fc.Submersible_List {
				total_subs += 1
				if sub.Return_Time < time.Now().Unix() {
					returned_subs += 1
					mark = "❌"
				}
			}
			value += mark
			if fc.Tanks < 4860 {
				value += "<:10155:1178500241847242832>"
			}
			if fc.Repairs < 1440 {
				value += "<:10373:1178500307169333319>"
			}
			value += fmt.Sprintf(" %s\n", fc.Name)
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   world.Name,
			Value:  value,
			Inline: true,
		})
	}

	embed.Description = fmt.Sprintf("Last Refreshed: <t:%d:R>\n%d/%d Submersibles Returned", db.Status_Last_Updated, returned_subs, total_subs)

	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: "https://raw.githubusercontent.com/luo-xuanxuan/MoneyGoblinUploader/master/Data/money_goblin.png",
	}

	return embed
}
