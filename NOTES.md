# KISIPAR NOTES

## LATEST PROBLEM: COMPLEXITY!

The code base is too complex, I'm losing track of what does what where.

How can we simplify?

    KEEP
    
    * GENERIC "PROVIDER" CONCEPT
    * LOADING FMD PAGES
    
        what are the MVP pathers?
        (what else to call them?)
        
        * File - serves file
        * Page - page w/HTML etc.
        * Content - raw content
        * Handler - handles it by itself
        
        type check on these, decide what to do
        
        only Page has a template
        
        site knows about templates (generally not from same provider)
        
        
        
        type Item interface {
            Path() string // normalized path, may not be lookup path
        }
        
        
    LOSE
    
    * COMPLEXITY IN INTERFACES
    * YAML PAGES (DON'T NEED THEM FOR ANYTHING)
    
    
    
---
ANNO...

    NO FUCKING ITERATORS FOR NOW
    
    LATEST PLAN:
    
    1. MVP FOR BLOGS & PLACEHOLDERS USING FILE SYSTEM
    2. WAIT & SEE WHAT ELSE
    3. CONSIDER EBP ON A DIFFERENT SYSTEM ALTOGETHER
        (but could still start it here)
    
    NEXT UP - TIME BASED SETS

    What is the MVP?
    
    * Frostopolis from files.
        * With RSS
    * Placeholders from files for all others planned.
    * All with TSL at the edge (maybe not in kisipar yet)
    * Single host DigitalOcean

    Low-hanging fruit after that:
    
    * Placeholder KFCOM
    * Placeholder old Migra from files / archive.org

## GOALS

### Abstract Data Source

Three types implemented at start:

1. Virtual (also used for testing)
2. FileSystem
3. Whatever I use for real data.

### Multiple Sites from One Executable

    Maybe something like GetHost or Get(r *http.Request)?
    
    Problem is that you might have a Provider that knows about
    a bunch of hosts, and can't tell them apart by path alone.
    (In that case it might also not work with the Request.)
    
    You can attach the host to the path but that fucks with the
    template search unless you also have templates for each
    host (you might).  Would need new default logic.
    
    The other thing is to have one Provider per Site but then
    build the db provider such that it knows which host is which.
    
    mp,err := NewPgMultiProvider(config)
    p, err := mp.HostProvider("example.com")
    
    Probably best that way right?

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
