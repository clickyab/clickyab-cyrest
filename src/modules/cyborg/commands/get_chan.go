package commands

// GetLastCommand is the command for last items and verify if they are available in the list or not
// also update the view in database
type GetChanCommand struct {
	// The channel name
	ChannelID int64
	// the count
	Count int
}

// GetTopic return this message topic
func (GetChanCommand) GetTopic() string {
	return "cy.rubik.getchan"
}

// GetQueue is the request queue
func (GetChanCommand) GetQueue() string {
	return "cy_rubik_getchan"
}
