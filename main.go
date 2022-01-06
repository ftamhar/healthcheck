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
	dateFormat = "2006-01-02 15:04:05 MST"
	idFormat   = "02-01-2006 15:04:05 MST"
)

var (
	myEmails emails

	server = new(mail.SMTPServer)
	client = new(mail.SMTPClient)

	email_host     = new(string)
	email_username = new(string)
	email_password = new(string)
	address        = new(string)

	email_port = new(int)

	delay  = new(time.Duration)
	bounce = new(time.Duration)

	//go:embed email.html
	emailHtml embed.FS
)

func init() {
	email_host = flag.String("h", "", "email host")
	email_username = flag.String("u", "", "email username")
	email_password = flag.String("p", "", "email password")
	address = flag.String("a", "", "alamat server yang akan dicek")

	email_port = flag.Int("port", 0, "email port")

	delay = flag.Duration("d", 3, "durasi di cek kembali setelah error (JAM)")
	bounce = flag.Duration("b", 10, "durasi di cek kembali setelah error (DETIK)")

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

	return err
}

type emails []string

func (e *emails) String() string {
	return "Oke deh"
}

func (i *emails) Set(s string) error {
	*i = append(*i, s)
	return nil
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

	email := mail.NewMSG()
	email.SetFrom("From Me <me@host.com>")
	email.AddTo(myEmails...)
	email.SetSubject("Notifikasi Server")

	t := time.Tick(*bounce * time.Second)

	log.Println("healthcheck started...")
	for {
		select {
		case <-t:
			_, err := http.Get(*address)
			now := time.Now()
			if err != nil {
				log.Printf("%s tidak dapat diakses, mengirim notifikasi...", *address)
				err = newClient()
				handleError(err)

				tmpl := template.Must(template.ParseFS(emailHtml, "email.html"))

				var s = new(strings.Builder)
				err := tmpl.Execute(s, server)
				if err != nil {
					panic(err.Error())
				}

				sNow := now.Format(dateFormat)
				email = email.SetDate(sNow)
				email.SetBody(mail.TextHTML, s.String())
				err = email.Send(client)
				handleError(err)

				time.Sleep(*delay * time.Hour)
				continue
			}
			sNow := now.Format(idFormat)
			server.LastCheck = sNow
		}
	}
}

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
