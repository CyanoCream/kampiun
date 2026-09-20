package kampiun.shared

import io.ktor.client.HttpClient
import io.ktor.client.call.body
import io.ktor.client.engine.HttpClientEngine
import io.ktor.client.plugins.HttpTimeout
import io.ktor.client.plugins.contentnegotiation.ContentNegotiation
import io.ktor.client.request.get
import io.ktor.client.request.header
import io.ktor.client.request.post
import io.ktor.client.request.setBody
import io.ktor.http.ContentType
import io.ktor.http.contentType
import io.ktor.serialization.kotlinx.json.json
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json

// DTO — sinkron dgn respons backend.
@Serializable
data class ScoreDto(
    val HomeUnits: List<Int>,
    val AwayUnits: List<Int>,
    val HomeWon: Int,
    val AwayWon: Int,
    val Finished: Boolean,
)

@Serializable
data class MatchDetailDto(
    val ID: String = "",
    val CompID: String = "",
)

@Serializable
data class AuthResp(val token: String)

// ApiClient = klien Ktor utk backend Kampiun. Base URL di-set per platform (Android lokal/iOS).
interface ApiClient {
    suspend fun login(email: String, password: String): String
    suspend fun fetchScore(matchId: String): ScoreDto
    suspend fun tap(matchId: String, side: String, delta: Int, token: String): ScoreDto
}

fun createApiClient(baseUrl: String, engine: HttpClientEngine? = null): ApiClient {
    val client: HttpClient
    if (engine != null) {
        client = HttpClient(engine) { configureApi() }
    } else {
        client = HttpClient { configureApi() }
    }
    return KtorApiClient(client, baseUrl.trimEnd('/'))
}

private fun io.ktor.client.HttpClientConfig<*>.configureApi() {
    install(ContentNegotiation) { json(Json { ignoreUnknownKeys = true }) }
    install(HttpTimeout) { requestTimeoutMillis = 15_000 }
}

class KtorApiClient(
    private val client: HttpClient,
    private val baseUrl: String,
) : ApiClient {

    override suspend fun login(email: String, password: String): String {
        val resp = client.post("$baseUrl/api/v1/auth/login") {
            contentType(ContentType.Application.Json)
            setBody(LoginRequest(email, password))
        }.body<AuthResp>()
        return resp.token
    }

    override suspend fun fetchScore(matchId: String): ScoreDto =
        client.get("$baseUrl/api/v1/matches/$matchId/score").body()

    override suspend fun tap(matchId: String, side: String, delta: Int, token: String): ScoreDto =
        client.post("$baseUrl/api/v1/matches/$matchId/tap") {
            contentType(ContentType.Application.Json)
            header("Authorization", "Bearer $token")
            setBody(TapRequest(side, delta))
        }.body()
}

@Serializable
private data class LoginRequest(val email: String, val password: String)

@Serializable
private data class TapRequest(val side: String, val delta: Int)
