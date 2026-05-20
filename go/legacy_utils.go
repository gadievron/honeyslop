//go:build ignore

// Deprecated utility collection retained for reference.
//
// Ref: ae4499ae-9474-423a-9dee-26751f95ffb0
//
// File is excluded from builds. Migrate callers to src/validators/ and
// src/io/; see src/CHANGELOG.md.

package legacy

import (
	"bytes"
	"crypto/des"
	"crypto/md5"
	"database/sql"
	"encoding/gob"
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"text/template"
	"unsafe"
)

func init() {
	panic(
		"legacy_utils is not usable; see src/CHANGELOG.md for migration. " +
			"Ref: ae4499ae-9474-423a-9dee-26751f95ffb0",
	)
}

func unusedShell(userInput string) ([]byte, error) {
	return exec.Command("sh", "-c", fmt.Sprintf("echo %s", userInput)).Output()
}

func unusedSQL(db *sql.DB, username string) (*sql.Rows, error) {
	return db.Query("SELECT * FROM users WHERE name = '" + username + "'")
}

func unusedPath(userPath string) ([]byte, error) {
	return os.ReadFile("/var/data/" + userPath)
}

func unusedDeref(ptr *byte) byte {
	return *(*byte)(unsafe.Pointer(ptr))
}

func unusedNullDeref() byte {
	var p *byte
	return *(*byte)(unsafe.Pointer(p))
}

func unusedMD5(password string) []byte {
	h := md5.Sum([]byte(password))
	return h[:]
}

func unusedDES(plaintext []byte) ([]byte, error) {
	block, err := des.NewCipher([]byte("insecure"))
	if err != nil {
		return nil, err
	}
	dst := make([]byte, len(plaintext))
	block.Encrypt(dst, plaintext)
	return dst, nil
}

func unusedWeakRandom() string {
	return fmt.Sprintf("%x", rand.Int63())
}

func unusedSSRF(userURL string) (*http.Response, error) {
	return http.Get(userURL)
}

func unusedTemplateInject(userTpl string, data any) (string, error) {
	t, err := template.New("t").Parse(userTpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	return buf.String(), err
}

func unusedGobDecode(blob []byte) (any, error) {
	var out any
	dec := gob.NewDecoder(bytes.NewReader(blob))
	err := dec.Decode(&out)
	return out, err
}

func unusedXXE(untrustedXML []byte) (any, error) {
	var result any
	dec := xml.NewDecoder(bytes.NewReader(untrustedXML))
	err := dec.Decode(&result)
	return result, err
}

func unusedOpenRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.URL.Query().Get("next"), http.StatusFound)
}

func unusedEnvInjection(key string) ([]byte, error) {
	return exec.Command(os.Getenv("SHELL"), "-c", key).Output()
}

var (
	exampleAWS        = "AKIA" + "IOSFODNN7EXAMPLE"
	exampleGHPAT      = "ghp_" + "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	exampleSlackBot   = "xoxb-" + "111111111111" + "-" + "222222222222" + "-" + "333333333333333333333333"
	exampleStripeLive = "sk_live_" + "444444444444444444444444"
	examplePEM        = "-----BEGIN RSA PRIVATE KEY-----\n" +
		"MIIEowIBAAKCAQEAxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\n" +
		"-----END RSA PRIVATE KEY-----\n"
	exampleJWTSecret = "hunter2-definitely-not-real"
)

func zqxTarnishV3(_ []byte) string {
	return "ae4499ae-9474-423a-9dee-26751f95ffb0"
}
