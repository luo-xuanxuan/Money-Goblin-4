package reminder

import (
	"MoneyGoblin4/db"
	"MoneyGoblin4/disc_util"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Reminder_Command struct {
	command *discordgo.ApplicationCommand
}

var reminder_time int64 = 0
var last_reminder int64 = 0

func (r *Reminder_Command) Preload(s *discordgo.Session, guild string) error {
	go send_resend(s)
	return nil
}

func (r *Reminder_Command) Reference() *discordgo.ApplicationCommand {
	if r.command == nil {
		r.command = &discordgo.ApplicationCommand{
			Name:        "reminder",
			Description: "Sets channel for reminders!",
		}
	}
	return r.command
}

func (*Reminder_Command) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	update_reminder_channel(i.ChannelID)

	err := disc_util.Respond_Ephemeral(s, i, "This channel has been set for submersible reminders.")
	if err != nil {
		return err
	}

	return nil
}

func get_final_return_time() int64 {

	var return_time int64 = 0

	for _, w := range db.World_Statuses {
		for _, fc := range w.Free_Company_List {
			for _, sub := range fc.Submersible_List {
				if sub.Return_Time > return_time {
					return_time = sub.Return_Time
				}
			}

		}
	}

	return return_time
}

func send_resend(s *discordgo.Session) {

	embed := &discordgo.MessageEmbed{
		Title: "Resend Reminder",
		Author: &discordgo.MessageEmbedAuthor{
			Name: "Money Goblin 4.0",
		},
		Color:       0xF8C8DC,
		Description: "All Submersibles have Returned!",
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://raw.githubusercontent.com/luo-xuanxuan/MoneyGoblinUploader/master/Data/money_goblin.png",
		},
	}

	for {

		reminder_time = get_final_return_time()

		log.Printf("Reminder time: %d", reminder_time)

		if reminder_time == last_reminder {

			sleep := reminder_time - time.Now().Unix()

			if sleep < 1 {
				log.Println("Sleep time: 60")
				time.Sleep(60 * time.Second)
				continue
			}
			log.Printf("Sleep time: %d", sleep)
			time.Sleep(time.Duration(sleep) * time.Second)
			continue
		}

		if time.Now().Unix() >= reminder_time {

			if channel_id == "" {
				continue
			}

			_, err := s.ChannelMessageSendEmbed(channel_id, embed)
			if err != nil {
				log.Fatal(err.Error())
				continue
			}

			last_reminder = reminder_time
		}

		sleep := reminder_time - time.Now().Unix()

		if sleep < 1 {
			log.Println("Sleep time: 60")
			time.Sleep(60 * time.Second)
			continue
		}
		log.Printf("Sleep time: %d", sleep)
		time.Sleep(time.Duration(sleep) * time.Second)

	}
}
