package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strings"

	"github.com/go-fed/activity/pub"
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/go-redis/redis/v8"
)

type database struct {
	client *redis.Client
}

func (db *database) Lock(c context.Context, id *url.URL) error {
	// letting redis handle this
	return nil
}

func (db *database) Unlock(c context.Context, id *url.URL) error {
	// letting redis handle this
	return nil
}

func (db *database) Get(c context.Context, id *url.URL) (value vocab.Type, err error) {
	entry, err := db.client.Get(c, id.String()).Result()
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}

	err = json.Unmarshal([]byte(entry), &data)
	if err != nil {
		return nil, err
	}
	return streams.ToType(c, data)
}

func (db *database) Create(c context.Context, asType vocab.Type) error {
	url, err := pub.GetId(asType)
	if err != nil {
		return err
	}
	// This could go through Exists() instead
	i, err := db.client.Exists(c, url.String()).Result()
	if err != nil {
		return err
	}
	if i == 0 {
		return db.createUpdate(c, asType)
	}

	return errors.New("key exists")
}

func (db *database) Update(c context.Context, asType vocab.Type) error {
	return db.createUpdate(c, asType)
}

func (db *database) Exists(c context.Context, id *url.URL) (exists bool, err error) {
	i, err := db.client.Exists(c, id.String()).Result()
	if i == 1 {
		return true, err
	}
	return false, err
}

func (db *database) Delete(c context.Context, id *url.URL) error {
	return db.client.Del(c, id.String()).Err()
}

func (db *database) NewID(c context.Context, t vocab.Type) (id *url.URL, err error) {
	return url.Parse("https://" + HOSTNAME + "/" + t.GetTypeName() + "/" + string(rand.Int()))
}

func (db *database) InboxContains(c context.Context, inbox, id *url.URL) (contains bool, err error) {
	inboxEntry, err := db.client.Get(c, inbox.String()).Result()
	if err != nil {
		return false, err
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(inboxEntry), &data)
	if err != nil {
		return false, err
	}

	for _, v := range data {
		if fmt.Sprintf("%v", v) == id.String() {
			return true, nil
		}
	}

	return false, nil
}

func (db *database) GetInbox(c context.Context, inboxIRI *url.URL) (inbox vocab.ActivityStreamsOrderedCollectionPage, err error) {
	inboxEntry, err := db.client.Get(c, inboxIRI.String()).Result()
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(inboxEntry), &data)
	page := streams.NewActivityStreamsOrderedCollectionPage()
	resolver, err := streams.NewJSONResolver(func(c context.Context, p vocab.ActivityStreamsOrderedCollectionPage) error {
		page = p
		return nil
	})
	if err != nil {
		return nil, err
	}

	err = resolver.Resolve(c, data)
	if err != nil {
		return nil, err
	}

	return page, nil
}

func (db *database) SetInbox(c context.Context, inbox vocab.ActivityStreamsOrderedCollectionPage) error {
	orderedItems := inbox.GetActivityStreamsOrderedItems()
	if orderedItems == nil || orderedItems.Len() == 0 {
		return nil
	}

	originalInboxURL, err := pub.GetId(inbox.GetActivityStreamsPartOf().GetActivityStreamsOrderedCollection())
	if err != nil {
		return err
	}

	inboxEntry, err := db.client.Get(c, originalInboxURL.String()).Result()
	if err != nil {
		return err
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(inboxEntry), &data)
	if err != nil {
		return err
	}
	var orderedCollection vocab.ActivityStreamsOrderedCollection
	resolver, err := streams.NewJSONResolver(func(c context.Context, o vocab.ActivityStreamsOrderedCollection) error {
		orderedCollection = o
		return nil
	})
	if err != nil {
		return err
	}

	err = resolver.Resolve(c, data)
	if err != nil {
		return err
	}

	inboxOrderedItems := orderedCollection.GetActivityStreamsOrderedItems()
	inboxOrderedItems.PrependActivityStreamsOrderedCollectionPage(inbox)
	orderedCollection.SetActivityStreamsOrderedItems(inboxOrderedItems)

	ocData, err := orderedCollection.Serialize()
	if err != nil {
		return err
	}

	entry, err := json.Marshal(&ocData)
	if err != nil {
		return err
	}
	return db.client.Set(c, originalInboxURL.String(), entry, 0).Err()
}

func (db *database) GetOutbox(c context.Context, outboxIRI *url.URL) (inbox vocab.ActivityStreamsOrderedCollectionPage, err error) {
	outboxEntry, err := db.client.Get(c, outboxIRI.String()).Result()
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(outboxEntry), &data)
	page := streams.NewActivityStreamsOrderedCollectionPage()
	resolver, err := streams.NewJSONResolver(func(c context.Context, p vocab.ActivityStreamsOrderedCollectionPage) error {
		page = p
		return nil
	})
	if err != nil {
		return nil, err
	}

	err = resolver.Resolve(c, data)
	if err != nil {
		return nil, err
	}

	return page, nil
}

