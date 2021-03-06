package core

import (
	"bytes"
	"reflect"
	"testing"
)

var (
	_fixtureCommitTree Sha1 = Sha1{
		0x93, 0x5e, 0x0a, 0x5c, 0x83, 0x61, 0xe5, 0x9f, 0x8b, 0xbc,
		0x01, 0xb2, 0xdb, 0xfb, 0xec, 0x3a, 0x44, 0xe2, 0x49, 0x04,
	}
	_fixtureCommitParents []Sha1 = []Sha1{{
		0x77, 0x5c, 0x72, 0x28, 0x62, 0x15, 0x59, 0x62, 0x34, 0x06,
		0x85, 0x7d, 0x18, 0x10, 0xa3, 0x15, 0x36, 0x16, 0x33, 0x6f,
	}}
	_fixtureCommitAuthor Person = NewPerson(
		"Kosuke Asami", "tfortress58@gmail.com", 1395160458, 9*3600,
	)
	_fixtureCommitCommitter Person = NewPerson(
		"Jack Nagel", "jacknagel@gmail.com", 1395293290, -5*3600,
	)
	_fixtureCommitMessage string = `byobu 5.75

This release includes fixes about prefix problem that is discussed
in #27045.

Closes #27667.

Signed-off-by: Jack Nagel <jacknagel@gmail.com>`
	_fixtureCommit *Commit = NewCommit(
		_fixtureCommitTree,
		_fixtureCommitParents,
		_fixtureCommitAuthor,
		_fixtureCommitCommitter,
		_fixtureCommitMessage,
	)
	_fixtureCommitString string = `tree 935e0a5c8361e59f8bbc01b2dbfbec3a44e24904
parent 775c7228621559623406857d1810a3153616336f
author Kosuke Asami <tfortress58@gmail.com> 1395160458 +0900
committer Jack Nagel <jacknagel@gmail.com> 1395293290 -0500

byobu 5.75

This release includes fixes about prefix problem that is discussed
in #27045.

Closes #27667.

Signed-off-by: Jack Nagel <jacknagel@gmail.com>
`
	_fixtureCommit2 *Commit = NewCommit(
		Sha1{
			0x2d, 0x95, 0x50, 0xa9, 0x1d, 0x46, 0x74, 0xe6, 0x4d, 0x03,
			0xe2, 0xe6, 0x52, 0x9f, 0xdf, 0xeb, 0xd3, 0x6d, 0x91, 0x37,
		},
		[]Sha1{{
			0x31, 0xbb, 0x0f, 0x62, 0x27, 0x5f, 0xf0, 0xae, 0xbe, 0xc0,
			0x2a, 0x93, 0xb9, 0xda, 0x79, 0xe7, 0x69, 0x04, 0x3d, 0xee,
		}},
		NewPerson(
			"Matthew Hawkins", "darthmdh@gmail.com", 1396097804, 11*3600,
		),
		NewPerson(
			"Adam Vandenberg", "flangy@gmail.com", 1396328662, -7*3600,
		),
		`vim: add luajit support

Add --with-luajit which passes through to the ./configure script
to tell Vim to link against luajit instead of lua.

Closes #27967.

Signed-off-by: Adam Vandenberg <flangy@gmail.com>`,
	)
)

func TestCommit_ParentsEqual(t *testing.T) {
	if !_fixtureCommit.ParentsEqual(_fixtureCommit) {
		t.Error("commit.Parents() should equal commit.Parents()")
	}

	if _fixtureCommit.ParentsEqual(_fixtureCommit2) {
		t.Errorf("%v should not share parents with %v", _fixtureCommit, _fixtureCommit2)
	}
}

func TestCommit_Equal(t *testing.T) {
	if !_fixtureCommit.Equal(_fixtureCommit) {
		t.Error("commit should equal self")
	}

	if _fixtureCommit.Equal(_fixtureCommit2) {
		t.Errorf("%v should not equal %v", _fixtureCommit, _fixtureCommit2)
	}
}

func TestCommit_Type(t *testing.T) {
	var actual string = _fixtureCommit.Type()
	var expected string = "commit"

	if actual != expected {
		t.Errorf("commit.Type() = %v, want %v", actual, expected)
	}
}

func TestCommit_Size(t *testing.T) {
	var actual int = _fixtureCommit.Size()
	var expected int = len(_fixtureCommitString)

	if actual != expected {
		t.Errorf("commit.Size() = %v, want %v", actual, expected)
	}
}

func TestCommit_Reader(t *testing.T) {
	var actual []byte
	var expected []byte = []byte(_fixtureCommitString)

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(_fixtureCommit.Reader())
	actual = buffer.Bytes()

	if !bytes.Equal(actual, expected) {
		t.Error("commit.Reader() did not generate same byte sequence")
	}
}

func TestCommit_Decode(t *testing.T) {
	var actual *Commit = &Commit{}
	var expected *Commit = _fixtureCommit

	reader := bytes.NewBuffer([]byte(_fixtureCommitString))
	err := actual.Decode(reader)

	if err != nil {
		t.Errorf("commit.Decode() returned error %v", err)
	}

	if !reflect.DeepEqual(*actual, *expected) {
		t.Errorf("commit.Decode() produced %v, want %v", actual, expected)
	}
}

func TestCommit_Tree(t *testing.T) {
	var actual Sha1 = _fixtureCommit.Tree()
	var expected Sha1 = _fixtureCommitTree

	if actual != expected {
		t.Errorf("commit.Tree() = %v, want %v", actual, expected)
	}
}

func TestCommit_Parents(t *testing.T) {
	var actual []Sha1 = _fixtureCommit.Parents()
	var expected []Sha1 = _fixtureCommitParents

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("commit.Parents() = %v, want %v", actual, expected)
	}
}

func TestCommit_Author(t *testing.T) {
	var actual Person = _fixtureCommit.Author()
	var expected Person = _fixtureCommitAuthor

	if !actual.Equal(expected) {
		t.Errorf("commit.Author() = %v, want %v", actual, expected)
	}
}

func TestCommit_Committer(t *testing.T) {
	var actual Person = _fixtureCommit.Committer()
	var expected Person = _fixtureCommitCommitter

	if !actual.Equal(expected) {
		t.Errorf("commit.Committer() = %v, want %v", actual, expected)
	}
}

func TestCommit_Message(t *testing.T) {
	var actual string = _fixtureCommit.Message()
	var expected string = _fixtureCommitMessage

	if actual != expected {
		t.Errorf("commit.Message() = %v, want %v", actual, expected)
	}
}
