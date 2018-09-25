package database

import (
	"errors"
	"fmt"

	gorm "github.com/jinzhu/gorm"
	"github.com/namsral/flag"
)

// Settings used to connect to the backing Postgres database
type Settings struct {
	Host        string
	Port        int
	User        string
	Password    string
	Name        string
	SSLMode     SSLMode
	SSLRootCert string
	SSLCert     string
	SSLKey      string
}

// DefaultPort is the default Postgres listen port
const DefaultPort = 5432

// GetSettings returns application configuration derived from command line
// options and environment variables.
func GetSettings() Settings {

	var s Settings
	var sslMode string

	flag.StringVar(&s.Name, "db-name", "", "DB name")
	flag.StringVar(&s.User, "db-user", "postgres", "DB user")
	flag.StringVar(&s.Password, "db-password", "", "DB password")
	flag.StringVar(&s.Host, "db-host", "127.0.0.1", "DB host address")
	flag.IntVar(&s.Port, "db-port", DefaultPort, "DB port")
	flag.StringVar(&sslMode, "db-ssl-mode", "disable", "DB SSL mode")
	flag.StringVar(&s.SSLRootCert, "db-ssl-root-cert", "", "DB SSL root certificate")
	flag.StringVar(&s.SSLCert, "db-ssl-cert", "", "DB SSL client certificate")
	flag.StringVar(&s.SSLKey, "db-ssl-key", "", "DB SSL client key")
	flag.Parse()

	s.SSLMode = SSLMode(sslMode)
	switch s.SSLMode {
	case SSLModeDisabled:
	case SSLModeRequired:
	case SSLModeVerifyCA:
	default:
		panic(fmt.Sprintf("Unknown database SSL mode: %s", sslMode))
	}

	return s
}

// SSLMode defines SSL settings used to connect to Postgres
type SSLMode string

const (
	// SSLModeDisabled disables SSL
	SSLModeDisabled SSLMode = "disable"

	// SSLModeRequired makes SSL required
	SSLModeRequired SSLMode = "require"

	// SSLModeVerifyCA enables SSL with server and client certificates
	SSLModeVerifyCA SSLMode = "verify-ca"
)

// Connect to a Postgres database using the given settings
// and returns a *gorm.DB handle.
func Connect(s Settings) (*gorm.DB, error) {

	if s.Host == "" {
		return nil, errors.New("Must specify database host")
	}
	if s.User == "" {
		return nil, errors.New("Must specify database user")
	}

	// Default the database name to the username
	if s.Name == "" {
		s.Name = s.User
	}

	args := fmt.Sprintf("user=%s dbname=%s sslmode=%s host=%s",
		s.User, s.Name, s.SSLMode, s.Host)

	if s.Password != "" {
		args += fmt.Sprintf(" password=%s", s.Password)
	}
	if s.SSLRootCert != "" {
		args += fmt.Sprintf(" sslrootcert=%s", s.SSLRootCert)
	}
	if s.SSLCert != "" {
		args += fmt.Sprintf(" sslcert=%s", s.SSLCert)
	}
	if s.SSLKey != "" {
		args += fmt.Sprintf(" sslkey=%s", s.SSLKey)
	}

	db, err := gorm.Open("postgres", args)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %s", err.Error())
	}

	db.LogMode(false)
	return db, nil
}
