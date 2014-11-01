package main

import (
	"database/sql"
	"log"
)

// Checks if given user is a channel moderator
func hasPerm(user string, permission string) bool {
	for _, u := range bot.Moderators {
		if u == user {
			return true
		}
	}
	_, err := db.SelectInt("select count(1) from permissions where user=? and perm=?", user, permission)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error getting permission '%s' for user '%s': %s\n", permission, user, err)
		return false
	} else if err == sql.ErrNoRows {
		return false
	} else {
		return true
	}
}
