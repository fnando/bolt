# gotestfmt

A better format for golang's tests when using
<https://github.com/stretchr/testify>.

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
Usage of gotestfmt:
  -fast-fail
      Fast fail
  -reporter string
      Choose report type (dot, json) (default "dot")
  -version
      Show version
```

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
