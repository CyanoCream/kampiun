package kampiun.shared

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertTrue

class ScoreEngineTest {

    private val voli = SportFormat("volleyball", "set", 25, 3, usesScore = true)

    @Test
    fun tapHomeIncrementsPoint() {
        var ev = EventFactory.initial("m1")
        ev = ScoreEngine.apply(voli, ev, "home", 1)
        assertEquals(1, ev.scores["home"])
        assertEquals(0, ev.unitsWon["home"] ?: 0)
    }

    @Test
    fun reach25WinsSet() {
        var ev = EventFactory.initial("m1")
        repeat(25) { ev = ScoreEngine.apply(voli, ev, "home", 1) }
        assertEquals(1, ev.unitsWon["home"])
        assertEquals(0, ev.scores["home"])  // reset ke 0 setelah set menang
    }

    @Test
    fun twoSetsFinish() {
        var ev = EventFactory.initial("m1")
        repeat(25) { ev = ScoreEngine.apply(voli, ev, "home", 1) }
        repeat(25) { ev = ScoreEngine.apply(voli, ev, "home", 1) }
        val sc = ev.toMatchScore(voli)
        assertEquals(2, sc.homeWon)
        assertTrue(sc.finished)
    }

    @Test
    fun undoNegativeClampsToZero() {
        var ev = EventFactory.initial("m1")
        ev = ScoreEngine.apply(voli, ev, "away", -1)
        assertEquals(0, ev.scores["away"])
    }
}
