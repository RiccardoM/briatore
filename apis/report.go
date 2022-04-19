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

		id := types.RandomReportID()
		go ComputeReport(cfg, id, addresses, date)

		c.String(http.StatusOK, "Report queued. Your id is %s", id)
	}
}

// ComputeReport computes the result of the report for the provided addresses and date,
// storing it associated with the given id.
func ComputeReport(cfg *types.Config, id types.ReportID, addresses []string, date time.Time) {
	result := report.GetReport(cfg, addresses, date)
	_ = StoreResults(id, result)
}
