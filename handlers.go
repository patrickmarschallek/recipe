package recipe

//todo rename package name....
import (
	"io/ioutil"
	"net/http"
	"os"
)

// CustomHandlerFunc function type for recipes handlers
type CustomHandlerFunc func(w http.ResponseWriter, r *http.Request) (data interface{}, err error)

// ServeHTTP fullfills the http.Handler interface
func (f CustomHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

// IndexHandler shows the root endpoint / Welcome message
func IndexHandler(w http.ResponseWriter, r *http.Request) (data interface{}, err error) {
	return "Hello chef", nil
}

// HealthHandler shows service help
func HealthHandler(w http.ResponseWriter, r *http.Request) (data interface{}, err error) {
	// TODO write some code which reflects the service health. (cpu, mem ...)
	return Health{
		Status: 100,
	}, nil
}

// SwaggerHandler returns the swagger definition content
func SwaggerHandler(w http.ResponseWriter, r *http.Request) (data interface{}, err error) {
	fileName := "swagger"
	var extension string

	if _, err := os.Stat(fileName + ".yaml"); os.IsNotExist(err) {
		extension = "yaml"
	}
	if _, err := os.Stat(fileName + ".json"); os.IsNotExist(err) {
		extension = "json"
	}

	if len(extension) != 0 {
		b, err := ioutil.ReadFile(fileName + "." + extension) // just pass the file name
		if err != nil {
			return nil, NewErrorResponse(err.Error(), http.StatusInternalServerError)
		}
		return b, nil
	}
	return nil, NewErrorResponse("no swagger file defined.", http.StatusNotFound)
}
