package main

import (
	"context"
	"net/http"
	"net/url"

	"github.com/go-fed/activity/pub"
	"github.com/go-fed/activity/streams/vocab"
)

type service struct{}

// CommonBehavior

func (*service) AuthenticateGetInbox(c context.Context,
	w http.ResponseWriter,
	r *http.Request) (out context.Context, authenticated bool, err error) {
	// TODO: Actually authenticate this.
	return r.Context(), true, nil
}

func (*service) AuthenticateGetOutbox(c context.Context,
	w http.ResponseWriter,
	r *http.Request) (out context.Context, authenticated bool, err error) {
	// TODO: Actually authenticate this.

	return r.Context(), true, nil
}

func (*service) GetOutbox(c context.Context,
	r *http.Request) (vocab.ActivityStreamsOrderedCollectionPage, error) {
	// TODO
	rc := r.Body
	r.URL
	return nil, nil
}

func (*service) NewTransport(c context.Context,
	actorBoxIRI *url.URL,
	gofedAgent string) (t pub.Transport, err error) {
	// TODO
	return
}

// FederatingProtocol

func (*service) PostInboxRequestBodyHook(c context.Context,
	r *http.Request,
	activity pub.Activity) (context.Context, error) {
	// TODO
	return nil, nil
}

func (*service) AuthenticatePostInbox(c context.Context,
	w http.ResponseWriter,
	r *http.Request) (out context.Context, authenticated bool, err error) {
	// TODO
	return
}

func (*service) Blocked(c context.Context,
	actorIRIs []*url.URL) (blocked bool, err error) {
	// TODO
	return
}

func (*service) FederatingCallbacks(c context.Context) (wrapped pub.FederatingWrappedCallbacks, other []interface{}, err error) {
	// Empty
	return
}

func (*service) DefaultCallback(c context.Context,
	activity pub.Activity) error {
	// TODO
	return nil
}

func (*service) MaxInboxForwardingRecursionDepth(c context.Context) int {
	// TODO
	return -1
}

func (*service) MaxDeliveryRecursionDepth(c context.Context) int {
	// TODO
	return -1
}

func (*service) FilterForwarding(c context.Context,
	potentialRecipients []*url.URL,
	a pub.Activity) (filteredRecipients []*url.URL, err error) {
	// TODO
	return
}

func (*service) GetInbox(c context.Context,
	r *http.Request) (vocab.ActivityStreamsOrderedCollectionPage, error) {
	// TODO
	return nil, nil
}
