package dal
import (
	"net/http"
)
// The DAL abstracts the underlying "database". A DAL must provide
// the typical CRUD operations.
// For most operations:
// - *http.Request is the request by user, used by some DB (eg: GAE)
// - []string is the key. (eg: [Esteban MyBlog MyAwesomePost] -> /user/Esteban/blog/MyBlog/post/MyAwesomePost)
// - interface{} is where the result will be writen, or the data to be writen.
// - ...string optional arguments. Eg: for pagination.
// The result will be always the key and if an error happens.
type DAL interface{
	// Lists the entities, where []string is the key. Puts the result in interface{}
	List(*http.Request, []string, interface{}, ...string) ([]string, error)
	Create(*http.Request, []string, interface{}, ...string) ([]string, error)
	Read(*http.Request, []string, interface{}, ...string) ([]string, error)
	Update(*http.Request, []string, interface{}, ...string) ([]string, error)
	Delete(*http.Request, []string, interface{}, ...string) ([]string, error)
}