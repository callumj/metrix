package handlers

import (
	"github.com/callumj/metrix/shared"
	"net/http"
)

func verifyAPIKey(c http.ResponseWriter, req *http.Request) bool {
	if !shared.Config.ApiKeyActive {
		return false
	}

	apiKey := req.FormValue("api_key")
	if apiKey != shared.Config.ApiKey {
		http.Error(c, "API Key invalid", http.StatusForbidden)
		return false
	}
	return true
}
