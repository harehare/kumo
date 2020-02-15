package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/harehare/kumo/middleware"
	"github.com/harehare/kumo/model"
)

func List(w http.ResponseWriter, r *http.Request) *httpError {
	var pageSize int
	var pageNo int
	pageSizeQuery, ok := r.URL.Query()["page_size"]

	if !ok || len(pageSizeQuery[0]) < 1 {
		pageSize = 30
	} else {
		v, err := strconv.Atoi(pageSizeQuery[0])

		if err != nil {
			return &httpError{Err: err, Msg: "Invalid pageSize parameter", StatusCode: http.StatusBadRequest}
		}

		pageSize = v
	}

	pageNoQuery, ok := r.URL.Query()["page_size"]

	if !ok || len(pageNoQuery[0]) < 1 {
		pageNo = 1
	} else {
		v, err := strconv.Atoi(pageNoQuery[0])

		if err != nil {
			return &httpError{Err: err, Msg: "Invalid pageNo parameter", StatusCode: http.StatusBadRequest}
		}

		pageNo = v
	}

	userID := r.Context().Value(middleware.UIDKey).(string)
	items, err := ItemService.List(r.Context(), userID, pageNo, pageSize)

	if err != nil {
		return &httpError{Err: err, Msg: "Error fetch item list", StatusCode: http.StatusInternalServerError}
	}

	if len(*items) == 0 {
		w.Write([]byte("[]"))
		return nil
	}

	res, err := json.Marshal(items)

	if err != nil {
		return &httpError{Err: err, Msg: "Error decode json", StatusCode: http.StatusInternalServerError}
	}

	w.Write(res)
	return nil
}

func Get(w http.ResponseWriter, r *http.Request) *httpError {
	vars := mux.Vars(r)
	userID := r.Context().Value(middleware.UIDKey).(string)
	itemID := vars["ID"]

	item, err := ItemService.Get(r.Context(), userID, itemID)

	if err != nil {
		return &httpError{Err: err, Msg: "Error fetch item", StatusCode: http.StatusInternalServerError}
	}

	res, err := json.Marshal(item)

	if err != nil {
		return &httpError{Err: err, Msg: "Error decode json", StatusCode: http.StatusInternalServerError}
	}

	w.Write(res)

	return nil
}

func Delete(w http.ResponseWriter, r *http.Request) *httpError {
	vars := mux.Vars(r)
	userID := r.Context().Value(middleware.UIDKey).(string)
	itemID := vars["ID"]

	err := ItemService.Delete(r.Context(), userID, itemID)

	if err != nil {
		return &httpError{Err: err, Msg: "Error delete item", StatusCode: http.StatusInternalServerError}
	}

	if err != nil {
		return &httpError{Err: err, Msg: "Error decode json", StatusCode: http.StatusInternalServerError}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func Save(w http.ResponseWriter, r *http.Request) *httpError {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return &httpError{Err: err, Msg: "Error decode body", StatusCode: http.StatusInternalServerError}
	}

	var item model.Item

	err = json.Unmarshal(b, &item)
	if err != nil {
		return &httpError{Err: err, Msg: "Error decode json", StatusCode: http.StatusInternalServerError}
	}

	userID := r.Context().Value(middleware.UIDKey).(string)
	item.UserID = userID

	err = ItemService.Save(r.Context(), &item)

	if err != nil {
		return &httpError{Err: err, Msg: "Error save item", StatusCode: http.StatusInternalServerError}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
