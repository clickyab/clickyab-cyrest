package commands

// IdentifyAD is the command for forward message
type IdentifyAD struct {
	// The ad ID
	AddID int64
}

// GetTopic return this message topic
func (IdentifyAD) GetTopic() string {
	return "cy.rubik.identifyAD"
}

// GetQueue is the request queue
func (IdentifyAD) GetQueue() string {
	return "cy_rubik_identifyAD"
}
