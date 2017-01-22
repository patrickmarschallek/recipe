package recipe

// TODO add secure flag.
// Route maps the url path to a handler
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc CustomHandlerFunc
}

// Routes a list of route
type Routes []Route

// BaseRoutes basic routes which should be included by every service.
var BaseRoutes = Routes{
	Route{"Index", "GET", "/", IndexHandler},
	Route{"Health", "GET", "/health", HealthHandler},
}

// AddRoute add a route to the route list.
func (r Routes) AddRoute(route Route) Routes {
	return Routes(append(r, route))
}

// AddRoutes add a list of routes to the route list.
func (r Routes) AddRoutes(routes Routes) Routes {
	return Routes(append(r, routes[:]...))
}
