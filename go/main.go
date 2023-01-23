package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

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
	if _, ok := a.issue.Load(id); ok {
		a.issue.Delete(id)

		attributes := make([]*agency.Protocol_IssuingAttributes_Attribute, 1)
		attributes[0] = &agency.Protocol_IssuingAttributes_Attribute{
			Name:  "foo",
			Value: "bar",
		}

		log.Printf("Propose credential, conn id: %s, credDefID: %s, attrs: %v", id, ourAgent.CredDefID, attributes)

		protocol := &agency.Protocol{
			ConnectionID: id,
			TypeID:       agency.Protocol_ISSUE_CREDENTIAL,
			Role:         agency.Protocol_INITIATOR,
			StartMsg: &agency.Protocol_IssueCredential{
				IssueCredential: &agency.Protocol_IssueCredentialMsg{
					CredDefID: ourAgent.CredDefID,
					AttrFmt: &agency.Protocol_IssueCredentialMsg_Attributes{
						Attributes: &agency.Protocol_IssuingAttributes{
							Attributes: attributes,
						},
					},
				},
			},
		}
		try.To1(ourAgent.Client.Conn.DoStart(context.TODO(), protocol))

	} else {
		a.verify.Delete(id)

		attributes := make([]*agency.Protocol_Proof_Attribute, 1)
		attributes[0] = &agency.Protocol_Proof_Attribute{
			CredDefID: ourAgent.CredDefID,
			Name:      "foo",
		}

		log.Printf("Request proof, conn id: %s, attrs: %v", id, attributes)

		protocol := &agency.Protocol{
			ConnectionID: id,
			TypeID:       agency.Protocol_PRESENT_PROOF,
			Role:         agency.Protocol_INITIATOR,
			StartMsg: &agency.Protocol_PresentProof{
				PresentProof: &agency.Protocol_PresentProofMsg{
					AttrFmt: &agency.Protocol_PresentProofMsg_Attributes{
						Attributes: &agency.Protocol_Proof{
							Attributes: attributes,
						},
					},
				},
			},
		}
		try.To1(ourAgent.Client.Conn.DoStart(context.TODO(), protocol))

	}
}

func (a *AgentListener) HandleNewCredential(id, connectionID string) {
	log.Printf(`Credential issued: %s`, id)

}

func (a *AgentListener) HandleNewProof(id, connectionID string) {
	log.Printf(`Proof verified: %s`, id)

}

func (a *AgentListener) HandleProofOnHold(id, connectionID string) {
	defer err2.Catch(func(err error) {
		log.Println(err)
	})

	log.Printf(`Proof paused: %s`, id)

	state := &agency.ProtocolState{
		ProtocolID: &agency.ProtocolID{
			TypeID: agency.Protocol_PRESENT_PROOF,
			Role:   agency.Protocol_RESUMER,
			ID:     id,
		},
		State: agency.ProtocolState_ACK,
	}

	try.To1(ourAgent.Client.ProtocolClient.Resume(
		context.TODO(),
		state,
	))
}

// Routes
func homeHandler(response http.ResponseWriter, r *http.Request) {
	defer err2.Catch(func(err error) {
		log.Println(err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	})
	try.To1(response.Write([]byte("Go example")))
}

func issueHandler(response http.ResponseWriter, r *http.Request) {
	defer err2.Catch(func(err error) {
		log.Println(err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	})
	id := try.To1(renderInvitation("Issue credential", response))
	agentListener.issue.Store(id, true)
}

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

	url := res.URL
	log.Printf("Created invitation\n %s\n", url)

	png, err := qrcode.Encode(url, qrcode.Highest, 256)
	imgSrc := "data:image/png;base64," + base64.StdEncoding.EncodeToString([]byte(png))

	html := `<html>
    <h1>` + header + `</h1>
    <p>Read the QR code with the wallet application:</p>
    <img src="` + imgSrc + `"/>
    <p>or copy-paste the invitation:</p>
    <textarea onclick="this.focus();this.select()" readonly="readonly" rows="10" cols="60">` + url + `</textarea>
</html>`

	try.To1(response.Write([]byte(html)))

	var invitationMap map[string]any
	try.To(json.Unmarshal([]byte(res.GetJSON()), &invitationMap))
	return invitationMap["@id"].(string), nil
}

func main() {
	defer err2.Catch(func(err error) {
		log.Fatal(err)
	})
	ourAgent = try.To1(agent.Init("go-example", agent.SchemaInfo{
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
