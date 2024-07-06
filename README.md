# go-nostrbuild

Go package for nostr.build

## Usage

```go
sign := func(ev *nostr.Event) error {
	ev.PubKey = pubKey
	return ev.Sign(sk)
}
result, err := nostrbuild.Upload(&buf, sign)
```

If you use nbcmd

```
$ export NBCMD_NSEC=nsec1xxxxxxxxxxxxxxxxxxxx

$ nbcmd upload my-picture.png 
https://image.nostr.build/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.png

$ nbcmd delete https://image.nostr.build/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.png
File deleted.
```

## Installation

go get github.com/mattn/go-nostrbuild@latest

## License

MIT

## Author

Yasuhiro Matsumoto
