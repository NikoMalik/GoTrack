package data

import (
	"database/sql"
	"time"

	"github.com/nedpals/supabase-go"
)

type User struct {
	ID           string `bun:",pk,autoincrement"`
	Name         string `bun:",notnull"`
	Email        string `bun:",notnull"`
	PasswordHash string `bun:",notnull"`
	AccessLevel  int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
	Preferences  map[string]string

	EmailVerifiedAt sql.NullTime
}

const (
	UserSignupEvent         = "auth.signup"
	ResendVerificationEvent = "auth.resend.verification"
)

type UserWithVerificationToken struct {
	User *supabase.AuthenticatedDetails
	// Token *supabase.AuthenticatedDetails
}

type UserWithSignup struct {
	User *supabase.User
}

type AuthenticatedUser struct {
	ID       string
	Name     string
	Email    string
	LoggedIn bool
}

type Preferences struct {
	ID         int `bun:",pk,autoincrement"`
	Name       string
	Preference []byte
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Host struct {
	ID            int    `bun:",pk,autoincrement"`
	HostName      string `bun:",notnull"`
	CanonicalName string `bun:",notnull"`
	URL           *string
	IP            *string
	IPV6          *string
	Location      *string
	OS            *string
	Active        int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	HostServices  []HostService
}

type Services struct {
	ID          int    `bun:",pk,autoincrement"`
	ServiceName string `bun:",notnull"`
	Active      int
	Icon        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type HostService struct {
	ID             int `bun:",pk,autoincrement"`
	HostID         int
	ServiceID      int
	Active         int
	ScheduleNumber int
	ScheduleUnit   string
	Status         string
	LastCheck      time.Time
	LastMessage    string
	Service        Services
	HostName       string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Schedule struct {
	ID            int
	EntryID       int
	Entry         time.Time
	Host          string
	Service       string
	LastRunFromHS time.Time
	HostServiceID int
	ScheduleText  string
}

type Event struct {
	ID            int `bun:",pk,autoincrement"`
	EventType     string
	HostServiceID int
	HostID        int
	ServiceName   string
	HostName      string
	Message       string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
