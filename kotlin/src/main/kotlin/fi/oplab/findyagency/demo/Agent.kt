package fi.oplab.findyagency.demo

import java.io.File
import java.util.*
import kotlinx.coroutines.runBlocking
import org.findy_network.findy_common_kt.*

class Agent {
  val userName = System.getenv("FCLI_USER") ?: "kotlin-example"
  public var connection: Connection
  var credDefId: String
  init {
    val apiPortStr = System.getenv("AGENCY_API_SERVER_PORT") ?: ""
    connection =
        Connection(
            authUrl = System.getenv("FCLI_URL") ?: "http://localhost:8088",
            authOrigin = System.getenv("FCLI_ORIGIN") ?: "http://localhost:3000",
            userName = userName,
            seed = "",
            key = System.getenv("FCLI_KEY")
                    ?: "15308490f1e4026284594dd08d31291bc8ef2aeac730d0daf6ff87bb92d4336c",
            server = System.getenv("AGENCY_API_SERVER") ?: "localhost",
            port = if (apiPortStr != "") Integer.parseInt(apiPortStr) else 50052,
            certFolderPath = System.getenv("FCLI_TLS_PATH")
        )
    credDefId = createCredentialDefinition()
    println("Credential definition ready ${credDefId}")
  }

  public fun createInvitation(): Invitation = runBlocking {
    connection.agentClient.createInvitation(label = userName)
  }

  private fun createCredentialDefinition(): String = runBlocking {
    var credDefId = ""
    try {
      credDefId = File("CRED_DEF_ID").readLines()[0]
    } catch (e: Exception) {}

    if (credDefId == "") {
      val schemaRes =
          connection.agentClient.createSchema(
              name = "foobar",
              attributes = listOf("foo"),
              version = "1.0"
          )
      do {
        var schemaCreated = false
        try {
          val schema = connection.agentClient.getSchema(id = schemaRes.id)
          println("Created schema ${schema.id}")
          schemaCreated = true
        } catch (e: Exception) {
          println("Waiting for the schema to be created...")
          Thread.sleep(1_000 * 1)
        }
      } while (!schemaCreated)

      println("Starting to create credential definition (may take a while)...")
      val credDef = connection.agentClient.createCredDef(schemaId = schemaRes.id, tag = userName)
      credDefId = credDef.id
      File("CRED_DEF_ID").writeText(credDefId)
    }
    credDefId
  }
}
