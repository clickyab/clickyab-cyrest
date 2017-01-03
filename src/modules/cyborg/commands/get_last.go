package commands

// GetLastCommand is the command for last items and verify if they are available in the list or not
// also update the view in database
type GetLastCommand struct {
	// The channel name
	Channel string
	// The Hash key to store data in
	HashKey string

	Count int
}

// GetTopic return this message topic
func (GetLastCommand) GetTopic() string {
	return "cy.rubik.getlast"
}

// GetQueue is the request queue
func (GetLastCommand) GetQueue() string {
	return "cy_rubik_getlast"
}
