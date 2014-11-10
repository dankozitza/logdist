package logdist

import (
   "fmt"
   "syscall"
	"net/http"
	"testing"
   "os"
   "io"
)

var file_path string = "logdist_test_log.log"

func TestAll(t *testing.T) {
	var ldh LogDistHandler
	http.Handle("/logdist", ldh)
	//go http.ListenAndServe("localhost:9000", nil)

	Message("", "first use of logdist.Message() in TestAll\n", true)

	Message(file_path, "first use of logdist.Message with "+
		"file_path set\n", true)
	var fldh LogDistHandler = LogDistHandler(file_path)
	http.Handle("/logdist2", fldh)

	Message(file_path, "THIS MESSAGE SHOULDNT BE SEEN IN STDOUT\n",
		false)

	//for {
	//	Message(file_path, "hi\n", true)
	//	time.Sleep(3 * time.Second)
	//}
}

func TestFile(t *testing.T) {

   dummyfile := "first use of logdist.Message with file_path set\n" +
      "THIS MESSAGE SHOULDNT BE SEEN IN STDOUT\n"

   fi, err := os.Open(file_path)
   if err != nil {
      fmt.Println("TestFave: could not open saved log file:", file_path)
      t.Fail()
      return
   }

   buff := make([]byte, 1024)

   n, err := fi.Read(buff)
   if err != nil && err != io.EOF {
      fmt.Println("TestSave: could not read from config file:", file_path,
         "err:", err)
      t.Fail()
      return
   }
   if string(buff[:n]) != dummyfile {
      fmt.Println("TestSave: config file does not match dummy file:")
      fmt.Println("saved file:", string(buff[:n]))
      fmt.Println("dummy file:", dummyfile)
      t.Fail()
   }
}


func TestClean(t *testing.T) {
   fmt.Println("TestClean: removing", file_path)
   syscall.Exec("/usr/bin/rm", []string{"rm", file_path}, nil)
}
