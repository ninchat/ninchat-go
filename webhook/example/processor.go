// Use "go run processor.go" to run this example program.

// +build ignore

package main

import (
	"crypto/ed25519" // golang.org/x/crypto/ed25519 supports old Go versions.
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ninchat/ninchat-go/webhook"
)

// HTTPS server program.
func main() {
	var (
		server = http.Server{
			Addr: ":443",
			TLSConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		}
		cert    = "cert.pem"
		certKey = "certkey.pem"
	)

	flag.StringVar(&server.Addr, "-addr", server.Addr, "Listening [host]:port")
	flag.StringVar(&cert, "-cert", cert, "TLS certificate filename")
	flag.StringVar(&certKey, "-certkey", certKey, "TLS certificate's private key filename")
	flag.Parse()

	http.HandleFunc("/ninchat-webhook", handleNinchatWebhook)
	log.Fatal(server.ListenAndServeTLS(cert, certKey))
}

// Add new key here when one is announced.  Remove the old one after it has
// been phased out.
var ninchatPublicKeys = map[string]string{
	"ninchat.com/ed25519-2019-04": "fyMJMmsTPDkp4YYxlcQcWrIwTvalCgfHMW370OE/y2I",
}

// Mitigate replay attacks by remembering X-Ninchat-Signatures until their
// expiration.  This is an optional security measure.  To be effective, the
// information should be shared among all nodes of a distributed system.
var (
	seenLock       sync.Mutex
	seenSignatures = map[string]int64{}
)

func rememberSignature(s string, exp int64) (duplicate bool) {
	seenLock.Lock()
	defer seenLock.Unlock()

	_, duplicate = seenSignatures[s]
	if !duplicate {
		seenSignatures[s] = exp
	}
	return
}

// A naive background job for removing expired signatures.
func init() {
	go func() {
		for {
			time.Sleep(time.Minute)

			seenLock.Lock()

			var expired []string

			for sig, exp := range seenSignatures {
				if exp < time.Now().Unix() {
					expired = append(expired, sig)
				}
			}

			for _, sig := range expired {
				delete(seenSignatures, sig)
			}

			seenLock.Unlock()
		}
	}()
}

// handleNinchatWebhook processes a HTTP request.
func handleNinchatWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if !strings.HasPrefix(r.Header.Get("User-Agent"), "ninchat-webhook/") {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Header.Get("Content-Type") != "application/json; charset=utf-8" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	if r.ContentLength < 0 {
		w.WriteHeader(http.StatusLengthRequired)
		return
	}
	if r.ContentLength > 256*1024*1024 { // Something huge...
		w.WriteHeader(http.StatusRequestEntityTooLarge) // Let Ninchat know.
		return
	}

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var document webhook.Webhook
	if json.Unmarshal(content, &document) != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if document.Aud != "realm:abcd1234" { // This Ninchat resource id must be yours.
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if document.Exp < time.Now().Unix() { // Expired?
		w.WriteHeader(http.StatusForbidden)
		return
	}

	publicKeyString, found := ninchatPublicKeys[document.Kid]
	if !found {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	publicKey, err := base64.RawURLEncoding.DecodeString(publicKeyString)
	if err != nil {
		panic(err)
	}

	signatureString := r.Header.Get("X-Ninchat-Signature")
	signature, err := hex.DecodeString(signatureString)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if !ed25519.Verify(publicKey, content, signature) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if rememberSignature(signatureString, document.Exp) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	switch document.Event {
	case webhook.EventWebhookVerification:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(document.WebhookVerificationResponse())
		return

	case webhook.EventAudienceRequested:
		if _, err := document.AudienceRequested(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Println("audience requested")

	case webhook.EventAudienceAccepted:
		if _, err := document.AudienceAccepted(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Println("audience accepted")

	case webhook.EventAudienceComplete:
		params, err := document.AudienceComplete()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		handleAudienceComplete(w, params)

	case webhook.EventMessageSent:
		if _, err := document.MessageSent(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	case webhook.EventDataAccess:
		params, err := document.DataAccess()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		handleDataAccess(w, params)

	default:
		// Misconfiguration at Ninchat end?  Or if event handling code is
		// missing from this end, Ninchat will retry later.
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// Below is a demo implementation of on-site data storage:

type myChannel struct {
	audienceMetadata      map[string]json.RawMessage
	memberMessageMetadata map[string]webhook.Metadata
	messages              []myMessage
}

type myMessage struct {
	id       string
	time     time.Time
	typ      string
	userID   string
	userName *string
	payload  json.RawMessage
}

var (
	storageLock sync.Mutex
	myStorage   = map[string]*myChannel{}
)

func handleAudienceComplete(w http.ResponseWriter, params webhook.AudienceComplete) {
	if params.ChannelID == "" {
		return
	}

	storageLock.Lock()
	_, duplicate := myStorage[params.ChannelID]
	storageLock.Unlock()
	if duplicate {
		// Ninchat resent an earlier (successful) event, probably because it
		// didn't see the response due to network problem.  Do nothing except
		// indicate success.
		w.WriteHeader(http.StatusOK)
		return
	}

	data := new(myChannel)
	data.audienceMetadata = params.Audience.Metadata
	data.memberMessageMetadata = params.MemberMessageMetadata

	for _, m := range params.Messages {
		data.messages = append(data.messages, myMessage{
			id:       m.ID,
			time:     time.Unix(int64(m.Time), int64(math.Mod(m.Time, 1)*float64(time.Nanosecond))),
			typ:      string(m.Type),
			userID:   m.UserID,
			userName: m.UserName,
			payload:  m.PayloadJSON,
		})
	}

	storageLock.Lock()
	myStorage[params.ChannelID] = data
	storageLock.Unlock()

	w.WriteHeader(http.StatusOK)
}

func handleDataAccess(w http.ResponseWriter, params webhook.DataAccess) {
	if params.ChannelID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	storageLock.Lock()
	data, found := myStorage[params.ChannelID]
	storageLock.Unlock()
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var response webhook.DataAccessResponse

	if params.Query.AudienceMetadata {
		response.Audience = &webhook.Audience{
			Metadata: data.audienceMetadata,
		}
	}

	if params.Query.MemberMessageMetadata {
		response.MemberMessageMetadata = data.memberMessageMetadata
	}

	if params.Query.Messages != nil {
		for _, m := range data.messages {
			if params.Query.Messages.MinID != "" && m.id < params.Query.Messages.MinID {
				continue
			}
			if params.Query.Messages.MaxID != "" && m.id > params.Query.Messages.MaxID {
				continue
			}

			response.Messages = append(response.Messages, webhook.Message{
				ID:          m.id,
				UserName:    m.userName,
				PayloadJSON: m.payload,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Print(err)
	}
}
