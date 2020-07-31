package main

import (
	"github.com/google/uuid"
	"github.com/imonke/monkebase"
	"github.com/imonke/monketype"

	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
)

const (
	nick  = "reporter"
	email = "report@imonke.io"
)

var (
	reporter monketype.User
)

func mustMarshal(it interface{}) (data []byte) {
	var err error
	if data, err = json.Marshal(it); err != nil {
		panic(err)
	}

	return
}

func reportOK(reporter monketype.User, report monketype.Report) (err error) {
	if reporter.ID != report.Reporter {
		err = fmt.Errorf("ID mismatch! have: %s, want: %s", reporter.ID, report.Reporter)
		return
	}

	return
}

func TestMain(main *testing.M) {
	monkebase.Connect(os.Getenv("MONKEBASE_CONNECTION"))
	reporter = monketype.NewUser(nick, "", email)

	var err error
	if err = monkebase.WriteUser(reporter.Map()); err != nil {
		panic(err)
	}

	var result int = main.Run()
	monkebase.DeleteUser(reporter.ID)
	os.Exit(result)
}

func Test_createReport(test *testing.T) {
	var set []byte
	var sets [][]byte = [][]byte{
		mustMarshal(map[string]interface{}{
			"reported": uuid.New().String(),
			"type":     "user",
			"reason":   "smelly",
		}),
		mustMarshal(map[string]interface{}{
			"reported": uuid.New().String(),
			"type":     "user",
		}),
	}

	var code int
	var r_map map[string]interface{}
	var err error

	var request *http.Request
	var vauled context.Context

	for _, set = range sets {
		vauled = context.WithValue(context.TODO(), "requester", reporter.ID)
		if request, err = http.NewRequestWithContext(vauled, "POST", "/", bytes.NewReader(set)); err != nil {
			test.Fatal(err)
		}

		if code, r_map, err = createReport(request); err != nil {
			test.Fatal(err)
		}

		if code != 200 {
			test.Errorf("got code %d", code)
		}

		if err = reportOK(reporter, r_map["report"].(monketype.Report)); err != nil {
			test.Fatal(err)
		}

	}
}

func Test_createReport_badrequest(test *testing.T) {
	var set []byte
	var sets [][]byte = [][]byte{
		mustMarshal(map[string]interface{}{
			"type":   "user",
			"reason": "I dunno didn't like them. Wait who don't I like?",
		}),
		mustMarshal(map[string]interface{}{
			"reported": uuid.New().String(),
		}),
		[]byte("Gastrodon is my favorite pokemon!"),
		nil,
	}

	var code int
	var r_map map[string]interface{}
	var err error

	var request *http.Request
	var valued context.Context

	for _, set = range sets {
		valued = context.WithValue(context.TODO(), "requester", reporter.ID)
		if request, err = http.NewRequestWithContext(valued, "POST", "/", bytes.NewReader(set)); err != nil {
			test.Fatal(err)
		}

		if code, r_map, err = createReport(request); err != nil {
			test.Fatal(err)
		}

		if code != 400 {
			test.Errorf("got code %d", code)
		}

		if r_map != nil {
			test.Errorf("%#v", r_map)
		}
	}
}

func Test_getReport(test *testing.T) {
	var reporter string = uuid.New().String()
	var reported string = uuid.New().String()
	var report monketype.Report = monketype.NewReport(reporter, reported, "user", "")

	var err error
	if err = monkebase.WriteReport(report.Map()); err != nil {
		test.Fatal(err)
	}

	var request *http.Request
	if request, err = http.NewRequest("GET", "/id/"+report.ID, nil); err != nil {
		test.Fatal(err)
	}

	var code int
	var r_map map[string]interface{}
	if code, r_map, err = getReport(request); err != nil {
		test.Fatal(err)
	}

	if code != 200 {
		test.Errorf("got code %d", code)
	}

	var mapped map[string]interface{} = report.Map()
	var key string
	var value interface{}
	for key, value = range r_map["report"].(map[string]interface{}) {
		if value != mapped[key] {
			test.Errorf("mismatch at %s! got: %#v, want: %#v", key, value, mapped[key])
		}
	}
}

func Test_getReport_notfound(test *testing.T) {
	var request *http.Request
	var err error
	if request, err = http.NewRequest("GET", "/id/"+uuid.New().String(), nil); err != nil {
		test.Fatal(err)
	}

	var code int
	var r_map map[string]interface{}
	if code, r_map, err = getReport(request); err != nil {
		test.Fatal(err)
	}

	if code != 404 {
		test.Errorf("got code %d", code)
	}

	if r_map["error"].(string) != "no_such_report" {
		test.Errorf("%#v", r_map)
	}
}

func Test_ReportQueue(test *testing.T) {
	monkebase.EmptyTable(monkebase.REPORT_TABLE)

	var size, index int = 20, 0
	var reports []monketype.Report = make([]monketype.Report, size)
	var report monketype.Report
	var err error
	for index != size {
		report = monketype.NewReport(reporter.ID, uuid.New().String(), "user", "")
		if err = monkebase.WriteReport(report.Map()); err != nil {
			test.Fatal(err)
		}

		reports[index] = report
		index++
	}

	var offset, count int = 2, 16
	var valued context.Context = context.WithValue(
		context.TODO(),
		"parsed_query",
		map[string]int{"offset": offset, "size": count},
	)

	var request *http.Request
	if request, err = http.NewRequestWithContext(valued, "GET", "/new", nil); err != nil {
		test.Fatal(err)
	}

	var code int
	var r_map map[string]interface{}
	if code, r_map, err = getReportQueue(request); err != nil {
		test.Fatal(err)
	}

	if code != 200 {
		test.Errorf("got code %d", code)
	}

	var fetched_size int = r_map["size"].(map[string]int)["reports"]
	if fetched_size != count {
		test.Errorf("got size %d\n%#v", fetched_size, r_map["reports"])
	}

	var fetched monketype.Report
	for index, fetched = range r_map["reports"].([]monketype.Report) {
		if err = reportOK(reporter, fetched); err != nil {
			test.Fatal(err)
		}
	}
}