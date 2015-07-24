[Go](https://golang.org) client package for
the [Ninchat](https://ninchat.com) [API](https://ninchat.com/api).
Compatible with [GopherJS](http://www.gopherjs.org)
(see also [ninchat-js](https://github.com/ninchat/ninchat-js)).
Uses [Gorilla WebSocket](http://www.gorillatoolkit.org/pkg/websocket) package
with native Go.

The implementation uses few standard packages directly in order to reduce
generated JavaScript code size.  Subsets of encoding/json, math/rand, net/http
and time packages are wrapped for native Go, and reimplemented using web
browser APIs for GopherJS.
