package tgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"

	"time"

	"github.com/Sirupsen/logrus"
)

// TelegramCli is the interface to handle the telegram cli
type TelegramCli interface {

	//ChannelInfo channel_info <channel>  Prints info about channel (id, members, admin, etc.)
	ChannelInfo(channelId string) (*ChannelInfo, error)

	//AddContact add_contact <phone> <first name> <last name>    Tries to add user to contact list
	AddContact(phone, firstName, lastName string) (*AddContact, error)

	//ContactList contact_list()
	ContactList() ([]Contact, error)

	//History history <peer> [limit] [offset] Prints messages with this peer (most recent message lower). Also marks messages as read
	History(peer string, limit, offset int) ([]History, error)

	// FwdMsg fwd <peer> <msg-id>+    Forwards message to peer. Forward to secret chats is forbidden
	FwdMsg(peer string, msg string) (*SuccessResp, error)

	// UserInfo user_info <user>        Prints info about user (id, last online, phone)
	UserInfo(user string) (*UserInfo, error)

	// ChannelJoin channel_join <channel>  Joins to channel
	ChannelJoin(channelId string) (*SuccessResp, error)

	//Msg msg <peer> <kbd> <text> Sends text message to peer with custom kbd
	Msg(peer string, msg string) (*SuccessResp, error)

	//Post post <peer> <text>      Sends text message to peer as admin
	Post(peer string, msg string) (*SuccessResp, error)

	//ChannelList channel_list [limit=100] [offset=0]     List of last channels
	ChannelList(limit, offset int) ([]ChannelList, error)

	//ResolveUsername resolve_username username       Searches user by username
	ResolveUsername(chUser string) (*ChannelUser, error)

	//ChannelInvite channel_invite <channel> <user> Invites user to channel
	ChannelInvite(channelId, user string) (*SuccessResp, error)

	//GetSelf get_self        Get our user info
	GetSelf() (*UserInfo, error)
}

