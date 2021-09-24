package form

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-kratos/kratos/v2/internal/testdata/encoding"
)

func TestPopulateRepeatedField(t *testing.T) {
	want := []string{"a", "b", "c"}
	require.Equal(t, populateValues("a,b,c"), want)
	require.Equal(t, populateValues("a|b|c"), want)
	require.Equal(t, populateValues("a b c"), want)
	require.Equal(t, populateValues("a\tb\tc"), want)
	require.Equal(t, populateValues("string_do_not_need_split"), []string{"string_do_not_need_split"})
}

func TestRepeatedFieldMapProto(t *testing.T) {
	tests := []string{
		"movie|game|reading|running",
		"movie,game,reading,running",
		"movie game reading running",
		"movie\tgame\treading\trunning",
	}
	want := []string{"movie", "game", "reading", "running"}
	for _, tt := range tests {
		message := new(encoding.TestModel)
		err := MapProto(message, map[string][]string{
			"hobby": {tt},
		})
		require.NoError(t, err)
		require.Equal(t, message.Hobby, want)
	}

	normal := "movie_game_reading_running"
	message := new(encoding.TestModel)
	err := MapProto(message, map[string][]string{
		"hobby": {normal},
	})
	require.NoError(t, err)
	require.Equal(t, message.Hobby, []string{normal})
}
