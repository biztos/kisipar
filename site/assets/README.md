# README: Kisipar site assets

Kisipar uses packaged data for various things such as default templates, demo
content, and the like.

This is much easier to manage when we use the `bindata` package:

https://github.com/jteeuwen/go-bindata

If you wish to change anything major here you should be familiar with that
package.  If you only wish to update the data files, then it should be
sufficient to simply run `./refresh.sh` in this directory.

You may need to first install `go-bindata`:

    https://github.com/jteeuwen/go-bindata

If any of this seems wrong, stupid, or annoying, please file a pull request
and it will be given due consideration.
