// IonMVC
// IonMVC is a experiment for developing high quality websites
// typing less code, especifically, eliminating a lot of boilter
// plate that is necesary to write when using Ion alone.
//
// For example, processing a form is always the same code: parse the form,
// load it to a struct, validate it and process it accordling.
//
// This packages provides a more complete, automatic MVC for Ion Framework.
// It considers the data struct as the basic unit of work, and generates
// the RESTful style actions based on them.
//
// The package isn't as flexible as the pure Ion Framework, but it does what
// most probably you would do in any case.
//
// It receives a struct and a path descriptor, that is used to generate the URL
// paths. IonMVC creates a RESTful path for the CRUD actions:
//
// - GET    /entity/            Lists all the entities
// - GET    /entity/new         Shows the create entity form
// - POST   /entity/            Creates a new entity
// - GET    /entity/:id         Shows the entity
// - GET    /entity/:id/edit    Shows the edit entity form
// - PUT    /entity/:id         Updates the entity
// - POST   /entity/:id         Updates the entity
// - POST   /entity/:id/delete  Deletes the entity
// - DELETE /entity/:id         Deletes the entity
//
// The package will create a Human HTML and JSON interface,
// but its possible to add more (XML, YAML, MessagePack, etc).
//
// The data struct can implement several interfaces, used by IonMVC to
// transform and guide the request flow. Specifically, IonMVC expects
// interfaces for the next categories:
// - Validation:    Validates the data.
// - Processing:    Allows to preprocess the data before executing actions.
// - Authorization: Determines if a operation is valid or not.
//
//
package ionmvc

import (
	"github.com/estebarb/ion"
	"github.com/estebarb/ion/ionmvc/dal"
	"github.com/estebarb/ion/ionvc"
	"net/http"
)

// Creates a path to be used by the router:
// Example:
// generatePath(["user", "blog"]) returns
// "/user/:user/blog/:blog"
func generatePath(ids []string) string{
	str := ""
	for _, value := range ids{
		str += "/"+value+"/:"+value
	}
	return str
}

func CRUD(r *ion.Router, DAL dal.DAL, ids []string, entity interface{}){

}

type CRUDer interface{
	Middleware()
}

// Creates the handlers that allow to create a entity.
// Expects that ionvc already have registered a template
// named <entity>_new
func Create(r *ion.Router, DAL dal.DAL, ids[]string, handler http.HandlerFunc){
	// GET      /entity/new
	r.GetFunc(generatePath(ids[0:-1])+"/new", ionvc.ControllerFunc(handler, ids[-1]+"_new"))
}

