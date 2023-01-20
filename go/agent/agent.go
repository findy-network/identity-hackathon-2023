package agent

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/findy-network/findy-agent-auth/acator/authn"
	"github.com/findy-network/findy-common-go/agency/client"
	agency "github.com/findy-network/findy-common-go/grpc/agency/v1"
	"github.com/google/uuid"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"google.golang.org/grpc"
)

const WaitTime = time.Millisecond * 500
const MaxWaitTime = time.Minute * 3

type AgencyClient struct {
	Conn           client.Conn
	AgentClient    agency.AgentServiceClient
	ProtocolClient agency.ProtocolServiceClient
}

type Agent struct {
	JWT        string
	Client     *AgencyClient
	AgencyHost string
	AgencyPort int
	CredDefID  string
	UserName   string
}

var authnCmd = authn.Cmd{
	SubCmd:   "",
	UserName: "go-example",
	Url:      os.Getenv("AGENCY_AUTH_URL"),
	AAGUID:   "12c85a48-4baf-47bd-b51f-f192871a1511",
	Key:      os.Getenv("AGENCY_KEY"),
	Counter:  0,
	Token:    "",
	Origin:   os.Getenv("AGENCY_AUTH_ORIGIN"),
}

func Init() (agent *Agent, err error) {
	defer err2.Handle(&err)

	// use default values if no environment configuration
	if authnCmd.Url == "" {
		authnCmd.Url = "http://localhost:8088"
	}
	if authnCmd.Origin == "" {
		authnCmd.Origin = "http://localhost:3000"
	}
	if authnCmd.Key == "" {
		authnCmd.Key = "15308490f1e4026284594dd08d31291bc8ef2aeac730d0daf6ff87bb92d4336c"
	}
	log.Printf("Auth url %s, origin %s, user %s", authnCmd.Url, authnCmd.Origin, authnCmd.UserName)

	serverAddress := os.Getenv("AGENCY_API_SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = "localhost"
	}
	serverPort, _ := strconv.Atoi(os.Getenv("AGENCY_API_SERVER_PORT"))
	if serverPort == 0 {
		serverPort = 50052
	}
	log.Printf("API server url %s, port %d", serverAddress, serverPort)

	// login (or register and login) to agency
	agent = &Agent{
		UserName:   authnCmd.UserName,
		AgencyHost: serverAddress,
		AgencyPort: serverPort,
	}
	try.To(agent.Login())

	log.Println("Agent login succeeded")

	// set up API connection
	conf := client.BuildClientConnBase(
		os.Getenv("AGENCY_API_SERVER_CERT_PATH"),
		agent.AgencyHost,
		agent.AgencyPort,
		[]grpc.DialOption{},
	)

	conn := client.TryAuthOpen(agent.JWT, conf)
	agent.Client = &AgencyClient{
		Conn:           conn,
		AgentClient:    agency.NewAgentServiceClient(conn),
		ProtocolClient: agency.NewProtocolServiceClient(conn),
	}

	agent.CredDefID = try.To1(agent.createCredDef())

	// start listening to events
	ch := try.To1(agent.Client.Conn.ListenStatus(context.TODO(), &agency.ClientID{ID: uuid.New().String()}))
	go func() {
		for {
			chRes, ok := <-ch
			if !ok {
				panic("Listening failed")
			}
			notification := chRes.GetNotification()
			log.Printf("Received agent notification %v\n", notification)

		}
	}()

	return agent, err
}

func (a *Agent) register() (err error) {
	defer err2.Handle(&err)

	myCmd := authnCmd
	myCmd.SubCmd = "register"

	try.To(myCmd.Validate())
	try.To1(myCmd.Exec(os.Stdout))
	return
}

func (a *Agent) login() (err error) {
	defer err2.Handle(&err)

	myCmd := authnCmd
	myCmd.SubCmd = "login"

	try.To(myCmd.Validate())
	r := try.To1(myCmd.Exec(os.Stdout))

	// TODO: renew token on expiry
	a.JWT = r.Token
	return
}

func (a *Agent) Login() (err error) {
	defer err2.Handle(&err)

	// first try to login
	err = a.login()
	if err != nil {
		// if login fails, try to register and relogin
		try.To(a.register())
		try.To(a.login())
	}

	return
}

func (a *Agent) createCredDef() (credDefID string, err error) {
	defer err2.Handle(&err)

	schemaRes := try.To1(a.Client.AgentClient.CreateSchema(
		context.TODO(),
		&agency.SchemaCreate{
			Name:       "foobar",
			Version:    "1.0",
			Attributes: []string{"foo"},
		},
	))

	// wait for schema to be readable before creating cred def
	schemaGet := &agency.Schema{
		ID: schemaRes.ID,
	}
	schemaFound := false
	for !schemaFound {
		if _, err := a.Client.AgentClient.GetSchema(context.TODO(), schemaGet); err == nil {
			schemaFound = true
		} else {
			time.Sleep(1 * time.Second)
		}
	}

	log.Printf("Schema %s created successfully", schemaRes.ID)

	res := try.To1(a.Client.AgentClient.CreateCredDef(
		context.TODO(),
		&agency.CredDefCreate{
			SchemaID: schemaRes.ID,
			Tag:      authnCmd.UserName,
		},
	))

	log.Printf("Credential definition %s created successfully", res.ID)

	return res.GetID(), nil
}
