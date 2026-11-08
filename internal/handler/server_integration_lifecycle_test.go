package handler

import (
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"
)

const (
	integrationProbeTimeout = 250 * time.Millisecond
	integrationPollInterval = 20 * time.Millisecond
)

func waitForIntegrationServer(t *testing.T, baseURL string, done <-chan error) {
	t.Helper()
	deadline := time.Now().Add(integrationWaitTimeout)
	client := &http.Client{Timeout: integrationProbeTimeout}
	for time.Now().Before(deadline) {
		if integrationServerReady(t, client, baseURL) {
			return
		}
		checkIntegrationServerRunning(t, done)
		time.Sleep(integrationPollInterval)
	}
	t.Fatal("framework server did not become ready")
}

func integrationServerReady(t *testing.T, client *http.Client, baseURL string) bool {
	t.Helper()
	response, err := client.Get(baseURL + "/health")
	if err != nil {
		return false
	}
	if closeErr := response.Body.Close(); closeErr != nil {
		t.Errorf("close readiness response: %v", closeErr)
	}
	return response.StatusCode == http.StatusOK
}

func checkIntegrationServerRunning(t *testing.T, done <-chan error) {
	t.Helper()
	select {
	case err := <-done:
		t.Fatalf("framework server exited before ready: %v", err)
	default:
	}
}

func stopIntegrationServer(t *testing.T, done <-chan error) {
	t.Helper()
	if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
		t.Errorf("stop framework server: %v", err)
		return
	}
	waitForIntegrationShutdown(t, done)
}

func waitForIntegrationShutdown(t *testing.T, done <-chan error) {
	t.Helper()
	select {
	case err := <-done:
		checkIntegrationShutdownError(t, err)
	case <-time.After(integrationWaitTimeout):
		t.Error("framework server did not stop")
	}
}

func checkIntegrationShutdownError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("framework server stopped with error: %v", err)
	}
}
