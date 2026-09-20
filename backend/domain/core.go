package domain

import (
	"errors"
	"time"
)

// Role pengguna.
type Role string

const (
	RoleSuperAdmin     Role = "superadmin"
	RoleOrganizer      Role = "organizer"
	RoleReferee        Role = "referee" // wasit (input skor)
	RoleCustomer       Role = "customer" // user umum yang membeli cup/langganan
)

// User = akun platform.
type User struct {
	ID           string
	Email        string
	Name         string
	Phone        string
	PasswordHash string
	Role         Role
	IsSuspended  bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Organization = penyelenggara (komunitas/kampung/panitia).
type Organization struct {
	ID        string
	Name      string
	Slug      string
	OwnerID   string // User.ID pemilik
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Competition = sebuah turnamen/kompetisi di dalam organisasi.
type Competition struct {
	ID             string
	OrgID          string
	Name           string
	Slug           string
	Sport          string // template pertandingan: "volleyball","badminton","fishing","football",...
	Format         SportFormat
	IsPublic       bool // publik dicari/disearch, private via kode
	AccessCode     string // kode unik pendek utk viewer private
	Status         CompStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// SportFormat = aturan skor template per jenis pertandingan.
type SportFormat struct {
	Sport          string
	Unit           string // "set","game","berat_kg","point",...
	TargetPoint    int    // poin utk menang 1 set (0=kustom/berat)
	BestOf         int    // best-of berapa set/game (0=langsung)
	UsesScore      bool   // true=+1 tap; false=isi angka (mancing berat)
}

type CompStatus string

const (
	CompActive    CompStatus = "active"
	CompFinished  CompStatus = "finished"
	CompCancelled CompStatus = "cancelled"
)

// Participant = tim/regu/orang dalam sebuah kompetisi.
type Participant struct {
	ID       string
	CompID   string
	Name     string
	Seed     int // unggulan (seeding)
	Color    string // warna UI blok tim
	Meta     string // json opsional (anggota, dsb)
	CreatedAt time.Time
}

// Match = satu pertandingan antara dua peserta.
type Match struct {
	ID          string
	CompID      string
	Round       int
	Position    int    // urutan kotak di bracket
	HomeID      string // peserta A ("" bila bye/belum)
	AwayID      string // peserta B
	Status      MatchStatus
	WinnerID    string
	StartedAt   *time.Time
	FinishedAt  *time.Time
	IsBye       bool
}

type MatchStatus string

const (
	MatchPending   MatchStatus = "pending"
	MatchLive      MatchStatus = "live"
	MatchFinished  MatchStatus = "finished"
)

// ScoreLine = skor akhir per unit (per set/game) untuk satu sisi.
type ScoreLine struct {
	Side string // "home" | "away"
	Units []int // poin per set/game
}

// MatchResult = hasil akhir + garis skor.
type MatchResult struct {
	ID        string
	WinnerID  string
	HomeLines []int
	AwayLines []int
}

// ScoreEvent = satu ketukan tap-tap skor (log/audit & sync).
type ScoreEvent struct {
	ID        string
	MatchID   string
	Side      string // home/away
	Delta     int    // +1 (tap) atau -1 (undo), atau nilai (mancing)
	UnitIndex int    // indeks set/game aktif
	ByUserID  string
	CreatedAt time.Time
}

var (
	ErrNotFound   = errors.New("data tidak ditemukan")
	ErrInvalid    = errors.New("input tidak valid")
	ErrForbidden  = errors.New("tidak berhak")
	ErrDuplicate  = errors.New("sudah ada")
	ErrUnauthorized = ErrForbidden
)
