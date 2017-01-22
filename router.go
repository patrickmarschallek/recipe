package recipe

import "github.com/gorilla/mux"

// TODO add decorator list methods, probably add a struct
// add route and position/order number
// var Decorators = []CustomHandlerFunc{
// 	ExecutionTimeHandler,
// 	FlowHandler,
// 	CustomHandler,
// }

// Router handles all defined routes and performs the handler chain
func NewRouter(basePath string) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range BaseRoutes {
		var cHandler CustomHandlerFunc

		// middleware decorators
		cHandler = route.HandlerFunc
		cHandler = ExecutionTimeHandler(cHandler, route.Name)
		cHandler = FlowHandler(cHandler)

		router.
			Methods(route.Method).
			Path(basePath + route.Pattern).
			Name(route.Name).
			Handler(CustomHandler(cHandler))
	}
	return router
}
