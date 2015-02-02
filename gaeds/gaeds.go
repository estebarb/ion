// This package provides a thin layer
// over GAE datastore and memcache.
//
// The operations made with gaeds automatically
// use memcache when possible, improving the performance
// and reducing the boilerplate and billing.
//
// When in doubt, please read the code of gaeds.
package gaeds

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
)

// Generates a identifier to be used with memcache
func genMemcacheKey(bucket, key string) string {
	return bucket + "_" + key
}

// Saves a struct on the Datastore and Memcache. Example:
//      err := gaeds.SaveAs(c, "framework", framework.Id, &fw)
func SaveAs(c appengine.Context, bucket, key string, object interface{}) error {
	ck := datastore.NewKey(c, bucket, key, 0, nil)
	_, err := datastore.Put(c, ck, object)
	if err != nil {
		return err
	}

	item := &memcache.Item{
		Key:    genMemcacheKey(bucket, key),
		Object: object,
	}
	memcache.Gob.Set(c, item)

	return err
}

// Gets a item from the datastore or memcache and saves it to
// a struct.
//      err := gaeds.Get(c, "framework", "ionframework", &fw)
func Get(c appengine.Context, bucket, key string, object interface{}) error {
	// First we try to get from memcache
	mcKey := genMemcacheKey(bucket, key)
	_, err := memcache.Gob.Get(c, mcKey, object)

	if err == nil {
		return nil
	}

	// If there is a problem then we try from datastore:
	ck := datastore.NewKey(c, bucket, key, 0, nil)
	err = datastore.Get(c, ck, object)
	if err == nil {
		// Save the retrieved object in memcache
		item := &memcache.Item{
			Key:    mcKey,
			Object: object,
		}
		memcache.Gob.Set(c, item)
	}
	return err
}

// Deletes a item from the datastore and memcache and saves it to
// a struct.
//      err := gaeds.Delete(c, "framework", "ionframework")
func Delete(c appengine.Context, bucket, key string) error {
	// First we delete it from datastore
	ck := datastore.NewKey(c, bucket, key, 0, nil)
	err := datastore.Delete(c, ck)

	// Delete the object from memcache
	memcache.Delete(c, genMemcacheKey(bucket, key))
	return err
}
