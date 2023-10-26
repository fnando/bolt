# gotestfmt

[![Latest release](https://img.shields.io/github/v/release/fnando/gotestfmt?label=version)](https://github.com/fnando/gotestfmt/releases/latest)
[![Tests](https://github.com/fnando/gotestfmt/actions/workflows/tests.yml/badge.svg)](https://github.com/fnando/gotestfmt/actions/workflows/tests.yml)

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

gotestfmt wraps `go test`. You can run it with:

```shell
$ gotestfmt run ./...
```

Options:

```shell
$ gotestfmt -h

gotestfmt is a golang test runner that has a nicer output.

  Usage: gotestfmt [command] [options]

  Options:
    --version                          Show version


  Commands:

    gotestfmt version                  Show gotestfmt version
    gotestfmt download-url             Output the latest binary download url
    gotestfmt run                      Run tests
    gotestfmt [command] --help         Display help on [command]


  Further information:
    https://github.com/fnando/gotestfmt


```

To get the latest download url for your binary, you can use
`gotestfmt download-url`.

### Reporters

gotestfmt comes with two different reporters:

### JSON

The JSON reporter outputs a nicer JSON format that can be used to do things that
require structured data.

```shell
$ gotestfmt run ./... --reporter json
```

### Progress

The progress reporter outputs a sequence of characters that represent the test's
status (fail, pass, skip). Once all tests have been executed, a summary with the
failing and skipped tests, plus a coverage list is printed.

```shell
$ gotestfmt run ./... --reporter progress
```

#### Overriding colors

You can override the colors by setting the following env vars:

```bash
export GOTESTFMT_TEXT_COLOR="30"
export GOTESTFMT_FAIL_COLOR="31"
export GOTESTFMT_PASS_COLOR="32"
export GOTESTFMT_SKIP_COLOR="33"
export GOTESTFMT_DETAIL_COLOR="34"
```

To disable color output completely, just set `NO_COLOR=1`.

#### Overriding symbols

To override the characters, you can set some env vars. The following example
shows how to use emojis instead:

```shell
export GOTESTFMT_FAIL_SYMBOL=❌
export GOTESTFMT_PASS_SYMBOL=🔥
export GOTESTFMT_SKIP_SYMBOL=😴
```

## Code of Conduct

Everyone interacting in the gotestfmt project’s codebases, issue trackers, chat
rooms and mailing lists is expected to follow the
[code of conduct](https://github.com/fnando/gotestfmt/blob/main/CODE_OF_CONDUCT.md).

## Developing

To generate new test replay files, you can use
`go test -cover -json ./reference/package > test/replays/[case].txt`.

To generate new benchmark replay files, you can use
`go test -json -fullpath -bench . ./reference/bench &> test/replays/benchmark.txt`.

Once files are exported, make sure you replace all paths to use `/home/test` as
the home directory, and `/home/test/gotestfmt` as the working directory.

You can run tests with `./bin/test`.

## Contributing

Bug reports and pull requests are welcome on GitHub at
https://github.com/fnando/gotestfmt. This project is intended to be a safe,
welcoming space for collaboration, and contributors are expected to adhere to
the [Contributor Covenant](http://contributor-covenant.org) code of conduct.

## License

The gem is available as open source under the terms of the
[MIT License](https://opensource.org/licenses/MIT).
