package debug

import (
	"context"
	"github.com/stretchr/testify/require"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestServe(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	server := Serve(rand.Intn(55536) + 10000)
	require.NotNil(t, server, "Debug request fail")
	runtime.Gosched()
	_, err := http.NewRequest("GET", "/debug/debug/", nil)
	require.NoError(t, err, "Debug request fail")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	require.NoError(t, server.Shutdown(ctx), "Shutdown fail")
	cancel()
}
