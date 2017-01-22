package bot

import (
	"errors"
	"modules/telegram/common/tgo"
	"time"
)

// KnownChannel is the list of all known channels for cyborg
// @Model {
//		table = known_channels
//		primary = true, id
//		find_by = id, name
// }
type KnownChannel struct {
	ID            int64            `db:"id" json:"id"`
	Name          string           `db:"name" json:"name"`
	Title         string           `db:"title" json:"title"`
	Info          string           `db:"info" json:"info"`
	UserCount     int64            `db:"user_count" json:"user_count"`
	CliTelegramID string           `db:"cli_telegram_id" json:"cli_telegram_id"`
	RawData       *tgo.ChannelInfo `db:"raw_data" json:"raw_data"`
	CreatedAt     time.Time        `db:"created_at" json:"created_at" sort:"true"`
	UpdatedAt     time.Time        `db:"updated_at" json:"updated_at" sort:"true"`
}

// CreateChannelByRawData try to create a record per channel for all channele we visiting
func (m *Manager) CreateChannelByRawData(c *tgo.ChannelInfo) (*KnownChannel, error) {
	if c.PeerType != "channel" {
		return nil, errors.New("this is not a channel")
	}

	kc := &KnownChannel{
		Name:          c.Username,
		Title:         c.PrintName,
		Info:          c.About,
		UserCount:     int64(c.ParticipantsCount),
		CliTelegramID: c.ID,
		RawData:       c,
	}

	err := m.CreateKnownChannel(kc)
	if err != nil {
		return nil, err
	}

	return kc, nil
}
