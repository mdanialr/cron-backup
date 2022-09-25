package logger_test

import (
	"testing"

	"github.com/mdanialr/go-cron-backup/pkg/logger"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestInitInfoLogger(t *testing.T) {
	testCases := []struct {
		name       string
		sampleConf *viper.Viper
		wantErr    bool
	}{
		{
			name:       "Should error when using inaccessible directory",
			sampleConf: fakeErrorViper,
			wantErr:    true,
		},
		{
			name:       "Should pass when using valid and accessible directory",
			sampleConf: fakeOKViper,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := logger.InitInfoLogger(tt.sampleConf)

			switch tt.wantErr {
			case false:
				require.NoError(t, err)
			case true:
				require.Error(t, err)
			}
		})
	}
}

func TestInitErrorLogger(t *testing.T) {
	testCases := []struct {
		name       string
		sampleConf *viper.Viper
		wantErr    bool
	}{
		{
			name:       "Should error when using inaccessible directory",
			sampleConf: fakeErrorViper,
			wantErr:    true,
		},
		{
			name:       "Should pass when using valid and accessible directory",
			sampleConf: fakeOKViper,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := logger.InitErrorLogger(tt.sampleConf)

			switch tt.wantErr {
			case false:
				require.NoError(t, err)
			case true:
				require.Error(t, err)
			}
		})
	}
}
