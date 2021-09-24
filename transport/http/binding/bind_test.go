package binding

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-kratos/kratos/v2/internal/testdata/encoding"
)

func TestBindFormRepeated(t *testing.T) {
	req := new(encoding.TestModel)
	err := BindQuery(url.Values{
		"hobby": []string{"movie,game,reading,running"},
	}, req)
	require.NoError(t, err)
	require.Equal(t, req.Hobby, []string{"movie", "game", "reading", "running"})
	req = new(encoding.TestModel)
	err = BindQuery(url.Values{
		"hobby": []string{"movie_game_reading_running"},
	}, req)
	require.NoError(t, err)
	require.Equal(t, req.Hobby, []string{"movie_game_reading_running"})
}
