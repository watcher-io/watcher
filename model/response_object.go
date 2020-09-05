package model

// LoginResponse is a object used for sending the generated JWT token to the client
// after successful user validation
type LoginResponse struct {
	Token string `json:"token"`
}

// PutKVResponse is used to send the PUT key-value response
type PutKVResponse struct {
	Revision int64  `json:"revision"`
	MemberID uint64 `json:"member_id"`
	RaftTerm uint64 `json:"raft_term"`
}
