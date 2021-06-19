// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/rs/cors"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	"example.com/open-api/models"
	operations "example.com/open-api/restapi/operations"
	todos "example.com/open-api/restapi/operations/todos"
)

//go:generate swagger generate server --target ../../open-api --name TodoList --spec ../swagger.yml --principal interface{}

func configureFlags(api *operations.TodoListAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

// ここから追加
var items = make(map[int64]*models.Item)
var lastID int64

var itemsLock = &sync.Mutex{}

func newItemID() int64 {
	return atomic.AddInt64(&lastID, 1)
}

func addItem(item *models.Item) error {
	if item == nil {
		return errors.New(500, "item must be present")
	}

	itemsLock.Lock()
	defer itemsLock.Unlock()

	newID := newItemID()
	item.ID = newID
	items[newID] = item

	return nil
}

func allItems(since int64, limit int32) (result []*models.Item) {
	result = make([]*models.Item, 0)
	for id, item := range items {
		if len(result) >= int(limit) {
			return
		}
		if since == 0 || id > since {
			result = append(result, item)
		}
	}
	return result
}

// ここまで追加

func configureAPI(api *operations.TodoListAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// ここから追加
	api.TodosAddOneHandler = todos.AddOneHandlerFunc(func(params todos.AddOneParams) middleware.Responder {
		if err := addItem(params.Body); err != nil {
			return todos.NewAddOneDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}
		return todos.NewAddOneCreated().WithPayload(params.Body)
	})

	api.TodosFindTodosHandler = todos.FindTodosHandlerFunc(func(params todos.FindTodosParams) middleware.Responder {
		mergedParams := todos.NewFindTodosParams()
		mergedParams.Since = swag.Int64(0)
		if params.Since != nil {
			mergedParams.Since = params.Since
		}
		if params.Limit != nil {
			mergedParams.Limit = params.Limit
		}
		return todos.NewFindTodosOK().WithPayload(allItems(*mergedParams.Since, *mergedParams.Limit))
	})

	api.TodosDestroyOneHandler = todos.DestroyOneHandlerFunc(func(params todos.DestroyOneParams) middleware.Responder {
		delete(items, params.ID)
		return todos.NewDestroyOneNoContent()
	})
	// ここまで追加

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	myCors := cors.New(cors.Options{
		AllowedOrigins: []string{"http://127.0.0.1:3001"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})
	handleCORS := myCors.Handler
	return handleCORS(handler)
}
