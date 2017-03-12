package bot

import (
	"modules/misc/t9n"
	"modules/misc/trans"
)

// ClearTrans is the command for
type ClearTrans struct {
}

// GetTopic return this message topic
func (ClearTrans) GetTopic() string {
	return "cy.rubik.clearTrans"
}

// GetQueue is the request queue
func (ClearTrans) GetQueue() string {
	return "cy_rubik_clearTrans"
}

// ClearTransAction worker
func ClearTransAction(in *ClearTrans) (bool, error) {
	trans.Clear()
	t9n.Clear()
	return false, nil
}
