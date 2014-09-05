package main

import (
	"errors"
	//"fmt"
	irc "github.com/thoj/go-ircevent"
	"log"
	"strings"
)

const (
	DEBUG = true
	PASS  = "oauth:6oppx20dk5tcvxhgwa5a16foj4cacri"
)

type IRCBot struct {
	Conn                                *irc.Connection
	Address, Nick, User, Password, Room string
	callbacks                           map[string]func(*IRCBot, *irc.Event)
}

// Launches the bot connecting it to the channel and listening for messages
func (bot *IRCBot) Run() {
	bot.Conn.Debug = DEBUG
	bot.Conn.Password = bot.Password

	if err := bot.Conn.Connect(bot.Address); err != nil {
		log.Fatalln("Could not connect", err)
	}
	defer bot.Conn.Disconnect()

	bot.Conn.Join(bot.Room)
	log.Printf("Connected to channel %s\n", bot.Room)

	bot.Conn.SendRaw("TWITCHCLIENT 3")
	log.Println("Subscribed to user events")

	bot.Conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		if !(strings.HasPrefix(e.Message(), "!") || e.Nick == "twitchnotify") {
			return
		}

		if e.Nick == "twitchnotify" && strings.Contains(e.Message(), "just subscribed!") {
			user := strings.Fields(e.Message())[0]
			log.Println("New sub: ", user)
		}

		for key, callback := range bot.callbacks {
			if strings.HasPrefix(e.Message(), key) {
				go callback(bot, e)
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

func (bot *IRCBot) RegisterCallback(command string, callback func(*irc.Event)) error {
	if _, ok := bot.callbacks[command]; !ok {
		bot.callbacks[command] = callback
		return nil
	} else {
		return errors.New("Callback for this command already exists")
	}
}

// Returns a new iRCBot instance
func NewIRCBot(address, nick, user, password, room string) *IRCBot {
	return &IRCBot{irc.IRC(nick, user), address, nick, user, password, room, make(map[string]func(*irc.Event))}
}

// Global IRC Bot definition
var bot = NewIRCBot("irc.twitch.tv:6667", "pejter95", "pejter95", PASS, "#levelbf")

func main() {
	bot.Run()
}
