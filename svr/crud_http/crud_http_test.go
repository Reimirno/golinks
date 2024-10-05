package crud_http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/reimirno/golinks/pkg/mapper"
	"github.com/reimirno/golinks/pkg/sanitizer"
	"github.com/reimirno/golinks/pkg/types"
	"github.com/stretchr/testify/assert"
)

var (
	fakePair = &types.PathUrlPair{
		Path: "fk",
		Url:  "https://fake.com",
	}
	fakePairAlt = &types.PathUrlPair{
		Path: "fk",
		Url:  "https://fakealt.com",
	}
	fakePair2 = &types.PathUrlPair{
		Path: "fk2",
		Url:  "https://fake2.com",
	}
	fakePair3 = &types.PathUrlPair{
		Path: "fk3",
		Url:  "https://fake3.com",
	}

	// When using it, please clone it first
	mockConfigurer = &mapper.MockMapperConfigurer{
		Name:        "mock",
		IsSingleton: false,
		IsReadOnly:  false,
		StarterPairs: types.PathUrlPairMap{
			"fk":  fakePair,
			"fk2": fakePair2,
		},
	}
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name          string
		configurers   []*mapper.MockMapperConfigurer
		persistorName string
		port          string
		wantErr       bool
	}{
		{
			name:          "happy path",
			configurers:   []*mapper.MockMapperConfigurer{mockConfigurer},
			persistorName: "mock",
			port:          "8082",
			wantErr:       false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm, err := mapper.NewMapperManager(test.persistorName, mapper.CloneConfigurers(test.configurers))
			assert.NoError(t, err)
			got, err := NewServer(mm, test.port)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

func TestServer_GetName(t *testing.T) {
	server := &Server{}
	assert.Equal(t, crudHttpServiceName, server.GetName())
}

func TestServer_GetUrl(t *testing.T) {
	tests := []struct {
		name          string
		configurers   []*mapper.MockMapperConfigurer
		persistorName string
		path          string
		want          *types.PathUrlPair
		wantStatus    int
	}{
		{
			name:          "happy path",
			configurers:   []*mapper.MockMapperConfigurer{mockConfigurer},
			persistorName: "mock",
			path:          "fk",
			want:          fakePair,
			wantStatus:    http.StatusOK,
		},
		{
			name:          "path not found",
			configurers:   []*mapper.MockMapperConfigurer{mockConfigurer},
			persistorName: "mock",
			path:          "invalid",
			wantStatus:    http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm, err := mapper.NewMapperManager(test.persistorName, mapper.CloneConfigurers(test.configurers))
			assert.NoError(t, err)
			server, err := NewServer(mm, "8082")
			assert.NoError(t, err)

			req, err := http.NewRequest("GET", fmt.Sprintf("/go/%s/", test.path), nil)
			assert.NoError(t, err)
			assert.NotNil(t, req)

			r := mux.NewRouter()
			r.HandleFunc("/go/{path}/", server.handleGetUrl).Methods("GET")

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, test.wantStatus, rr.Code)
			if test.wantStatus == http.StatusOK {
				resp := rr.Result()
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)
				var gotPair types.PathUrlPair
				err = json.Unmarshal(body, &gotPair)
				assert.NoError(t, err)
				canonicalPath, err := sanitizer.CanonicalizePath(gotPair.Path)
				assert.NoError(t, err)
				assert.Equal(t, canonicalPath, gotPair.Path)
				assert.Equal(t, test.want.Url, gotPair.Url)
			}
		})
	}
}

