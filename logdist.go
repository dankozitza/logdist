package logdist

import (
	"bytes"
	"fmt"
	"github.com/dankozitza/seestack"
	"github.com/dankozitza/shiftlist"
	"log"
	"net/http"
	"os"
	"regexp"
)

type LogDist struct {
	Log  *log.Logger
	Tail *shiftlist.ShiftList
}

var logs map[string]*LogDist = make(map[string]*LogDist)
var nologbuf bytes.Buffer
var default_MaxIndex int = 100

func init() {
	// logs is a map of LogDist objects that is mostly keyed by file paths
	logs["stdout"] = &LogDist{
		log.New(os.Stdout, "", 0),
		shiftlist.New(default_MaxIndex)}
}

func Message(file_path string, to_stdout bool, msg ...interface{}) {

	// distribute the message using various methods. currently print to stdout,
	// print to file, and store in shiftlist

	if file_path == "" {
		file_path = "stdout"
	}

	// create a new log object if we don't have one
	if _, ok := logs[file_path]; !ok {

		fo, err := os.Create(file_path)
		if err != nil {
			panic(seestack.Short() + ": " + err.Error())
		}
		logs[file_path] = &LogDist{
			log.New(fo, "", 0),
			shiftlist.New(default_MaxIndex)}
	}

	strmsg := fmt.Sprint(msg...)

	logs[file_path].Log.Print(strmsg)
	logs[file_path].Tail.Add(strmsg)

	if file_path != "stdout" && to_stdout {
		logs["stdout"].Log.Print(strmsg)
		logs["stdout"].Tail.Add(strmsg)
	}

	return
}

type HTTPHandler string

func (l HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if l == "" {
		l = "stdout"
	}

	if _, ok := logs[string(l)]; !ok {

		fo, err := os.Create(string(l))
		if err != nil {
			panic(seestack.Short() + ": " + err.Error())
		}
		logs[string(l)] = &LogDist{
			log.New(fo, "", 0),
			shiftlist.New(default_MaxIndex)}
	}

	// print the strings in logs[string(l)].Tail
	for i := 0; i < logs[string(l)].Tail.NumEntries; i++ {

		switch v := logs[string(l)].Tail.Get(i).(type) {

		case string:
			re, _ := regexp.Compile("\n")
			newline := re.ReplaceAllString(v, "<br>\n")
			fmt.Fprint(w, newline)

		default:
			// this should never happen. panic?
			fmt.Fprint(w, v)
			fmt.Fprint(w, "<br>\n")
		}
	}
}
