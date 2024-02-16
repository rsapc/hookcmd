package librenms

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

func GetLibreError(resp *resty.Response) (LibreErrorResponse, error) {
	obj := LibreErrorResponse{}
	err := json.Unmarshal(resp.Body(), &obj)
	return obj, err
}
