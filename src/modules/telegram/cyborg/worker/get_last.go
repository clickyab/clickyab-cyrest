package worker

import (
	"common/assert"
	"common/redis"
	"encoding/json"
	"modules/telegram/cyborg/bot"
	"modules/telegram/cyborg/commands"
	"time"
)

func (mw *MultiWorker) getLast(in *commands.GetLastCommand) (bool, error) {
	// first try to resolve the channel
	m := bot.NewBotManager()
	c, err := m.FindKnownChannelByName(in.Channel)
	if err != nil {
		ch, err := mw.discoverChannel(in.Channel)
		if err != nil {
			// Oh crap. can not resolve this :/
			assert.Nil(aredis.StoreHashKey(in.HashKey, "STATUS", "failed", time.Hour))
			return false, err
		}

		c, err = bot.NewBotManager().CreateChannelByRawData(ch)
		if err != nil {
			assert.Nil(aredis.StoreHashKey(in.HashKey, "STATUS", "failed", time.Hour))
			return false, err
		}
	}
	h, err := mw.getLastMessages(c.CliTelegramID, 99, 0)
	if len(h) > in.Count {
		h = h[:in.Count]
	}
	if err != nil {
		assert.Nil(aredis.StoreHashKey(in.HashKey, "STATUS", "failed", time.Hour))
		return false, err
	}
	b, err := json.Marshal(h)
	assert.Nil(err)
	assert.Nil(aredis.StoreHashKey(in.HashKey, "DATA", string(b), time.Hour))
	assert.Nil(aredis.StoreHashKey(in.HashKey, "STATUS", "done", time.Hour))
	return false, nil
}
