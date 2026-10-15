# junit-munger

Merge RSpec JUnit reports and regroup them by spec file, sorted slowest first.

RSpec's JUnit output is one `<testsuite>` per run (or per parallel worker). To split a test suite into
equal-time groups you need timings per **file**, not per run. `junit-munger` reads any number of JUnit
XML files, groups every test case by its `file` attribute, and writes a single JUnit document where each
`<testsuite>` is one spec file with its total duration.

It pairs well with [split-tests](https://github.com/scruplelesswizard/split-tests), which can use the
output as its `junit-path` to balance parallel CI runners by real durations.

## Install

```bash
go install github.com/scruplelesswizard/junit-munger@latest
```

or download a binary for Linux or macOS from the
[releases page](https://github.com/scruplelesswizard/junit-munger/releases), or build from source:

```bash
go build -o junit-munger .
```

## Usage

```bash
junit-munger --in 'results/rspec-*.xml' > rspec.xml
```

| Flag | Description |
|------|-------------|
| `--in` | Glob of JUnit XML files to merge (required). `**` is supported. Quote it so your shell does not expand it. |

The merged report is written to stdout. Progress (which files were loaded) is logged to stderr, so
redirecting stdout is safe.

### Example

Given two worker reports where `orders_spec.rb` ran across both:

```
results/rspec-0.xml   user_spec.rb (2 tests), orders_spec.rb (1 test)
results/rspec-1.xml   orders_spec.rb (1 test), mailer_spec.rb (1 test)
```

the output has one suite per file, slowest first, with durations summed:

```xml
<testsuites tests="5" failures="1" skipped="1">
  <testsuite name="./spec/api/orders_spec.rb" tests="2" time="3.750000">...</testsuite>
  <testsuite name="./spec/models/user_spec.rb" tests="2" failures="1" time="0.750000">...</testsuite>
  <testsuite name="./spec/jobs/mailer_spec.rb" tests="1" skipped="1" time="0.200000">...</testsuite>
</testsuites>
```

Failures, errors and skips are preserved on each test case. Test cases inside each suite are also sorted
slowest first.

## Requirements

Each `<testcase>` needs a `file` attribute. The
[rspec_junit_formatter](https://github.com/sj26/rspec_junit_formatter) gem emits it by default. Test
cases without one are grouped together under an empty suite name.

## License

BSD 2-Clause. See [LICENSE](LICENSE).
