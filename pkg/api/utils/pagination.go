package utils

import (
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const MAX_PAGE_SIZE = 50
const DEFAULT_PAGE = 1
const DEFAULT_PAGE_SIZE = 10

func ParsePage(w http.ResponseWriter, r *http.Request) int {
	pageString := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		log.Debug(err)
		page = DEFAULT_PAGE
		return page
	}
	if page < 1 {
		page = DEFAULT_PAGE
	}
	return page
}

func ParsePageSize(w http.ResponseWriter, r *http.Request) int {
	pageSizeString := r.URL.Query().Get("pageSize")
	pageSize, err := strconv.Atoi(pageSizeString)
	if err != nil {
		log.Debug(err)
		pageSize = DEFAULT_PAGE_SIZE
		return pageSize
	}
	if pageSize > MAX_PAGE_SIZE {
		pageSize = DEFAULT_PAGE_SIZE
	}
	return pageSize
}
