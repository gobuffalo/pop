package pop

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_PostgreSQL_ConnectionDetails_Values_Finalize(t *testing.T) {
	cd := &ConnectionDetails{
		Dialect:  "postgres",
		Database: "database",
		Host:     "host",
		Port:     "1234",
		User:     "user",
		Password: "pass#",
	}
	err := cd.Finalize()
	require.NoError(t, err)

	p := &postgresql{commonDialect: commonDialect{ConnectionDetails: cd}}

	require.Equal(t, "postgres://user:pass%23@host:1234/database?", p.URL())
}

func Test_PostgreSQL_Connection_String(t *testing.T) {
	url := "host=host port=1234 dbname=database user=user password=pass#"
	cd := &ConnectionDetails{
		Dialect: "postgres",
		URL:     url,
	}
	err := cd.Finalize()
	require.NoError(t, err)

	require.Equal(t, url, cd.URL)
	require.Equal(t, "postgres", cd.Dialect)
	require.Equal(t, "host", cd.Host)
	require.Equal(t, "pass#", cd.Password)
	require.Equal(t, "1234", cd.Port)
	require.Equal(t, "user", cd.User)
	require.Equal(t, "database", cd.Database)
}

func genPrivateKey(tb testing.TB, caKeyPath string) *rsa.PrivateKey {
	tb.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(tb, err)

	f, err := os.Create(caKeyPath)
	require.NoError(tb, err)
	require.NoError(tb, pem.Encode(f, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}))
	require.NoError(tb, f.Close())
	return key
}

func genCertificate(
	tb testing.TB,
	template, parent *x509.Certificate,
	pub *rsa.PublicKey,
	priv *rsa.PrivateKey,
	path string,
) {
	tb.Helper()

	caBytes, err := x509.CreateCertificate(rand.Reader, template, parent, pub, priv)
	require.NoError(tb, err)

	f, err := os.Create(path)
	require.NoError(tb, err)
	defer f.Close()
	require.NoError(tb, pem.Encode(f, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	}))
}

func setupCerts(tb testing.TB, caKeyPath, caCertPath, serverKeyPath, serverCertPath string) {
	tb.Helper()
	snLimit := new(big.Int).Lsh(big.NewInt(1), 128) // 128-bit serial number limit

	caKey := genPrivateKey(tb, caKeyPath)

	caKeyBytes := x509.MarshalPKCS1PublicKey(&caKey.PublicKey)
	caKeyHash := sha256.Sum256(caKeyBytes)
	snCA, err := rand.Int(rand.Reader, snLimit)
	require.NoError(tb, err)

	caCert := &x509.Certificate{
		SerialNumber:   snCA,
		Subject:        pkix.Name{CommonName: "Test CA"},
		SubjectKeyId:   caKeyHash[:],
		AuthorityKeyId: caKeyHash[:],
		NotBefore:      time.Now(),
		NotAfter:       time.Now().Add(365 * 24 * time.Hour),
		IsCA:           true,
	}

	genCertificate(tb, caCert, caCert, &caKey.PublicKey, caKey, caCertPath)

	serverKey := genPrivateKey(tb, serverKeyPath)
	serverKeyBytes := x509.MarshalPKCS1PublicKey(&serverKey.PublicKey)
	serverKeyHash := sha256.Sum256(serverKeyBytes)
	snServer, err := rand.Int(rand.Reader, snLimit)
	require.NoError(tb, err)

	serverCert := &x509.Certificate{
		SerialNumber:   snServer,
		Subject:        pkix.Name{CommonName: "Test DB"},
		SubjectKeyId:   serverKeyHash[:],
		AuthorityKeyId: caKeyHash[:],
		IPAddresses:    []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		DNSNames:       []string{"localhost"},
		NotBefore:      time.Now(),
		NotAfter:       time.Now().Add(365 * 24 * time.Hour),
	}
	genCertificate(tb, serverCert, caCert, &serverKey.PublicKey, caKey, serverCertPath)
}

func Test_PostgreSQL_Connection_String_Options(t *testing.T) {
	tempDir := t.TempDir()
	caKeyPath := filepath.Join(tempDir, "ca.key")
	caCertPath := filepath.Join(tempDir, "ca.crt")
	serverKeyPath := filepath.Join(tempDir, "server.key")
	serverCertPath := filepath.Join(tempDir, "server.crt")
	setupCerts(t, caKeyPath, caCertPath, serverKeyPath, serverCertPath)

	url := fmt.Sprintf(
		"host=host port=1234 dbname=database user=user password=pass# sslmode=disable "+
			"fallback_application_name=test_app connect_timeout=10 sslcert=%s sslkey=%s sslrootcert=%s",
		filepath.ToSlash(serverCertPath),
		filepath.ToSlash(serverKeyPath),
		filepath.ToSlash(caCertPath),
	)
	cd := &ConnectionDetails{
		Dialect: "postgres",
		URL:     url,
	}
	require.NoError(t, cd.Finalize())

	require.Equal(t, url, cd.URL)

	require.Equal(t, "disable", cd.Options["sslmode"])
	require.Equal(t, "test_app", cd.Options["fallback_application_name"])
}

func Test_PostgreSQL_Connection_String_Without_User(t *testing.T) {
	url := "dbname=database"
	cd := &ConnectionDetails{
		Dialect: "postgres",
		URL:     url,
	}
	err := cd.Finalize()
	require.NoError(t, err)

	uc := os.Getenv("PGUSER")
	if uc == "" {
		c, err := user.Current()
		if err == nil {
			uc = c.Username
		}
	}

	require.Equal(t, url, cd.URL)
	require.Equal(t, "postgres", cd.Dialect)

	var foundHost bool
	if slices.Contains([]string{
		"/var/run/postgresql",
		"/private/tmp",
		"/tmp",
		"localhost",
	}, cd.Host) {
		foundHost = true
	}
	require.True(t, foundHost, `Got host: "%s"`, cd.Host)

	require.Equal(t, os.Getenv("PGPASSWORD"), cd.Password)
	require.Equal(t, portPostgreSQL, cd.Port) // fallback
	require.Equal(t, uc, cd.User)
	require.Equal(t, "database", cd.Database)
}

func Test_PostgreSQL_Connection_String_Failure(t *testing.T) {
	url := "abc"
	cd := &ConnectionDetails{
		Dialect: "postgres",
		URL:     url,
	}
	err := cd.Finalize()
	require.Error(t, err)
	require.Equal(t, "postgres", cd.Dialect)
}

func Test_PostgreSQL_Quotable(t *testing.T) {
	p := postgresql{}

	require.Equal(t, `"table_name"`, p.Quote("table_name"))
	require.Equal(t, `"schema"."table_name"`, p.Quote("schema.table_name"))
	require.Equal(t, `"schema"."table name"`, p.Quote(`"schema"."table name"`))
}
