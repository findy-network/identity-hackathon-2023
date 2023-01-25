package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	client "github.com/findy-network/findy-common-go/agency/client/async"
	agency "github.com/findy-network/findy-common-go/grpc/agency/v1"
	"github.com/findy-network/identity-hackathon-2023/go/agent"
	"github.com/gorilla/mux"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	qrcode "github.com/skip2/go-qrcode"
)

var (
	ourAgent      *agent.Agent
	agentListener *AgentListener = &AgentListener{}
)

type AgentListener struct {
	issue  sync.Map
	verify sync.Map
}

func (a *AgentListener) HandleNewConnection(id string) {
	defer err2.Catch(func(err error) {
		log.Println(err)
	})

	log.Printf(`New connection: %s`, id)

	pw := client.NewPairwise(ourAgent.Client.Conn, id)

	// If connection was for issuing, continue by issuing the "foobar" credential
	if _, ok := a.issue.Load(id); ok {
		a.issue.Delete(id)

		attributes := make([]*agency.Protocol_IssuingAttributes_Attribute, 1)
		attributes[0] = &agency.Protocol_IssuingAttributes_Attribute{
			Name:  "foo",
			Value: "bar",
		}

		log.Printf("Offer credential, conn id: %s, credDefID: %s, attrs: %v", id, ourAgent.CredDefID, attributes)

		res := try.To1(pw.IssueWithAttrs(
			context.TODO(),
			ourAgent.CredDefID,
			&agency.Protocol_IssuingAttributes{
				Attributes: attributes,
			}),
		)

		log.Printf("Credential offered: %s", res.GetID())

		// If connection was for verifying, continue by verifying the "foobar" credential
	} else {
		a.verify.Delete(id)

		attributes := make([]*agency.Protocol_Proof_Attribute, 1)
		attributes[0] = &agency.Protocol_Proof_Attribute{
			CredDefID: ourAgent.CredDefID,
			Name:      "foo",
		}

		log.Printf("Request proof, conn id: %s, attrs: %v", id, attributes)

		res := try.To1(pw.ReqProofWithAttrs(context.TODO(), &agency.Protocol_Proof{
			Attributes: attributes,
		}))

		log.Printf("Proof verified: %s", res.GetID())

	}
}

func (a *AgentListener) HandleNewCredential(id, connectionID string) {
	log.Printf(`Credential issued: %s`, id)

}

func (a *AgentListener) HandleNewProof(id, connectionID string) {
	log.Printf(`Proof verified: %s`, id)

}

// This function is called after proof is verified cryptographically.
// The application can execute its business logic and reject the proof
// if the attribute values are not valid.
func (a *AgentListener) HandleProofOnHold(id, connectionID string) {
	defer err2.Catch(func(err error) {
		log.Println(err)
	})

	log.Printf("Proof paused: %s", id)

	pw := client.NewPairwise(ourAgent.Client.Conn, id)

	// we have no special logic here - accept all received values
	res := try.To1(pw.Resume(
		context.TODO(),
		id,
		agency.Protocol_PRESENT_PROOF,
		agency.ProtocolState_ACK,
	))

	log.Printf("Proof continued: %s", res.GetID())

}

// Routes
func homeHandler(response http.ResponseWriter, r *http.Request) {
	defer err2.Catch(func(err error) {
		log.Println(err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	})
	try.To1(response.Write([]byte("Go example")))
}

// Show pairwise invitation. Once connection is established, issue credential.
func issueHandler(response http.ResponseWriter, r *http.Request) {
	defer err2.Catch(func(err error) {
		log.Println(err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	})
	id := try.To1(renderInvitation("Issue credential", response))
	agentListener.issue.Store(id, true)
}

// Show pairwise invitation. Once connection is established, verify credential.
func verifyHandler(response http.ResponseWriter, r *http.Request) {
	defer err2.Catch(func(err error) {
		log.Println(err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	})
	id := try.To1(renderInvitation("Verify proof", response))
	agentListener.verify.Store(id, true)
}

func renderInvitation(header string, response http.ResponseWriter) (invitationID string, err error) {
	defer err2.Handle(&err)

	res := try.To1(ourAgent.Client.AgentClient.CreateInvitation(
		context.TODO(),
		&agency.InvitationBase{Label: ourAgent.UserName},
	))

	var invitationMap map[string]any
	try.To(json.Unmarshal([]byte(res.GetJSON()), &invitationMap))

	url := res.URL
	log.Printf("Created invitation\n %s\n", url)

	png, err := qrcode.Encode(url, qrcode.Medium, 512)
	imgSrc := "data:image/png;base64," + base64.StdEncoding.EncodeToString([]byte(png))

	html := `<html>
    <h1>` + header + `</h1>
    <p>Read the QR code with the wallet application:</p>
    <img src="` + imgSrc + `"/>
    <p>or copy-paste the invitation:</p>
    <textarea onclick="this.focus();this.select()" readonly="readonly" rows="10" cols="60">` + url + `</textarea>
</html>`

	try.To1(response.Write([]byte(html)))

	return invitationMap["@id"].(string), nil
}

func main() {
	defer err2.Catch(func(err error) {
		log.Fatal(err)
	})
	agentName := os.Getenv("AGENCY_USER_NAME")
	if agentName == "" {
		agentName = "go-example"
	}
	ourAgent = try.To1(agent.Init(agentName, agent.SchemaInfo{
		Name:       "foobar",
		Attributes: []string{"foo"},
	}, agentListener))
	// TODO: renew token on expiry
	ourAgent.Login()

	router := mux.NewRouter()

	router.HandleFunc("/", homeHandler).Methods(http.MethodGet)
	router.HandleFunc("/issue", issueHandler).Methods(http.MethodGet)
	router.HandleFunc("/verify", verifyHandler).Methods(http.MethodGet)

	http.Handle("/", router)

	addr := ":3001"
	log.Printf("Starting server at %s", addr)

	server := http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	try.To(server.ListenAndServe())
}
