package utils

import (
	"common/assert"
	"io/ioutil"
	"os"
	"strings"
)

// ChangeInFile is the function for replace a string in file
func ChangeInFile(path string, from, to string) error {
	o, err := os.Open(path)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(o)
	if err != nil {
		return err
	}
	assert.Nil(o.Close())
	target := strings.Replace(string(data), from, to, -1)
	o, err = os.Create(path)
	if err != nil {
		return err
	}

	_, err = o.WriteString(target)
	return err

}
