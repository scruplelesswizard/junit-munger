package main

import (
	"strings"
	"testing"
)

func TestMergeGroupsByFileAndSortsSlowestFirst(t *testing.T) {
	suites := getSuitesFromJUnitXML("testdata/rspec-*.xml")
	out := buildReportOutput(buildSuitesFromFiles(suites))

	type want struct {
		name  string
		tests int
		time  string
	}
	expected := []want{
		{"./spec/api/orders_spec.rb", 2, "3.750000"},
		{"./spec/models/user_spec.rb", 2, "0.750000"},
		{"./spec/jobs/mailer_spec.rb", 1, "0.200000"},
	}

	if len(out.Suites) != len(expected) {
		t.Fatalf("expected %d suites, got %d", len(expected), len(out.Suites))
	}
	for i, w := range expected {
		got := out.Suites[i]
		if got.Name != w.name || got.Tests != w.tests || got.Time != w.time {
			t.Errorf("suite %d: got {%s %d %s}, want {%s %d %s}",
				i, got.Name, got.Tests, got.Time, w.name, w.tests, w.time)
		}
	}
}

func TestTestcasesSortedSlowestFirst(t *testing.T) {
	suites := getSuitesFromJUnitXML("testdata/rspec-*.xml")
	out := buildReportOutput(buildSuitesFromFiles(suites))

	orders := out.Suites[0]
	if len(orders.Testcases) != 2 {
		t.Fatalf("expected 2 testcases, got %d", len(orders.Testcases))
	}
	if orders.Testcases[0].Name != "creates an order" || orders.Testcases[1].Name != "lists orders" {
		t.Errorf("testcases not sorted slowest first: %q, %q",
			orders.Testcases[0].Name, orders.Testcases[1].Name)
	}
}

func TestFailuresAndSkipsPreserved(t *testing.T) {
	suites := getSuitesFromJUnitXML("testdata/rspec-*.xml")
	out := buildReportOutput(buildSuitesFromFiles(suites))

	var failed, skipped int
	for _, s := range out.Suites {
		for _, tc := range s.Testcases {
			if tc.Failure != nil {
				failed++
				if tc.Failure.Message != "expected invalid" {
					t.Errorf("unexpected failure message %q", tc.Failure.Message)
				}
			}
			if tc.Skipped != nil {
				skipped++
			}
		}
	}
	if failed != 1 || skipped != 1 {
		t.Errorf("got %d failures and %d skips, want 1 and 1", failed, skipped)
	}
	if out.Failures != 1 || out.Skipped != 1 {
		t.Errorf("totals: failures=%d skipped=%d, want 1 and 1", out.Failures, out.Skipped)
	}
}

func TestLoadJUnitXML(t *testing.T) {
	const xml = `<testsuite name="s"><testcase classname="c" name="n" file="a.rb" time="1.5"/></testsuite>`
	suites := loadJUnitXML(strings.NewReader(xml))
	if len(suites) != 1 || len(suites[0].Tests) != 1 {
		t.Fatalf("unexpected suites: %+v", suites)
	}
	if got := suites[0].Tests[0].Properties["file"]; got != "a.rb" {
		t.Errorf("file property = %q, want a.rb", got)
	}
}

func TestNoMatchingFilesYieldsNoSuites(t *testing.T) {
	if got := getSuitesFromJUnitXML("testdata/does-not-exist-*.xml"); len(got) != 0 {
		t.Errorf("expected no suites, got %d", len(got))
	}
}
