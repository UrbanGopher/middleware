/*

Package middleware - is an easy, yet powerful, tool for adding middleware to
your standard http or gorilla/mux routes.

This package is as simple and lightweight as it gets. Just take a look at the
source code, it's only 29 lines (not counting the unit tests and documentation,
of course) and has no 3rd party dependencies.

Middleware is used to enhance the functionality of route handlers without
bloating your handler functions or duplicating code across multiple handlers.
For those familiar with aspect oriented programming or cross-cutting concerns,
middleware is roughly analogous.

The canonical example of middleware is adding request logging across all of your
http handlers. You you may want it to do some logging before the handler is
invoked and some additional logging after the response is sent to the client.
This is not logic that you want to include in the handler method and have it
detract from the actual business logic, nor is it code that you want to repeat
in every handler method that it applies to.

Another common example is adding authentication/authorization across multiple
routes. As with logging, you want to keep this logic isolated and easily apply
it to each route that requires it.

The best way to understand is by example, so here is a fleshed out, fully
operational example of an http service that exposes two endpoints: "/" and
"/secret". logging is applied across both routes, however security is only
applied to the "/secret" route.


Example:

  package main

  import (
  	"fmt"
  	"net/http"
  	"time"

  	"github.com/UrbanGopher/middleware"
  )

  func main() {
  	mw := &middleware.Wrapper{}

  	mw.Append(loggingMiddleware) // add logging middleware
  	http.HandleFunc("/", mw.WrapHandler(handler))

  	mw.Append(securityMiddleware) // add some security middleware
  	http.HandleFunc("/secret", mw.WrapHandler(secretHandler))

  	http.ListenAndServe(":8080", nil)
  }

  func handler(w http.ResponseWriter, r *http.Request) {
  	w.Write([]byte("Hello, world!"))
  }

  func secretHandler(w http.ResponseWriter, r *http.Request) {
  	w.Write([]byte("You got in!"))
  }

  // A function that adheres to the middleware.Middleware specification
  func loggingMiddleware(next middleware.RouteHandler) middleware.RouteHandler {
  	return func(w http.ResponseWriter, r *http.Request) {
  		fmt.Println("Handling the request")
  		defer func(begin time.Time) {
  			fmt.Println("Request took: ", time.Since(begin))
  		}(time.Now())

  		next(w, r)
  	}
  }

  // Another function that adheres to middleware.Middleware
  func securityMiddleware(next middleware.RouteHandler) middleware.RouteHandler {
  	return func(w http.ResponseWriter, r *http.Request) {
  		qp := r.URL.Query()

  		if len(qp) > 0 && len(qp["word"]) > 0 && qp["word"][0] == "please" {
  			next(w, r)
  		} else {
  			w.WriteHeader(403)
  			w.Write([]byte("You didn't say the magic word!"))
  		}
  	}
  }

This middleware package can easily be applied to gorilla/mux routes as well
since they rely on handler functions with signatures that are identical to the
standard go http.HandleFunc handler functions. To demonstrate, the above
example's main method could be rewritten as follows:

  import (
    ...
    "github.com/gorilla/mux"
  )

  func main() {
    router := mux.NewRouter()
  	mw := &middleware.Wrapper{}

  	mw.Append(loggingMiddleware) // add logging middleware
  	router.HandleFunc("/", mw.WrapHandler(handler))

  	mw.Append(securityMiddleware) // add some security middleware
  	router.HandleFunc("/secret", mw.WrapHandler(secretHandler))

    http.Handle("/", router)
  	http.ListenAndServe(":8080", nil)
  }

*/
package middleware
