# Streamlet

Streamlet is originally a NodeJS-based stream-oriented database optimal for small scale projects that need a high speed method for storing data. This version of streamlet has been rewritten in Go for maximum speed and portability.

This version of Streamlet used to use `os.Write`, but after extensive testing, it was determined that it was about 2000x faster to use `bufio`. Before the change to `bufio`, the database clocked in at `~3500 ns/op` for bulk inserts. It now takes `~1600 ns/op`.

On average, Streamlet can achieve `450,000 ops/sec` write speed on a Dell XPS 13. Streamlet can achieve `1000ns/doc` read speed on the same device.

## Installation

```bash

GO111MODULE=on
go get github.com/SuperAuguste/streamlet

```

## License

This project is licensed under the MIT license.
