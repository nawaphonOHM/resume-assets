package handler

import (
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
)

const (
	integrationWaitTimeout = 5 * time.Second
	integrationHost        = "127.0.0.1"
)

func startFrameworkServer(t *testing.T, blueprint *rest.BluePrint) string {
	t.Helper()
	port := reserveIntegrationPort(t)
	configureIntegrationServer(t, port)
	done := make(chan error, 1)
	go func() { done <- rest.StartREST(blueprint) }()
	baseURL := "http://" + net.JoinHostPort(integrationHost, strconv.Itoa(port))
	waitForIntegrationServer(t, baseURL, done)
	t.Cleanup(func() { stopIntegrationServer(t, done) })
	return baseURL
}

func reserveIntegrationPort(t *testing.T) int {
	t.Helper()
	listener, err := net.Listen("tcp", net.JoinHostPort(integrationHost, "0"))
	if err != nil {
		t.Fatalf("reserve server port: %v", err)
	}
	port := integrationListenerPort(t, listener)
	if closeErr := listener.Close(); closeErr != nil {
		t.Fatalf("release server port: %v", closeErr)
	}
	return port
}

func integrationListenerPort(t *testing.T, listener net.Listener) int {
	t.Helper()
	address, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatal("reserved listener address is not TCP")
	}
	return address.Port
}

func configureIntegrationServer(t *testing.T, port int) {
	t.Helper()
	t.Setenv("OHM9996_SERVER_API_GROUP", "raw")
	t.Setenv("OHM9996_GIN_MODE", "test")
	t.Setenv("OHM9996_SERVER_HOST", integrationHost)
	t.Setenv("OHM9996_SERVER_PORT", strconv.Itoa(port))
	t.Setenv("OHM9996_SERVER_CORS_PATH", "")
}
