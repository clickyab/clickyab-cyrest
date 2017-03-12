package trans

import (
	"fmt"
	"modules/misc/t9n"
	"sync"

	"github.com/Sirupsen/logrus"
)

var (
	translations = make(map[string]map[string]string)
	lock         = &sync.RWMutex{}
	_            = new(baseTranslated) // Make the stupid unused stop the warning
)

const (
	// DefaultLang is the default application language
	DefaultLang = "en_US"
	// PersianLang PersianLang
	PersianLang = "fa_IR"
	// EnglishLang EnglishLang
	EnglishLang = "en_US"
)

type baseTranslated interface {
	// Text return the unformatted text
	GetText() string
	// Params is the parameters
	GetParams() []interface{}
	// Translate is the actual translation
	Translate(string) string
}

// Translated is the interface to handle the translation
type Translated interface {
	fmt.Stringer
	baseTranslated
}

// TranslatedError is the error type translator
type TranslatedError interface {
	error
	baseTranslated
}

type t9Base struct {
	Text   string        `json:"text"`
	Params []interface{} `json:"params"`
}

// T9String is the Translation string
type T9String struct {
	t9Base
}

// T9Error is translation error
type T9Error struct {
	t9Base
}

func (t9 t9Base) GetText() string {
	return t9.Text
}

func (t9 t9Base) GetParams() []interface{} {
	return t9.Params
}

func (t9 t9Base) Translate(lang string) string {
	lock.RLock()
	t, ok := translations[lang]
	if !ok {
		lock.RUnlock()
		lock.Lock()
		translations[lang] = t9n.NewT9nManager().LoadAllInMap(lang)
		t = translations[lang]
		lock.Unlock()
		lock.RLock()
	}
	tmp := t[t9.Text]
	return fmt.Sprintf(tmp, t9.Params...)
}

func (t9 T9String) String() string {
	return fmt.Sprintf(t9.GetText(), t9.GetParams()...)
}

func (t9 T9Error) Error() string {
	return fmt.Sprintf(t9.GetText(), t9.GetParams()...)
}

// T is the universal translate function
func T(translationID string, args ...interface{}) (res T9String) {
	defer func() {
		if e := recover(); e != nil {
			logrus.Warn("Use trans module before Initialize models? maybe in global scope or init")
			res = T9String{
				t9Base{
					Text:   translationID,
					Params: args,
				},
			}
		}
	}()
	lock.RLock()

	if translations[DefaultLang] == nil {
		lock.RUnlock()
		lock.Lock()
		m := t9n.NewT9nManager()
		translations[DefaultLang] = m.LoadAllInMap(DefaultLang)
		lock.Unlock()
		lock.RLock()
	}

	_, ok := translations[DefaultLang][translationID]
	lock.RUnlock()
	if !ok {
		var err error
		lock.Lock()
		m := t9n.NewT9nManager()
		err = m.AddMissing(translationID)
		if err == nil {
			logrus.Debugf("[ADD TRANSLATION] %s", translationID)
			translations[DefaultLang][translationID] = translationID
		}
		lock.Unlock()
	}

	return T9String{
		t9Base{
			Text:   translationID,
			Params: args,
		},
	}
}

// E is the error version of the T
func E(translationID string, args ...interface{}) T9Error {
	text := T(translationID, args...)
	return T9Error{
		t9Base: text.t9Base,
	}
}

// EE try to translate an already generated error
func EE(e error) T9Error {
	if e == nil {
		return T9Error{}
	}
	switch t9 := e.(type) {
	case T9Error:
		return t9
	default:
		return E(e.Error())
	}
}

// Clear Clear
func Clear() {
	lock.Lock()
	defer lock.Unlock()
	translations = make(map[string]map[string]string)

}
