# gotestfmt

[![Latest release](https://img.shields.io/github/v/release/fnando/gotestfmt?label=version)](https://github.com/fnando/gotestfmt/releases/latest)

A better format output for golang's tests.

We've all been there... a bunch of tests failing and you have no idea where to
start, because everything looks the same. And that one test failing in the sea
of passing tests? So frustrating. Not anymore!

With `gotestfmt` you'll see a progress output while tests are being executed.
Once tests are done, you'll see an output with only the tests that failed.
Simple and easy!

Here's the before and after:

![An image showing the comparison between the native output versus gotestfmt's](https://github.com/fnando/gotestfmt/raw/main/gotestfmt.png)

## Install

Download the binary for your system from the
[latest release](https://github.com/fnando/gotestfmt/releases/latest) and place
it anywhere on your path.

## Usage

Pipe test results in JSON format into this tool:

```shell
$ go test -json . | gotestfmt
```

Options:

```shell
$ gotestfmt -h
gotestfmt is a tool that generates a better output format for golang tests.

Usage: gotestfmt [OPTIONS]

  -cover
      Show module coverage (default true)
  -cover-count int
      Number of coverage items to display (default 10)
  -cover-threshold float
      Only show module coverage below this threshold (default 100)
  -fastfail
      Fast fail
  -reporter string
      Choose report type (dot, json) (default "dot")

Other commands:

  gotestfmt download-url
      display the latest binary download url

  gotestfmt update
      download the latest binary and replace the running one

  gotestfmt version
      display the version


For more info, visit https://github.com/fnando/gotestfmt
```

To get the latest download url for your binary, you can use
`gotestfmt download-url`.

### Overriding colors

You can override the colors by setting the following env vars:

```bash
export GOTESTFMT_TEXT_COLOR="30"
export GOTESTFMT_FAIL_COLOR="31"
export GOTESTFMT_PASS_COLOR="32"
export GOTESTFMT_SKIP_COLOR="33"
export GOTESTFMT_DETAIL_COLOR="34"
export GOTESTFMT_COVERAGE_BAD_COLOR="31"  # coverage < 60%
export GOTESTFMT_COVERAGE_GOOD_COLOR="32" # coverage > 70%
export GOTESTFMT_COVERAGE_OK_COLOR="33"   # coverage between 100-70%
```

To disable color output completely, just set `NO_COLOR=1`.

## Code of Conduct

Everyone interacting in the gotestfmt project’s codebases, issue trackers, chat
rooms and mailing lists is expected to follow the
[code of conduct](https://github.com/fnando/gotestfmt/blob/main/CODE_OF_CONDUCT.md).

## Contributing

Bug reports and pull requests are welcome on GitHub at
https://github.com/fnando/gotestfmt. This project is intended to be a safe,
welcoming space for collaboration, and contributors are expected to adhere to
the [Contributor Covenant](http://contributor-covenant.org) code of conduct.

## License

The gem is available as open source under the terms of the
[MIT License](https://opensource.org/licenses/MIT).
