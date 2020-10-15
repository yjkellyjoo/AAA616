# AAA616

AAA616 Projects

### Installation
1. Install Go (https://golang.org)
2. Install llir (https://github.com/llir/llvm)
```sh
$ go get -u github.com/llir/llvm/...
```

### Running an analysis
```sh
$ cd sign
$ go build
$ ./sign ../test/toy/ll/fact.ll
```

