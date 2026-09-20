// Package email = pengirim email via SMTP (stdlib net/smtp). Mengimplementasikan
// notify.Notifier tanpa mengubah kontrak service (Lines => paragraf HTML).
//
// Kredensial SMTP diisi belakangan (SMTP_HOST/PORT/USER/PASS). Selama kosong
// (Host ""), Mailer menjadi Nop — tidak melempari log error.
package email

import (
	"context"
	"fmt"
	"log/slog"
	"net/smtp"
	"strings"

	"kampiun/kernel/notify"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Pass     string
	From     string
	FromName string
}

func (c Config) Enabled() bool { return c.Host != "" }

// Mailer mengirim email; bila SMTP tak dikonfigurasi, Notify jadi no-op.
type Mailer struct {
	cfg Config
	log *slog.Logger
	To  string // penerima notifikasi (admin/dev) sampai ada relasi user->email
}

func New(cfg Config, to string, log *slog.Logger) *Mailer { return &Mailer{cfg: cfg, log: log, To: to} }

// Notify mengirim email utk notifikasi; subject dari Kind.
func (m *Mailer) Notify(ctx context.Context, n notify.Notification) {
	if !m.cfg.Enabled() || m.To == "" {
		return
	}
	if err := m.send(subjectFor(n.Kind, n.Title), n.Lines); err != nil {
		m.log.Error("kirim email", "to", m.To, "err", err)
	}
}

func (m *Mailer) send(subject string, lines []string) error {
	body := renderHTML(subject, m.cfg.FromName, lines)
	addr := fmt.Sprintf("%s:%d", m.cfg.Host, m.cfg.Port)
	var auth smtp.Auth
	if m.cfg.User != "" {
		auth = smtp.PlainAuth("", m.cfg.User, m.cfg.Pass, m.cfg.Host)
	}
	from := m.cfg.From
	if m.cfg.FromName != "" {
		from = fmt.Sprintf("%s <%s>", m.cfg.FromName, m.cfg.From)
	}
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		from, m.To, subject, body)
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if auth != nil {
		if err := c.Auth(auth); err != nil {
			return err
		}
	}
	if err := c.Mail(from); err != nil {
		return err
	}
	if err := c.Rcpt(m.To); err != nil {
		return err
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err = w.Write([]byte(msg)); err != nil {
		return err
	}
	return w.Close()
}

func renderHTML(subject, fromName string, lines []string) string {
	var b strings.Builder
	b.WriteString("<!DOCTYPE html><html><body style=\"font-family:sans-serif;color:#1a1a2e\">")
	fmt.Fprintf(&b, "<h2 style=\"font-weight:300\">%s</h2>", htmlEscape(subject))
	for _, l := range lines {
		fmt.Fprintf(&b, "<p style=\"color:#333\">%s</p>", htmlEscape(l))
	}
	b.WriteString("</body></html>")
	return b.String()
}

func subjectFor(kind, fallback string) string {
	switch kind {
	case notify.KindUserRegister:
		return "Selamat datang di Kampiun"
	}
	return fallback
}

func htmlEscape(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&#34;")
	return r.Replace(s)
}
