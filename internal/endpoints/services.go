package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/redhatinsights/platform-changelog-go/internal/db"
	l "github.com/redhatinsights/platform-changelog-go/internal/logging"
	"github.com/redhatinsights/platform-changelog-go/internal/metrics"
	"github.com/redhatinsights/platform-changelog-go/internal/models"
)

func GetServicesAll(w http.ResponseWriter, r *http.Request) {
	metrics.IncRequests(r.URL.Path, r.Method, r.UserAgent())

	q, err := initQuery(r)

	if err != nil {
		writeResponse(w, http.StatusBadRequest, "Invalid query")
		return
	}

	servicesWithTimelines, count, err := db.GetServicesAll(db.DB, q.Offset, q.Limit)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	servicesList := models.ExpandedServicesList{Count: count, Data: servicesWithTimelines}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(servicesList)
}

func GetServiceByName(w http.ResponseWriter, r *http.Request) {
	metrics.IncRequests(r.URL.Path, r.Method, r.UserAgent())

	serviceName := chi.URLParam(r, "service")
	service, _, err := db.GetServiceByName(db.DB, serviceName)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	/**
	 * the service doesn't exist
	 */
	if service.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Service not found"))
		return
	}

	l.Log.Debugf("URL Param: %s", serviceName)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(service)
}
