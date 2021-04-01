package warehouse

import (
	"testing"

	"github.com/fullstorydev/hauser/config"
)

func makeConf(databaseSchema string) *config.RedshiftConfig {
	return &config.RedshiftConfig{
		DatabaseSchema: databaseSchema,
		VarCharMax:     20,
		ExportTable:    "exportTable",
		SyncTable:      "syncTable",
	}
}

func TestRedshiftValueToString(t *testing.T) {
	wh := &Redshift{
		conf: makeConf("some_schema"),
	}

	var testCases = []struct {
		input    interface{}
		isTime   bool
		expected string
	}{
		{"short string", false, "short string"},
		{"I'm too long, truncate me", false, "I'm too long, trunc"},
		{"no\nnew\nlines", false, "no new lines"},
		{"no\x00null\x00chars", false, "nonullchars"},
		{5, false, "5"},
		{"2009-11-10T23:00:00.000Z", true, "2009-11-10 23:00:00 +0000 UTC"},
	}

	for _, testCase := range testCases {
		if got := wh.ValueToString(testCase.input, testCase.isTime); got != testCase.expected {
			t.Errorf("Expected value %q, got %q", testCase.expected, got)
		}
	}
}

func TestValidateSchemaConfig(t *testing.T) {

	testCases := []struct {
		conf       *config.RedshiftConfig
		hasError   bool
		errMessage string
	}{
		{
			conf:       makeConf(""),
			hasError:   true,
			errMessage: "DatabaseSchema definition missing from Redshift configuration. More information: https://github.com/fullstorydev/hauser/blob/master/Redshift.md#database-schema-configuration",
		},
		{
			conf:       makeConf("test"),
			hasError:   false,
			errMessage: "",
		},
		{
			conf:       makeConf("search_path"),
			hasError:   false,
			errMessage: "",
		},
	}

	for _, tc := range testCases {
		wh := NewRedshift(tc.conf)
		err := wh.validateSchemaConfig()
		if tc.hasError && err == nil {
			t.Errorf("expected Redshift.validateSchemaConfig() to return an error when config.Config.Redshift.DatabaseSchema is empty")
		}
		if tc.hasError && err.Error() != tc.errMessage {
			t.Errorf("expected Redshift.validateSchemaConfig() to return \n%s \nwhen config.Config.Redshift.DatabaseSchema is empty, returned \n%s \ninstead", tc.errMessage, err)
		}
		if !tc.hasError && err != nil {
			t.Errorf("unexpected error thrown for DatabaseSchema %s: %s", tc.conf.DatabaseSchema, err)
		}
	}
}
