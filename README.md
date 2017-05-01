# Sider example key value in-memory storage with rest api.

There are two separated applications:
1. siderd - long run daemon with rest api.
2. sider - command line client to rest api.

## Build

Install [govendor](https://github.com/kardianos/govendor):
```
$ go get -u github.com/kardianos/govendor
```

Build applications separately

Daemon:
```
$ cd $GOPATH/src/github.com/aandryashin/siderd
```
Sync dependencies:
```
$ govendor sync
```
Build source:
```
$ go build
```
Run daemon:
```
$ ./siderd --help
```

Command line client:
```
$ cd $GOPATH/src/github.com/aandryashin/sider
```
Sync dependencies:
```
$ govendor sync
```
Build source:
```
$ go build
```
Run command line client:
```
$ ./sider --help
```

Keys are strings, values are objects in json notation.

Simple session:
```
$ sider keys
[]

$ sider set 1 '["one"]'
$ sider set 2 '["two"]'
$ sider set 3 '["three"]'

$ sider keys
[
    "3",
    "1",
    "2"
]

$ sider get 1
[
    "one"
]
```

By default keys are not expired, to set keys with expiration timeout provide ttl option in golang time.Duration notation:
```
$ sider set key '{}' --ttl 10s
```
