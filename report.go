package main

import (
	"github.com/brane-app/librane/database"
	"github.com/brane-app/librane/tools"
	"github.com/brane-app/librane/types"
	"github.com/gastrodon/groudon/v2"

	"net/http"
)

func getReportQueue(request *http.Request) (code int, r_map map[string]interface{}, err error) {
	var query map[string]interface{} = request.Context().Value("query").(map[string]interface{})
	var before string = query["before"].(string)
	var size int = query["size"].(int)
	var reports []types.Report
	if reports, size, err = database.ReadManyUnresolvedReport(before, size); err != nil {
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
	var parts []string = tools.SplitPath(request.URL.Path)

	var report types.Report
	var exists bool
	if report, exists, err = database.ReadSingleReport(parts[len(parts)-1]); err != nil {
		return
	}

	if !exists {
		code = 404
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
	var report types.Report = types.NewReport(
		reporter, body.Reported,
		body.Type,
		body.Reason,
	)

	err = database.WriteReport(report.Map())
	code = 200
	r_map = map[string]interface{}{"report": report}
	return
}

func patchReport(body PatchReportBody, report types.Report) (patched types.Report, changed bool) {
	patched = report

	if body.Resolved != nil {
		changed = true
		patched.Resolved = *body.Resolved
	}

	if body.Resolution != nil && *body.Resolution != "" {
		changed = true
		patched.Resolution = *body.Resolution
	}

	return
}

func updateReport(request *http.Request) (code int, r_map map[string]interface{}, err error) {
	var body PatchReportBody
	var external error
	if err, external = groudon.SerializeBody(request.Body, &body); err != nil || external != nil {
		code = 400
		return
	}

	var parts []string = tools.SplitPath(request.URL.Path)
	var id string = parts[len(parts)-1]

	var report types.Report
	var ok bool
	if report, ok, err = database.ReadSingleReport(id); !ok || err != nil {
		code = 404
		r_map = map[string]interface{}{"error": "no_such_report"}
		return
	}

	if report, ok = patchReport(body, report); !ok {
		code = 400
		return
	}

	if err = database.WriteReport(report.Map()); err != nil {
		panic(err)
		return
	}

	code = 200
	r_map = map[string]interface{}{"report": report}
	return
}
