# KISIPAR NOTES


**HOW THE FUCK DO WE HANDLE LISTS?**

Maybe there is a List item type?

Then it's up to the provider.  But isn't it better to have that logic
in the templates?


    
```go

// tbd

```

## GOALS

### Abstract Data Source

Three types implemented at start:

1. Virtual (also used for testing)
2. FileSystem
3. Whatever I use for real data.

### Multiple Sites from One Executable


### Caching of Known Paths

### Still Support "folder full of Markdown files" story

## Layout & Code &c.

Use "lots of files in package" as that seems to be the norm.

Try to reduce the code footprint!

Use sync.WaitGroup for the multi-site logic

Use standard things with overrides where possible, so e.g. DefaultMux can
be used, we just set it to our own if needed (or better yet, don't, and
let someone else if they need to).

Logger, ditto.  What else?

## Generations

Bindoc I think in any case, right?

If we end up with enums:

https://godoc.org/golang.org/x/tools/cmd/stringer

## DS Interaction


    if !ds.Has(rpath) {
        // Send 404 etc.
    }
    // We now have either a Page or a File or (maybe) raw data.
    if page, err := ds.Page(rpath); err == nil {
        // Serve Page
    }
    if page := ds.Page(have-page) {
        ..serve-page..
    } else {
        ... get it as a file or what?
    }

## Executable

* Option to read from a YAML file, presumably -y or --yaml FILE.
