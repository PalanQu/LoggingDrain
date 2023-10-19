# LoggingDrain3

## Introduce

LoggingDrain3 is go implemented for [drain3](http://jiemingzhu.github.io/pub/pjhe_icws2017.pdf)


## Example

``` bash
go run examples/stdin.go
```

## Test

run unittest

``` bash
go test .
```

run benchmark

``` bash
go test -bench=.
```

```
goos: darwin
goarch: arm64
pkg: github.com/palanqu/loggingdrain3
BenchmarkBuildTree-8       	 1632832	       699.0 ns/op
BenchmarkMatchTree-8       	 3431210	       349.4 ns/op
BenchmarkUnmarshalJson-8   	  378432	      3172 ns/op
PASS
ok  	github.com/palanqu/loggingdrain3	4.840s
```