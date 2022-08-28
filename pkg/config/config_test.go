package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupDefault(t *testing.T) {
	testCases := []struct {
		name    string
		sample  *viper.Viper
		expect  string
		wantErr bool
	}{
		{
			name:    "Should error without root",
			sample:  viperWithoutSecret,
			expect:  "`root` for this app root directories is required",
			wantErr: true,
		},
		{
			name:   "Should pass and other fields has default value as expected",
			sample: viperComplete,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := SetupDefault(tc.sample)

			switch tc.wantErr {
			case true:
				require.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			case false:
				require.NoError(t, err)
				assert.Equal(t, "/tmp/logs", tc.sample.GetString("root"))
				assert.Equal(t, 6, tc.sample.GetInt("max_days"))
				assert.Equal(t, 1, tc.sample.GetInt("db.max_worker"))
				assert.Equal(t, 6, tc.sample.GetInt("db.max_days"))
				assert.Equal(t, 1, tc.sample.GetInt("app.max_worker"))
				assert.Equal(t, 6, tc.sample.GetInt("app.max_days"))
			}
		})
	}
}

func TestInitConfig(t *testing.T) {
	testCases := []struct {
		name    string
		sample  string
		wantErr bool
	}{
		{
			name:    "Should error when config file not found",
			sample:  "/fake/path",
			wantErr: true,
		},
		{
			name:   "Should pass when config file found",
			sample: "/tmp/",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := InitConfig(tc.sample)

			switch tc.wantErr {
			case true:
				require.Error(t, err)
			case false:
				require.NoError(t, err)
			}
		})
	}
}
