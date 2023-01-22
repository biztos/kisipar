# Kisipar funcmap!

## TODO

* Godoc and testing links once public.

## What?

A collection of functions useful to have in Go templates.

## Why?

Go has an interesting and quite powerful templating system.  However, it has
very meagre function support out of the box.

This is by design: Go makes it is easy to define your own functions, and
thus not waste any time or memory on functions you don't plan to use.

So far so good, until you want a general-purpose templating system for, say,
a web site, and you want to give template authors a broad palette of
functions.

Hence this collection, which will make use of other collections as well.

## How?

```go
TODO
```

## Who?

Kevin Frost, for functions defined in this package directly.

Kyoung-chan Lee, for functions defined in **gtf**:
[https://github.com/leekchan/gtf](https://github.com/leekchan/gtf)

## Whither?

This *may* grow into some kind of monster collection, but it may also not.

In fact I rather hope it doesn't: other people have more time to spend on
this problem, and better reasons to spend it; and I'm more than happy to send
them pull requests once I'm convinced my functions have sufficient merit.

In the meantime, `funcmap` is a subpackage of `kisipar`, and probably not
very interesting outside that context.
