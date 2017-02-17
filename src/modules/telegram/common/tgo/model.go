package tgo

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// ChannelInfo is the channel informmation from cli
type ChannelInfo struct {
	ID                string `json:"id"`
	PeerType          string `json:"peer_type"`
	PeerID            int    `json:"peer_id"`
	PrintName         string `json:"print_name"`
	Flags             int    `json:"flags"`
	Title             string `json:"title"`
	ParticipantsCount int    `json:"participants_count"`
	AdminsCount       int    `json:"admins_count"`
	KickedCount       int    `json:"kicked_count"`
	About             string `json:"about"`
	Username          string `json:"username"`
}

// Scan convert the json array ino string slice
func (ci *ChannelInfo) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}

	return json.Unmarshal(b, ci)
}

// Value try to get the string slice representation in database
func (ci ChannelInfo) Value() (driver.Value, error) {
	return json.Marshal(ci)
}

// History is a single history record
type History struct {
	Event string `json:"event"`
	ID    string `json:"id"`
	Flags int    `json:"flags"`
	From  struct {
		ID                string `json:"id"`
		PeerType          string `json:"peer_type"`
		PeerID            int    `json:"peer_id"`
		PrintName         string `json:"print_name"`
		Flags             int    `json:"flags"`
		Title             string `json:"title"`
		ParticipantsCount int    `json:"participants_count"`
		AdminsCount       int    `json:"admins_count"`
		KickedCount       int    `json:"kicked_count"`
		Username          string `json:"username"`
	} `json:"from"`
	To struct {
		ID                string `json:"id"`
		PeerType          string `json:"peer_type"`
		PeerID            int    `json:"peer_id"`
		PrintName         string `json:"print_name"`
		Flags             int    `json:"flags"`
		Title             string `json:"title"`
		ParticipantsCount int    `json:"participants_count"`
		AdminsCount       int    `json:"admins_count"`
		KickedCount       int    `json:"kicked_count"`
		Username          string `json:"username"`
	} `json:"to"`
	Out     bool   `json:"out"`
	Unread  bool   `json:"unread"`
	Service bool   `json:"service"`
	Date    int    `json:"date"`
	Views   int    `json:"views"`
	PostID  int    `json:"post_id"`
	Link    string `json:"link"`
	Media   *struct {
		Type    string `json:"type"`
		Caption string `json:"caption"`
	} `json:"media,omitempty"`
	FwdFrom *struct {
		ID                string `json:"id"`
		PeerType          string `json:"peer_type"`
		PeerID            int    `json:"peer_id"`
		PrintName         string `json:"print_name"`
		Flags             int    `json:"flags"`
		Title             string `json:"title"`
		ParticipantsCount int    `json:"participants_count"`
		AdminsCount       int    `json:"admins_count"`
		KickedCount       int    `json:"kicked_count"`
		Username          string `json:"username"`
	} `json:"fwd_from"`
	FwdDate int    `json:"fwd_date,omitempty"`
	Text    string `json:"text,omitempty"`
}

// UserInfo is the user information
type UserInfo struct {
	ID        string `json:"id"`
	PeerType  string `json:"peer_type"`
	PeerID    int    `json:"peer_id"`
	PrintName string `json:"print_name"`
	Flags     int    `json:"flags"`
	FirstName string `json:"first_name"`
	When      string `json:"when"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Username  string `json:"username"`
}

//AddContact - duplicate userInfo struct
type AddContact struct {
	ID        string `json:"id"`
	PeerType  string `json:"peer_type"`
	PeerID    int    `json:"peer_id"`
	PrintName string `json:"print_name"`
	Flags     int    `json:"flags"`
	FirstName string `json:"first_name"`
	When      string `json:"when"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Username  string `json:"username"`
}

// ChannelList is the list of channel
type ChannelList struct {
	ID                string `json:"id"`
	PeerType          string `json:"peer_type"`
	PeerID            int    `json:"peer_id"`
	PrintName         string `json:"print_name"`
	Flags             int    `json:"flags"`
	Title             string `json:"title"`
	ParticipantsCount int    `json:"participants_count"`
	AdminsCount       int    `json:"admins_count"`
	KickedCount       int    `json:"kicked_count"`
}

// SuccessResp is the message for success status
type SuccessResp struct {
	Result string `json:"result"`
}

// MsgResp is the message response
type MsgResp struct {
	Result string `json:"result"`
}

// Contact is the single contact
type Contact struct {
	ID        string `json:"id"`
	PeerType  string `json:"peer_type"`
	PeerID    int    `json:"peer_id"`
	PrintName string `json:"print_name"`
	Flags     int    `json:"flags"`
	FirstName string `json:"first_name,omitempty"`
	When      string `json:"when,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Username  string `json:"username,omitempty"`
}

// ChannelUser is the user for a channel
type ChannelUser struct {
	ID                string `json:"id"`
	PeerType          string `json:"peer_type"`
	PeerID            int    `json:"peer_id"`
	PrintName         string `json:"print_name,omitempty"`
	FirstName         string `json:"first_name,omitempty"`
	When              string `json:"when"`
	LastName          string `json:"last_name"`
	Flags             int    `json:"flags"`
	Phone             string `json:"phone,omitempty"`              //friend
	Username          string `json:"username,omitempty"`           //user
	Title             string `json:"title,omitempty"`              //channel
	ParticipantsCount int    `json:"participants_count,omitempty"` //channel
	AdminsCount       int    `json:"admins_count,omitempty"`       //channel
	KickedCount       int    `json:"kicked_count,omitempty"`       //channel
}
