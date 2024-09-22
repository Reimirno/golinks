package redirector

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/reimirno/golinks/pkg/mapper"
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

	mockConfigurer = &mapper.MockMapperConfigurer{
		Name:        "mock",
		IsSingleton: false,
		IsReadOnly:  false,
		StarterPairs: types.PathUrlPairMap{
			"fk":  fakePair,
			"fk2": fakePair2,
		},
	}
	mockConfigurerAlt = &mapper.MockMapperConfigurer{
		Name:        "mockAlt",
		IsSingleton: false,
		IsReadOnly:  false,
		StarterPairs: types.PathUrlPairMap{
			"fk": fakePairAlt,
		},
	}
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name          string
		configurers   []mapper.MapperConfigurer
		persistorName string
		port          string
		wantErr       bool
	}{
		{
			name:          "happy path",
			configurers:   []mapper.MapperConfigurer{mockConfigurer},
			persistorName: "mock",
			port:          "8080",
			wantErr:       false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm, err := mapper.NewMapperManager(test.persistorName, test.configurers)
			assert.NoError(t, err)
			assert.NotNil(t, mm)
			server, err := NewServer(mm, test.port)
			if test.wantErr {
				assert.Error(t, err)
				assert.Nil(t, server)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, server)
				assert.Equal(t, test.port, server.port)
				assert.Equal(t, mm, server.manager)
			}
		})
	}
}

func TestServer_GetName(t *testing.T) {
	server := &Server{}
	assert.Equal(t, redirectorServiceName, server.GetName())
}

func TestServer_handleRedirect(t *testing.T) {
	tests := []struct {
		name          string
		configurers   []mapper.MapperConfigurer
		persistorName string
		path          string
		statusCode    int
		redirectUrl   string
	}{
		{
			name:          "happy path",
			configurers:   []mapper.MapperConfigurer{mockConfigurer},
			persistorName: "mock",
			path:          "fk",
			redirectUrl:   fakePair.Url,
			statusCode:    http.StatusFound,
		},
		{
			name:          "not found",
			configurers:   []mapper.MapperConfigurer{mockConfigurer},
			persistorName: "mock",
			path:          "invalid",
			statusCode:    http.StatusNotFound,
		},
		{
			name:          "precedence",
			configurers:   []mapper.MapperConfigurer{mockConfigurer, mockConfigurerAlt},
			persistorName: "mock",
			path:          "fk",
			redirectUrl:   fakePair.Url,
			statusCode:    http.StatusFound,
		},
		{
			name:          "precedence 2",
			configurers:   []mapper.MapperConfigurer{mockConfigurerAlt, mockConfigurer},
			persistorName: "mock",
			path:          "fk",
			redirectUrl:   fakePairAlt.Url,
			statusCode:    http.StatusFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm, err := mapper.NewMapperManager(test.persistorName, test.configurers)
			assert.NoError(t, err)
			assert.NotNil(t, mm)
			server, err := NewServer(mm, "8080")
			assert.NoError(t, err)
			assert.NotNil(t, server)

			req, err := http.NewRequest("GET", "/"+test.path, nil)
			assert.NoError(t, err)
			assert.NotNil(t, req)

			r := mux.NewRouter()
			r.HandleFunc("/{path}", server.handleRedirect).Methods("GET")

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, test.statusCode, rr.Code)
			if test.statusCode == http.StatusFound {
				assert.Equal(t, test.redirectUrl, rr.Header().Get("Location"))
			}
		})
	}
}
