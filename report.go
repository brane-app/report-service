package main

import (
	"github.com/gastrodon/groudon"
	"github.com/imonke/monkebase"
	"github.com/imonke/monketype"

	"net/http"
	"strings"
)

func pathSplit(it rune) (ok bool) {
	ok = it == '/'
	return
}

func getReportQueue(request *http.Request) (code int, r_map map[string]interface{}, err error) {
	var query map[string]int = request.Context().Value("parsed_query").(map[string]int)

	var reports []monketype.Report
	var size int
	if reports, size, err = monkebase.ReadManyUnresolvedReport(query["offset"], query["size"]); err != nil {
		panic(err)
		return
	}

	code = 200
	r_map = map[string]interface{}{
		"reports": reports,
		"size":    map[string]int{"reports": size},
	}
	return
}

func getReport(request *http.Request) (code int, r_map map[string]interface{}, err error) {
	var parts []string = strings.FieldsFunc(request.URL.Path, pathSplit)

	var report monketype.Report
	var exists bool
	if report, exists, err = monkebase.ReadSingleReport(parts[len(parts)-1]); err != nil {
		return
	}

	if !exists {
		code = 400
		r_map = map[string]interface{}{"error": "no_such_report"}
		return
	}

	code = 200
	r_map = map[string]interface{}{"report": report.Map()}
	return
}

func createReport(request *http.Request) (code int, r_map map[string]interface{}, err error) {
	var body CreateReportBody
	var external error
	if err, external = groudon.SerializeBody(request.Body, &body); err != nil || external != nil {
		code = 400
		return
	}

	var reporter string = request.Context().Value("requester").(string)
	var report monketype.Report = monketype.NewReport(
		reporter, body.Reported,
		body.Type,
		body.Reason,
	)

	err = monkebase.WriteReport(report.Map())
	code = 200
	r_map = map[string]interface{}{"report": report}
	return
}