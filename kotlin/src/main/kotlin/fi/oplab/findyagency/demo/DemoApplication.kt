package fi.oplab.findyagency.demo

import com.google.zxing.BarcodeFormat
import com.google.zxing.qrcode.QRCodeWriter
import java.awt.Color
import java.awt.image.BufferedImage
import java.io.ByteArrayOutputStream
import java.util.Base64
import java.util.Collections
import javax.imageio.ImageIO
import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.launch
import kotlinx.serialization.*
import kotlinx.serialization.json.*
import org.findy_network.findy_common_kt.*
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RestController

@SpringBootApplication class DemoApplication

fun main(args: Array<String>) {
  runApplication<DemoApplication>(*args)
}

@Serializable data class InvitationData(@SerialName("@id") val id: String)

@RestController
class AppController {

  var agent = Agent()
  var issue: MutableList<String> = Collections.synchronizedList(mutableListOf())
  var verify: MutableList<String> = Collections.synchronizedList(mutableListOf())

  @GetMapping("/") fun index(): String = "Kotlin example"
  @GetMapping("/issue") fun issue(): String = createInvitationPage("Issue credential", issue)
  @GetMapping("/verify") fun verify(): String = createInvitationPage("Verify credential", verify)

  init {
    GlobalScope.launch {
      agent.connection.agentClient.listen().collect {
        println("Received from Agency:\n$it")
        val status = it.notification
        when (status.typeID) {
          Notification.Type.STATUS_UPDATE -> {
            // info contains the protocol related information
            val info = agent.connection.protocolClient.status(status.protocolID)
            val getType =
                fun(): Protocol.Type =
                    if (info.state.state == ProtocolState.State.OK) status.protocolType
                    else Protocol.Type.NONE

            when (getType()) {
              // New connection established
              Protocol.Type.DIDEXCHANGE -> {
                println("New connection ${status.protocolID} established")

                // If connection was for issuing, continue by issuing the "foobar"
                // credential
                if (issue.contains(status.connectionID)) {
                  issue.remove(status.connectionID)

                  agent.connection.protocolClient.sendCredentialOffer(
                      status.connectionID,
                      mapOf("foo" to "bar"),
                      agent.credDefId
                  )

                  // If connection was for verifying, continue by verifying the
                  // "foobar" credential
                } else {
                  verify.remove(status.connectionID)

                  agent.connection.protocolClient.sendProofRequest(
                      status.connectionID,
                      listOf(ProofRequestAttribute("foo", agent.credDefId)),
                  )
                }
              }

              // Credential issued
              Protocol.Type.ISSUE_CREDENTIAL -> {
                println("Credential ${status.protocolID} issued")
              }

              // Verification ready
              Protocol.Type.PRESENT_PROOF -> {
                println("Proof ${status.protocolID} verified")
              }
              else -> println("no handler for protocol type: ${status.protocolType}")
            }
            // Proof on hold
          }
          Notification.Type.PROTOCOL_PAUSED -> {
            // the cryptographic proof is done, we don't care about the values, so
            // accept always
            agent.connection.protocolClient.resumeProofRequest(status.protocolID, true)
          }
          else -> println("no handler for notification type: ${status.typeID}")
        }
      }
    }
  }

  fun createInvitationPage(header: String, idList: MutableList<String>): String {
    val invitation: Invitation = agent.createInvitation()

    val data =
        Json { ignoreUnknownKeys = true }.decodeFromString<InvitationData>(invitation.getJSON())
    println("Created invitaion with id ${data.id}")
    idList.add(data.id)

    // Create QR code for invitation URL
    val content = invitation.url
    val writer = QRCodeWriter()
    val bitMatrix = writer.encode(content, BarcodeFormat.QR_CODE, 512, 512)
    val width = bitMatrix.width
    val height = bitMatrix.height
    val bitmap = BufferedImage(width, height, BufferedImage.TYPE_USHORT_565_RGB)
    for (x in 0 until width) {
      for (y in 0 until height) {
        bitmap.setRGB(x, y, if (bitMatrix.get(x, y)) Color.BLACK.getRGB() else Color.WHITE.getRGB())
      }
    }
    val out = ByteArrayOutputStream()
    ImageIO.write(bitmap, "PNG", out)
    val imgSrc = "data:image/png;base64," + Base64.getEncoder().encodeToString(out.toByteArray())

    // render simple HTML page
    return """<html>
    <h1>${header}</h1>
    <p>Read the QR code with the wallet application:</p>
    <img src="${imgSrc}"/>
    <p>or copy-paste the invitation:</p>
    <textarea onclick="this.focus();this.select()" readonly="readonly" rows="10" cols="60">${invitation.url}</textarea>
</html>"""
  }
}
