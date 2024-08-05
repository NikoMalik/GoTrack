package data

import (
	"time"

	"github.com/nedpals/supabase-go"
)

const (
	defaultLimit    = 15
	UserSignupEvent = "auth.signup"
)

type UserWithVerificationToken struct {
	User    *supabase.AuthenticatedDetails
	Message any
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

type DomainTrackingInfo struct {
	ServerIP      string
	Issuer        string
	Port          string
	SignatureAlgo string
	PublicKeyAlgo string
	EncodedPEM    string
	PublicKey     string
	Signature     string
	DNSNames      string
	KeyUsage      string
	ExtKeyUsages  []string `bun:",array"`
	Expires       time.Time
	Status        string
	LastPollAt    time.Time
	Latency       int
	Error         string
}

type DomainTracking struct {
	ID         int64 `bun:"id,pk,autoincrement"`
	User       *supabase.User
	DomainName string

	DomainTrackingInfo
}

type Host struct {
	ID            int    `bun:"id,pk,autoincrement,notnull"`
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
