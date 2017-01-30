#!/bin/bash
#
# refresh.sh - refresh kisipar content and LAST_REFRESH date.
# ----------
# The version is the date, since every update needs its own version and there
# end up being a lot of them; however it must match this regex:
# (?:^(?!-)[a-z\d\-]{0,62}[a-z\d]$)
go-bindata -pkg kisipar -prefix bindata bindata/... \
    && sed -i .bak "s/^LAST_REFRESH =.*/LAST_REFRESH = \"`date +%F-%H%M%S`\"/" version.go \
    && /bin/rm version.go.bak \
    && echo REFRESHED || echo FAILED
