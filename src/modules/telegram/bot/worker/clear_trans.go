package bot

import "modules/misc/trans"

// ClearTrans is the command for
type ClearTrans struct {
}

// GetTopic return this message topic
func (ClearTrans) GetTopic() string {
	return "cy.rubik.sendWarn"
}

// GetQueue is the request queue
func (ClearTrans) GetQueue() string {
	return "cy_rubik_sendWarn"
}

// ClearTransAction worker
func ClearTransAction(in *ClearTrans) (bool, error) {
	trans.Clear()
	return false, nil
}