func TestServer_ListUrls(t *testing.T) {
	tests := []struct {
		name          string
		configurers   []*mapper.MockMapperConfigurer
		persistorName string
		wantStatus    int
		numPairs      int
	}{
		{
			name:          "happy path",
			configurers:   []*mapper.MockMapperConfigurer{mockConfigurer},
			persistorName: "mock",
			wantStatus:    http.StatusOK,
			numPairs:      2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm, err := mapper.NewMapperManager(test.persistorName, mapper.CloneConfigurers(test.configurers))
			assert.NoError(t, err)
			server, err := NewServer(mm, "8082")
			assert.NoError(t, err)

			req, err := http.NewRequest("GET", "/go/", nil)
			assert.NoError(t, err)
			assert.NotNil(t, req)

			r := mux.NewRouter()
			r.HandleFunc("/go/", server.handleListUrls).Methods("GET")

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, test.wantStatus, rr.Code)
			if test.wantStatus == http.StatusOK {
				resp := rr.Result()
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)
				var got []types.PathUrlPair
				err = json.Unmarshal(body, &got)
				assert.NoError(t, err)
				assert.Equal(t, test.numPairs, len(got))
			}
		})
	}
}

func TestServer_PutUrl(t *testing.T) {
	tests := []struct {
		name          string
		configurers   []*mapper.MockMapperConfigurer
		persistorName string
		pair          *types.PathUrlPair
		wantStatus    int
		finalPair     *types.PathUrlPair
	}{
		{
			name:          "happy path update",
			configurers:   []*mapper.MockMapperConfigurer{mockConfigurer},
			persistorName: "mock",
			pair:          fakePairAlt,
			wantStatus:    http.StatusAccepted,
			finalPair:     fakePairAlt,
		},
		{
			name:          "happy path create",
			configurers:   []*mapper.MockMapperConfigurer{mockConfigurer},
			persistorName: "mock",
			pair:          fakePair3,
			wantStatus:    http.StatusAccepted,
			finalPair:     fakePair3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm, err := mapper.NewMapperManager(test.persistorName, mapper.CloneConfigurers(test.configurers))
			assert.NoError(t, err)
			server, err := NewServer(mm, "8082")
			assert.NoError(t, err)

			body, err := json.Marshal(test.pair)
			assert.NoError(t, err)
			req, err := http.NewRequest("POST", "/go/", bytes.NewBuffer(body))
			assert.NoError(t, err)
			assert.NotNil(t, req)

			r := mux.NewRouter()
			r.HandleFunc("/go/", server.handlePutUrl).Methods("POST")

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, test.wantStatus, rr.Code)
			if test.wantStatus == http.StatusAccepted {
				resp := rr.Result()
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)
				var got types.PathUrlPair
				err = json.Unmarshal(body, &got)
				assert.NoError(t, err)
				canonicalPath, err := sanitizer.CanonicalizePath(test.finalPair.Path)
				assert.NoError(t, err)
				assert.Equal(t, canonicalPath, got.Path)
				assert.Equal(t, test.finalPair.Url, got.Url)
			}
		})
	}
}

func TestServer_DeleteUrl(t *testing.T) {
	tests := []struct {
		name          string
		configurers   []*mapper.MockMapperConfigurer
		persistorName string
		path          string
		wantStatus    int
	}{
		{
			name:          "happy path",
			configurers:   []*mapper.MockMapperConfigurer{mockConfigurer},
			persistorName: mockConfigurer.Name,
			path:          "fk",
			wantStatus:    http.StatusNoContent,
		},
		{
			name:          "path not found is fine",
			configurers:   []*mapper.MockMapperConfigurer{mockConfigurer},
			persistorName: mockConfigurer.Name,
			path:          "invalid",
			wantStatus:    http.StatusNoContent,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm, err := mapper.NewMapperManager(test.persistorName, mapper.CloneConfigurers(test.configurers))
			assert.NoError(t, err)
			server, err := NewServer(mm, "8082")
			assert.NoError(t, err)

			req, err := http.NewRequest("DELETE", fmt.Sprintf("/go/%s/", test.path), nil)
			assert.NoError(t, err)
			assert.NotNil(t, req)

			r := mux.NewRouter()
			r.HandleFunc("/go/{path}/", server.handleDeleteUrl).Methods("DELETE")

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, test.wantStatus, rr.Code)
		})
	}
}
