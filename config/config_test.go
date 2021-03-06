package config

import (
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

type configFixture struct {
	Config
	String           string
	NormalizedString string
}

var (
	_fixtureConfig configFixture = configFixture{
		Config: Config{
			"user": {"user", Dict{
				"name":  "Jane Doe",
				"email": "jane@example.com",
			}},
			"core": {"core", Dict{
				"repositoryformatversion": int64(0),
				"filemode":                true,
				"diff":                    "auto",
				"bare":                    false,
			}},
		},
		String: `
[user]
	name =   Jane Doe
	email =     jane@example.com

	[core]

repositoryformatversion = 0
	filemode = true
diff=auto
	bare=false
`,
		NormalizedString: `[core]
	bare = false
	diff = auto
	filemode = true
	repositoryformatversion = 0

[user]
	email = jane@example.com
	name = "Jane Doe"
`,
	}
)

func TestConfig_String(t *testing.T) {
	var actual string = _fixtureConfig.Config.String()
	var expected string = _fixtureConfig.NormalizedString

	if actual != expected {
		t.Errorf("config.String() = %v, want %v", actual, expected)
	}
}

func TestConfig_Decode(t *testing.T) {
	var actual *Config = &Config{}
	var expected *Config = &_fixtureConfig.Config

	actual.Decode(strings.NewReader(_fixtureConfig.String))

	if !reflect.DeepEqual(*actual, *expected) {
		t.Errorf("config.Decode() produced %v, want %v", *actual, *expected)
	}
}

func ExampleConfig_Reader() {
	defaultConfig := Config{
		"core": {"core", Dict{
			"repositoryformatversion": int64(0),
			"filemode":                true,
			"bare":                    false,
			"logallrefupdates":        true,
			"ignorecase":              true,
			"precomposeunicode":       false,
		}},
	}

	io.Copy(os.Stdout, defaultConfig.Reader())
	// Output:
	// [core]
	// 	bare = false
	// 	filemode = true
	// 	ignorecase = true
	// 	logallrefupdates = true
	// 	precomposeunicode = false
	// 	repositoryformatversion = 0
}
