package report

import (
	"MoneyGoblin4/db"
	"fmt"
	"sort"
	"time"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Report_Command struct {
	command *discordgo.ApplicationCommand
}

func (r *Report_Command) Reference() *discordgo.ApplicationCommand {
	if r.command == nil {
		r.command = &discordgo.ApplicationCommand{
			Name:        "report",
			Description: "Generates a report!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "days",
					Description: "How many days to include in the report.",
					Required:    false,
				},
			},
		}
	}
	return r.command
}

func (*Report_Command) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	data := i.ApplicationCommandData()
	options := data.Options

	days := 1

	for _, v := range options {
		if v.Name == "days" {
			days = int(v.IntValue())
		}
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{build_report(days)},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

var salvage_values = map[int]int{
	22500: 8000,
	22501: 9000,
	22502: 10000,
	22503: 13000,
	22504: 27000,
	22505: 28500,
	22506: 30000,
	22507: 34500,
}

var salvage_emotes = map[int]string{
	22500: "<:22500:1176734113319882853>",
	22501: "<:22501:1176734115614171226>",
	22502: "<:22502:1176734116717264988>",
	22503: "<:22503:1176734117824561213>",
	22504: "<:22504:1176734118722158693>",
	22505: "<:22505:1176734120731217980>",
	22506: "<:22506:1176734121905635430>",
	22507: "<:22507:1176734123289751572>",
}

func build_report(duration int) *discordgo.MessageEmbed {
	response, err := db.Fetch_Report(duration)
	if err != nil {
		return nil
	}

	report := make(map[string]map[int]int)
	totals := make(map[string]int)

	for _, row := range response {
		world := row.World

		// Initialize map for this world if it doesn't exist
		if _, ok := report[world]; !ok {
			report[world] = map[int]int{}
		}

		// Increment the quantity of the item directly
		report[world][row.Item_id] += row.Quantity

		// Initialize totals for this world if it doesn't exist
		if _, ok := totals[world]; !ok {
			totals[world] = 0
		}

		// Increment the total value for the world
		totals[world] += salvage_values[row.Item_id] * row.Quantity
	}

	sorted_worlds := make([]string, 0, len(totals))
	for world := range totals {
		sorted_worlds = append(sorted_worlds, world)
	}
	sort.Slice(sorted_worlds, func(i, j int) bool {
		return totals[sorted_worlds[i]] > totals[sorted_worlds[j]]
	})

	total := 0
	for _, value := range totals {
		total += value
	}

	p := message.NewPrinter(language.English)

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("Salvage Report: <t:%d:D>-<t:%d:D>", time.Now().Unix()-(int64(duration)*86400), time.Now().Unix()),
		Author: &discordgo.MessageEmbedAuthor{
			Name: "Money Goblin 4.0",
		},
		Color:       0xF8C8DC,
		Description: p.Sprintf("Total: %dg", total),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://raw.githubusercontent.com/luo-xuanxuan/MoneyGoblinUploader/master/Data/money_goblin.png",
		},
		Fields: []*discordgo.MessageEmbedField{},
	}

	for _, v := range sorted_worlds {

		loot := ""

		for item, quantity := range report[v] {
			loot += p.Sprintf("%sx%d ", salvage_emotes[item], quantity)
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   v,
			Value:  p.Sprintf("%dg\n%s", totals[v], loot),
			Inline: true,
		})
	}

	return embed
}
