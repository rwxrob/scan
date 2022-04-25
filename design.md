# Design Considerations

* **One job, one function/method**

  Instead of `Peek(a any)` that takes a string and a regular expression
  break them out into two, `Peek(a string)` and `Match(a
  *regexp.Regexp)`. We want to keep this scanner the fastest possible of
  its class.

* **No need for Unicode method since in Go Regexp**

  Go now fully supports the Unicode character class notation (`\p{L}`)
  so the `Match(a *regexp.Regexp)` covers everything needed.
