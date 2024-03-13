package main

import (
	"encoding/json"
	"errors"
		"net/http"

	"github.com/marioromandono/supplementapp/internal/supplement"
)

type ErrorResponseBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func addRoutes(mux *http.ServeMux, service *supplement.SupplementService) {
	mux.HandleFunc("GET /supplement/{gtin}", getSupplementHandler(service))
	mux.HandleFunc("POST /supplement", createSupplementHandler(service))
	mux.HandleFunc("PUT /supplement/{gtin}", updateSupplementHandler(service))
	mux.HandleFunc("PATCH /supplement/{gtin}", updateSupplementHandler(service))
	mux.HandleFunc("DELETE /supplement/{gtin}", deleteSupplementHandler(service))
}

func getSupplementHandler(service *supplement.SupplementService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gtin := r.PathValue("gtin")

		supplement, err := service.FindByGtin(r.Context(), gtin)

		if err != nil {
			handleError(err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(supplement)
	}
}

func createSupplementHandler(service *supplement.SupplementService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var supplement supplement.Supplement
		json.NewDecoder(r.Body).Decode(&supplement)

		err := service.Create(r.Context(), supplement)

		if err != nil {
			handleError(err, w)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Location", "/supplement/"+supplement.Gtin)
	}
}

func updateSupplementHandler(service *supplement.SupplementService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gtin := r.PathValue("gtin")

		var supplement supplement.UpdatableSupplement
		json.NewDecoder(r.Body).Decode(&supplement)

		err := service.Update(r.Context(), gtin, supplement)

		if err != nil {
			handleError(err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func deleteSupplementHandler(service *supplement.SupplementService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gtin := r.PathValue("gtin")

		err := service.Delete(r.Context(), gtin)

		if err != nil {
			handleError(err, w)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func handleError(err error, w http.ResponseWriter) {
	message := err.Error()
	var code int
	
	switch {
	case errors.Is(err, supplement.ErrNotFound):
		code = http.StatusNotFound
	case errors.Is(err, supplement.ErrAlreadyExists):
		code = http.StatusConflict
	default:
		code = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(ErrorResponseBody{Code: code, Message: message})
}
