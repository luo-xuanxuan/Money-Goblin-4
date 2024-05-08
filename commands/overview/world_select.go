package overview

import (
	"MoneyGoblin4/db"
	"MoneyGoblin4/structs"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func create_world_select() discordgo.SelectMenu {

	options := []discordgo.SelectMenuOption{
		{
			Label:       "Overview",
			Value:       "Overview",
			Description: "Switches back to the Overview.",
		},
	}

	for _, world := range db.World_Statuses {
		options = append(options, discordgo.SelectMenuOption{
			Label:       world.Name,
			Value:       world.Name,
			Description: fmt.Sprintf("Views %s", world.Name),
		})
	}

	menu := discordgo.SelectMenu{
		CustomID:    "world_select",
		Placeholder: "Select a World",
		Options:     options,
	}

	return menu
}

func world_select_handler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	component := i.MessageComponentData()

	if component.CustomID != "world_select" {
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		fmt.Println("Failed to respond to interaction:", err)
		return
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				create_world_select(),
			},
		},
	}

	message := &discordgo.MessageEdit{
		ID:      i.Message.ID,
		Channel: i.Message.ChannelID,
	}

	if component.Values[0] == "Overview" {
		message.Embeds = &[]*discordgo.MessageEmbed{build_overview()}
		components = append(components, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				create_sync_button(),
			},
		})
		message.Components = &components
	} else {
		message.Embeds = &[]*discordgo.MessageEmbed{build_world_embed(component.Values[0])}
		components = append(components, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				create_fc_select(component.Values[0]),
			},
		})
		components = append(components, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				create_sync_button(),
			},
		})
		message.Components = &components
	}

	s.ChannelMessageEditComplex(message)
}

func build_world_embed(world string) *discordgo.MessageEmbed {
	embed := discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name: "Money Goblin 4.0",
		},
		Color: 0xF8C8DC,
	}

	var world_status *structs.World_Status

	for _, v := range db.World_Statuses {
		if v.Name == world {
			world_status = v
			break
		}
	}

	for _, fc := range world_status.Free_Company_List {
		value := ""
		for _, sub := range fc.Submersible_List {
			mark := "✅"
			if sub.Return_Time < time.Now().Unix() {
				mark = "❌"
			}
			value += fmt.Sprintf("%s <t:%d:R>\n", mark, sub.Return_Time)
		}
		if fc.Tanks < 4860 {
			value += "<:10155:1178500241847242832>"
		}
		if fc.Repairs < 1440 {
			value += "<:10373:1178500307169333319>"
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   fc.Name,
			Value:  value,
			Inline: true,
		})
	}

	embed.Description = fmt.Sprintf("<t:%d:F>", time.Now().Unix())

	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: "https://raw.githubusercontent.com/luo-xuanxuan/MoneyGoblinUploader/master/Data/money_goblin.png",
	}

	return &embed
}
