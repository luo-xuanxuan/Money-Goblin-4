package overview

import (
	"MoneyGoblin4/db"
	"MoneyGoblin4/structs"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func create_fc_select(world string) discordgo.SelectMenu {
	options := make([]discordgo.SelectMenuOption, 0)

	for _, w := range db.World_Statuses {
		if w.Name == world {
			for _, fc := range w.Free_Company_List {
				options = append(options, discordgo.SelectMenuOption{
					Label:       fc.Name,
					Value:       fc.Name,
					Description: fmt.Sprintf("Views %s", fc.Name),
				})
			}
		}

	}

	menu := discordgo.SelectMenu{
		CustomID:    "fc_select",
		Placeholder: "Select a Free Company",
		Options:     options,
	}

	return menu
}

func fc_select_handler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	component := i.MessageComponentData()

	if component.CustomID != "fc_select" {
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		fmt.Println("Failed to respond to interaction:", err)
		return
	}

	world, embed := build_fc_view(component.Values[0])

	message := &discordgo.MessageEdit{
		ID:      i.Message.ID,
		Channel: i.Message.ChannelID,
		Embeds:  &[]*discordgo.MessageEmbed{embed},
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				create_world_select(),
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				create_fc_select(world),
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				create_sync_button(),
			},
		},
	}

	message.Components = &components

	s.ChannelMessageEditComplex(message)
}

func build_fc_view(fc string) (string, *discordgo.MessageEmbed) {

	embed := discordgo.MessageEmbed{
		Title: fc,
		Author: &discordgo.MessageEmbedAuthor{
			Name: "Money Goblin 4.0",
		},
		Color:       0xF8C8DC,
		Description: "",
	}

	var status *structs.Free_Company_Status = nil
	var world string

	for _, w := range db.World_Statuses {
		for _, f := range w.Free_Company_List {
			if f.Name == fc { //this line is bugged if two FCs have the same name, i dont wanna fix it rn, but the fix would be to change values in the fc_select to match ID instead of name
				status = f
				world = w.Name
				break
			}
		}
		if status != nil {
			break
		}
	}

	for _, sub := range status.Submersible_List {
		mark := "✅"
		if sub.Return_Time < time.Now().Unix() {
			mark = "❌"
		}
		embed.Description += fmt.Sprintf("%s %s: <t:%d:R>\n<t:%d:F>\n", mark, sub.Name, sub.Return_Time, sub.Return_Time)
	}

	value := fmt.Sprintf("<:10155:1178500241847242832>x%d\n<:10373:1178500307169333319>x%d\n", status.Tanks, status.Repairs)
	if status.Tanks < 4860 {
		value += "<:10155:1178500241847242832>"
	}
	if status.Repairs < 1440 {
		value += "<:10373:1178500307169333319>"
	}

	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   fc,
		Value:  value,
		Inline: true,
	})

	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: "https://raw.githubusercontent.com/luo-xuanxuan/MoneyGoblinUploader/master/Data/money_goblin.png",
	}

	return world, &embed

}
