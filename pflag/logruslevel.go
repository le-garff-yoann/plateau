package pflag

import "github.com/sirupsen/logrus"

// LogrusLevel is a log level value.
type LogrusLevel logrus.Level

// Type implements the `pflag.Value` interface.
func (s *LogrusLevel) Type() string {
	return "logruslevel"
}

// Set implements the `pflag.Value` interface.
func (s *LogrusLevel) Set(val string) error {
	logrusLevel, err := logrus.ParseLevel(val)

	*s = LogrusLevel(logrusLevel)

	return err
}

func (s LogrusLevel) String() string {
	return logrus.Level(logrus.Level(s)).String()
}
