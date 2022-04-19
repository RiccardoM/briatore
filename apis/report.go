package apis

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/riccardom/briatore/report"
	"github.com/riccardom/briatore/types"
)

const (
	addressesParam = "addresses"
	dateParam      = "date"
	outputParam    = "output"
)

// GetReportHandler returns the APIs handler to get a report
func GetReportHandler(cfg *types.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		addresses := strings.Split(c.Query(addressesParam), ",")
		if len(addresses) == 0 {
			c.String(http.StatusBadRequest, "No addresses provided")
			return
		}

		date, err := time.Parse(time.RFC3339, c.Query(dateParam))
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid date. Must be in RFC3339 format")
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

		bz, err := report.GetReportBytes(cfg, addresses, date, output)
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
