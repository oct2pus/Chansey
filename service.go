package main

import (
	"net/http"
	"net/url"
)

type service struct{}

// CommonBehavior

func (*service) AuthenticateGetInbox(c context.Contect,
	w http.ResponseWriter,
	r *http.Request) (out context.Context, authenticated bool, err error) {
	// TODO
	return
}

func (*myService) AuthenticateGetOutbox(c context.Context,
	w http.ResponseWriter,
	r *http.Request) (out context.Context, authenticated bool, err error) {
	// TODO
	return
}

func (*myService) GetOutbox(c context.Context,
	r *http.Request) (vocab.ActivityStreamsOrderedCollectionPage, error) {
	// TODO
	return nil, nil
}

func (*myService) NewTransport(c context.Context,
	actorBoxIRI *url.URL,
	gofedAgent string) (t pub.Transport, err error) {
	// TODO
	return
}

// FederatingProtocol

func (*myService) PostInboxRequestBodyHook(c context.Context,
	r *http.Request,
	activity Activity) (context.Context, error) {
	// TODO
	return nil, nil
}

func (*myService) AuthenticatePostInbox(c context.Context,
	w http.ResponseWriter,
	r *http.Request) (out context.Context, authenticated bool, err error) {
	// TODO
	return
}

func (*myService) Blocked(c context.Context,
	actorIRIs []*url.URL) (blocked bool, err error) {
	// TODO
	return
}

func (*myService) FederatingCallbacks(c context.Context) (wrapped FederatingWrappedCallbacks, other []interface{}, err error) {
	// TODO
	return
}

func (*myService) DefaultCallback(c context.Context,
	activity Activity) error {
	// TODO
	return nil
}

func (*myService) MaxInboxForwardingRecursionDepth(c context.Context) int {
	// TODO
	return -1
}

func (*myService) MaxDeliveryRecursionDepth(c context.Context) int {
	// TODO
	return -1
}

func (*myService) FilterForwarding(c context.Context,
	potentialRecipients []*url.URL,
	a Activity) (filteredRecipients []*url.URL, err error) {
	// TODO
	return
}

func (*myService) GetInbox(c context.Context,
	r *http.Request) (vocab.ActivityStreamsOrderedCollectionPage, error) {
	// TODO
	return nil, nil
}
