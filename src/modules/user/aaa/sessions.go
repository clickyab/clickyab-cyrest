package aaa

import (
	"common/redis"
	"fmt"
	"time"

	"github.com/mssola/user_agent"
)

// Session is the active session for a user
type Session struct {
	Value       string        `json:"-"`
	Current     bool          `json:"current"`
	Key         string        `json:"key"`
	OS          string        `json:"os"`
	Browser     string        `json:"browser"`
	Version     string        `json:"version"`
	IP          string        `json:"ip"`
	ExpireAfter time.Duration `json:"expire_after"`
	CreatedAt   time.Time     `json:"created_at"`
}

// Sessions is slice of session
type Sessions []Session

func getSession(st string) Session {
	// TODO : HGETALL
	s := Session{
		Key: st,
	}
	s.Value, _ = aredis.GetHashKey(s.Key, "token", false, 0)
	agent, _ := aredis.GetHashKey(s.Key, "ua", false, 0)
	ua := user_agent.New(agent)
	browser, version := ua.Browser()
	s.Browser = browser
	s.Version = version
	s.OS = ua.OS()
	s.IP, _ = aredis.GetHashKey(s.Key, "ip", false, 0)
	tmp, _ := aredis.GetHashKey(s.Key, "date", false, 0)
	s.CreatedAt, _ = time.Parse(time.RFC3339, tmp)
	s.ExpireAfter, _ = aredis.GetExpire(s.Key)

	return s
}

// GetSessions return active sessions for the current user
func (m *Manager) GetSessions(u *User, current string, count int64) (Sessions, bool) {
	res := make(Sessions, 1)
	res[0] = getSession(current)
	res[0].Current = true
	scanner := aredis.Client.Scan(0, fmt.Sprintf("%d:*", u.ID), 0).Iterator()
	var cnt int64
	for scanner.Next() {
		if scanner.Val() == current {
			continue
		}
		s := getSession(scanner.Val())
		if s.Value == res[0].Value {
			res = append(res, s)
			cnt++
			if cnt >= count {
				break
			}
		}
	}

	return res, scanner.Next()
}
