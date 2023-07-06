package main

import (
	"github.com/emicklei/go-restful/v3"
	"io"
	"log"
	"net/http"
)

// This example shows how to have a program with 2 WebServices containers
// each having a http server listening on its own port.
//
// The first "hello" is added to the restful.DefaultContainer (and uses DefaultServeMux)
// For the second "hello", a new container and ServeMux is created
// and requires a new http.Server with the container being the Handler.
// This first server is spawn in its own go-routine such that the program proceeds to create the second.
//
// GET http://localhost:8080/hello
// GET http://localhost:8081/hello

func main() {
	Ws := new(restful.WebService)
	Ws.Route(Ws.GET("/hello").To(hello1))
	restful.Add(Ws)
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	container2 := restful.NewContainer()
	Ws2 := new(restful.WebService)
	Ws2.Route(Ws2.GET("/hello").To(hello2))
	container2.Add(Ws2)
	server := &http.Server{Addr: ":8081", Handler: container2}
	log.Fatal(server.ListenAndServe())
	//

}
func hello1(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "Default world")
}
func hello2(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "second world")
}
