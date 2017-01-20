package commands

// IdentifyAD is the command for forward message
type UpdateMessage struct {
	// The ad ID
	CLiChannelID string
}

// GetTopic return this message topic
func (UpdateMessage) GetTopic() string {
	return "cy.rubik.updateMessage"
}

// GetQueue is the request queue
func (UpdateMessage) GetQueue() string {
	return "cy_rubik_updateMessageD"
}
