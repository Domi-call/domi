package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/filebrowser/filebrowser/v2/settings"
)

type QuotaRequest struct {
	FilesetName string `json:"filesetName"`
	QuotaLimit  int    `json:"quotaLimmit"`
	QuotaMax    int    `json:"quotaMax"`
}

var createFilesetHandler = withAdmin(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return http.StatusBadRequest, err
	}
	m := make(map[string]interface{})
	if err := json.Unmarshal(body, &m); err != nil {
		return http.StatusBadRequest, err
	}
	username := m["username"].(string)
	err = settings.CreateFileset(username, d.settings.Gpfs)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	_ = settings.SetQuota(username, d.settings.Gpfs.QuotaLimmit, d.settings.Gpfs.QuotaMax, d.settings.Gpfs)
	return renderJSON(w, r, "")
})

var setUserQuotaHandler = withAdmin(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	var req QuotaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return http.StatusBadRequest, err
	}

	err := settings.SetQuota(req.FilesetName, uint16(req.QuotaLimit), uint16(req.QuotaMax), d.settings.Gpfs)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return renderJSON(w, r, "")
})

var queryQuotaHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	queryParams := r.URL.Query()
	objectName := queryParams.Get("objectName")
	err, quotas := settings.QueryFilesetQuota(objectName, d.settings.Gpfs)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return renderJSON(w, r, quotas)
})

var queryQuotaDefaultHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	m := make(map[string]int)
	m["soft"] = int(d.settings.Gpfs.QuotaLimmit)
	m["hard"] = int(d.settings.Gpfs.QuotaMax)
	return renderJSON(w, r, m)
})

var queryUserUsageHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	err, usage := settings.QueryFilesetUsage(d.user.Username, d.settings.Gpfs)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	m := make(map[string]interface{})
	m["usage"] = usage
	return renderJSON(w, r, m)
})

var queryUserQuotaHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	err, quota := settings.QueryFilesetQuota(d.user.Username, d.settings.Gpfs)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	m := make(map[string]interface{})
	m["soft"] = int64(0)
	m["hard"] = int64(0)
	m["used"] = int64(0)
	if len(quota) != 0 {
		m["soft"] = quota[0].BlockQuota * 1024 //B
		m["hard"] = quota[0].BlockLimit * 1024
		m["used"] = quota[0].BlockUsage * 1024
	}
	return renderJSON(w, r, m)
})

var queryFilesetUsageHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	queryParams := r.URL.Query()
	fileset := queryParams.Get("fileset")
	err, usage := settings.QueryFilesetUsage(fileset, d.settings.Gpfs)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	m := make(map[string]interface{})
	m["usage"] = usage
	return renderJSON(w, r, m)
})
