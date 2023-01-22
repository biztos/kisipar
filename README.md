# Kisipar: An UNFINISHED (and abandoned) Markdown-centric web server for small sites.

This was a fun project and I rewrote it several times -- twice in Perl, thrice
in Go, and even once in Javascript/nodejs. But in the end, I never quite got
it where I wanted it, and other things solved my own problems better. These
days I use Zola.

Still, there are some useful bits of code and ideas in here, so I keep it
around for reference, and you are welcome to use any of it subject to the
MIT LICENSE.

---

    LATEST IDEAS - 10 MAY 2016

    Site Layout

        config.yaml / config.json (as you prefer)
        static/
        pages/
        templates/

    TheDot
        Site
            Name
            Description
            Domain
            Owner
            BaseURL
            Pageset -- all pages for the site
            Config  -- your full config; handle with care!
            StaticInfo(path) -> FileInfo
            StaticList(prefix,suffix,extension) --> []FileInfo
        Page -- the current page if applicable, not present for created index
        Pageset -- the pageset if applicable (current level down)
        Static  -- interface to lookup static files, if desired.

    PAGE FINDING HIERARCHY

        1. STATIC (so you can override on disk)
        2. VIRTUAL
        3. REGULAR FULL MATCH
        4. REGULAR INDEX MATCH


    STATIC ASSET HANDLING

    Exact match in static dir (configurable, or "static") trumps all else,
    and is never cached.  No approximate matches are attempted.

    No auto-indexing of static pages, for instance.  If you want to store
    things in /static and

    INDEX HANDLING

    Specific page trumps index, so for `/foo` a match on `/foo.md` will always
    beat `/foo/index.md` and will *not* result in a Pageset automatically
    present for the Dot.

        IF exact static match THEN serve that file
        ELSE IF exact md THEN serve that as Page
        ELSE IF index THEN serve that as Page with Pageset for dir
        ELSE IF dir matches THEN serve generic Index with Pageset for dir
        ELSE serve not-found

    TEMPLATE SELECTION

    Templates can e overridden in the config, or in the individual file's
    meta.  Config overrides are either specific by cleaned URL, or general
    covering a partial match; in the latter case, the most complete match
    is used, i.e. `/foo/bar/baz` matches `/foo/bar` before `/foo`.
    Config templates take precedence over file-meta templates.
    Regexp template overrides are not (yet) supported, but might be a good
    idea.

    If a specified template is not found, a 500 error results.

        IF match for
        IF exact match for clean URL in config THEN use that
        ELSE IF set in file meta THEN use that
        ELSE IF partial match in config THEN use that
    First: any page-specific override
    First: template specified in the file, if available (500 if not).
    Second: any override template specified in config
    USE CASES
        /foo
        /foo/index.html


"Kisipar" is Hungarian for "Kleingewerbe," which is German for small business.

Kisipar is also an opinionated, Markdown-centric web server written in the Go
language.

## WORK IN PROGRESS

CAVEAT LECTOR! This is a work in _early_ progress and may change radically at
any moment. Do not try to use it before it at least reaches version 1.0.
Thank you.

    CURRENT VERSION: 0.0.0 (i.e. just getting started here folks...)

NOTE: now a private repo until it does something useful (but it might go
public before 1.0.0, just not while it's this embarrassingly raw).

## Motivation

I want an off-the-shelf web server for my content-based projects, with a
minimum of configuration required for the simplest use cases. I want it to be
easily expandable to handle more complicated sites. I want the textual
content itself to be written in Markdown.

Since others might reasonably want the same thing, let's make it an Open
Source project!

## Goals

- Zero-config serving of simple sites.
  - `kisipar /path/to/your/site`
- High performance.
- Minimum number of dependencies outside the standard library.
- Easy/obvious way to extend (in code) for new routes.
- Flexible, within reasonable use cases.
- Support multiple languages.
- Placeholder sites just run the binary w/ cli args for config.
- Basic "real" sites just run the binary but w/various customizations.
- More complicated sites `import "kisipar"` and go from there.

### Stretc (or 2.0) Goals

- Run multiple sites from one executable.
- Sane multilingual options.
- SSL (normally you'd let Nginx do this for you).

## NOTES TO SELF

- Caching rendered content defeats the purpose of greater control in templates.
  - You can't cache the request object, and it may be meaningful.
  - If you want real caching for performance reasons then do it in Nginx etc.
- Follow blackfriday example for packaging: https://github.com/russross/blackfriday
  - Except we want an executable, so I guess it needs a main package.
- https://travis-ci.org
- docopt in main
- vendor in deps or what?
- frosted markdown as option (once ready); otherwise just regular
  - then extract title via node parsing or what?
  - first line is title?
- what to do for config?
- Use the executable, or fork and customize
- What about using one instance of kisipar to serve a bunch of sites?
  - Shared static content for example...

## Features

- Content can be Markdown, text, or HTML.
- Content is rendered in the browser.
- Uses Frosted Markdown for defining metadata.
- Templates are used.
- Zero config option: can use command-line or env.
- Contact Form available (simple mail sender).
  - Or not? Security risk?
- MAYBE SSL?

## Site Layout:

    config.toml
    templates/
        wrapper.html
        ...anything else you want, it's within the templates...
    static/
        js/
        css/

        (anything)
    content/

Order of checking:

1. static (trumps all else; not cached)
2. cached, rendered (respect TTL)
3. content directory: md,txt,html as fragments; others served straight.
4. special URLs, eg "contact" and "rss"

Anything not found this way is 404.
