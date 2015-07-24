[Go](https://golang.org) client package for the [Ninchat](https://ninchat.com)
API.  Compatible with [GopherJS](http://www.gopherjs.org) (see also
[ninchat-js](https://github.com/ninchat/ninchat-js)).

The implementation uses few standard packages directly in order to reduce
generated JavaScript code size.  Subsets of encoding/json, math/rand, net/http
and time packages are wrapped for native Go, and reimplemented using web
browser APIs for GopherJS.
