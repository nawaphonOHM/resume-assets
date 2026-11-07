package handler

import (
	"net/http"
	"testing"
)

func TestFrameworkAssetRouter(t *testing.T) {
	service := &routerAssetService{}
	baseURL := startFrameworkServer(t, RegisterRoutes(NewAssetHandler(service)))
	client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
	assertIntegrationAsset(t, client, baseURL, service)
	assertIntegrationRoutes(t, client, baseURL, service)
	assertIntegrationProbes(t, client, baseURL)
}
