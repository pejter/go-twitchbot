package main

import (
	"code.google.com/p/gcfg"
	"database/sql"
	"errors"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	irc "github.com/thoj/go-ircevent"
	"log"
	"os"
	"strings"
)

type Config struct {
	General struct {
		Debug                      bool
		User, Nick, Password, Room string
	}
	Database struct {
		Filename, Handler string
	}
}

type IRCBot struct {
	conn          *irc.Connection
	Moderators    []string
	Address, Room string
	callbacks     map[string]func(SimpleMessage)
}

type SimpleMessage struct {
	User    string
	Content string
}

// Launches the bot connecting it to the channel and listening for messages
func (bot *IRCBot) Run() {
	bot.conn.Debug = cfg.General.Debug
	bot.conn.Password = cfg.General.Password

	if err := bot.conn.Connect(bot.Address); err != nil {
		log.Fatalln("Could not connect", err)
	}
	defer bot.conn.Disconnect()

	bot.conn.Join(bot.Room)
	log.Printf("Connected to channel %s\n", bot.Room)

	bot.conn.SendRaw("TWITCHCLIENT 3")
	log.Println("Subscribed to user events")

	bot.Message("/mods")

	bot.conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		m := SimpleMessage{e.User, e.Message()}

		if !(strings.HasPrefix(m.Content, "!") || m.User == "twitchnotify" || m.User == "jtv") {
			checkSpam(m)
			return
		}

		if m.User == "twitchnotify" {
			user := strings.Fields(m.Content)[0]
			log.Println("New sub: ", user)
		} else if m.User == "jtv" {
			if strings.HasPrefix(m.Content, "The moderators of this room are: ") {
				list := strings.TrimPrefix(m.Content, "The moderators of this room are: ")
				mods := strings.Split(list, ", ")
				bot.Moderators = mods
				log.Println("Moderator list updated. Mods:")
				for _, m := range bot.Moderators {
					log.Println(m)
				}
			}
		}

		for key, callback := range bot.callbacks {
			if strings.HasPrefix(m.Content, key) {
				strings.TrimPrefix(m.Content, key)
				go callback(m)
				return
			}
		}
	})

	//bot.Message("Levelbot joined channel")
	bot.conn.Loop()
}

// Sends message to channel
func (bot *IRCBot) Message(s string) {
	bot.conn.Privmsg(bot.Room, s)
}

// Sends formatted message to channel
func (bot *IRCBot) Messagef(s string, v ...interface{}) {
	bot.conn.Privmsgf(bot.Room, s, v...)
}

// Adds new callback to the dispather
func (bot *IRCBot) RegisterCallback(command string, callback func(SimpleMessage)) error {
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
func NewIRCBot(address, nick, user, room string) *IRCBot {
	return &IRCBot{irc.IRC(nick, user), []string{user}, address, "#" + room, make(map[string]func(SimpleMessage))}
}

func initDb() {
	dbHandle, err := sql.Open(cfg.Database.Handler, cfg.Database.Filename)
	if err != nil {
		log.Fatalln("sql.Open failed")
	}

	// construct a gorp DbMap
	db = &gorp.DbMap{Db: dbHandle, Dialect: gorp.SqliteDialect{}}

	// Enable gorp logging
	if cfg.General.Debug {
		db.TraceOn("[db]", log.New(os.Stdout, "", log.LstdFlags))
	}
}

// Global IRC Bot & Config definition
var cfg Config
var bot *IRCBot
var db *gorp.DbMap

func main() {
	err := gcfg.ReadFileInto(&cfg, "config.ini")
	if err != nil {
		log.Fatalln("Error while reading config: ", err)
	}
	bot = NewIRCBot("irc.twitch.tv:6667", cfg.General.Nick, cfg.General.User, cfg.General.Room)
	initDb()
	initInfo()
	initPoll()
	bot.Run()
}
