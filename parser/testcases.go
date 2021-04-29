package parser

import (
	"fmt"

	"github.com/leodido/go-conventionalcommits"
	cctesting "github.com/leodido/go-conventionalcommits/testing"
)

type testCase struct {
	title        string
	input        []byte
	ok           bool
	value        conventionalcommits.Message
	partialValue conventionalcommits.Message
	errorString  string
}

var testCases = []testCase{
	// INVALID / empty
	{
		"empty",
		[]byte(""),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEmpty+ColumnPositionTemplate, 0),
	},
	// INVALID / invalid type (1 char)
	{
		"invalid-type-1-char",
		[]byte("f"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "f", 1),
	},
	// INVALID / invalid type (2 char)
	{
		"invalid-type-2-char",
		[]byte("fx"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "x", 1),
	},
	// INVALID / invalid type (2 char) with almost valid type
	{
		"invalid-type-2-char-feat",
		[]byte("fe"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "e", 2),
	},
	// INVALID / invalid type (3 char)
	{
		"invalid-type-3-char",
		[]byte("fit"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "t", 2),
	},
	// INVALID / invalid type (3 char) again
	{
		"invalid-type-3-char-feat",
		[]byte("fei"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "i", 2),
	},
	// INVALID / invalid type (3 char) with almost valid type
	{
		"invalid-type-3-char-feat",
		[]byte("fea"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "a", 3),
	},
	// INVALID / invalid type (4 char)
	{
		"invalid-type-4-char",
		[]byte("feax"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "x", 3),
	},
	// INVALID / missing colon after type fix
	{
		"invalid-after-valid-type-fix",
		[]byte("fix"),
		false,
		nil,
		nil, // no partial result because it is not a minimal valid commit message
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "x", 2),
	},
	// INVALID / missing colon after type feat
	{
		"invalid-after-valid-type-feat",
		[]byte("feat"),
		false,
		nil,
		nil, // no partial result because it is not a minimal valid commit message
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "t", 3),
	},
	// INVALID / invalid type (2 char) + colon
	{
		"invalid-type-2-char-colon",
		[]byte("fi:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, ":", 2),
	},
	// INVALID / invalid type (3 char) + colon
	{
		"invalid-type-3-char-colon",
		[]byte("fea:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, ":", 3),
	},
	// VALID / minimal commit message
	{
		"valid-minimal-commit-message",
		[]byte("fix: x"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
		},
		"",
	},
	// INVALID / missing colon after valid commit message type
	{
		"missing-colon-after-type-3-chars",
		[]byte("fix>"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, ">", 3),
	},
	// INVALID / missing colon after valid commit message type
	{
		"missing-colon-after-type-4-chars",
		[]byte("feat?"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "?", 4),
	},
	// INVALID / invalid after valid type and scope
	{
		"invalid-after-valid-type-and-scope",
		[]byte("fix(scope)"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, ")", 9),
	},
	// VALID / type + scope + description
	{
		"valid-with-scope",
		[]byte("fix(aaa): bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
		},
		"",
	},
	// VALID / type + scope + multiple whitespaces + description
	{
		"valid-with-scope-multiple-whitespaces",
		[]byte("fix(aaa):          bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
		},
		"",
	},
	// VALID / type + scope + breaking + description
	{
		"valid-breaking-with-scope",
		[]byte("fix(aaa)!: bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			Exclamation: true,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			Exclamation: true,
		},
		"",
	},
	// VALID / empty scope is ignored
	{
		"valid-empty-scope-is-ignored",
		[]byte("fix(): bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
		},
		"",
	},
	// VALID / type + empty scope + breaking + description
	{
		"valid-breaking-with-empty-scope",
		[]byte("fix()!: bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			Exclamation: true,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			Exclamation: true,
		},
		"",
	},
	// VALID / type + breaking + description
	{
		"valid-breaking-without-scope",
		[]byte("fix!: bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			Exclamation: true,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			Exclamation: true,
		},
		"",
	},
	// INVALID / missing whitespace after colon (with breaking)
	{
		"invalid-missing-ws-after-colon-with-breaking",
		[]byte("fix!:a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 5),
	},
	// INVALID / missing whitespace after colon with scope
	{
		"invalid-missing-ws-after-colon-with-scope",
		[]byte("fix(x):a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 7),
	},
	// INVALID / missing whitespace after colon with empty scope
	{
		"invalid-missing-ws-after-colon-with-empty-scope",
		[]byte("fix():a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 6),
	},
	// INVALID / missing whitespace after colon
	{
		"invalid-missing-ws-after-colon",
		[]byte("fix:a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 4),
	},
	// INVALID / invalid initial character
	{
		"invalid-initial-character",
		[]byte("(type: a description"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "(", 0),
	},
	// INVALID / invalid second character
	{
		"invalid-second-character",
		[]byte("f description"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, " ", 1),
	},
	// INVALID / invalid after valid type, scope, and breaking
	{
		"invalid-after-valid-type-scope-and-breaking",
		[]byte("fix(scope)!"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "!", 10),
	},
	// INVALID / invalid after valid type, scope, and colon
	{
		"invalid-after-valid-type-scope-and-colon",
		[]byte("fix(scope):"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, ":", 10),
	},
	// INVALID / invalid after valid type, scope, breaking, and colon
	{
		"invalid-after-valid-type-scope-breaking-and-colon",
		[]byte("fix(scope)!:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, ":", 11),
	},
	// INVALID / invalid after valid type, scope, breaking, colon, and white-space
	{
		"invalid-after-valid-type-scope-breaking-colon-and-space",
		[]byte("fix(scope)!: "),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescription+ColumnPositionTemplate, " ", 13),
	},
	// INVALID / invalid after valid type, scope, breaking, colon, and white-spaces
	{
		"invalid-after-valid-type-scope-breaking-colon-and-spaces",
		[]byte("fix(scope)!:  "),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescription+ColumnPositionTemplate, " ", 14),
	},
	// INVALID / double left parentheses in scope
	{
		"invalid-double-left-parentheses-scope",
		[]byte("fix(("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrMalformedScope+ColumnPositionTemplate, "(", 4),
	},
	// INVALID / double left parentheses in scope after valid character
	{
		"invalid-double-left-parentheses-scope-after-valid-character",
		[]byte("fix(a("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrMalformedScope+ColumnPositionTemplate, "(", 5),
	},
	// INVALID / double right parentheses in place of an exclamation, or a colon
	{
		"invalid-double-right-parentheses-scope",
		[]byte("fix(a))"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, ")", 6),
	},
	// INVALID / new left parentheses after valid scope
	{
		"invalid-new-left-parentheses-after-valid-scope",
		[]byte("feat(az)("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "(", 8),
	},
	// INVALID / newline rather than whitespace in description
	{
		"invalid-newline-rather-than-whitespace-description",
		[]byte("feat(az):\x0A description on newline"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "\n", 9),
	},
	// INVALID / newline after whitespace in description
	{
		"invalid-newline-after-whitespace-description",
		[]byte("feat(az): \x0Adescription on newline"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrNewline+ColumnPositionTemplate, 11),
	},
	// INVALID / newline in the description
	// VALID / until the newline
	{
		"invalid-newline-in-description",
		[]byte("feat(az): new\x0Aline"),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Scope:       cctesting.StringAddress("az"),
			Description: "new",
		},
		fmt.Sprintf(ErrMissingBlankLineAtBeginning+ColumnPositionTemplate, 14),
	},
	// INVALID / newline in the description
	// VALID / until the newline
	{
		"invalid-newline-in-description-2",
		[]byte("feat(az)!: bla\x0Al"),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Scope:       cctesting.StringAddress("az"),
			Exclamation: true,
			Description: "bla",
		},
		fmt.Sprintf(ErrMissingBlankLineAtBeginning+ColumnPositionTemplate, 15),
	},
	// INVALID / newline in the description
	// VALID / until the newline
	{
		"description-ending-with-single-newline",
		[]byte("feat(az)!: bla\x0A"),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Scope:       cctesting.StringAddress("az"),
			Exclamation: true,
			Description: "bla",
		},
		fmt.Sprintf(ErrMissingBlankLineAtBeginning+ColumnPositionTemplate, 15),
	},
	// VALID / multi-line body is valid (after a blank line)
	{
		"valid-with-multi-line-body",
		[]byte(`fix: x

see the issue for details

on typos fixed.`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			Body:        cctesting.StringAddress("see the issue for details\n\non typos fixed."),
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			Body:        cctesting.StringAddress("see the issue for details\n\non typos fixed."),
		},
		"",
	},
	// VALID / multi-line body ending with multiple blank lines (they gets discarded) is valid
	{
		"valid-with-multi-line-body-ending-extras-blank-lines",
		[]byte(`fix: x

see the issue for details

on typos fixed.

`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			Body:        cctesting.StringAddress("see the issue for details\n\non typos fixed."),
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			Body:        cctesting.StringAddress("see the issue for details\n\non typos fixed."),
		},
		"",
	},
	// VALID / multi-line body starting with many extra blank lines is valid
	{
		"valid-with-multi-line-body-after-two-extra-blank-lines",
		[]byte(`fix: magic



see the issue for details

on typos fixed.`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "magic",
			Body:        cctesting.StringAddress("\n\nsee the issue for details\n\non typos fixed."),
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "magic",
			Body:        cctesting.StringAddress("\n\nsee the issue for details\n\non typos fixed."),
		},
		"",
	},
	// VALID / multi-line body starting and ending with many extra blank lines is valid
	{
		"valid-with-multi-line-body-with-extra-blank-lines-before-and-after",
		[]byte(`fix: magic



see the issue for details

on typos fixed.


`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "magic",
			Body:        cctesting.StringAddress("\n\nsee the issue for details\n\non typos fixed."),
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "magic",
			Body:        cctesting.StringAddress("\n\nsee the issue for details\n\non typos fixed."),
		},
		"",
	},
	// VALID / single line body (after blank line) is valid
	{
		"valid-with-single-line-body",
		[]byte(`fix: correct minor typos in code

see the issue for details.`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct minor typos in code",
			Body:        cctesting.StringAddress("see the issue for details."),
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct minor typos in code",
			Body:        cctesting.StringAddress("see the issue for details."),
		},
		"",
	},
	// VALID / empty body is okay (it's optional)
	{
		"valid-with-empty-body",
		[]byte(`fix: correct something

`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct something",
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct something",
		},
		"",
	},
	// VALID / multiple blank lines body is okay (it's considered empty)
	{
		"valid-with-multiple-blank-lines-body",
		[]byte(`fix: descr





`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "descr",
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "descr",
		},
		"",
	},
	// VALID / only footer
	{
		"valid-with-footer-only",
		[]byte(`fix: only footer

Fixes #3
Signed-off-by: Leo`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "only footer",
			Footers: map[string][]string{
				"fixes":         {"3"},
				"signed-off-by": {"Leo"},
			},
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "only footer",
			Footers: map[string][]string{
				"fixes":         {"3"},
				"signed-off-by": {"Leo"},
			},
		},
		"",
	},
	// VALID / only footer after many blank lines (that gets ignored)
	{
		"valid-with-footer-only-after-many-blank-lines",
		[]byte(`fix: only footer




Fixes #3
Signed-off-by: Leo`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "only footer",
			Footers: map[string][]string{
				"fixes":         {"3"},
				"signed-off-by": {"Leo"},
			},
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "only footer",
			Footers: map[string][]string{
				"fixes":         {"3"},
				"signed-off-by": {"Leo"},
			},
		},
		"",
	},
	// VALID / only footer ending with many blank lines (that gets ignored)
	{
		"valid-with-footer-only-ending-with-many-blank-lines",
		[]byte(`fix: only footer

Fixes #3
Signed-off-by: Leo


`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "only footer",
			Footers: map[string][]string{
				"fixes":         {"3"},
				"signed-off-by": {"Leo"},
			},
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "only footer",
			Footers: map[string][]string{
				"fixes":         {"3"},
				"signed-off-by": {"Leo"},
			},
		},
		"",
	},
	// VALID / only footer containing repetitions
	{
		"valid-with-footer-containing-repetitions",
		[]byte(`fix: only footer

Fixes #3
Fixes #4
Fixes #5`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "only footer",
			Footers: map[string][]string{
				"fixes": {"3", "4", "5"},
			},
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "only footer",
			Footers: map[string][]string{
				"fixes": {"3", "4", "5"},
			},
		},
		"",
	},
	// VALID / Multi-line body with extras blank lines after and footer with multiple trailers
	{
		"valid-with-multi-line-body-containing-extra-blank-lines-inside-and-after-plus-footer-many-trailers",
		[]byte(`fix: sarah

FUCK

COVID-19.
This is the only message I have in my mind

right now.



Fixes #22
Co-authored-by: My other personality <persona@email.com>
Signed-off-by: Leonardo Di Donato <some@email.com>`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "sarah",
			Body:        cctesting.StringAddress("FUCK\n\nCOVID-19.\nThis is the only message I have in my mind\n\nright now."),
			Footers: map[string][]string{
				"fixes":          {"22"},
				"co-authored-by": {"My other personality <persona@email.com>"},
				"signed-off-by":  {"Leonardo Di Donato <some@email.com>"},
			},
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "sarah",
			Body:        cctesting.StringAddress("FUCK\n\nCOVID-19.\nThis is the only message I have in my mind\n\nright now."),
			Footers: map[string][]string{
				"fixes":          {"22"},
				"co-authored-by": {"My other personality <persona@email.com>"},
				"signed-off-by":  {"Leonardo Di Donato <some@email.com>"},
			},
		},
		"",
	},
	// VALID / Multi-line body with newlines inside and many blank lines after and footer with multiple trailers
	{
		"valid-with-multi-line-body-and-extra-blank-lines-after-plus-footer-many-trailers",
		[]byte(`fix: sarah

FUCK
COVID-19.
This is the only message I have in my mind
right
now.



Fixes #22
Co-authored-by: My other personality <persona@email.com>
Signed-off-by: Leonardo Di Donato <some@email.com>`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "sarah",
			Body:        cctesting.StringAddress("FUCK\nCOVID-19.\nThis is the only message I have in my mind\nright\nnow."),
			Footers: map[string][]string{
				"fixes":          {"22"},
				"co-authored-by": {"My other personality <persona@email.com>"},
				"signed-off-by":  {"Leonardo Di Donato <some@email.com>"},
			},
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "sarah",
			Body:        cctesting.StringAddress("FUCK\nCOVID-19.\nThis is the only message I have in my mind\nright\nnow."),
			Footers: map[string][]string{
				"fixes":          {"22"},
				"co-authored-by": {"My other personality <persona@email.com>"},
				"signed-off-by":  {"Leonardo Di Donato <some@email.com>"},
			},
		},
		"",
	},
	// VALID / Multi-line body with newlines inside and many blank lines before it, plus footer with multiple trailers
	{
		"valid-with-multi-line-body-and-extra-blank-lines-before-plus-footer-many-trailers",
		[]byte(`fix: sarah



FUCK
COVID-19.
This is the only message I have in my mind
right
now.



Fixes #22
Co-authored-by: My other personality <persona@email.com>
Signed-off-by: Leonardo Di Donato <some@email.com>`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "sarah",
			// First blank line ("\n\n") gets ignored
			Body: cctesting.StringAddress("\n\nFUCK\nCOVID-19.\nThis is the only message I have in my mind\nright\nnow."),
			Footers: map[string][]string{
				"fixes":          {"22"},
				"co-authored-by": {"My other personality <persona@email.com>"},
				"signed-off-by":  {"Leonardo Di Donato <some@email.com>"},
			},
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "sarah",
			// First blank line ("\n\n") gets ignored
			Body: cctesting.StringAddress("\n\nFUCK\nCOVID-19.\nThis is the only message I have in my mind\nright\nnow."),
			Footers: map[string][]string{
				"fixes":          {"22"},
				"co-authored-by": {"My other personality <persona@email.com>"},
				"signed-off-by":  {"Leonardo Di Donato <some@email.com>"},
			},
		},
		"",
	},
}

var testCasesForFalcoTypes = []testCase{
	// INVALID / empty
	{
		"empty",
		[]byte(""),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEmpty+ColumnPositionTemplate, 0),
	},
	// INVALID / invalid type (1 char)
	{
		"invalid-type-1-char",
		[]byte("c"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "c", 1),
	},
	// INVALID / invalid type (2 char)
	{
		"invalid-type-2-char",
		[]byte("bx"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "x", 1),
	},
	// INVALID / invalid type (2 char) with almost valid type
	{
		"invalid-type-2-char-feat",
		[]byte("fe"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "e", 2),
	},
	// INVALID / invalid type (2 char) with almost valid type
	{
		"invalid-type-2-char-revert",
		[]byte("re"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "e", 2),
	},
	// INVALID / invalid type (3 char)
	{
		"invalid-type-3-char",
		[]byte("net"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "t", 2),
	},
	// INVALID / invalid type (3 char) again
	{
		"invalid-type-3-char-feat",
		[]byte("fei"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "i", 2),
	},
	// INVALID / invalid type (3 char) with almost valid type
	{
		"invalid-type-3-char-feat",
		[]byte("bui"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "i", 3),
	},
	// INVALID / invalid type (4 char)
	{
		"invalid-type-4-char",
		[]byte("docx"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "x", 3),
	},
	// INVALID / invalid type (4 char)
	{
		"invalid-type-4-char",
		[]byte("perz"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "z", 3),
	},
	// INVALID / missing colon after type fix
	{
		"invalid-after-valid-type-fix",
		[]byte("fix"),
		false,
		nil,
		nil, // no partial result because it is not a minimal valid commit message
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "x", 2),
	},
	// INVALID / missing colon after type feat
	{
		"invalid-after-valid-type-feat",
		[]byte("test"),
		false,
		nil,
		nil, // no partial result because it is not a minimal valid commit message
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "t", 3),
	},
	// INVALID / invalid type (2 char) + colon
	{
		"invalid-type-2-char-colon",
		[]byte("ch:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, ":", 2),
	},
	// INVALID / invalid type (3 char) + colon
	{
		"invalid-type-3-char-colon",
		[]byte("upd:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, ":", 3),
	},
	// VALID / minimal commit message
	{
		"valid-minimal-commit-message",
		[]byte("fix: w"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "w",
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "w",
		},
		"",
	},
	// VALID / minimal commit message
	{
		"valid-minimal-commit-message-rule",
		[]byte("rule: super secure rule"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "rule",
			Description: "super secure rule",
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "rule",
			Description: "super secure rule",
		},
		"",
	},
	// INVALID / missing colon after valid commit message type
	{
		"missing-colon-after-type-3-chars",
		[]byte("new>"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, ">", 3),
	},
	// INVALID / missing colon after valid commit message type
	{
		"missing-colon-after-type-4-chars",
		[]byte("perf?"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "?", 4),
	},
	// INVALID / missing colon after valid commit message type
	{
		"missing-colon-after-type-5-chars",
		[]byte("build?"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "?", 5),
	},
	// VALID / type + scope + description
	{
		"valid-with-scope",
		[]byte("new(xyz): ccc"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "new",
			Scope:       cctesting.StringAddress("xyz"),
			Description: "ccc",
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "new",
			Scope:       cctesting.StringAddress("xyz"),
			Description: "ccc",
		},
		"",
	},
	// VALID / type + scope + multiple whitespaces + description
	{
		"valid-with-scope-multiple-whitespaces",
		[]byte("fix(aaa):          bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
		},
		"",
	},
	// VALID / type + scope + breaking + description
	{
		"valid-breaking-with-scope",
		[]byte("fix(aaa)!: bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			Exclamation: true,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			Exclamation: true,
		},
		"",
	},
	// VALID / type + scope + breaking + description
	{
		"valid-breaking-with-scope-feat",
		[]byte("feat(aaa)!: bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			Exclamation: true,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			Exclamation: true,
		},
		"",
	},
	// VALID / empty scope is ignored
	{
		"valid-empty-scope-is-ignored",
		[]byte("fix(): bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
		},
		"",
	},
	// VALID / type + empty scope + breaking + description
	{
		"valid-breaking-with-empty-scope",
		[]byte("fix()!: bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			Exclamation: true,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			Exclamation: true,
		},
		"",
	},
	// VALID / type + breaking + description
	{
		"valid-breaking-without-scope",
		[]byte("fix!: bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			Exclamation: true,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			Exclamation: true,
		},
		"",
	},
	// INVALID / missing whitespace after colon (with breaking)
	{
		"invalid-missing-ws-after-colon-with-breaking",
		[]byte("fix!:a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 5),
	},
	// INVALID / missing whitespace after colon with scope
	{
		"invalid-missing-ws-after-colon-with-scope",
		[]byte("fix(x):a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 7),
	},
	// INVALID / missing whitespace after colon with empty scope
	{
		"invalid-missing-ws-after-colon-with-empty-scope",
		[]byte("fix():a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 6),
	},
	// INVALID / missing whitespace after colon
	{
		"invalid-missing-ws-after-colon",
		[]byte("fix:a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 4),
	},
	// INVALID / invalid after valid type and scope
	{
		"invalid-after-valid-type-and-scope",
		[]byte("new(scope)"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, ")", 9),
	},
	// INVALID / invalid initial character
	{
		"invalid-initial-character",
		[]byte("(type: a description"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "(", 0),
	},
	// INVALID / invalid second character
	{
		"invalid-second-character",
		[]byte("c description"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, " ", 1),
	},
	// INVALID / invalid after valid type, scope, and breaking
	{
		"invalid-after-valid-type-scope-and-breaking",
		[]byte("new(scope)!"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "!", 10),
	},
	// INVALID / invalid after valid type, scope, and colon
	{
		"invalid-after-valid-type-scope-and-colon",
		[]byte("fix(scope):"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, ":", 10),
	},
	// INVALID / invalid after valid type, scope, breaking, and colon
	{
		"invalid-after-valid-type-scope-breaking-and-colon",
		[]byte("new(scope)!:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, ":", 11),
	},
	// INVALID / invalid after valid type, scope, breaking, colon, and white-space
	{
		"invalid-after-valid-type-scope-breaking-colon-and-space",
		[]byte("revert(scope)!: "),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescription+ColumnPositionTemplate, " ", 16),
	},
	// INVALID / invalid after valid type, scope, breaking, colon, and white-spaces
	{
		"invalid-after-valid-type-scope-breaking-colon-and-spaces",
		[]byte("ci(scope)!:  "),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescription+ColumnPositionTemplate, " ", 13),
	},
	// INVALID / double left parentheses in scope
	{
		"invalid-double-left-parentheses-scope",
		[]byte("chore(("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrMalformedScope+ColumnPositionTemplate, "(", 6),
	},
	// INVALID / double left parentheses in scope after valid character
	{
		"invalid-double-left-parentheses-scope-after-valid-character",
		[]byte("perf(a("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrMalformedScope+ColumnPositionTemplate, "(", 6),
	},
	// INVALID / double right parentheses in place of an exclamation, or a colon
	{
		"invalid-double-right-parentheses-scope",
		[]byte("fix(a))"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, ")", 6),
	},
	// INVALID / new left parentheses after valid scope
	{
		"invalid-new-left-parentheses-after-valid-scope",
		[]byte("new(az)("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "(", 7),
	},
	// INVALID / newline rather than whitespace in description
	{
		"invalid-newline-rather-than-whitespace-description",
		[]byte("perf(ax):\x0A description on newline"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "\n", 9),
	},
	// INVALID / newline after whitespace in description
	{
		"invalid-newline-after-whitespace-description",
		[]byte("feat(az): \x0Adescription on newline"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrNewline+ColumnPositionTemplate, 11),
	},
	// INVALID / newline in the description
	// VALID / until the newline
	{
		"invalid-newline-in-description",
		[]byte("feat(ae): new\x0Aline"),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Scope:       cctesting.StringAddress("ae"),
			Description: "new",
		},
		fmt.Sprintf(ErrMissingBlankLineAtBeginning+ColumnPositionTemplate, 14),
	},
	// INVALID / newline in the description
	// VALID / until the newline
	{
		"invalid-newline-in-description-2",
		[]byte("docs(az)!: bla\x0Al"),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "docs",
			Scope:       cctesting.StringAddress("az"),
			Exclamation: true,
			Description: "bla",
		},
		fmt.Sprintf(ErrMissingBlankLineAtBeginning+ColumnPositionTemplate, 15),
	},
	// INVALID / newline in the description
	// VALID / until the newline
	{
		"description-ending-with-single-newline",
		[]byte("docs(az)!: bla\x0A"),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "docs",
			Scope:       cctesting.StringAddress("az"),
			Exclamation: true,
			Description: "bla",
		},
		fmt.Sprintf(ErrMissingBlankLineAtBeginning+ColumnPositionTemplate, 15),
	},
	// VALID
	{
		"valid-with-multiline-body",
		[]byte(`fix: correct minor typos in code

see the issue for details

on typos fixed.`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct minor typos in code",
			Body:        cctesting.StringAddress("see the issue for details\n\non typos fixed."),
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct minor typos in code",
			Body:        cctesting.StringAddress("see the issue for details\n\non typos fixed."),
		},
		"",
	},
	// VALID
	{
		"valid-with-singleline-body",
		[]byte(`fix: correct minor typos in code

see the issue for details.`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct minor typos in code",
			Body:        cctesting.StringAddress("see the issue for details."),
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct minor typos in code",
			Body:        cctesting.StringAddress("see the issue for details."),
		},
		"",
	},
	// VALID
	{
		"valid-with-empty-body",
		[]byte(`fix: correct something

`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct something",
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct something",
		},
		"",
	},
	// VALID
	{
		"valid-with-multiple-blank-lines-body",
		[]byte(`fix: correct something



`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct something",
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct something",
		},
		"",
	},
}

var testCasesForConventionalTypes = []testCase{
	// INVALID / empty
	{
		"empty",
		[]byte(""),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEmpty+ColumnPositionTemplate, 0),
	},
	// INVALID / invalid type (1 char)
	{
		"invalid-type-1-char",
		[]byte("c"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "c", 1),
	},
	// INVALID / invalid type (2 char)
	{
		"invalid-type-2-char",
		[]byte("bx"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "x", 1),
	},
	// INVALID / invalid type (2 char) with almost valid type
	{
		"invalid-type-2-char-feat",
		[]byte("fe"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "e", 2),
	},
	// INVALID / invalid type (2 char) with almost valid type
	{
		"invalid-type-2-char-revert",
		[]byte("re"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "e", 2),
	},
	// INVALID / invalid type (3 char)
	{
		"invalid-type-3-char",
		[]byte("fit"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "t", 2),
	},
	// INVALID / invalid type (3 char) again
	{
		"invalid-type-3-char-feat",
		[]byte("fei"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "i", 2),
	},
	// INVALID / invalid type (3 char) with almost valid type
	{
		"invalid-type-3-char-feat",
		[]byte("bui"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "i", 3),
	},
	// INVALID / invalid type (4 char)
	{
		"invalid-type-4-char",
		[]byte("tesx"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "x", 3),
	},
	// INVALID / invalid type (4 char) with almost valid type
	{
		"invalid-type-4-char-refactor",
		[]byte("refa"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "a", 4),
	},
	// INVALID / invalid type (4 char)
	{
		"invalid-type-4-char",
		[]byte("perz"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "z", 3),
	},
	// INVALID / invalid type (5 char) with almost valid type
	{
		"invalid-type-5-char-refactor",
		[]byte("refac"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "c", 5),
	},
	// INVALID / invalid type (6 char) with almost valid type
	{
		"invalid-type-6-char-refactor",
		[]byte("refact"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "t", 6),
	},
	// INVALID / invalid type (7 char) with almost valid type
	{
		"invalid-type-7-char-refactor",
		[]byte("refacto"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "o", 7),
	},
	// INVALID / missing colon after type fix
	{
		"invalid-after-valid-type-fix",
		[]byte("fix"),
		false,
		nil,
		nil, // no partial result because it is not a minimal valid commit message
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "x", 2),
	},
	// INVALID / missing colon after type test
	{
		"invalid-after-valid-type-test",
		[]byte("test"),
		false,
		nil,
		nil, // no partial result because it is not a minimal valid commit message
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "t", 3),
	},
	// INVALID / invalid type (2 char) + colon
	{
		"invalid-type-2-char-colon",
		[]byte("ch:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, ":", 2),
	},
	// INVALID / invalid type (3 char) + colon
	{
		"invalid-type-3-char-colon",
		[]byte("sty:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, ":", 3),
	},
	// VALID / minimal commit message
	{
		"valid-minimal-commit-message",
		[]byte("fix: w"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "w",
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "w",
		},
		"",
	},
	// VALID / minimal commit message
	{
		"valid-minimal-commit-message-style",
		[]byte("style: CSS skillz"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "style",
			Description: "CSS skillz",
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "style",
			Description: "CSS skillz",
		},
		"",
	},
	// INVALID / missing colon after valid commit message type
	{
		"missing-colon-after-type-3-chars",
		[]byte("fix>"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, ">", 3),
	},
	// INVALID / missing colon after valid commit message type
	{
		"missing-colon-after-type-4-chars",
		[]byte("perf?"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "?", 4),
	},
	// INVALID / missing colon after valid commit message type
	{
		"missing-colon-after-type-5-chars",
		[]byte("build?"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "?", 5),
	},
	// VALID / type + scope + description
	{
		"valid-with-scope",
		[]byte("refactor(xyz): ccc"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "refactor",
			Scope:       cctesting.StringAddress("xyz"),
			Description: "ccc",
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "refactor",
			Scope:       cctesting.StringAddress("xyz"),
			Description: "ccc",
		},
		"",
	},
	// VALID / type + scope + multiple whitespaces + description
	{
		"valid-with-scope-multiple-whitespaces",
		[]byte("fix(aaa):          bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
		},
		"",
	},
	// VALID / type + scope + breaking + description
	{
		"valid-breaking-with-scope",
		[]byte("fix(aaa)!: bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			Exclamation: true,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			Exclamation: true,
		},
		"",
	},
	// VALID / type + scope + breaking + description
	{
		"valid-breaking-with-scope-feat",
		[]byte("feat(aaa)!: bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			Exclamation: true,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			Exclamation: true,
		},
		"",
	},
	// VALID / empty scope is ignored
	{
		"valid-empty-scope-is-ignored",
		[]byte("fix(): bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
		},
		"",
	},
	// VALID / type + empty scope + breaking + description
	{
		"valid-breaking-with-empty-scope",
		[]byte("fix()!: bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			Exclamation: true,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			Exclamation: true,
		},
		"",
	},
	// VALID / type + breaking + description
	{
		"valid-breaking-without-scope",
		[]byte("fix!: bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			Exclamation: true,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			Exclamation: true,
		},
		"",
	},
	// INVALID / missing whitespace after colon (with breaking)
	{
		"invalid-missing-ws-after-colon-with-breaking",
		[]byte("fix!:a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 5),
	},
	// INVALID / missing whitespace after colon with scope
	{
		"invalid-missing-ws-after-colon-with-scope",
		[]byte("fix(x):a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 7),
	},
	// INVALID / missing whitespace after colon with empty scope
	{
		"invalid-missing-ws-after-colon-with-empty-scope",
		[]byte("fix():a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 6),
	},
	// INVALID / missing whitespace after colon
	{
		"invalid-missing-ws-after-colon",
		[]byte("fix:a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 4),
	},
	// INVALID / invalid after valid type and scope
	{
		"invalid-after-valid-type-and-scope",
		[]byte("test(scope)"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, ")", 10),
	},
	// INVALID / invalid initial character
	{
		"invalid-initial-character",
		[]byte("(type: a description"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "(", 0),
	},
	// INVALID / invalid second character
	{
		"invalid-second-character",
		[]byte("c description"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, " ", 1),
	},
	// INVALID / invalid after valid type, scope, and breaking
	{
		"invalid-after-valid-type-scope-and-breaking",
		[]byte("test(scope)!"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "!", 11),
	},
	// INVALID / invalid after valid type, scope, and colon
	{
		"invalid-after-valid-type-scope-and-colon",
		[]byte("fix(scope):"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, ":", 10),
	},
	// INVALID / invalid after valid type, scope, breaking, and colon
	{
		"invalid-after-valid-type-scope-breaking-and-colon",
		[]byte("ci(scope)!:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, ":", 10),
	},
	// INVALID / invalid after valid type, scope, breaking, colon, and white-space
	{
		"invalid-after-valid-type-scope-breaking-colon-and-space",
		[]byte("revert(scope)!: "),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescription+ColumnPositionTemplate, " ", 16),
	},
	// INVALID / invalid after valid type, scope, breaking, colon, and white-spaces
	{
		"invalid-after-valid-type-scope-breaking-colon-and-spaces",
		[]byte("ci(scope)!:  "),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescription+ColumnPositionTemplate, " ", 13),
	},
	// INVALID / double left parentheses in scope
	{
		"invalid-double-left-parentheses-scope",
		[]byte("chore(("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrMalformedScope+ColumnPositionTemplate, "(", 6),
	},
	// INVALID / double left parentheses in scope after valid character
	{
		"invalid-double-left-parentheses-scope-after-valid-character",
		[]byte("perf(a("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrMalformedScope+ColumnPositionTemplate, "(", 6),
	},
	// INVALID / double right parentheses in place of an exclamation, or a colon
	{
		"invalid-double-right-parentheses-scope",
		[]byte("fix(a))"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, ")", 6),
	},
	// INVALID / new left parentheses after valid scope
	{
		"invalid-new-left-parentheses-after-valid-scope",
		[]byte("build(az)("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "(", 9),
	},
	// INVALID / newline rather than whitespace in description
	{
		"invalid-newline-rather-than-whitespace-description",
		[]byte("perf(ax):\x0A description on newline"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "\n", 9),
	},
	// INVALID / newline after whitespace in description
	{
		"invalid-newline-after-whitespace-description",
		[]byte("feat(az): \x0Adescription on newline"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrNewline+ColumnPositionTemplate, 11),
	},
	// INVALID / newline in the description
	// VALID / newline in description ignored in best effort mode
	{
		"invalid-newline-in-description",
		[]byte("feat(ap): new\x0Aline"),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Scope:       cctesting.StringAddress("ap"),
			Description: "new",
		},
		fmt.Sprintf(ErrMissingBlankLineAtBeginning+ColumnPositionTemplate, 14),
	},
	// INVALID / newline in the description
	// VALID / newline in description ignored in best effort mode
	{
		"invalid-newline-in-description-2",
		[]byte("perf(at)!: rrr\x0Al"),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "perf",
			Scope:       cctesting.StringAddress("at"),
			Exclamation: true,
			Description: "rrr",
		},
		fmt.Sprintf(ErrMissingBlankLineAtBeginning+ColumnPositionTemplate, 15),
	},
	// INVALID / newline in the description
	// VALID / until the newline
	{
		"description-ending-with-single-newline",
		[]byte("perf(at)!: rrr\x0A"),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "perf",
			Scope:       cctesting.StringAddress("at"),
			Exclamation: true,
			Description: "rrr",
		},
		fmt.Sprintf(ErrMissingBlankLineAtBeginning+ColumnPositionTemplate, 15),
	},
	// VALID
	{
		"valid-with-multiline-body",
		[]byte(`fix: correct minor typos in code

see the issue for details

on typos fixed.`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct minor typos in code",
			Body:        cctesting.StringAddress("see the issue for details\n\non typos fixed."),
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct minor typos in code",
			Body:        cctesting.StringAddress("see the issue for details\n\non typos fixed."),
		},
		"",
	},
	// VALID
	{
		"valid-with-singleline-body",
		[]byte(`fix: correct minor typos in code

see the issue for details.`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct minor typos in code",
			Body:        cctesting.StringAddress("see the issue for details."),
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct minor typos in code",
			Body:        cctesting.StringAddress("see the issue for details."),
		},
		"",
	},
	// VALID
	{
		"valid-with-empty-body",
		[]byte(`fix: correct something

`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct something",
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct something",
		},
		"",
	},
	// VALID
	{
		"valid-with-multiple-blank-lines-body",
		[]byte(`fix: correct something



`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct something",
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct something",
		},
		"",
	},
}