# KISIPAR NOTES

```go

itemType := ds.Has(r.URL.Path)
if (itemType == kisipar.NoItem) {
    http.NotFound(w,r)
    return
}
switch itemType {
    case kisipar.PageItem:
        p, err := ds.Page(r.URL.Path)
        // serve page or error
    default:
        panic("unknown item type: " + itemType)
}

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

## Layout

Use "lots of files in package" as that seems to be the norm.

Try to reduce the code footprint!



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