# bolt

[![Latest release](https://img.shields.io/github/v/release/fnando/bolt?label=version)](https://github.com/fnando/bolt/releases/latest)
[![Tests](https://github.com/fnando/bolt/actions/workflows/tests.yml/badge.svg)](https://github.com/fnando/bolt/actions/workflows/tests.yml)

A better format output for golang's tests.

We've all been there... a bunch of tests failing and you have no idea where to
start, because everything looks the same. And that one test failing in the sea
of passing tests? So frustrating. Not anymore!

With `bolt` you'll see a progress output while tests are being executed. Once
tests are done, you'll see an output with only the tests that failed. Simple and
easy!

Here's the before and after:

![An image showing the comparison between the native output versus bolt's](https://github.com/fnando/bolt/raw/main/bolt.png)

Features:

- Colored output
- Dotenv files support
- Coverage output
- Slowest tests output
- Benchmark output

## Install

Download the binary for your system from the
[latest release](https://github.com/fnando/bolt/releases/latest) and place it
anywhere on your path.

## Usage

bolt wraps `go test`. You can run it with:

```shell
$ bolt run ./...
```

Options:

```shell
$ bolt -h

bolt is a golang test runner that has a nicer output.

  Usage: bolt [command] [options]

  Commands:

    bolt version                  Show bolt version
    bolt run                      Run tests
    bolt update                   Update to the latest released version
    bolt [command] --help         Display help on [command]


  Further information:
    https://github.com/fnando/bolt
```

To get the latest download url for your binary, you can use `bolt download-url`.

### Reporters

bolt comes with two different reporters:

### JSON

The JSON reporter outputs a nicer JSON format that can be used to do things that
require structured data.

```shell
$ bolt run ./... --reporter json
```

### Progress

The progress reporter outputs a sequence of characters that represent the test's
status (fail, pass, skip). Once all tests have been executed, a summary with the
failing and skipped tests, plus a coverage list is printed.

```shell
$ bolt run ./... --reporter progress
```

#### Overriding colors

You can override the colors by setting the following env vars:

```bash
export BOLT_TEXT_COLOR="30"
export BOLT_FAIL_COLOR="31"
export BOLT_PASS_COLOR="32"
export BOLT_SKIP_COLOR="33"
export BOLT_DETAIL_COLOR="34"
```

To disable color output completely, just set `NO_COLOR=1`.

#### Overriding symbols

To override the characters, you can set some env vars. The following example
shows how to use emojis instead:

```shell
export BOLT_FAIL_SYMBOL=❌
export BOLT_PASS_SYMBOL=⚡️
export BOLT_SKIP_SYMBOL=😴
```

## Code of Conduct

Everyone interacting in the bolt project’s codebases, issue trackers, chat rooms
and mailing lists is expected to follow the
[code of conduct](https://github.com/fnando/bolt/blob/main/CODE_OF_CONDUCT.md).

## Developing

To generate new test replay files, you can use
`go test -cover -json -tags=reference ./test/reference/package > test/replays/[case].txt`.

To generate new benchmark replay files, you can use
`go test -json -fullpath -tags=reference -bench . ./test/reference/bench &> test/replays/benchmark.txt`.

Once files are exported, make sure you replace all paths to use `/home/test` as
the home directory, and `/home/test/bolt` as the working directory.

You can run tests with `./bin/test`.

## Contributing

Bug reports and pull requests are welcome on GitHub at
https://github.com/fnando/bolt. This project is intended to be a safe, welcoming
space for collaboration, and contributors are expected to adhere to the
[Contributor Covenant](http://contributor-covenant.org) code of conduct.

## License

The gem is available as open source under the terms of the
[MIT License](https://opensource.org/licenses/MIT).
