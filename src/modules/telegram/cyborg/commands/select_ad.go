package commands

// SelectAd is the command for channel ad
type SelectAd struct {
	// channel ID
	ChannelID int64
	// chat ID
	ChatID int64
}

// GetTopic return this message topic
func (SelectAd) GetTopic() string {
	return "cy.rubik.selectAd"
}

// GetQueue is the request queue
func (SelectAd) GetQueue() string {
	return "cy_rubik_selectAd"
}
