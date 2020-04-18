package resource

import (
	"github.com/estebarb/ion"
	"github.com/estebarb/ion/components/router"
	"net/http"
)

type Resource struct {
	Root         string
	Name         string
	Middleware   []ion.Middleware
	Subresources map[string]Resource
	Get          http.Handler
	Post         http.Handler
	GetID        http.Handler
	PostID       http.Handler
	DeleteID     http.Handler
}

func (res Resource) Attach(router *router.Router) {
	mid := ion.Chain(res.Middleware)

	if res.Get != nil {
		router.Get(res.Root, mid.Then(res.Get))
	}

	if res.Post != nil {
		router.Get(res.Root, mid.Then(res.Post))
	}

	if res.GetID != nil {
		router.Get(res.Root+":"+res.Name+"_id", mid.Then(res.GetID))
	}

	if res.PostID != nil {
		router.Get(res.Root+":"+res.Name+"_id", mid.Then(res.PostID))
	}

	if res.DeleteID != nil {
		router.Get(res.Root+":"+res.Name+"_id", mid.Then(res.DeleteID))
	}

	for subresourceName, subresource := range res.Subresources {
		subresource.Root = res.Root + subresourceName + "/"
		subresource.Middleware = append(res.Middleware, subresource.Middleware...)
		subresource.Attach(router)
	}
}
