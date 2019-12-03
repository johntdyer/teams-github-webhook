package main

import "time"

type webexTeamsSpace struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Type         string    `json:"type"`
	IsLocked     bool      `json:"isLocked"`
	LastActivity time.Time `json:"lastActivity"`
	TeamID       string    `json:"teamId"`
	CreatorID    string    `json:"creatorId"`
	Created      time.Time `json:"created"`
	OwnerID      string    `json:"ownerId"`
}

type webexTeamsUser struct {
	ID           string   `json:"id"`
	Emails       []string `json:"emails"`
	PhoneNumbers []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"phoneNumbers"`
	DisplayName  string    `json:"displayName"`
	NickName     string    `json:"nickName"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	Avatar       string    `json:"avatar"`
	OrgID        string    `json:"orgId"`
	Created      time.Time `json:"created"`
	LastActivity time.Time `json:"lastActivity"`
	Status       string    `json:"status"`
	Type         string    `json:"type"`
}

type attachmentData struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	MessageID string `json:"messageId"`
	Inputs    struct {
		PullRequestID string `json:"pullRequestId"`
		RepoName      string `json:"repoName"`
		Comments      string `json:"comments"`
		ReviewStatus  string `json:"reviewStatus"`
	} `json:"inputs"`
	PersonID string    `json:"personId"`
	RoomID   string    `json:"roomId"`
	Created  time.Time `json:"created"`
}

// actionEvent
type webhookData struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	TargetURL string    `json:"targetUrl"`
	Resource  string    `json:"resource"`
	Event     string    `json:"event"`
	OrgID     string    `json:"orgId"`
	CreatedBy string    `json:"createdBy"`
	AppID     string    `json:"appId"`
	OwnedBy   string    `json:"ownedBy"`
	Status    string    `json:"status"`
	Created   time.Time `json:"created"`
	ActorID   string    `json:"actorId"`
	Data      struct {
		ID        string    `json:"id"`
		Type      string    `json:"type"`
		MessageID string    `json:"messageId"`
		PersonID  string    `json:"personId"`
		RoomID    string    `json:"roomId"`
		Created   time.Time `json:"created"`
	} `json:"data"`
}

type outputEvent struct {
	Body        string            `json:"Body"`
	Headers     map[string]string `json:"Headers"`
	Method      string            `json:"Method"`
	QueryParams string            `json:"QueryParams"`
}
