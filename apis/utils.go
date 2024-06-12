package apis

import (
	"encoding/json"
	"os"
	"path"

	"github.com/riccardom/briatore/types"
)

var (
	resultsFolder = path.Join(types.HomePath, "results")
)

func init() {
	_ = os.MkdirAll(resultsFolder, 0755)
}

// StoreResults stores the results for the report having the given id
func StoreResults(id types.ReportID, result *types.ReportResult) error {
	bz, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return os.WriteFile(path.Join(resultsFolder, id.String()), bz, 0600)
}

// GetResult returns the result for the report having the given id
func GetResult(id types.ReportID) (result *types.ReportResult, found bool, err error) {
	bz, err := os.ReadFile(path.Join(resultsFolder, id.String()))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, err
	}

	var res types.ReportResult
	err = json.Unmarshal(bz, &res)
	if err != nil {
		return nil, false, err
	}

	return &res, true, nil
}
