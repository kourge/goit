package config

import (
	"reflect"
	"strings"
	"testing"
)

type sectionFixture struct {
	Section
	String           string
	NormalizedString string
}

var (
	_fixtureSection1 sectionFixture = sectionFixture{
		Section: Section{"core", Dict{
			"repositoryformatversion": int64(0),
			"filemode":                true,
			"diff":                    "auto",
			"bare":                    false,
			"name":                    "John Doe",
		}},
		String: `[core]
	repositoryformatversion = 0
	# comment 1
	filemode = true

	diff = auto
; comment 2
	bare = false
	name = "John Doe"`,
		NormalizedString: `[core]
	bare = false
	diff = auto
	filemode = true
	name = "John Doe"
	repositoryformatversion = 0
`,
	}
)

func TestSection_String(t *testing.T) {
	var actual string = _fixtureSection1.Section.String()
	var expected string = _fixtureSection1.NormalizedString

	if actual != expected {
		t.Errorf("section.String() = %v, want %v", actual, expected)
	}
}

func TestSection_Decode(t *testing.T) {
	var actual *Section = &Section{}
	var expected *Section = &_fixtureSection1.Section

	actual.Decode(strings.NewReader(_fixtureSection1.String))

	if !reflect.DeepEqual(*actual, *expected) {
		t.Errorf("section.Decode() produced %v, want %v", *actual, *expected)
	}
}
