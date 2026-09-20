package kampiun.shared

// SportFormat dari backend (tetap sinkron dgn domain/core.go).
data class SportFormat(
    val sport: String,
    val unit: String,
    val targetPoint: Int,
    val bestOf: Int,
    val usesScore: Boolean,
)

// MatchScore hasil evaluasi skor (identik dgn service/score.go MatchScore).
data class MatchScore(
    val homeUnits: List<Int>,
    val awayUnits: List<Int>,
    val homeWon: Int,
    val awayWon: Int,
    val finished: Boolean,
) {
    val isDraw: Boolean get() = homeWon == awayWon
}

// ScoreEngine = logika murni tap-skore, dipakai UI & (saat offline) disimpan lokal.
object ScoreEngine {

    // Terapkan satu aksi tap; hasilkan MatchScore berikutnya.
    fun apply(
        format: SportFormat,
        current: ScoreEvent,
        side: String,
        delta: Int,
    ): ScoreEvent {
        require(side == "home" || side == "away") { "side harus home/away" }
        val scores = current.scores.toMutableMap()
        scores[side] = (scores[side] ?: 0) + delta
        if (scores[side]!! < 0) scores[side] = 0

        // set tercapai? (hanya mode usesScore + targetPoint)
        val home = scores["home"] ?: 0
        val away = scores["away"] ?: 0
        var newWon = current.unitsWon
        if (format.usesScore && format.targetPoint > 0 && home >= format.targetPoint && home != away) {
            newWon = addWin(format, current.unitsWon, "home")
            scores["home"] = 0
            scores["away"] = 0
        } else if (format.usesScore && format.targetPoint > 0 && away >= format.targetPoint && home != away) {
            newWon = addWin(format, current.unitsWon, "away")
            scores["home"] = 0
            scores["away"] = 0
        }
        return EventFactory.next(current, scores, newWon)
    }

    private fun addWin(format: SportFormat, won: Map<String, Int>, side: String): Map<String, Int> {
        val needed = format.bestOf / 2 + 1
        val out = won.toMutableMap()
        out[side] = (out[side] ?: 0) + 1
        // sisa tak perlu; finished dicek dari MapWon
        return out
    }

    fun isFinished(format: SportFormat, won: Map<String, Int>): Boolean {
        if (format.bestOf <= 0) return false
        val needed = format.bestOf / 2 + 1
        return (won["home"] ?: 0) >= needed || (won["away"] ?: 0) >= needed
    }
}

// ScoreEvent = state mutable sederhana per match (dipakai UI offline-first).
data class ScoreEvent(
    val matchId: String,
    val scores: Map<String, Int>, // home/away -> angka saat ini di set aktif
    val unitsWon: Map<String, Int>, // home/away -> jumlah set menang
) {
    fun toMatchScore(format: SportFormat): MatchScore {
        val needed = format.bestOf / 2 + 1
        val h = unitsWon["home"] ?: 0
        val a = unitsWon["away"] ?: 0
        // catat skor set aktif sebagai unit terakhir
        val homeUnits = listOf(scores["home"] ?: 0)
        val awayUnits = listOf(scores["away"] ?: 0)
        return MatchScore(
            homeUnits = homeUnits,
            awayUnits = awayUnits,
            homeWon = h,
            awayWon = a,
            finished = h >= needed || a >= needed,
        )
    }
}

object EventFactory {
    fun initial(matchId: String) = ScoreEvent(matchId, emptyMap(), emptyMap())

    fun next(cur: ScoreEvent, scores: Map<String, Int>, won: Map<String, Int>) =
        ScoreEvent(cur.matchId, scores, won)
}
