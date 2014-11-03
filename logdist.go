package logdist

import (
	"bytes"
	"fmt"
	"github.com/dankozitza/seestack"
	"github.com/dankozitza/shiftlist"
	"net/http"
	"log"
	"os"
)

type LogDist struct {
	Log *log.Logger
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

func Message(file_path string, msg string) {

	// distribute the message using various methods

	if file_path == "" {
		file_path = "stdout"
	}

	if _, ok := logs[file_path]; !ok {
		fmt.Println(seestack.Short(), file_path, logs[file_path], ok)

		fo, err := os.Create(file_path)
		if err != nil {
			panic(seestack.Short() + ": " + err.Error())
		}
		defer func() {
			if err := fo.Close(); err != nil {
				panic(seestack.Short() + ": " + err.Error())
			}
		}()
		logs[file_path] = &LogDist{
			log.New(fo, "", 0),
			shiftlist.New(default_MaxIndex)}
	}

	logs[file_path].Log.Print(msg)
	logs[file_path].Tail.Add(msg)

	if (file_path != "stdout") {
		logs["stdout"].Log.Print(msg)
		logs["stdout"].Tail.Add(msg)
	}

	return
}

type LogDistHandler string

func (l LogDistHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if l == "" {
		l = "stdout"
	}

	for i := 0; i < logs[string(l)].Tail.NumEntries; i++ {
		fmt.Fprint(w, logs[string(l)].Tail.Get(i))
	}
	//fmt.Fprint(w, logs[string(l)].Tail)
}