type TelegramCliFull interface {
	TelegramCli

	//accept_secret_chat <secret chat>        Accepts secret chat. Only useful with -E option
	AcceptSecretChat()

	//block_user <user>       Blocks user
	BlockUser()

	//broadcast <user>+ <text>        Sends text to several users at once
	Broadcast()

	//channel_get_admins <channel> [limit=100] [offset=0]     Gets channel admins
	ChannelGetAdmins()

	//channel_get_members <channel> [limit=100] [offset=0]    Gets channel members
	ChannelGetMembers()

	//channel_kick <channel> <user>   Kicks user from channel
	Channel_kick()

	//channel_leave <channel> Leaves from channel
	Channel_leave()

	//channel_list [limit=100] [offset=0]     List of last channels
	Channel_list()

	//channel_set_about <channel> <about>     Sets channel about info.
	Channel_set_about()

	//channel_set_admin <channel> <admin> <type>      Sets channel admin. 0 - not admin, 1 - moderator, 2 - editor
	Channel_set_admin()

	//channel_set_username <channel> <username>       Sets channel username info.
	Channel_set_username()

	//channel_set_photo <channel> <filename>  Sets channel photo. Photo will be cropped to square
	Channel_set_photo()

	//chat_add_user <chat> <user> [msgs-to-forward]   Adds user to chat. Sends him last msgs-to-forward message from this chat. Default 100
	Chat_add_user()

	//chat_del_user <chat> <user>     Deletes user from chat
	Chat_del_user()

	//chat_info <chat>        Prints info about chat (id, members, admin, etc.)
	Chat_info()

	//chat_set_photo <chat> <filename>        Sets chat photo. Photo will be cropped to square
	Chat_set_photo()

	//chat_upgrade <chat>     Upgrades chat to megagroup
	Chat_upgrade()

	//chat_with_peer <peer>   Interface option. All input will be treated as messages to this peer. Type /quit to end this mode
	ChatWithPeer()

	//clear   Clears all data and exits. For debug.
	Clear()

	//contact_search username Searches user by username
	ContactSearch()

	//create_channel <name> <about> <user>+   Creates channel with users
	CreateChannel()

	//create_group_chat <name> <user>+        Creates group chat with users
	CreateGroupChat()

	//create_secret_chat <user>       Starts creation of secret chat
	CreateSecretChat()

	//del_contact <user>      Deletes contact from contact list
	DelContact()

	//delete_msg <msg-id>     Deletes message
	DeleteMsg()

	//dialog_list [limit=100] [offset=0]      List of last conversations
	DialogList()

	//export_card     Prints card that can be imported by another user with import_card method
	ExportCard()

	//export_channel_link     Prints channel link that can be used to join to channel
	ExportChannelLink()

	//export_chat_link        Prints chat link that can be used to join to chat
	ExportChatLink()

	//fwd_media <peer> <msg-id>       Forwards message media to peer. Forward to secret chats is forbidden. Result slightly differs from fwd
	FwdMedia()

	//get_terms_of_service    Prints telegram's terms of service
	GetTermsOfService()

	//get_message <msg-id>    Get message by id
	GetMessage()

	//help [command]  Prints this help
	Help()

	//import_card <card>      Gets user by card and prints it name. You can then send messages to him as usual
	ImportCard()

	//import_chat_link <hash> Joins to chat by link
	ImportChatLink()

	//import_channel_link <hash>      Joins to channel by link
	ImportChannelLink()

	//load_audio <msg-id>     Downloads file to downloads dirs. Prints file name after download end
	LoadAudio()

	//load_channel_photo <channel>    Downloads file to downloads dirs. Prints file name after download end
	LoadChannelPhoto()

	//load_chat_photo <chat>  Downloads file to downloads dirs. Prints file name after download end
	LoadChatPhoto()

	//load_document <msg-id>  Downloads file to downloads dirs. Prints file name after download end
	LoadDocument()

	//load_document_thumb <msg-id>    Downloads file to downloads dirs. Prints file name after download end
	LoadDocumentThumb()

	//load_file <msg-id>      Downloads file to downloads dirs. Prints file name after download end
	LoadFile()

	//load_file_thumb <msg-id>        Downloads file to downloads dirs. Prints file name after download end
	LoadFileThumb()

	//load_photo <msg-id>     Downloads file to downloads dirs. Prints file name after download end
	LoadPhoto()

	//load_user_photo <user>  Downloads file to downloads dirs. Prints file name after download end
	LoadUserPhoto()

	//load_video <msg-id>     Downloads file to downloads dirs. Prints file name after download end
	LoadVideo()

	//load_video_thumb <msg-id>       Downloads file to downloads dirs. Prints file name after download end
	LoadVideoThumb()

	//main_session    Sends updates to this connection (or terminal). Useful only with listening socket
	MainSession()

	//mark_read <peer>        Marks messages with peer as read
	MarkRead()

	//post_audio <peer> <file>        Posts audio to peer
	PostAudio()

	//post_document <peer> <file>     Posts document to peer
	PostDocument()

	//post_file <peer> <file> Sends document to peer
	PostFile()

	//post_location <peer> <latitude> <longitude>     Sends geo location
	PostLocation()

	//post_photo <peer> <file> [caption]      Sends photo to peer
	PostPhoto()

	//post_text <peer> <file> Sends contents of text file as plain text message
	PostText()

	//post_video <peer> <file> [caption]      Sends video to peer
	PostVideo()

	//quit    Quits immediately
	Quit()

	//rename_channel <channel> <new name>     Renames channel
	RenameChannel()

	//rename_chat <chat> <new name>   Renames chat
	RenameChat()

	//rename_contact <user> <first name> <last name>  Renames contact
	RenameContact()

	//reply <msg-id> <text>   Sends text reply to message
	Reply()

	//reply_audio <msg-id> <file>     Sends audio to peer
	ReplyAudio()

	//reply_contact <msg-id> <phone> <first-name> <last-name> Sends contact (not necessary telegram user)
	ReplyContact()

	//reply_document <msg-id> <file>  Sends document to peer
	ReplyDocument()

	//reply_file <msg-id> <file>      Sends document to peer
	ReplyFile()

	//reply_location <msg-id> <latitude> <longitude>  Sends geo location
	ReplyLocation()

	//reply_photo <msg-id> <file> [caption]   Sends photo to peer
	ReplyPhoto()

	//reply_video <msg-id> <file>     Sends video to peer
	ReplyVideo()

	//safe_quit       Waits for all queries to end, then quits
	SafeQuit()

	//search [peer] [limit] [from] [to] [offset] pattern      Search for pattern in messages from date from to date to (unixtime) in messages with peer (if peer not present, in all messages)
	Search()

	//send_audio <peer> <file>        Sends audio to peer
	SendAudio()

	//send_contact <peer> <phone> <first-name> <last-name>    Sends contact (not necessary telegram user)
	SendContact()

	//send_document <peer> <file>     Sends document to peer
	SendDocument()

	//send_file <peer> <file> Sends document to peer
	SendFile()

	//send_location <peer> <latitude> <longitude>     Sends geo location
	SendLocation()

	//send_photo <peer> <file> [caption]      Sends photo to peer
	SendPhoto()

	//send_text <peer> <file> Sends contents of text file as plain text message
	SendText()

	//send_typing <peer> [status]     Sends typing notification. You can supply a custom status (range 0-10): none, typing, cancel, record video, upload video, record audio, upload audio, upload photo, upload document, geo, choose contact.
	SendTyping()

	//send_typing_abort <peer>        Sends typing notification abort
	SendTypingAbort()

	//send_video <peer> <file> [caption]      Sends video to peer
	SendVideo()

	//set <param> <value>     Sets value of param. Currently available: log_level, debug_verbosity, alarm, msg_num
	Set()

	//set_password <hint>     Sets password
	SetPassword()

	//set_profile_name <first-name> <last-name>       Sets profile name.
	SetProfileName()

	//set_profile_photo <filename>    Sets profile photo. Photo will be cropped to square
	SetProfilePhoto()

	//set_ttl <secret chat>   Sets secret chat ttl. Client itself ignores ttl
	SetTtl()

	//set_username <name>     Sets username.
	SetUsername()

	//set_phone_number <phone>        Changes the phone number of this account
	SetPhoneNumber()

	//show_license    Prints contents of GPL license
	ShowLicense()

	//start_bot <bot> <chat> <data>   Adds bot to chat
	StartBot()

	//stats   For debug purpose
	Stats()

	//status_online   Sets status as online
	StatusOnline()

	//status_offline  Sets status as offline
	StatusOffline()

	//unblock_user <user>     Unblocks user
	UnblockUser()

	//version Prints client and library version
	Version()

	//view_audio <msg-id>     Downloads file to downloads dirs. Then tries to open it with system default action
	ViewAudio()

	//view_channel_photo <channel>    Downloads file to downloads dirs. Then tries to open it with system default action
	ViewChannelPhoto()

	//view_chat_photo <chat>  Downloads file to downloads dirs. Then tries to open it with system default action
	ViewChatPhoto()

	//view_document <msg-id>  Downloads file to downloads dirs. Then tries to open it with system default action
	ViewDocument()

	//view_document_thumb <msg-id>    Downloads file to downloads dirs. Then tries to open it with system default action
	ViewDocumentThumb()

	//view_file <msg-id>      Downloads file to downloads dirs. Then tries to open it with system default action
	ViewFile()

	//view_file_thumb <msg-id>        Downloads file to downloads dirs. Then tries to open it with system default action
	ViewFileThumb()

	//view_photo <msg-id>     Downloads file to downloads dirs. Then tries to open it with system default action
	ViewPhoto()

	//view_user_photo <user>  Downloads file to downloads dirs. Then tries to open it with system default action
	ViewUserPhoto()

	//view_video <msg-id>     Downloads file to downloads dirs. Then tries to open it with system default action
	ViewVideo()

	//view_video_thumb <msg-id>       Downloads file to downloads dirs. Then tries to open it with system default action
	ViewVideoThumb()

	//view <msg-id>   Tries to view message contents
	View()

	//visualize_key <secret chat>     Prints visualization of encryption key (first 16 bytes sha1 of it in fact)
	VisualizeKey()
}

