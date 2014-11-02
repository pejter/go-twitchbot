package main

import (
	"fmt"
	"log"
)

// Checks if given user is a channel moderator
func hasPerm(user string, perm ...string) bool {
	for _, u := range bot.Moderators {
		if u == user {
			return true
		}
	}
	query := fmt.Sprintf("SELECT count(1) FROM permissions WHERE user='%s'", user)
	for _, p := range perm {
		query += fmt.Sprintf(" AND perm LIKE '%s'", p)
	}
	result, err := db.SelectInt(query)
	log.Printf("Result : %v", result)
	if err != nil {
		log.Printf("Error getting permission '%s' for user '%s': %s\n", perm, user, err)
		return false
	} else if result == 0 {
		return false
	} else {
		return true
	}
}
