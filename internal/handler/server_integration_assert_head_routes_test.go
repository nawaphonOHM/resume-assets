package handler

import (
	"net/http"
	"testing"
)

func assertHeadMetadataOmittedOnError(t *testing.T, response *http.Response, test integrationRouteCase) {
	t.Helper()
	if test.status == http.StatusOK {
		return
	}
	for _, name := range []string{
		handlerETagHeader,
		integrationHeaderLastModified,
		integrationHeaderContentType,
		handlerContentLengthHeader,
		integrationHeaderMode,
	} {
		assertHeadErrorHeaderOmitted(t, response, test, name)
	}
}

func assertHeadErrorHeaderOmitted(t *testing.T, response *http.Response, test integrationRouteCase, name string) {
	t.Helper()
	if response.Header.Get(name) != "" {
		t.Errorf(
			"HEAD error %s included asset metadata header %s=%q",
			test.path, name, response.Header.Get(name),
		)
	}
}
