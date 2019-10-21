package handler

import (
	"encoding/json"
	"net/http"

	"github.com/skygeario/skygear-server/pkg/core/db"
	"github.com/skygeario/skygear-server/pkg/core/skyerr"
	xskyerr "github.com/skygeario/skygear-server/pkg/core/xskyerr"
)

type APIHandler interface {
	DecodeRequest(request *http.Request, resp http.ResponseWriter) (RequestPayload, error)
	WithTx() bool
	Handle(requestPayload interface{}) (interface{}, error)
}

type APITxHandler interface {
	WillCommitTx() error
	DidCommitTx()
}

// TODO(error): use new APIError
type APIResponse struct {
	Result interface{}       `json:"result,omitempty"`
	Err    skyerr.Error      `json:"error,omitempty"`
	Error  *xskyerr.APIError `json:",omitempty"`
}

func APIHandlerToHandler(apiHandler APIHandler, txContext db.TxContext) http.Handler {
	txHandler, _ := apiHandler.(APITxHandler)

	handleAPICall := func(r *http.Request, resp http.ResponseWriter) (response APIResponse) {
		payload, err := apiHandler.DecodeRequest(r, resp)
		if err != nil {
			response.Error = xskyerr.AsAPIError(err)
			return
		}

		if err := payload.Validate(); err != nil {
			response.Error = xskyerr.AsAPIError(err)
			return
		}

		defer func() {
			if err != nil {
				response.Error = xskyerr.AsAPIError(err)
			}
		}()

		if apiHandler.WithTx() {
			// assume txContext != nil if apiHandler.WithTx() is true
			if err := txContext.BeginTx(); err != nil {
				response.Error = xskyerr.AsAPIError(err)
				return
			}

			defer func() {
				err = db.EndTx(txContext, err)
				if err == nil && txHandler != nil {
					txHandler.DidCommitTx()
				}
			}()
		}

		responsePayload, err := apiHandler.Handle(payload)

		if err == nil && txHandler != nil {
			err = txHandler.WillCommitTx()
		}

		if err == nil {
			response.Result = responsePayload
		}

		return
	}

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		response := handleAPICall(r, rw)
		WriteResponse(rw, response)
	})
}

func WriteResponse(rw http.ResponseWriter, response APIResponse) {
	httpStatus := http.StatusOK
	encoder := json.NewEncoder(rw)

	if response.Err != nil {
		httpStatus = skyerr.ErrorDefaultStatusCode(response.Err)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(httpStatus)
	encoder.Encode(response)
}

// Transactional runs f within a transaction.
// If err is non-nil, the transaction is rolled back.
// Otherwise the transaction is committed.
// It is a lightweight and flexible alternative to APIHandler
// because it is not coupled with http.
func Transactional(txContext db.TxContext, f func() (interface{}, error)) (result interface{}, err error) {
	err = db.WithTx(txContext, func() error {
		result, err = f()
		return err
	})
	return
}
