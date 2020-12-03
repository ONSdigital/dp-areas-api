package api

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/ONSdigital/dp-authorisation/auth"
	dphandlers "github.com/ONSdigital/dp-net/handlers"
	dprequest "github.com/ONSdigital/dp-net/request"
	"github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/config"
	"github.com/ONSdigital/dp-topic-api/store"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

var (
	trueStringified = strconv.FormatBool(true)

	createPermission = auth.Permissions{Create: true}
	readPermission   = auth.Permissions{Read: true}
	updatePermission = auth.Permissions{Update: true}
	deletePermission = auth.Permissions{Delete: true}
)

//API provides a struct to wrap the api around
type API struct {
	Router                 *mux.Router
	dataStore              store.DataStore
	topicPermissions       AuthHandler
	permissions            AuthHandler
	enablePrivateEndpoints bool
}

//Setup function sets up the api and returns an api
func Setup(ctx context.Context, cfg *config.Configuration, router *mux.Router, dataStore store.DataStore, topicPermissions AuthHandler, permissions AuthHandler) *API {
	api := &API{
		Router:                 router,
		dataStore:              dataStore,
		topicPermissions:       topicPermissions,
		permissions:            permissions,
		enablePrivateEndpoints: cfg.EnablePrivateEndpoints,
	}

	if cfg.EnablePrivateEndpoints {
		// create publishing related endpoints ...
		log.Event(ctx, "enabling private endpoints for topic api", log.INFO)

		api.enablePrivateTopicEndpoints(ctx)
	} else {
		// create web related endpoints ...

		log.Event(ctx, "enabling only public endpoints for dataset api", log.INFO)
		api.enablePublicEndpoints(ctx)
	}

	router.HandleFunc("/hello", HelloHandler()).Methods("GET")
	return api
}

// enablePublicEndpoints register only the public GET endpoints.
func (api *API) enablePublicEndpoints(ctx context.Context) {
	api.get("/topics/{id}", api.getTopicPublicHandler)
}

// enablePrivateTopicEndpoints register the topics endpoints with the appropriate authentication and authorisation
// checks required when running the dataset API in publishing (private) mode.
func (api *API) enablePrivateTopicEndpoints(ctx context.Context) {
	api.get(
		"/topics/{id}",
		api.isAuthenticated(
			api.isAuthorised(readPermission, api.getTopicPrivateHandler)), //...ask Eleanor if its this line that should be used ?
		/*api.isAuthorisedForTopics(readPermission, api.getTopicPrivateHandler*/
	)
}

// isAuthenticated wraps a http handler func in another http handler func that checks the caller is authenticated to
// perform the requested action. handler is the http.HandlerFunc to wrap in an
// authentication check. The wrapped handler is only called if the caller is authenticated
func (api *API) isAuthenticated(handler http.HandlerFunc) http.HandlerFunc {
	return dphandlers.CheckIdentity(handler)
}

// isAuthorised wraps a http.HandlerFunc another http.HandlerFunc that checks the caller is authorised to perform the
// requested action. required is the permissions required to perform the action, handler is the http.HandlerFunc to
// apply the check to. The wrapped handler is only called if the caller has the required permissions.
func (api *API) isAuthorised(required auth.Permissions, handler http.HandlerFunc) http.HandlerFunc {
	return api.permissions.Require(required, handler)
}

// isAuthorised wraps a http.HandlerFunc another http.HandlerFunc that checks the caller is authorised to perform the
// requested topics action. This authorisation check is specific to topics. required is the permissions required to
// perform the action, handler is the http.HandlerFunc to apply the check to. The wrapped handler is only called if the
// caller has the required topic permissions.
func (api *API) isAuthorisedForTopics(required auth.Permissions, handler http.HandlerFunc) http.HandlerFunc {
	return api.topicPermissions.Require(required, handler)
}

// get register a GET http.HandlerFunc.
func (api *API) get(path string, handler http.HandlerFunc) {
	api.Router.HandleFunc(path, handler).Methods("GET")
}

// get register a PUT http.HandlerFunc.
func (api *API) put(path string, handler http.HandlerFunc) {
	api.Router.HandleFunc(path, handler).Methods("PUT")
}

// get register a POST http.HandlerFunc.
func (api *API) post(path string, handler http.HandlerFunc) {
	api.Router.HandleFunc(path, handler).Methods("POST")
}

// get register a DELETE http.HandlerFunc.
func (api *API) delete(path string, handler http.HandlerFunc) {
	api.Router.HandleFunc(path, handler).Methods("DELETE")
}

