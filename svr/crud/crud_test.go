package crud

import (
	"context"
	"testing"

	"github.com/reimirno/golinks/pkg/mapper"
	"github.com/reimirno/golinks/pkg/pb"
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
		debug         bool
		wantErr       bool
	}{
		{
			name:          "happy path",
			configurers:   []*mapper.MockMapperConfigurer{mockConfigurer},
			persistorName: "mock",
			port:          "8081",
			debug:         true,
			wantErr:       false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm, err := mapper.NewMapperManager(test.persistorName, mapper.CloneConfigurers(test.configurers))
			assert.NoError(t, err)
			got, err := NewServer(mm, test.port, test.debug)
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
	tests := []struct {
		name string
		want string
	}{
		{
			name: "happy path",
			want: "crud",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NewServer(nil, "", false)
			assert.NoError(t, err)
			assert.Equal(t, test.want, got.GetName())
		})
	}
}

func TestServer_GetUrl(t *testing.T) {
	tests := []struct {
		name          string
		configurers   []*mapper.MockMapperConfigurer
		persistorName string
		path          string
		want          *types.PathUrlPair
		wantErr       bool
	}{
		{
			name:          "happy path",
			configurers:   []*mapper.MockMapperConfigurer{mockConfigurer},
			persistorName: "mock",
			path:          "fk",
			want:          fakePair,
		},
		{
			name:          "path not found",
			configurers:   []*mapper.MockMapperConfigurer{mockConfigurer},
			persistorName: "mock",
			path:          "invalid",
			wantErr:       true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm, err := mapper.NewMapperManager(test.persistorName, mapper.CloneConfigurers(test.configurers))
			assert.NoError(t, err)
			server, err := NewServer(mm, "8081", false)
			assert.NoError(t, err)

			resp, err := server.GetUrl(context.Background(), &pb.GetUrlRequest{Path: test.path})
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				canonicalPath, err := sanitizer.CanonicalizePath(test.want.Path)
				assert.NoError(t, err)
				assert.Equal(t, canonicalPath, resp.GetPath())
				assert.Equal(t, test.want.Url, resp.GetUrl())
			}
		})
	}
}

func TestServer_ListUrls(t *testing.T) {
	tests := []struct {
		name          string
		configurers   []*mapper.MockMapperConfigurer
		persistorName string
		wantErr       bool
		numPairs      int
	}{
		{
			name:          "happy path",
			configurers:   []*mapper.MockMapperConfigurer{mockConfigurer},
			persistorName: "mock",
			wantErr:       false,
			numPairs:      2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm, err := mapper.NewMapperManager(test.persistorName, mapper.CloneConfigurers(test.configurers))
			assert.NoError(t, err)
			server, err := NewServer(mm, "8081", false)
			assert.NoError(t, err)

			resp, err := server.ListUrls(context.Background(), &pb.ListUrlsRequest{
				Offset: 0,
				Limit:  100,
			})
			assert.NoError(t, err)
			assert.Equal(t, test.numPairs, len(resp.GetPairs()))
		})
	}
}

func TestServer_PutUrl(t *testing.T) {
	tests := []struct {
		name          string
		configurers   []*mapper.MockMapperConfigurer
		persistorName string
		pair          *types.PathUrlPair
		wantErr       bool
		want          *types.PathUrlPair
	}{
		{
			name:          "happy path update",
			configurers:   []*mapper.MockMapperConfigurer{mockConfigurer},
			persistorName: "mock",
			pair:          fakePairAlt,
			wantErr:       false,
			want:          fakePairAlt,
		},
		{
			name:          "happy path create",
			configurers:   []*mapper.MockMapperConfigurer{mockConfigurer},
			persistorName: "mock",
			pair:          fakePair3,
			wantErr:       false,
			want:          fakePair3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm, err := mapper.NewMapperManager(test.persistorName, mapper.CloneConfigurers(test.configurers))
			assert.NoError(t, err)
			server, err := NewServer(mm, "8081", false)
			assert.NoError(t, err)

			resp, err := server.PutUrl(context.Background(), &pb.PathUrlPair{
				Path: test.pair.Path,
				Url:  test.pair.Url,
			})
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				canonicalPath, err := sanitizer.CanonicalizePath(test.want.Path)
				assert.NoError(t, err)
				assert.Equal(t, canonicalPath, resp.GetPath())
				assert.Equal(t, test.want.Url, resp.GetUrl())
			}

			resp, err = server.GetUrl(context.Background(), &pb.GetUrlRequest{Path: test.pair.Path})
			assert.NoError(t, err)
			canonicalPath, err := sanitizer.CanonicalizePath(test.want.Path)
			assert.NoError(t, err)
			assert.Equal(t, canonicalPath, resp.GetPath())
			assert.Equal(t, test.want.Url, resp.GetUrl())
		})
	}
}

func TestServer_DeleteUrl(t *testing.T) {
	tests := []struct {
		name          string
		configurers   []*mapper.MockMapperConfigurer
		persistorName string
		path          string
		wantErr       bool
		want          string
	}{
		{
			name:          "happy path",
			configurers:   []*mapper.MockMapperConfigurer{mockConfigurer},
			persistorName: mockConfigurer.Name,
			path:          "fk",
			wantErr:       false,
			want:          "",
		},
		{
			name:          "path not found is fine",
			configurers:   []*mapper.MockMapperConfigurer{mockConfigurer},
			persistorName: mockConfigurer.Name,
			path:          "invalid",
			wantErr:       false,
			want:          "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm, err := mapper.NewMapperManager(test.persistorName, mapper.CloneConfigurers(test.configurers))
			assert.NoError(t, err)
			server, err := NewServer(mm, "8081", false)
			assert.NoError(t, err)

			_, err = server.DeleteUrl(context.Background(), &pb.DeleteUrlRequest{Path: test.path})
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			resp, err := server.GetUrl(context.Background(), &pb.GetUrlRequest{Path: test.path})
			if test.want != "" {
				assert.NoError(t, err)
				assert.Equal(t, test.want, resp.GetUrl())
			} else {
				assert.Error(t, err)
			}
		})
	}
}
