package cmangos

import (
	"io/ioutil"
	golog "log"
)

var (
	log *golog.Logger
)

func init() {
	// discard all logs by default
	log = golog.New(ioutil.Discard, "", 0)
}