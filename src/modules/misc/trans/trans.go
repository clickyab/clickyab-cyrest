package trans

import (
	"errors"
	"modules/misc/t9n"
	"sync"

	"fmt"

	"github.com/Sirupsen/logrus"
)

var (
	translations map[string]t9n.Translation
	lock         = &sync.RWMutex{}
)

// T is the translate function
func T(translationID string, args ...interface{}) string {
	return translationID

	lock.RLock()
	if translations == nil {
		lock.RUnlock()
		lock.Lock()
		m := t9n.NewT9nManager()
		translations = m.LoadAllInMap()
		lock.Unlock()
		lock.RLock()
	}

	tt, ok := translations[translationID]
	lock.RUnlock()
	if !ok {
		var err error
		lock.Lock()
		m := t9n.NewT9nManager()
		tt, err = m.AddMissing(translationID)
		if err == nil {
			translations[translationID] = tt
		}
		lock.Unlock()
	}

	if tt.Single.Valid {
		return fmt.Sprintf(tt.Single.String, args...)
	}
	logrus.Debugf("NOT TRANSLATED : %s ", translationID)
	return fmt.Sprintf(translationID, args...)
}

// E is the error version of the T
func E(translationID string, args ...interface{}) error {
	return errors.New(T(translationID, args...))
}
