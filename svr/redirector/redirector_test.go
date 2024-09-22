package redirector

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gorilla/mux"
// 	"github.com/reimirno/golinks/pkg/mapper"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"

// 	"github.com/reimirno/golinks/pkg/types"
// )

// // MockMapperManager is a mock implementation of mapper.MapperManager
// type MockMapperManager struct {
// 	mock.Mock
// }

// func (m *MockMapperManager) GetUrl(path string, incrementUseCount bool) (*types.PathUrlPair, error) {
// 	args := m.Called(path, incrementUseCount)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*types.PathUrlPair), args.Error(1)
// }

// func (m *MockMapperManager) DeleteUrl(path string) error {
// 	args := m.Called(path)
// 	return args.Error(0)
// }

// func (m *MockMapperManager) GetName() string {
// 	return "MockMapperManager"
// }

// func (m *MockMapperManager) GetType() string {
// 	return "MockMapperManager"
// }

// func (m *MockMapperManager) ListUrls() ([]*types.PathUrlPair, error) {
// 	args := m.Called()
// 	return args.Get(0).([]*types.PathUrlPair), args.Error(1)
// }

// var _ mapper.MapperManager = (*MockMapperManager)(nil)

// func TestNewServer(t *testing.T) {
// 	m := &MockMapperManager{}
// 	server, err := NewServer(m, "8080")

// 	assert.NoError(t, err)
// 	assert.NotNil(t, server)
// 	assert.Equal(t, ":8080", server.server.Addr)
// 	assert.Equal(t, "8080", server.port)
// 	assert.Equal(t, m, server.manager)
// }

// func TestServer_GetName(t *testing.T) {
// 	server := &Server{}
// 	assert.Equal(t, redirectorServiceName, server.GetName())
// }

// func TestServer_handleRedirect(t *testing.T) {
// 	tests := []struct {
// 		name           string
// 		path           string
// 		mockReturn     *types.PathUrlPair
// 		mockError      error
// 		expectedStatus int
// 		expectedURL    string
// 	}{
// 		{
// 			name:           "Successful redirect",
// 			path:           "example",
// 			mockReturn:     &types.PathUrlPair{Url: "https://example.com"},
// 			mockError:      nil,
// 			expectedStatus: http.StatusFound,
// 			expectedURL:    "https://example.com",
// 		},
// 		{
// 			name:           "Mapping not found",
// 			path:           "nonexistent",
// 			mockReturn:     nil,
// 			mockError:      nil,
// 			expectedStatus: http.StatusNotFound,
// 			expectedURL:    "",
// 		},
// 		{
// 			name:           "Error occurred",
// 			path:           "error",
// 			mockReturn:     nil,
// 			mockError:      assert.AnError,
// 			expectedStatus: http.StatusInternalServerError,
// 			expectedURL:    "",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			m := &MockMapperManager{}
// 			m.On("GetUrl", tt.path, true).Return(tt.mockReturn, tt.mockError)

// 			server, _ := NewServer(m, "8080")

// 			req, _ := http.NewRequest("GET", "/"+tt.path, nil)
// 			rr := httptest.NewRecorder()

// 			router := mux.NewRouter()
// 			router.HandleFunc("/{path}", server.handleRedirect)
// 			router.ServeHTTP(rr, req)

// 			assert.Equal(t, tt.expectedStatus, rr.Code)
// 			if tt.expectedURL != "" {
// 				assert.Equal(t, tt.expectedURL, rr.Header().Get("Location"))
// 			}

// 			m.AssertExpectations(t)
// 		})
// 	}
// }
