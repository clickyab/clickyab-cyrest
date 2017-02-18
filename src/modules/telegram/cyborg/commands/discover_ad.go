package commands

// DiscoverAd is the command to discover an ad in channel
type DiscoverAd struct {
	// ChannelID is the channel id to use
	Channel int64
	ChatID  int64
	Reshot  bool
}

// GetTopic return this message topic
func (DiscoverAd) GetTopic() string {
	return "cy.rubik.discover"
}

// GetQueue is the request queue
func (DiscoverAd) GetQueue() string {
	return "cy_rubik_discover"
}
