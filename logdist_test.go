package logdist

import (
	"net/http"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	var ldh LogDistHandler = ""
	http.Handle("/stdout", ldh)
	go http.ListenAndServe("localhost:9000", nil)

	Message("", "first use of logdist.Message() in TestAll\n", true)
	time.Sleep(3 * time.Second)

	Message("logdist_test_log.log", "first use of logdist.Message with "+
		"file_path set\n", true)
	var fldh LogDistHandler = "logdist_test_log.log"
	http.Handle("/ltl", fldh)

	Message("logdist_test_log.log", "THIS MESSAGE SHOULDNT BE SEEN IN STDOUT\n",
		false)

	for {
		Message("logdist_test_log.log", "hi\n", true)
		time.Sleep(3 * time.Second)
	}
}
