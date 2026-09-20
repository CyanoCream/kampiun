package notify

import "context"

const (
	KindUserRegister  = "user.registered"
	KindOrderCreated  = "order.created"
	KindOrderProof    = "order.proof"
	KindOrderApproved = "order.approved"
	KindOrderRejected = "order.rejected"
)

type Action struct {
	Label string
	Data  string
}

type Notification struct {
	Kind    string
	Title   string
	Lines   []string
	Link    string
	Actions []Action
}

type Notifier interface {
	Notify(ctx context.Context, n Notification)
}

// Nop = notifier kosong (bila Telegram tak dikonfigurasi).
type Nop struct{}

func (Nop) Notify(context.Context, Notification) {}
