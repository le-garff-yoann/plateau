package pflag

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func TestLogrusLevel(t *testing.T) {
	t.Parallel()

	require.Implements(t, (*pflag.Value)(nil), new(LogrusLevel))

	lvl := LogrusLevel(logrus.WarnLevel)

	require.Equal(t, logrus.WarnLevel.String(), lvl.String())

	require.NoError(t, lvl.Set("info"))
	require.Equal(t, LogrusLevel(logrus.InfoLevel), lvl)

	require.Error(t, lvl.Set("hopeitdoesnotexist"))
}

func TestLogrusLevelType(t *testing.T) {
	t.Parallel()

	lvl := LogrusLevel(logrus.WarnLevel)
	require.Equal(t, "logruslevel", lvl.Type())
}
