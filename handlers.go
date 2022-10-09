package main

import (
	"github.com/brane-app/librane/tools/middleware"
	"github.com/gastrodon/groudon/v2"

	"os"
)

const (
	uuid_regex = `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`
)

var (
	prefix = os.Getenv("PATH_PREFIX")

	routeAny                 = "^" + prefix + "/?.*"
	routeRoot                = "^" + prefix + "/?$"
	routeNew                 = "^" + prefix + "/new/?$"
	routeId                  = "^" + prefix + "/id/" + uuid_regex + "/?$"
	routeModeratorRestricted = "^" + prefix + "/(id|new).*$"

	forbidden = map[string]interface{}{"error": "forbidden"}
)

func register_handlers() {
	groudon.AddCodeResponse(403, forbidden)

	groudon.AddMiddleware("GET", routeAny, middleware.MustAuth)
	groudon.AddMiddleware("POST", routeAny, middleware.MustAuth)
	groudon.AddMiddleware("PATCH", routeAny, middleware.MustAuth)

	groudon.AddMiddleware("GET", routeNew, middleware.PaginationParams)
	groudon.AddMiddleware("GET", routeModeratorRestricted, middleware.MustModerator)
	groudon.AddMiddleware("PATCH", routeModeratorRestricted, middleware.MustModerator)

	groudon.AddHandler("POST", routeRoot, createReport)
	groudon.AddHandler("GET", routeNew, getReportQueue)
	groudon.AddHandler("GET", routeId, getReport)
	groudon.AddHandler("PATCH", routeId, updateReport)
}
