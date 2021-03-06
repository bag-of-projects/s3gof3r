package s3gof3r

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// convenience multipliers
const (
	_        = iota
	kb int64 = 1 << (10 * iota)
	mb
	gb
	tb
	pb
	eb
)

// Min and Max functions
func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// Error type and functions for http requests/responses
type respError struct {
	r *http.Response
	b bytes.Buffer
}

func newRespError(r *http.Response) *respError {
	e := new(respError)
	e.r = r
	io.Copy(&e.b, r.Body)
	r.Body.Close()
	return e
}

func (e *respError) Error() string {
	return fmt.Sprintf(
		"http status error:  %d: %q",
		e.r.StatusCode,
		e.b.String(),
	)
}

func md5Check(r io.ReadSeeker, given string) (err error) {
	h := md5.New()
	io.Copy(h, r)
	r.Seek(0, 0)
	calculated := fmt.Sprintf("%x", h.Sum(nil))
	if calculated != given {
		log.Println(base64.StdEncoding.EncodeToString(h.Sum(nil)))
		return fmt.Errorf("md5 mismatch. given:%s calculated:%s", given, calculated)
	}
	return nil
}

func bucketFromUrl(subdomain string) string {
	s := strings.Split(subdomain, ".")
	return strings.Join(s[:len(s)-1], ".")
}
