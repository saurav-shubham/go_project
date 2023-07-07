package main

import (
	"fmt"
	restful "github.com/emicklei/go-restful/v3"
	"log"
	"net/http"
	"strings"
)

// UserResource is the REST layer to the User domain
type UserResource struct {
	// normally one would use DAO (data access object)
	users map[string]User
}
type abc int //this is the new type of the data type, "abc data type"

// WebService creates a new service that can handle REST requests for User resources.
func setRootPath(ws *restful.WebService, path string) {
	log.Println("setting root path:", path, " for webservice:", ws.Documentation())
	ws.Path(path)
}
func RegisterApi(ws *restful.WebService, writes interface{}, httpMethodName, apiPath, doc string, handler func(request *restful.Request, response *restful.Response)) {
	log.Println("Registering api:", apiPath, "with http method:", httpMethodName)
	//ws.Doc(doc).Method(strings.ToUpper(httpMethodName)).Path(apiPath).To(handler).Writes(writes)
	ws.Route(ws.Doc(doc).Method(strings.ToUpper(httpMethodName)).Path(apiPath).To(handler).Writes(writes))
}
func (a abc) testMethod() {
	fmt.Println(a)
}
func (u UserResource) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Doc("Employee")
	setRootPath(ws, "/user")
	ws.Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well
	RegisterApi(ws, []User{}, "GET", "", "get all user list", u.findAllUsers)
	RegisterApi(ws, User{}, "GET", "/{user-id}", "get a user", u.findUser)

	ws.Route(ws.PUT("/{user-id}").To(u.updateUser).
		// docs
		Doc("update a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		//Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(User{})) // from the request

	ws.Route(ws.POST("/{user-id}").To(u.createUser).
		// docs
		Doc("create a user").
		//Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(User{})) // from the request

	ws.Route(ws.DELETE("/{user-id}").To(u.removeUser).
		// docs
		Doc("delete a user").
		//	Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")))

	return ws
}

// GETr
func (u UserResource) findAllUsers(request *restful.Request, response *restful.Response) {
	list := []User{}
	for _, each := range u.users {
		list = append(list, each)
	}
	response.WriteEntity(list)
}

// GET http://localhost:8080/users/1
func (u UserResource) findUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	usr := u.users[id]
	if len(usr.ID) == 0 {
		response.WriteErrorString(http.StatusNotFound, "User could not be found.")
	} else {
		response.WriteEntity(usr)
	}
}

// PUT http://localhost:8080/users/1
// <User><Id>1</Id><Name>Melissa Raspberry</Name></User>
func (u *UserResource) updateUser(request *restful.Request, response *restful.Response) {
	usr := new(User)
	err := request.ReadEntity(&usr)
	if err == nil {
		u.users[usr.ID] = *usr
		response.WriteEntity(usr)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

// PUT http://localhost:8080/users/1
// <User><Id>1</Id><Name>Melissa</Name></User>
func (u *UserResource) createUser(request *restful.Request, response *restful.Response) {
	usr := User{ID: request.PathParameter("user-id")}
	roll := request.QueryParameters("roll")
	fmt.Println(roll)
	err := request.ReadEntity(&usr)
	if err == nil {
		u.users[usr.ID] = usr
		response.WriteHeaderAndEntity(http.StatusCreated, usr)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

// DELETE http://localhost:8080/users/1
func (u *UserResource) removeUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	delete(u.users, id)
}

func main() {
	u := UserResource{map[string]User{}}
	log.Printf("\n ##################Registering api############")
	restful.DefaultContainer.Add(u.WebService())
	port := ":8081"
	log.Printf("start listening on localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// User is just a sample type
type User struct {
	ID   string `json:"id" description:"identifier of the user"`
	Name string `json:"name" description:"name of the user" default:"john"`
	Age  int    `json:"age" description:"age of the user" default:"21"`
}
