Chme
====

[![CircleCI](https://circleci.com/gh/tomocy/chme.svg?style=svg)](https://circleci.com/gh/tomocy/chme)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Chme provides middlewares which enable you to do requests with other methods than GET and POST from HTML form.  
Chme stands for "change methods".   

## Description
HTML form does not support methods except GET and POST.   
But you sometimes want it to request with more proper methods such as when we provide some edit feature.   

## Usage
The signature of middlewares in this repository is the standard one.  
```go
func(next http.Handler) http.Handler
```

### ChangePostToHiddenMethod
The usage is very simple.   
All you have to do is add input named "_method" in POST form.  

```html
<form action="/articles/1" method="POST">
    <input type="hidden" name="_method" value="PUT">
</form>
```

And in Go, you can handle the request as if it is done with the method set in the input.  ([chi](https://github.com/go-chi/chi) is used as an example.)
```go
r := chi.NewRouter()
r.Use(ChangePostToHiddenMethod)
r.Put("/articles/{id}", func(w http.ResponseWriter, r *http.Request) {
    // do something
})
```

You can change POST to PUT, PATCH, and DELETE.   
Make sure to specify html form method as POST because the default method of HTML form is GET and this function ignores requests with GET.  

You can change default input name "_method" to anything you want.  
```go
// Change to "other" for the input
r.Use(NewChme("other").ChangePostToHiddenMethod)
```

## Install
```
go get github.com/tomocy/chme
```

## Licence
[MIT](https://github.com/tcnksm/tool/blob/master/LICENCE)

## Author
tomocy
