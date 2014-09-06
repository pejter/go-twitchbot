package main

import (
	"errors"
	//"fmt"
	"code.google.com/p/gcfg"
	irc "github.com/thoj/go-ircevent"
	"log"
	"strings"
)

type Config struct {
	General struct {
		Debug                      bool
		User, Nick, Password, Room string
	}
}

type IRCBot struct {
	Conn                                *irc.Connection
	Moderators                          []string
	Address, Nick, User, Password, Room string
	callbacks                           map[string]func(*IRCBot, SimpleMessage)
}

type SimpleMessage struct {
	User    string
	Content string
}

// Launches the bot connecting it to the channel and listening for messages
func (bot *IRCBot) Run() {
	bot.Conn.Debug = cfg.General.Debug
	bot.Conn.Password = bot.Password

	if err := bot.Conn.Connect(bot.Address); err != nil {
		log.Fatalln("Could not connect", err)
	}
	defer bot.Conn.Disconnect()

	bot.Conn.Join(bot.Room)
	log.Printf("Connected to channel %s\n", bot.Room)

	bot.Conn.SendRaw("TWITCHCLIENT 3")
	log.Println("Subscribed to user events")

	bot.Message("/mods")

	bot.Conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		if !(strings.HasPrefix(e.Message(), "!") || e.Nick == "twitchnotify" || e.Nick == "jtv") {
			return
		}

		if e.Nick == "twitchnotify" {
			user := strings.Fields(e.Message())[0]
			log.Println("New sub: ", user)
		} else if e.Nick == "jtv" {
			if strings.HasPrefix(e.Message(), "The moderators of this room are: ") {
				list := strings.TrimPrefix(e.Message(), "The moderators of this room are: ")
				mods := strings.Split(list, ", ")
				bot.Moderators = mods
				log.Println("Moderator list updated. Mods:")
				for _, m := range bot.Moderators {
					log.Println(m)
				}
			}
		}

		for key, callback := range bot.callbacks {
			if strings.HasPrefix(e.Message(), key) {
				m := SimpleMessage{e.Nick, e.Message()}
				go callback(bot, m)
				return
			}
		}
	})

	bot.Conn.Loop()
}

// Sends message to channel
func (bot *IRCBot) Message(s string) {
	bot.Conn.Privmsg(bot.Room, s)
}

// Sends formatted message to channel
func (bot *IRCBot) Messagef(s string, v ...interface{}) {
	bot.Conn.Privmsgf(bot.Room, s, v...)
}

// Adds new callback to the dispather
func (bot *IRCBot) RegisterCallback(command string, callback func(*IRCBot, SimpleMessage)) error {
	if _, ok := bot.callbacks[command]; !ok {
		bot.callbacks[command] = callback
		return nil
	} else {
		return errors.New("Callback for this command already exists")
	}
}

// Removes callback from the dispatcher
func (bot *IRCBot) RemoveCallback(command string) {
	delete(bot.callbacks, command)
}

// Returns a new iRCBot instance
func NewIRCBot(address, nick, user, password, room string) *IRCBot {
	return &IRCBot{irc.IRC(nick, user), []string{user}, address, nick, user, password, "#" + room, make(map[string]func(*IRCBot, SimpleMessage))}
}

// Checks if given user is a channel moderator
func isMod(bot *IRCBot, user string) bool {
	for _, u := range bot.Moderators {
		if u == user {
			return true
		}
	}
	return false
}

// Global IRC Bot & Config definition
var cfg Config
var bot *IRCBot

func main() {
	err := gcfg.ReadFileInto(&cfg, "config.ini")
	if err != nil {
		log.Fatalln("Error while reading config: ", err)
	}
	bot = NewIRCBot("irc.twitch.tv:6667", cfg.General.Nick, cfg.General.User, cfg.General.Password, cfg.General.Room)
	initInfo(bot)
	initPoll(bot)
	bot.Run()
}
