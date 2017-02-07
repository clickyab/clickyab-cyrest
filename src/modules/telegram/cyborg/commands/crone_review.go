package commands

// CronReview is the command for last items and crone review if they are available in the list or not
// also update the view in database
type CronReview struct {
	// The channel name
	ChannelID int64
	// the count
	Count int
}

// GetTopic return this message topic
func (CronReview) GetTopic() string {
	return "cy.rubik.cronreview"
}

// GetQueue is the request queue
func (CronReview) GetQueue() string {
	return "cy_rubik_cronreview"
}
