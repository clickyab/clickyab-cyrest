package commands

// ExistChannelAd is the command for channel ad
type ExistChannelAd struct {
	// channel ID
	ChannelID int64
	// ad ID
	AdID []int64
	// chat ID
	ChatID int64
}

// GetTopic return this message topic
func (ExistChannelAd) GetTopic() string {
	return "cy.rubik.existChannelAd"
}

// GetQueue is the request queue
func (ExistChannelAd) GetQueue() string {
	return "cy_rubik_existChannelAd"
}