func (api *API) authenticate(r *http.Request, logData log.Data) bool {
	var authenticated bool

	if api.enablePrivateEndpoints {
		var hasCallerIdentity, hasUserIdentity bool

		// NOTE:
		// If the identity exists then the user has been authenticated.
		// There is an earlier step in the middleware which will call off to zebedee to
		// authenticate the request (user/service) and this will add the identity to the
		// request context for later use in the application ...
		// ... which happens to be here:

		callerIdentity := dprequest.Caller(r.Context())
		if callerIdentity != "" {
			logData["caller_identity"] = callerIdentity
			hasCallerIdentity = true
		}

		userIdentity := dprequest.User(r.Context())
		if userIdentity != "" {
			logData["user_identity"] = userIdentity
			hasUserIdentity = true
		}

		if hasCallerIdentity || hasUserIdentity {
			authenticated = true
		}
		logData["authenticated"] = authenticated

		return authenticated
	}
	return authenticated
}

// WriteJSONBody marshals the provided interface into json, and writes it to the response body.
func WriteJSONBody(ctx context.Context, v interface{}, w http.ResponseWriter, data log.Data) error {

	// Set headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Marshal provided model
	payload, err := json.Marshal(v)
	if err != nil {
		handleError(ctx, w, apierrors.ErrInternalServer, data)
		return err
	}

	// Write payload to body
	if _, err := w.Write(payload); err != nil {
		handleError(ctx, w, apierrors.ErrInternalServer, data)
		return err
	}
	return nil
}

// ReadJSONBody reads the bytes from the provided body, and marshals it to the provided model interface.
func ReadJSONBody(ctx context.Context, body io.ReadCloser, v interface{}, w http.ResponseWriter, data log.Data) error {
	defer body.Close()

	// Set headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Get Body bytes
	payload, err := ioutil.ReadAll(body)
	if err != nil {
		handleError(ctx, w, apierrors.ErrUnableToReadMessage, data)
		return err
	}

	// Unmarshal body bytes to model
	if err := json.Unmarshal(payload, v); err != nil {
		handleError(ctx, w, apierrors.ErrUnableToParseJSON, data)
		return err
	}

	return nil
}

// handleError is a utility function that maps api errors to an http status code and sets the provided responseWriter accordingly
func handleError(ctx context.Context, w http.ResponseWriter, err error, data log.Data) {
	var status int
	if err != nil {
		switch err {
		//!!! fix this list for this service as the application develops to all features
		case apierrors.ErrTopicNotFound,
			apierrors.ErrVariantNotFound:
			status = http.StatusNotFound
		case apierrors.ErrUnableToReadMessage,
			apierrors.ErrUnableToParseJSON,
			apierrors.ErrImageFilenameTooLong,
			apierrors.ErrImageNoCollectionID,
			apierrors.ErrTopicInvalidState,
			apierrors.ErrImageDownloadTypeMismatch,
			apierrors.ErrImageDownloadInvalidState,
			apierrors.ErrImageIDMismatch,
			apierrors.ErrVariantIDMismatch:
			status = http.StatusBadRequest
		case apierrors.ErrImageAlreadyPublished,
			apierrors.ErrTopicAlreadyCompleted,
			apierrors.ErrTopicStateTransitionNotAllowed,
			apierrors.ErrImageBadInitialState,
			apierrors.ErrImageNotImporting,
			apierrors.ErrImageNotPublished,
			apierrors.ErrVariantAlreadyExists,
			apierrors.ErrVariantStateTransitionNotAllowed,
			apierrors.ErrImageDownloadBadInitialState:
			status = http.StatusForbidden
		default:
			status = http.StatusInternalServerError
		}
	}

	if data == nil {
		data = log.Data{}
	}

	switch status {
	case http.StatusNotFound, http.StatusForbidden:
		data["response_status"] = status
		data["user_error"] = err.Error()
		log.Event(ctx, "request unsuccessful", log.ERROR, data)
		http.Error(w, err.Error(), status)
	default:
		// a stack trace is added for Non User errors
		data["response_status"] = status
		log.Event(ctx, "request unsuccessful", log.ERROR, log.Error(err), data)
		http.Error(w, apierrors.ErrInternalServer.Error(), status)
	}

	//	data["response_status"] = status
	//	log.Event(ctx, "request unsuccessful", log.ERROR, log.Error(err), data)
	//	http.Error(w, err.Error(), status)
}
