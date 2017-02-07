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

// Message type const
const Message string = "message"

// Photo type const
const Photo string = "photo"

// TelegramCli is the interface to handle the telegram cli
type TelegramCli interface {

	//ChannelInfo channel_info <channel>  Prints info about channel (id, members, admin, etc.)
	ChannelInfo(channelID string) (*ChannelInfo, error)

	//AddContact add_contact <phone> <first name> <last name>    Tries to add user to contact list
	AddContact(phone, firstName, lastName string) (*AddContact, error)

	//ContactList contact_list()
	ContactList() ([]Contact, error)

	//History history <peer> [limit] [offset] Prints messages with this peer
	// (most recent message lower). Also marks messages as read
	History(peer string, limit, offset int) ([]History, error)

	// history of just event with message
	MessageHistory(peer string, limit, offset int) ([]History, error)

	// FwdMsg fwd <peer> <msg-id>+    Forwards message to peer. Forward to secret chats is forbidden
	FwdMsg(peer string, msg string) (*SuccessResp, error)

	// UserInfo user_info <user>        Prints info about user (id, last online, phone)
	UserInfo(user string) (*UserInfo, error)

	// ChannelJoin channel_join <channel>  Joins to channel
	ChannelJoin(channelID string) (*SuccessResp, error)

	//Msg msg <peer> <kbd> <text> Sends text message to peer with custom kbd
	Msg(peer string, msg string) (*SuccessResp, error)

	//Post post <peer> <text>      Sends text message to peer as admin
	Post(peer string, msg string) (*SuccessResp, error)

	//ChannelList channel_list [limit=100] [offset=0]     List of last channels
	ChannelList(limit, offset int) ([]ChannelList, error)

	//ResolveUsername resolve_username username       Searches user by username
	ResolveUsername(chUser string) (*ChannelUser, error)

	//ChannelInvite channel_invite <channel> <user> Invites user to channel
	ChannelInvite(channelID, user string) (*SuccessResp, error)

	//GetSelf get_self        Get our user info
	GetSelf() (*UserInfo, error)

	// Close the connection
	Close() error
}

type telegram struct {
	lock *sync.Mutex
	Conn *net.TCPConn
	IP   net.IP
	Port int
}

// NewTelegramCli return a new instance of telegram cli
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

func (t *telegram) Close() error {
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

func (t *telegram) ChannelInfo(channelID string) (*ChannelInfo, error) {

	var data ChannelInfo
	cmd := fmt.Sprintf("channel_info %s", channelID)
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
	if err != nil {
		return nil, err
	}
	fmt.Println(string(x))
	err = json.Unmarshal(x, &data)
	if err != nil {
		return nil, err
	}
	return data, nil

}

func (t *telegram) MessageHistory(user string, limit, offset int) ([]History, error) {
	var data []History
	var res []History

	cmd := fmt.Sprintf("history %s %d %d", user, limit, offset)
	x, err := t.exec(cmd)

	if err != nil {
		return nil, err
	}
	fmt.Println(string(x))
	err = json.Unmarshal(x, &data)
	if err != nil {
		return nil, err
	}

	for i := range data {
		if data[i].Event == Message {
			res = append(res, data[i])
		}
	}

	return res, nil
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
func (t *telegram) ChannelJoin(channelID string) (*SuccessResp, error) {

	var data SuccessResp
	cmd := fmt.Sprintf("channel_join %s", channelID)
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

func (t *telegram) ChannelInvite(channelID, user string) (*SuccessResp, error) {
	var data SuccessResp
	cmd := fmt.Sprintf("channel_invite %s %s", channelID, user)
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
