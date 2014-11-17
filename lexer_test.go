package main

import (
	"fmt"
	"testing"
)

type lexTest struct {
	name  string // Name of the sub-test.
	input string // Input string to be lexed.
	want  []item // Expected items produced by lexer.
}

const (
	snip1 = "echo $PATH\n" +
		"echo $GOPATH"
	snip2 = "kill -9 $pid"
)

var (
	tEOF = item{itemEOF, ""}
)

var lexTests = []lexTest{
	{"empty", "", []item{tEOF}},
	{"spaces", " \t\n", []item{tEOF}},
	{"text", "blah blah blinkity blah", []item{tEOF}},
	{"comment1", "<!-- -->", []item{tEOF}},
	{"comment2", "a <!-- --> b", []item{tEOF}},
	{"block1", "aa <!-- @1 -->\n" +
		"```\n" + snip1 + "```\n bbb",
		[]item{{itemThreadLabel, "1"},
			{itemSnippet, snip1},
			tEOF}},
	{"block2", "aa <!-- @1 @2-->\n" +
		"```\n" + snip1 + "```\n bb cc\n" +
		"dd <!-- @3 @4-->\n" +
		"```\n" + snip2 + "```\n ee ff\n",
		[]item{
			{itemThreadLabel, "1"},
			{itemThreadLabel, "2"},
			{itemSnippet, snip1},
			{itemThreadLabel, "3"},
			{itemThreadLabel, "4"},
			{itemSnippet, snip2},
			tEOF}},
}

// collect gathers the emitted items into a slice.
func collect(t *lexTest) (items []item) {
	l := newLex(t.input)
	for {
		item := l.nextItem()
		items = append(items, item)
		if item.typ == itemEOF || item.typ == itemError {
			break
		}
	}
	return
}

func equal(i1, i2 []item) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			fmt.Printf("types not equal - got : %s\n", i1[k].typ)
			fmt.Printf("types not equal - want: %s\n", i2[k].typ)
			fmt.Printf("\n")
			return false
		}
		if i1[k].val != i2[k].val {
			fmt.Printf("vals not equal - got : %q\n", i1[k].val)
			fmt.Printf("vals not equal - want: %q\n", i2[k].val)
			fmt.Printf("\n")
			return false
		}
	}
	return true
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		got := collect(&test)
		if !equal(got, test.want) {
			t.Errorf("%s:\ngot\n\t%+v\nexpected\n\t%v", test.name, got, test.want)
		}
	}
}