type telegram struct {
	lock *sync.Mutex
	Conn *net.TCPConn
	IP   net.IP
	Port int
}

func NewTelegramCli(ip net.IP, port int) (TelegramCli, error) {
	t := &telegram{
		lock: &sync.Mutex{},
		IP:   ip,
		Port: port,
	}

	if err := t.connect(); err != nil {
		return nil, err
	}
	return t, nil

}

func (t *telegram) connect() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	var err error
	t.Conn, err = net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   t.IP,
		Port: t.Port,
	})
	return err
}

func (t *telegram) disConnect() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	err := t.Conn.Close()
	t.Conn = nil
	return err
}

func (t *telegram) exec(cmd string) (d []byte, e error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	defer time.Sleep(time.Second)

	if t.Conn == nil {
		return nil, errors.New("not connected")
	}
	defer func() {
		logrus.Debugf("[TCLI] %s , err : %v result : %d", cmd, e, len(d))
	}()
	_, err := t.Conn.Write([]byte(cmd + "\n"))
	if err != nil {
		return nil, err
	}
	l := len("ANSWER ")
	buf := make([]byte, l)
	i, err := t.Conn.Read(buf)
	if err != nil || i != l {
		return nil, errors.New("can not read header")
	}

	var f []byte
	for {
		b := make([]byte, 1)
		i, err := t.Conn.Read(b)
		if err != nil || i != 1 {
			return nil, errors.New("can not read header")
		}
		if string(b) == "\n" {
			break
		}
		f = append(f, b...)
	}
	realLen, err := strconv.Atoi(string(f))
	if err != nil {
		return nil, errors.New("can not read header")
	}
	buffer := make([]byte, realLen+1)
	read := 0
	for read < realLen {
		i, err = t.Conn.Read(buffer[read:])
		read += i
		if err != nil {
			return nil, errors.New("can not read data")
		}
	}
	return buffer, nil
}

