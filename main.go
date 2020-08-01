package main

import (
	"github.com/gastrodon/groudon"
	"github.com/imonke/monkebase"
	"github.com/imonke/monkelib/middleware"

	"log"
	"net/http"
	"os"
)

const (
	uuid_regex = `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`
)

var (
	forbidden = map[string]interface{}{"error": "forbidden"}
)

func main() {
	monkebase.Connect(os.Getenv("MONKEBASE_CONNECTION"))

	groudon.RegisterCatch(403, forbidden)
	groudon.RegisterMiddleware(middleware.MustAuth)
	groudon.RegisterMiddleware(middleware.RangeQueryParams)

	groudon.RegisterMiddlewareRoute([]string{"GET", "PATCH"}, `^/(id|new).*$`, MustModerator)

	groudon.RegisterHandler("POST", `^/$`, createReport)
	groudon.RegisterHandler("GET", `^/new/?$`, getReportQueue)
	groudon.RegisterHandler("GET", `^/id/`+uuid_regex+`/?$`, getReport)
	// groudon.RegisterHandler("PATCH", `^/id/`+uuid_regex+`/?$`, updateReport)
	// groudon.RegisterHandler("GET", `^/user/`+uuid_regex+`/?$`, getUserReports)

	http.Handle("/", http.HandlerFunc(groudon.Route))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
