# gotestfmt

A better format for golang's tests when using
<https://github.com/stretchr/testify>.

## Install

Download binary for your system from the
[Releases page](https://github.com/fnando/gotestfmt/releases) and place it
anywhere on your path.

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
