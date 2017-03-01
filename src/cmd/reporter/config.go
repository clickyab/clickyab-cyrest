package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	path       string
	tagPattern = regexp.MustCompile(`(ref|fix)?\s*#([0-9]+)`)
	dev        bool
)

func config() {
	help()
	devmode()

}
func devmode() {
	_, e := argValue("test")
	if e == nil {
		fmt.Println("Test mode")
	}
}

func limit() int {
	count := 100
	cs, e := argValue("l")
	if e != nil || cs == "true" {
		return count
	}
	c, e := strconv.Atoi(cs)
	if e != nil || c < 1 {
		fError(errors.New("Wrong value"), "%v is not a valid value for -l parameter! it sould be integer and bigger then 0.", c)
	}
	return c
}

func redmineEndpoint() (string, error) {
	v, e := argValue("re")
	if e != nil {
		return "", e
	}
	return v, nil
}

func redmineApikey() (string, error) {
	v, e := argValue("ra")
	if e != nil {
		return "", e
	}
	return v, nil
}

func help() {
	_, e := argBool("h")
	if e == nil {
		println(HELP)
		os.Exit(0)
	}
	return

}

func argBool(key string) (bool, error) {
	v, e := argValue(key)
	if e != nil {
		return false, e
	}
	if b, e := strconv.ParseBool(v); e == nil {
		return b, e
	}
	fError(errors.New("Wrong value"), "Wrong value has been set for '%s'!! it's sould be true or false.", keyer(key))
	return false, nil
}

func keyer(key string) string {
	var prefix string
	if len(key) == 1 {
		prefix = "-"
	} else {
		prefix = "--"
	}
	return prefix + key
}

func getPath() string {
	dir, err := argValue("i")
	fError(err, "you missed the -i flag!!")
	path, _ := filepath.Abs(dir)
	res, err := os.Stat(path)
	if err != nil || !res.IsDir() {
		fError(err, "There is o directory at %s. use -h for more information", path)
	}

	return path

}

func argValue(key string) (string, error) {

	for i, a := range os.Args {
		if !strings.HasPrefix(a, keyer(key)) {
			continue
		}
		if keyer(key) == a {
			if len(os.Args) < i+2 {
				return "true", nil
			}
			if strings.HasPrefix(os.Args[i+1], "-") {
				return "true", nil
			}
			return os.Args[i+1], nil
		}
	}
	return "", errors.New("flag does not exists")
}

func fError(err interface{}, message ...interface{}) {
	if err != nil {
		fmt.Println(fmt.Sprintf(message[0].(string), message[1:]...))
		os.Exit(1)
	}
}
