package plank

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBuildMetadata(t *testing.T) {
	handler := &testHandler{}
	srv := httptest.NewServer(handler)
	defer srv.Close()
	client := NewLogkeeperClient(NewLogkeeperClientOptions{BaseURL: srv.URL})

	build := Build{
		ID:            "build",
		Builder:       "builder",
		BuildNum:      1,
		TaskID:        "task",
		TaskExecution: 2,
		Tests: []Test{
			{
				ID:            "test0",
				Name:          "Test_0",
				BuildID:       "build",
				TaskID:        "task",
				TaskExecution: 2,
				Phase:         "phase",
				Command:       "cmd",
			},
		},
	}
	buildData, err := json.Marshal(&build)
	require.NoError(t, err)

	t.Run("DoRequestFails", func(t *testing.T) {
		canceledCtx, cancel := context.WithCancel(context.Background())
		cancel()
		handler.statusCode = http.StatusOK
		handler.body = buildData

		_, err := client.GetBuildMetadata(canceledCtx, "id")
		assert.Error(t, err)
	})
	t.Run("BadStatus", func(t *testing.T) {
		handler.statusCode = http.StatusNotFound
		handler.body = nil

		_, err := client.GetBuildMetadata(context.Background(), "id")
		assert.Error(t, err)
	})
	t.Run("ExpectedBuildPayload", func(t *testing.T) {
		handler.statusCode = http.StatusOK
		handler.body = buildData

		out, err := client.GetBuildMetadata(context.Background(), "id")
		require.NoError(t, err)
		assert.Equal(t, build, out)
	})
}

func TestGetTestMetadata(t *testing.T) {
	handler := &testHandler{}
	srv := httptest.NewServer(handler)
	defer srv.Close()
	client := NewLogkeeperClient(NewLogkeeperClientOptions{BaseURL: srv.URL})

	test := Test{
		ID:            "test0",
		Name:          "Test_0",
		BuildID:       "build",
		TaskID:        "task",
		TaskExecution: 2,
		Phase:         "phase",
		Command:       "cmd",
	}
	testData, err := json.Marshal(&test)
	require.NoError(t, err)

	t.Run("DoRequestFails", func(t *testing.T) {
		canceledCtx, cancel := context.WithCancel(context.Background())
		cancel()
		handler.statusCode = http.StatusOK
		handler.body = testData

		_, err := client.GetTestMetadata(canceledCtx, "build", "test")
		assert.Error(t, err)
	})
	t.Run("BadStatus", func(t *testing.T) {
		handler.statusCode = http.StatusNotFound
		handler.body = nil

		_, err := client.GetTestMetadata(context.Background(), "build", "test")
		assert.Error(t, err)
	})
	t.Run("ExpectedBuildPayload", func(t *testing.T) {
		handler.statusCode = http.StatusOK
		handler.body = testData

		out, err := client.GetTestMetadata(context.Background(), "build", "test")
		require.NoError(t, err)
		assert.Equal(t, test, out)
	})
}

type testHandler struct {
	statusCode int
	body       []byte
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.statusCode > 0 {
		w.WriteHeader(h.statusCode)
	}
	if h.body != nil {
		_, _ = w.Write(h.body)
	}
}
