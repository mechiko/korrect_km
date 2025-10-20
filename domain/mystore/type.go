package mystore

import (
	"bytes"
	"crypto/sha1"
	"crypto/x509"
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"go.uber.org/zap"
)

const layout = "2006-01-02"

const modError = "mystore"

func List(logger *zap.SugaredLogger) (out map[string]string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s panic %v", modError, r)
		}
	}()
	out = make(map[string]string)
	store, err := syscall.UTF16PtrFromString("MY")
	if err != nil {
		// logger.Errorf("%s %w", modError, syscall.GetLastError())
		return out, syscall.GetLastError()
	}
	storeHandle, err := syscall.CertOpenSystemStore(0, store)
	if err != nil {
		return out, syscall.GetLastError()
	}
	defer syscall.CertCloseStore(storeHandle, 0)

	var cert *syscall.CertContext
	for {
		cert, err = syscall.CertEnumCertificatesInStore(storeHandle, cert)
		if err != nil {
			return out, syscall.GetLastError()
		}
		if cert == nil {
			break
		}
		// Copy the buf, since ParseCertificate does not create its own copy.
		buf := (*[1 << 20]byte)(unsafe.Pointer(cert.EncodedCert))[:]
		buf2 := make([]byte, cert.Length)
		copy(buf2, buf)
		if c, err := x509.ParseCertificate(buf2); err == nil {
			out[fing(c)] = fmt.Sprintf("%s (%s) %s ", c.Subject.CommonName, c.Issuer.CommonName, c.NotAfter.Local().Format(layout))
		}
	}
	return out, nil
}

func fing(c *x509.Certificate) string {
	fingerprint := sha1.Sum(c.Raw)
	var buf bytes.Buffer
	for _, f := range fingerprint {
		fmt.Fprintf(&buf, "%02X", f)
	}
	return strings.ToLower(buf.String())
}
