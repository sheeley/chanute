package chanute

import "github.com/aws/aws-sdk-go/aws/session"

type NamedSession struct {
	Name    string
	Session *session.Session
}
