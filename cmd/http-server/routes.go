package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/marioromandono/supplementapp/internal/supplement"
)

type ErrorResponseBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func addRoutes(mux *http.ServeMux, service *supplement.SupplementService) {
	mux.HandleFunc("GET /supplement/{gtin}", getSupplementHandler(service))
	mux.HandleFunc("GET /supplement", listAllSupplementsHandler(service))
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
		err = json.NewEncoder(w).Encode(supplement)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func listAllSupplementsHandler(service *supplement.SupplementService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		supplements, err := service.ListAll(r.Context())

		if err != nil {
			handleError(err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(supplements)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func createSupplementHandler(service *supplement.SupplementService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var supplement supplement.Supplement
		err := json.NewDecoder(r.Body).Decode(&supplement)
		defer r.Body.Close()

		if err != nil {
			handleError(err, w)
			return
		}

		err = service.Create(r.Context(), supplement)

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
		err := json.NewDecoder(r.Body).Decode(&supplement)
		defer r.Body.Close()

		if err != nil {
			handleError(err, w)
			return
		}

		err = service.Update(r.Context(), gtin, supplement)

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

	var syntaxErr *json.SyntaxError
	var unmarshalTypeErr *json.UnmarshalTypeError
	var invalidUnmarshalErr *json.InvalidUnmarshalError
	var unsupportedTypeError *json.UnsupportedTypeError
	var unsupportedValueErr *json.UnsupportedValueError

	switch {
	case errors.Is(err, supplement.ErrNotFound):
		code = http.StatusNotFound
	case errors.Is(err, supplement.ErrAlreadyExists):
		code = http.StatusConflict
	case
		errors.Is(err, io.EOF),
		errors.As(err, &syntaxErr),
		errors.As(err, &unmarshalTypeErr),
		errors.As(err, &invalidUnmarshalErr),
		errors.As(err, &unsupportedTypeError),
		errors.As(err, &unsupportedValueErr),
		errors.Is(err, supplement.ErrInvalidSupplement):
		code = http.StatusBadRequest
	default:
		code = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	encodeErr := json.NewEncoder(w).Encode(ErrorResponseBody{Code: code, Message: message})
	if encodeErr != nil {
		http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
		return
	}
}
