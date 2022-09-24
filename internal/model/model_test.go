package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabase_setDefault(t *testing.T) {
	testCases := []struct {
		name   string
		sample Database
		expect Database
	}{
		{
			name:   "Given empty host should set to localhost",
			sample: Database{},
			expect: Database{Host: "localhost"},
		},
		{
			name:   "Given empty port should set to 3306 if the type is either my or md",
			sample: Database{Type: "my"},
			expect: Database{Type: "my", Host: "localhost", Port: 3306, User: "root"},
		},
		{
			name:   "Given empty user should set to root if the type is either my or md",
			sample: Database{Type: "my"},
			expect: Database{Type: "my", Host: "localhost", Port: 3306, User: "root"},
		},
		{
			name:   "Given empty port should set to 5432 if the type is pg",
			sample: Database{Type: "pg"},
			expect: Database{Type: "pg", Host: "localhost", Port: 5432, User: "postgres"},
		},
		{
			name:   "Given empty user should set to postgres if the type is pg",
			sample: Database{Type: "pg"},
			expect: Database{Type: "pg", Host: "localhost", Port: 5432, User: "postgres"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.sample.setDefault()
			assert.Equal(t, tc.expect, tc.sample)
		})
	}
}

func TestDatabase_buildID(t *testing.T) {
	testCases := []struct {
		name   string
		sample Database
		expect string
	}{
		{
			name:   "Given 'type' pg empty 'host' and 'port' also 'name' sample should has expected generated id",
			sample: Database{Type: "pg", Name: "sample"},
			expect: "pg_localhost_5432_sample",
		},
		{
			name:   "Given 'type' my empty 'host' 'port' 3376 also 'name' db should has expected generated id",
			sample: Database{Type: "my", Port: 3376, Name: "db"},
			expect: "my_localhost_3376_db",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.sample.setDefault()
			tc.sample.buildID()
			assert.Equal(t, tc.expect, tc.sample.ID)
		})
	}
}

func TestDatabase_buildCMD(t *testing.T) {
	testCases := []struct {
		name   string
		sample Database
		expect string
	}{
		{
			name:   "Given 'type' pg empty 'host' and 'port' also 'name' sample should has expected generated cmd",
			sample: Database{Type: "pg", Name: "sample"},
			expect: "pg_dump postgresql://postgres:@localhost:5432/sample?",
		},
		{
			name:   "Given 'type' my empty 'host' and 'pass' but has 'port' 3376 also 'name' db should has expected generated cmd",
			sample: Database{Type: "my", Port: 3376, Name: "db"},
			expect: "mysqldump -h localhost -P 3376 -uroot   db",
		},
		{
			name:   "Given 'type' md 'host' 42.123.23.1 'port' 3376 'pass' secret also 'name' db should has expected generated cmd",
			sample: Database{Type: "md", Host: "42.123.23.1", Port: 3376, Pass: "secret", Name: "db"},
			expect: "mariadb-dump -h 42.123.23.1 -P 3376 -uroot -psecret  db",
		},
		{
			name:   "Given sample 'docker' is not empty should has prefix `docker exec -t docker-name`, then followed by the generated cmd",
			sample: Database{Type: "pg", Name: "sample", Docker: "docker-name"},
			expect: "docker exec -t docker-name pg_dump postgresql://postgres:@localhost:5432/sample?",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.sample.setDefault()
			cmd := tc.sample.buildCMD()
			assert.Equal(t, tc.expect, cmd)
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	appSample := APP{Apps: []*App{{Name: "app1"}, {Name: "app2"}}}
	dbSample := DB{Databases: []*Database{
		{Type: "pg", Name: "db1", BackupName: "db1"},
		{Type: "my", Name: "db2", BackupName: "db2"},
		{Type: "md", Name: "db3", BackupName: "db3"},
		{Type: "pg", Name: "db4", Docker: "docker-name", BackupName: "db4"},
	}}

	testCases := []struct {
		name      string
		sample    Config
		wantErr   bool
		expectMsg string
	}{
		{
			name:      "Should fail when given Config has duplicate database backup name",
			sample:    Config{DB: DB{Databases: []*Database{{BackupName: "sample"}, {BackupName: "sample"}}}},
			wantErr:   true,
			expectMsg: "found duplicate database backup name: duplicate backup name (sample)",
		},
		{
			name:      "Should fail when given Config has duplicate app backup name",
			sample:    Config{APP: APP{Apps: []*App{{Name: "sample"}, {Name: "sample"}}}},
			wantErr:   true,
			expectMsg: "found duplicate app backup name: duplicate app name (sample)",
		},
		{
			name:      "Should fail when given Config has database unsupported prefix type, like should be 'md' instead of 'mariadb'",
			sample:    Config{DB: DB{Databases: []*Database{{BackupName: "sample", Type: "mariadb"}}}},
			wantErr:   true,
			expectMsg: "unsupported database type: 'mariadb'. currently supported types are [pg,my,md]",
		},
		{
			name:   "Should pass when there are no duplicate in either databases or apps, and the database type is supported",
			sample: Config{APP: appSample, DB: dbSample},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.sample.Validate()
			if tc.wantErr {
				require.Error(t, err)
				assert.Equal(t, tc.expectMsg, err.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}
