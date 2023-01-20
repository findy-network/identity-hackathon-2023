import express, { Express, Request, Response } from 'express';
import { createAcator, openGRPCConnection, agencyv1 } from '@findy-network/findy-common-ts'
import QRCode from 'qrcode'

import prepareIssuer from './prepare';


const app: Express = express();
const port = process.env.PORT || 3001;
const invitations: { issue: string[], verify: string[] } = {
    issue: [],
    verify: []
}
const userName = process.env.AGENCY_USER_NAME || 'ts-example'

const setupFindyAgency = async () => {
    const acatorProps = {
        authUrl: process.env.AGENCY_AUTH_URL || 'http://localhost:8088',
        authOrigin: process.env.AGENCY_AUTH_ORIGIN || 'http://localhost:3000',
        userName,
        key: process.env.AGENCY_KEY || '15308490f1e4026284594dd08d31291bc8ef2aeac730d0daf6ff87bb92d4336c',
    }
    const authenticator = createAcator(acatorProps)

    const grpcProps = {
        serverAddress: process.env.AGENCY_API_SERVER_ADDRESS || 'localhost',
        serverPort: parseInt(process.env.AGENCY_API_SERVER_PORT || '50052', 10),
        // NOTE: make sure cert path is defined when using localhost and self-issued certificate.
        // e.g. ../tools/local-env/cert
        certPath: process.env.AGENCY_API_SERVER_CERT_PATH || '',
    }

    // Authenticate and open GRPC connection to agency
    return openGRPCConnection(grpcProps, authenticator)
}

const runApp = async () => {
    const agencyConnection = await setupFindyAgency()
    const { createAgentClient, createProtocolClient } = agencyConnection
    const agentClient = await createAgentClient()
    const protocolClient = await createProtocolClient()
    // Credential definition is created on server startup.
    // We need it to be able to issue credentials.
    const credDefId = await prepareIssuer(agentClient, userName)

    // Listening callback handles agent events
    await agentClient.startListeningWithHandler(
        {
            // New connection is established
            DIDExchangeDone: async (info) => {
                console.log(`New connection: ${info.connectionId}`)

                // If connection was for issuing, continue by issuing the "foobar" credential
                if (invitations.issue.includes(info.connectionId)) {
                    const attributes = new agencyv1.Protocol.IssuingAttributes()
                    const attr = new agencyv1.Protocol.IssuingAttributes.Attribute()
                    attr.setName("foo")
                    attr.setValue("bar")
                    attributes.addAttributes(attr)

                    const credential = new agencyv1.Protocol.IssueCredentialMsg()
                    credential.setCredDefid(credDefId)
                    credential.setAttributes(attributes)

                    await protocolClient.sendCredentialOffer(info.connectionId, credential)

                    // If connection was for verifying, continue by verifying the "foobar" credential
                } else {
                    const attributes = new agencyv1.Protocol.Proof()
                    const attr = new agencyv1.Protocol.Proof.Attribute()
                    attr.setName("foo")
                    attr.setCredDefid(credDefId)
                    attributes.addAttributes(attr)

                    const proofRequest = new agencyv1.Protocol.PresentProofMsg()
                    proofRequest.setAttributes(attributes)

                    await protocolClient.sendProofRequest(info.connectionId, proofRequest)
                }
            },
            IssueCredentialDone: (info) => {
                console.log(`Credential issued: ${info.protocolId}`)
                invitations.issue = invitations.issue.filter(item => item !== info.connectionId)
            },

            // This function is called after proof is verified cryptographically.
            // The application can execute its business logic and reject the proof
            // if the attribute values are not valid.
            PresentProofPaused: async (info, presentProof) => {
                console.log(`Proof paused: ${info.protocolId}`)
                presentProof.getProof()?.getAttributesList().forEach((value, index) => {
                    console.log(`Proof attribute ${index} ${value.getName()}: ${value.getValue()}`)
                })
                const protocolID = new agencyv1.ProtocolID()
                protocolID.setId(info.protocolId)
                protocolID.setTypeid(agencyv1.Protocol.Type.PRESENT_PROOF)
                protocolID.setRole(agencyv1.Protocol.Role.RESUMER)
                const msg = new agencyv1.ProtocolState()
                msg.setProtocolid(protocolID)
                // we have no special logic here - accept all received values
                msg.setState(agencyv1.ProtocolState.State.ACK)
                await protocolClient.resume(msg)
            },
            PresentProofDone: (info) => {
                console.log(`Proof verified: ${info.protocolId}`)
                invitations.verify = invitations.verify.filter(item => item !== info.connectionId)
            },
        },
        {
            protocolClient,
            retryOnError: true,
        }
    )

    const renderInvitation = async (header: string, res: Response) => {
        const msg = new agencyv1.InvitationBase()
        msg.setLabel(userName)

        const invitation = await agentClient.createInvitation(msg)

        console.log(`Created invitation with Findy Agency: ${invitation.getUrl()}`)
        const qrData = await QRCode.toDataURL(invitation.getUrl())

        res.send(`<html>
    <h1>${header}</h1>
    <p>Read the QR code with the wallet application:</p>
    <img src="${qrData}"/>
    <p>or copy-paste the invitation:</p>
    <textarea onclick="this.focus();this.select()" readonly="readonly" rows="10" cols="60">${invitation.getUrl()}</textarea>
</html>`);

        return JSON.parse(invitation.getJson())["@id"]
    }

    // Show pairwise invitation. Once connection is established, verify credential.
    app.get('/verify', async (req: Request, res: Response) => {
        const id = await renderInvitation("Verify proof", res)
        invitations.verify = [...invitations.verify, id]
    });

    // Show pairwise invitation. Once connection is established, issue credential.
    app.get('/issue', async (req: Request, res: Response) => {
        const id = await renderInvitation("Issue credential", res)
        invitations.issue = [...invitations.issue, id]
    });

    app.get('/', (req: Request, res: Response) => {
        res.send('Typescript example');
    });

    app.listen(port, async () => {
        console.log(`⚡️[server]: Server is running at http://localhost:${port}`);
    });
}

runApp()