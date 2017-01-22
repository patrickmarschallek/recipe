package recipe

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/context"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/uber-go/zap"
)

// CustomHandler starts the internal handlers and writes
// the return to the standard put in given format.
func CustomHandler(fn CustomHandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Logger.Debug("start CustomHandler",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
		)
		data, err := fn(w, r)
		Logger.Debug("inner handler returned",
			zap.Bool("data", data != nil),
			zap.Bool("err", err != nil),
		)
		flowID := context.Get(r, "flowID").(string)
		if err != nil {
			Logger.Error("write error response",
				zap.Error(err),
			)
			writeResponse(w, flowID, err, r.URL.Query().Get("type"))
		}
		if data != nil {
			writeResponse(w, flowID, data, r.URL.Query().Get("type"))
		}
	})
}

// FlowHandler adds a flow id to the context to identify one request flow.
func FlowHandler(inner CustomHandlerFunc) CustomHandlerFunc {
	return CustomHandlerFunc(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		Logger.Debug("Set flowID")
		flowID, err := uuid.NewV4()
		if err != nil {
			Logger.Fatal("UUID generating error",
				zap.Error(err),
			)
		}
		context.Set(r, "flowID", flowID.String())
		// set flowID to the logger context
		Logger = Logger.With(
			zap.String("flowID", flowID.String()),
		)
		Logger.Debug("flowID was set")
		return inner(w, r)
	})
}

// ExecutionTimeHandler logs the execution time for this request.
func ExecutionTimeHandler(inner CustomHandlerFunc, name string) CustomHandlerFunc {
	return CustomHandlerFunc(func(w http.ResponseWriter, r *http.Request) (data interface{}, err error) {
		start := time.Now()
		data, err = inner(w, r)

		duration := time.Since(start)
		Logger.Info("request time execution",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.String("name", name),
			zap.Duration("execution_time", duration),
		)
		return
	})
}

// oAuthHandler TODO check out clientToken etc..
func oAuthHandler(inner CustomHandlerFunc) CustomHandlerFunc {
	return CustomHandlerFunc(func(w http.ResponseWriter, r *http.Request) (data interface{}, err error) {
		if len(r.URL.Query().Get("key")) == 0 {
			return nil, NewErrorResponse("missing key", http.StatusUnauthorized)
		}
		return inner(w, r)
	})
}

// //-- mandetory params
// func MustParams(fn http.HandlerFunc, params ...string) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		for _, param := range params {
// 			if len(r.URL.Query().Get(param)) == 0 {
// 				http.Error(w, "missing "+param, http.StatusBadRequest)
// 				return
// 			}
// 		}
// 		fn(w, r) // success - call handler
// 	}
// }

// //-- validation
// func Validate(inner http.Handler, name string, obj *interface{}) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		decoder := json.NewDecoder(r.Body)
// 		err := decoder.Decode(&obj)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		context.Set(r, "validatedInput", *obj)
// 		inner.ServeHTTP(w, r)
// 	})
// }

func writeResponse(w http.ResponseWriter, flowID string, data interface{}, ct string) {
	var err error
	var payload []byte

	switch ct {
	case "xml":
		payload, err = xml.Marshal(data)
	default:
		payload, err = json.Marshal(data)
	}
	if err != nil {
		Logger.Info("error during response writing",
			zap.Error(err),
		)
	}
	Logger.Debug("write response",
		zap.String("type", ct),
		zap.String("payload", string(payload)),
	)
	w.Header().Set("Content-Type", fmt.Sprintf("application/%s; charset=utf-8", ct))
	w.Write(payload)
}
