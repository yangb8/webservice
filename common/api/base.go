package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

type Handler struct {
	*mux.Router
	RequestIDGenFunc func(req *http.Request) string
}

func defaultRequestIDGenFunc(req *http.Request) string {
	return uuid.NewV4().String()
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var requestID string
	if h.RequestIDGenFunc == nil {
		requestID = defaultRequestIDGenFunc(r)
	} else {
		requestID = h.RequestIDGenFunc(r)
	}

	//	defer func() {
	//		// recover here for panic
	//		if p := recover(); p != nil {
	//			w.WriteHeader(500)
	//
	//			// report to sentry
	//			if sentry.SentryClient != nil {
	//				err, ok := p.(error)
	//				if !ok {
	//					err = errors.New(p.(string))
	//				}
	//
	//				packet := raven.NewPacket(err.Error(),
	//					raven.NewException(err, raven.NewStacktrace(2, 3, sentry.SentryClient.IncludePaths())),
	//					raven.NewHttp(r))
	//				packet.AddTags(map[string]string{"RequestID": requestID})
	//
	//				sentry.SentryClient.Capture(packet, nil)
	//			}
	//		}
	//	}()

	log.Printf("%s %s %s", requestID, r.Method, r.URL)
	w.Header().Set("X-Request-ID", requestID)
	h.Router.ServeHTTP(w, r)
}
