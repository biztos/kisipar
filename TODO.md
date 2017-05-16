# KISIPAR TODO

## STOP EXPOSING TEMPLATES

http://localhost:8080/style

Instead, what?

1. We `Get` or `GetSince` from the Provider for all things not already done.
    * Ergo static, specials like contact or news, etc.
2. Provider takes the Request so it can branch on get/post logic if it wants.
    * Ergo we could make a read-write kisipar wiki, it's up to the provider.
    * Nice use-case for a provider: write shit using git?
    * Make use of context maybe?
3. One thing the provider can return (new interface) is an Error Pather.
    * It has a Code, Message, PublicDetail, PrivateDetail
    * Those explicit names are to keep us honest when writing templates.
4. We still look up templates separately so we can if needed do that.

So the `FSP` does:

* Reject any non-GET out of hand.
* Exclude templates matching a regex.
* Never sees templates in `/shared`, never serves "default.*"
* Climbs up the tree, that's really very useful.
* Maybe keeps track of template names so type-defaulting is better.
    * Want "foo.any" to be a valid template.
    * Want to have a configurable list of ordered types.
    * I guess you can require `.any` to be included.

    
Can't do `s.Provider.TemplateForPath(path)` because that's the same logic we
use for pages.

    s.Provider.RequestTemplate(r) - template for *request*
    
    Thus the provider can do any fancy shit they want with it.
    
    ...or maybe the provider tells you what to do...?
    
    
    
    Maybe walk up after all? (This is up to the provider.)
    
        GET /foo/bar/baz
        
        TEMPLATE
            /foo/bar/baz.html
            /foo/bar.html
            /foo.html
            /default.html
            
    s.Provider.PathTemplate(path) - returns template to handle req.
    s.Provider.

## INDEX / LIST HANDLERS

* Use templates, make some examples -- need to e.g. get latest 20.

## EXPOSE PROVIDER+SITE IN PAGE SOMEHOW

Back to the Dot?

## TRULY MINIMALIST FILE SERVER (WITH FMD + YAML)

## Code Audits / Refactors

* Audit NewThing pattern, make sure all return (*Thing,err) if taking args.

## Allow non-HTML templates

Use the template extension to set the content type.

## Cache StaticDir content on load; configurable TTL.

## Write templates to vars then to output.

So we can catch errors without partial fucking results eh.
