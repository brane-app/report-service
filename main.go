package main

import (
	"github.com/gastrodon/groudon"
	"git.gastrodon.io/imonke/monkebase"
	"git.gastrodon.io/imonke/monkelib/middleware"

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

	groudon.RegisterMiddlewareRoute([]string{"GET"}, "^/new/?", middleware.PaginationParams)
	groudon.RegisterMiddlewareRoute([]string{"GET", "PATCH"}, `^/(id|new).*$`, middleware.MustModerator)

	groudon.RegisterHandler("POST", `^/$`, createReport)
	groudon.RegisterHandler("GET", `^/new/?$`, getReportQueue)
	groudon.RegisterHandler("GET", `^/id/`+uuid_regex+`/?$`, getReport)
	groudon.RegisterHandler("PATCH", `^/id/`+uuid_regex+`/?$`, updateReport)

	http.Handle("/", http.HandlerFunc(groudon.Route))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
