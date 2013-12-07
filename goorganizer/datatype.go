package goorganizer

import (
	"time"
	"appengine/datastore"
)

type User struct {
	Id datastore.Key
	Email string
	Nickname string
	Verified bool
	Password string
}

type Thread struct {
	Id datastore.Key
	Author User
	Title string
	Text string
	Time time.Time
	Posts []Post
	Participant []User
}

type Post struct {
    Author User
    Text string
    Time time.Time
}