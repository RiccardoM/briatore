package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/riccardom/briatore/report"
	"github.com/riccardom/briatore/types"
)

const (
	idParam     = "id"
	outputParam = "output"
)

// GetResultHandler returns the handler used to get the results of a report
func GetResultHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		result, found, err := GetResult(types.ParseReportID(c.Query(idParam)))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		if !found {
			c.String(http.StatusNotFound, "Report not found. Please try calling this endpoint later")
			return
		}

		if result.IsError() {
			c.String(http.StatusBadRequest, result.Error)
			return
		}

		outParam := c.Query(outputParam)
		if outParam == "" {
			outParam = types.OutText.String()
		}

		output, err := types.ParseOutput(outParam)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}

		bz, err := report.MarshalAmounts(result.GetAmounts(), output)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		var contentType string
		switch output {
		case types.OutJSON:
			contentType = "application/json"
		case types.OutCSV:
			contentType = "text/csv"
		case types.OutText:
			contentType = "text/plain"
		}

		c.Data(http.StatusOK, contentType, bz)
	}
}