func (db *database) SetOutbox(c context.Context, outbox vocab.ActivityStreamsOrderedCollectionPage) error {
	orderedItems := outbox.GetActivityStreamsOrderedItems()
	if orderedItems == nil || orderedItems.Len() == 0 {
		return nil
	}

	originaloutboxURL, err := pub.GetId(outbox.GetActivityStreamsPartOf().GetActivityStreamsOrderedCollection())
	if err != nil {
		return err
	}

	outboxEntry, err := db.client.Get(c, originaloutboxURL.String()).Result()
	if err != nil {
		return err
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(outboxEntry), &data)
	if err != nil {
		return err
	}
	var orderedCollection vocab.ActivityStreamsOrderedCollection
	resolver, err := streams.NewJSONResolver(func(c context.Context, o vocab.ActivityStreamsOrderedCollection) error {
		orderedCollection = o
		return nil
	})
	if err != nil {
		return err
	}

	err = resolver.Resolve(c, data)
	if err != nil {
		return err
	}

	outboxOrderedItems := orderedCollection.GetActivityStreamsOrderedItems()
	outboxOrderedItems.PrependActivityStreamsOrderedCollectionPage(outbox)
	orderedCollection.SetActivityStreamsOrderedItems(outboxOrderedItems)

	ocData, err := orderedCollection.Serialize()
	if err != nil {
		return err
	}

	entry, err := json.Marshal(&ocData)
	if err != nil {
		return err
	}
	return db.client.Set(c, originaloutboxURL.String(), entry, 0).Err()
}

func (db *database) Owns(c context.Context, id *url.URL) (owns bool, err error) {
	// TODO: this seems like a naive implimentation
	return strings.Contains(id.String(), HOSTNAME), nil
}

func (db *database) ActorForOutbox(c context.Context, outboxIRI *url.URL) (actorIRI *url.URL, err error) {
	return url.Parse(strings.TrimSuffix(outboxIRI.String(), "/out"))
}

func (db *database) ActorForInbox(c context.Context, inboxIRI *url.URL) (actorIRI *url.URL, err error) {
	return url.Parse(strings.TrimSuffix(inboxIRI.String(), "/in"))
}

func (db *database) OutboxForInbox(c context.Context, inboxIRI *url.URL) (outboxIRI *url.URL, err error) {
	outbox := strings.TrimSuffix(inboxIRI.String(), "/in")
	outbox += "/out"
	return url.Parse(outbox)
}

func (db *database) Followers(c context.Context, actorIRI *url.URL) (followers vocab.ActivityStreamsCollection, err error) {
	person, err := db.getPersonFromIRI(c, actorIRI)
	if err != nil {
		return nil, err
	}

	followersProperty := person.GetActivityStreamsFollowers()
	if followersProperty.IsActivityStreamsCollection() {
		return followersProperty.GetActivityStreamsCollection(), err
	}

	// May not be an orderedcollection unfortunately :()
	fID, err := pub.GetId(followersProperty.GetType())
	if err != nil {
		return nil, err
	}
	collection, err := db.getCollectionFromID(c, fID)
	return collection, err
}

func (db *database) Following(c context.Context, actorIRI *url.URL) (following vocab.ActivityStreamsCollection, err error) {
	person, err := db.getPersonFromIRI(c, actorIRI)
	if err != nil {
		return nil, err
	}
	followingProperty := person.GetActivityStreamsFollowing()
	if followingProperty.IsActivityStreamsCollection() {
		return followingProperty.GetActivityStreamsCollection(), err
	}

	// May not be an orderedcollection unfortunately :()
	fID, err := pub.GetId(followingProperty.GetType())
	if err != nil {
		return nil, err
	}

	collection, err := db.getCollectionFromID(c, fID)
	return collection, err
}

func (db *database) Liked(c context.Context, actorIRI *url.URL) (liked vocab.ActivityStreamsCollection, err error) {
	person, err := db.getPersonFromIRI(c, actorIRI)
	if err != nil {
		return nil, err
	}

	likedProperty := person.GetActivityStreamsLiked()
	if likedProperty.IsActivityStreamsCollection() {
		return likedProperty.GetActivityStreamsCollection(), err
	}

	// May not be an orderedcollection unfortunately :()
	lID, err := pub.GetId(likedProperty.GetType())
	if err != nil {
		return nil, err
	}
	collection, err := db.getCollectionFromID(c, lID)
	return collection, err
}

// helper funcs

func (db *database) createUpdate(c context.Context, asType vocab.Type) error {
	data, err := streams.Serialize(asType)
	if err != nil {
		return err
	}

	entry, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	url, err := pub.GetId(asType)
	if err != nil {
		return err
	}
	return db.client.Set(c, url.String(), entry, 0).Err()
}

func (db *database) getCollectionFromID(c context.Context, ID *url.URL) (collection vocab.ActivityStreamsCollection, err error) {
	entry, err := db.client.Get(c, ID.String()).Result()
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(entry), &data)
	if err != nil {
		return nil, err
	}
	resolver2, err := streams.NewJSONResolver(func(c context.Context, col vocab.ActivityStreamsCollection) error {
		collection = col

		return nil
	})
	if err != nil {
		return nil, err
	}

	err = resolver2.Resolve(c, data)
	return collection, err
}

func (db *database) getPersonFromIRI(c context.Context, IRI *url.URL) (person vocab.ActivityStreamsPerson, err error) {
	entry, err := db.client.Get(c, IRI.String()).Result()
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(entry), &data)
	if err != nil {
		return nil, err
	}

	resolver, err := streams.NewJSONResolver(func(c context.Context, p vocab.ActivityStreamsPerson) error {
		person = p
		return nil
	})
	if err != nil {
		return nil, err
	}

	err = resolver.Resolve(c, data)
	return person, err
}

// new func

func databasebNew() *database {
	var db *database
	db.client = redis.NewClient(&redis.Options{})
	return db
}
