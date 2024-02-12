package utils

import (
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const MAX_PAGE_SIZE = 50

func ParsePage(w http.ResponseWriter, r *http.Request) (int, error) {
	pageString := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		HttpJSONError(w, "invalid page", http.StatusBadRequest)
		log.Error(err)
		return -1, err
	}
	return page, nil
}

func ParsePageSize(w http.ResponseWriter, r *http.Request) (int, error) {
	pageSizeString := r.URL.Query().Get("pageSize")
	pageSize, err := strconv.Atoi(pageSizeString)
	if err != nil {
		HttpJSONError(w, "invalid page size", http.StatusBadRequest)
		log.Error(err)
		return -1, err
	}
	if pageSize > MAX_PAGE_SIZE {
		pageSize = MAX_PAGE_SIZE
	}
	return pageSize, nil
}
