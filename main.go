package main

import (
	"embed"
	"flag"
	"net/http"
	"strings"
	"text/template"
	"time"

	"log"

	mail "github.com/xhit/go-simple-mail/v2"
)

const (
	emailFormat      = "2006-01-02 15:04:05 MST"
	indonesianFormat = "02-01-2006 15:04:05 MST"
)

type emails []string

func (e *emails) String() string {
	return "Oke deh"
}

func (i *emails) Set(s string) error {
	*i = append(*i, s)
	return nil
}

var (
	myEmails emails

	server = new(mail.SMTPServer)
	client = new(mail.SMTPClient)

	email_host     = flag.String("h", "", "host email")
	email_username = flag.String("u", "", "username email")
	email_password = flag.String("p", "", "password email")
	address        = flag.String("a", "", "alamat server yang akan dicek")

	email_port = flag.Int("port", 0, "port email")

	delay  = flag.Duration("d", 3, "durasi di cek kembali setelah error (JAM)")
	bounce = flag.Duration("b", 10, "durasi pengecekan (DETIK)")

	//go:embed email.html
	emailHtml embed.FS
)

func init() {
	flag.Var(&myEmails, "email", "list email yang akan dinotifikasi")
}

func newClient() (err error) {
	server = mail.NewSMTPClient()
	server.Host = *email_host
	server.Port = *email_port
	server.Username = *email_username
	server.Password = *email_password
	server.Encryption = mail.EncryptionTLS

	client, err = server.Connect()

	return
}

type Check struct {
	Server    string
	LastCheck string
}

func main() {
	flag.Parse()

	if *address == "" {
		log.Fatalln("alamat harus diisi, -a <alamat server>. cek --help untuk pemakaian")
	}

	if len(myEmails) == 0 {
		log.Fatalln("email tujuan harus di isi (boleh banyak). contoh: -email admin@local.com -email user@local.com")
	}

	server := Check{Server: *address}
	var now time.Time
	var err error
	var tmpl = new(template.Template)

	email := mail.NewMSG()
	email.SetFrom("From Me <me@host.com>")
	email.AddTo(myEmails...)
	email.SetSubject("Notifikasi Server")

	t := time.Tick(*bounce * time.Second)

	log.Println("healthcheck started...")
	for range t {
		_, err = http.Get(*address)
		now = time.Now()
		if err != nil {
			err = newClient()
			handleError(err)

			var s = new(strings.Builder)
			tmpl = template.Must(template.ParseFS(emailHtml, "email.html"))

			err = tmpl.Execute(s, server)
			handleError(err)

			email = email.SetDate(now.Format(emailFormat))
			email.SetBody(mail.TextHTML, s.String())
			err = email.Send(client)
			handleError(err)

			log.Printf("%s tidak dapat diakses, notifikasi telah dikirim...", *address)

			time.Sleep(*delay * time.Hour)
			continue
		}
		server.LastCheck = now.Format(indonesianFormat)
	}
}

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
