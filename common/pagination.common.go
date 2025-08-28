package common

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

const (
	DEFAULT_PAGE  = 1
	DEFAULT_LIMIT = 10
	MAX_LIMIT     = 100
)

func GetPagination(ctx *gin.Context) *Pagination {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", strconv.Itoa(DEFAULT_PAGE)))
	if err != nil || page < 1 {
		page = DEFAULT_PAGE
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", strconv.Itoa(DEFAULT_LIMIT)))
	if err != nil || limit < 1 {
		limit = DEFAULT_LIMIT
	}

	if limit > MAX_LIMIT {
		limit = MAX_LIMIT
	}

	return &Pagination{
		Page:  page,
		Limit: limit,
	}
}
func PaginatedResponse(ctx *gin.Context, status int, data interface{}, totalDocs int, p *Pagination) {
	totalPages := (totalDocs + p.Limit - 1) / p.Limit

	nextPage := 0
	if p.Page < totalPages {
		nextPage = p.Page + 1
	}

	prevPage := 0
	if p.Page > 1 {
		prevPage = p.Page - 1
	}

	ctx.JSON(status, gin.H{
		"data": data,
		"meta": gin.H{
			"page":       p.Page,
			"limit":      p.Limit,
			"totalDocs":  totalDocs,
			"totalPages": totalPages,
			"nextPage":   nextPage,
			"prevPage":   prevPage,
		},
	})
}
