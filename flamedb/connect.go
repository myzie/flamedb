package flamedb

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
	flag.StringVar(&s.User, "db-user", "", "DB user")
	flag.StringVar(&s.Password, "db-password", "", "DB password")
	flag.StringVar(&s.Host, "db-host", "", "DB host address")
	flag.IntVar(&s.Port, "db-port", DefaultPort, "DB port")
	flag.StringVar(&sslMode, "db-ssl-mode", "", "DB SSL mode")
	flag.StringVar(&s.SSLRootCert, "db-ssl-root-cert", "", "DB SSL root certificate")
	flag.StringVar(&s.SSLCert, "db-ssl-cert", "", "DB SSL client certificate")
	flag.StringVar(&s.SSLKey, "db-ssl-key", "", "DB SSL client key")
	flag.Parse()

	s.SSLMode = SSLMode(sslMode)

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
		s.Host = "127.0.0.1"
	}
	if s.Name == "" {
		return nil, errors.New("Must specify database name")
	}
	if s.User == "" {
		return nil, errors.New("Must specify database user")
	}
	if s.Password == "" {
		return nil, errors.New("Must specify database password")
	}

	args := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s",
		s.User, s.Password, s.Name, s.SSLMode, s.Host)

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
