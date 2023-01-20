package main

import (
	"context"
	"log"
	"net/http"
	"time"

	agency "github.com/findy-network/findy-common-go/grpc/agency/v1"
	"github.com/findy-network/identity-hackathon-2023/go/agent"
	"github.com/gorilla/mux"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

var ourAgent *agent.Agent

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
	try.To1(response.Write([]byte("Go example")))
}

func verifyHandler(response http.ResponseWriter, r *http.Request) {
	defer err2.Catch(func(err error) {
		log.Println(err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	})
	try.To1(response.Write([]byte("Go example")))
}

func renderInvitation(header string, response http.ResponseWriter) (err error) {
	defer err2.Handle(&err)

	res := try.To1(ourAgent.Client.AgentClient.CreateInvitation(
		context.TODO(),
		&agency.InvitationBase{Label: ourAgent.UserName},
	))

	url := res.URL
	log.Printf("Created invitation\n %s\n", url)

}

func main() {
	defer err2.Catch(func(err error) {
		log.Fatal(err)
	})
	ourAgent = try.To1(agent.Init())
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
