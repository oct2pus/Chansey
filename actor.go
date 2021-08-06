package main

// NewFederatingActor builds a new Actor concept that handles only the Federating
// Protocol part of ActivityPub.
//
// This Actor can be created once in an application and reused to handle
// multiple requests concurrently and for different endpoints.
//
// It leverages as much of Go-Fed as possible to ensure the implementation is
// compliant with the ActivityPub specification, while providing enough freedom
// to be productive without shooting one's self in the foot.
//
// Do not try to use NewSocialActor and NewFederatingActor together to cover
// both the Social and Federating parts of the protocol. Instead, use NewActor.
func NewFederatingActor(c CommonBehavior,
	s2s FederatingProtocol,
	db Database,
	clock Clock) FederatingActor