func (t *telegram) ChannelInfo(channelId string) (*ChannelInfo, error) {

	var data ChannelInfo
	cmd := fmt.Sprintf("channel_info %s", channelId)
	x, err := t.exec(cmd)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(x, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil

}
func (t *telegram) History(user string, limit, offset int) ([]History, error) {
	var data []History
	cmd := fmt.Sprintf("history %s %d %d", user, limit, offset)
	x, err := t.exec(cmd)
	//fmt.Print(string(x))
	if err != nil {
		return nil, err
	}
	fmt.Println(string(x))
	err = json.Unmarshal(x, &data)
	if err != nil {
		return nil, err
	}
	//fmt.Println(data)
	return data, nil

}
func (t *telegram) UserInfo(user string) (*UserInfo, error) {
	var data UserInfo
	cmd := fmt.Sprintf("user_info %s", user)
	x, err := t.exec(cmd)
	//fmt.Print(string(x))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(x, &data)
	if err != nil {
		return nil, err
	}
	//fmt.Println(data)
	return &data, nil
}
func (t *telegram) AddContact(phone, firstName, lastName string) (*AddContact, error) {
	var data AddContact

	cmd := fmt.Sprintf("add_contact %s %s %s", phone, firstName, lastName)
	x, err := t.exec(cmd)
	//fmt.Print(string(x))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(x, &data)
	if err != nil {
		return nil, err
	}
	//fmt.Println(data)
	return &data, nil

}
func (t *telegram) FwdMsg(peer string, msg string) (*SuccessResp, error) {
	var data SuccessResp

	cmd := fmt.Sprintf("fwd %s %s", peer, msg)
	x, err := t.exec(cmd)
	//fmt.Print(string(x))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(x, &data)
	if err != nil {
		return nil, err
	}
	//fmt.Println(data)
	return &data, nil

}
func (t *telegram) Msg(peer string, msg string) (*SuccessResp, error) {
	var data SuccessResp

	cmd := fmt.Sprintf("msg %s %s", peer, msg)
	x, err := t.exec(cmd)
	//fmt.Print(string(x))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(x, &data)
	if err != nil {
		return nil, err
	}
	//fmt.Println(data)
	return &data, nil

}
func (t *telegram) ContactList() ([]Contact, error) {
	var data []Contact
	cmd := fmt.Sprint("contact_list ")
	x, err := t.exec(cmd)
	//fmt.Print(string(x))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(x, &data)
	if err != nil {
		return nil, err
	}
	//fmt.Println(data)
	return data, nil
}
func (t *telegram) Post(peer string, msg string) (*SuccessResp, error) {
	var data SuccessResp
	cmd := fmt.Sprintf("msg %s %s", peer, msg)
	x, err := t.exec(cmd)
	//fmt.Print(string(x))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(x, &data)
	if err != nil {
		return nil, err
	}
	//fmt.Println(data)
	return &data, nil
}
func (t *telegram) ChannelJoin(channelId string) (*SuccessResp, error) {

	var data SuccessResp
	cmd := fmt.Sprintf("channel_join %s", channelId)
	x, err := t.exec(cmd)
	//fmt.Print(string(x))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(x, &data)
	if err != nil {
		return nil, err
	}
	//fmt.Println(data)
	return &data, nil

}
func (t *telegram) ChannelList(limit, offset int) ([]ChannelList, error) {
	var data []ChannelList
	cmd := fmt.Sprintf("channel_list %d %d", limit, offset)
	x, err := t.exec(cmd)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(x, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (t *telegram) ResolveUsername(chUser string) (*ChannelUser, error) {
	var data ChannelUser
	cmd := fmt.Sprintf("resolve_username %s", chUser)
	x, err := t.exec(cmd)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(x, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
func (t *telegram) ChannelInvite(channelId, user string) (*SuccessResp, error) {
	var data SuccessResp
	cmd := fmt.Sprintf("channel_invite %s %s", channelId, user)
	x, err := t.exec(cmd)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(x, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil

}

func (t *telegram) GetSelf() (*UserInfo, error) {
	var data UserInfo
	cmd := "get_self"
	x, err := t.exec(cmd)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(x, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
