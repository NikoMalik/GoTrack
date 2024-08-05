package handlers

import (
	"bytes"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/NikoMalik/GoTrack/data"
	"github.com/gofiber/fiber/v2"
)

func PollDomain(c *fiber.Ctx, domain string) (*data.DomainTrackingInfo, error) {
	currentTime := time.Now()
	resultch := make(chan data.DomainTrackingInfo)
	defer close(resultch)

	config := &tls.Config{InsecureSkipVerify: true}

	go func() {
		conn, err := tls.Dial("tcp", domain+":443", config)
		if err != nil {
			handleError(currentTime, err, resultch)

			return
		}

		defer conn.Close()

		state := conn.ConnectionState()
		cert := state.PeerCertificates[0]
		names := make([]string, 0, len(cert.DNSNames)+1)

		if IsDomainName(cert.Subject.CommonName) {
			names = append(names, cert.Subject.CommonName)
		}
		for _, name := range cert.DNSNames {
			if IsDomainName(name) {
				names = append(names, name)
			}
		}

		keyUsages := make([]string, len(cert.ExtKeyUsage))
		for i, usage := range cert.ExtKeyUsage {
			keyUsages[i] = extKeyUsageToString(usage)
		}

		host, port := splitHostPort(conn.RemoteAddr().String())
		var ports []string
		if port == "" {

			ports = []string{"80", "443"}
		} else {
			ports = []string{port}
		}

		// CIDR?
		if isCIDR(host) {
			// expand CIDR
			ips, err := expandCIDR(host)
			if err != nil {
				resultch <- data.DomainTrackingInfo{ServerIP: host, Error: err.Error()}
				return
			}

			for ip := range ips {
				for _, port := range ports {
					handleIPPort(ip, port, resultch)
				}
			}
		} else {

			for _, port := range ports {
				handleIPPort(host, port, resultch)
			}
		}

		resultch <- data.DomainTrackingInfo{
			ServerIP:      host,
			PublicKeyAlgo: cert.PublicKeyAlgorithm.String(),
			SignatureAlgo: cert.SignatureAlgorithm.String(),
			Signature:     sha1Hex(cert.Signature),
			KeyUsage:      keyUsageToString(cert.KeyUsage),
			ExtKeyUsages:  keyUsages,
			DNSNames:      strings.Join(names, ", "),
			Expires:       cert.NotAfter,
			PublicKey:     publicKeyFromCert(cert),
			Issuer:        cert.Issuer.Organization[0],
			EncodedPEM:    encodedPEMFromCert(cert),
			Port:          port,

			LastPollAt: currentTime,
			Latency:    int(time.Since(currentTime).Milliseconds()),
			Status:     getStatus(cert.NotAfter),
		}

	}()

	select {
	case <-c.Context().Done():
		return &data.DomainTrackingInfo{
			Error:      c.Context().Err().Error(),
			LastPollAt: currentTime,
			Status:     data.StatusUnresponsive,
		}, nil

	case result := <-resultch:
		return &result, nil

	}

}

func handleError(currentTime time.Time, err error, resultch chan data.DomainTrackingInfo) {
	info := data.DomainTrackingInfo{
		LastPollAt: currentTime,
		Error:      err.Error(),
		Latency:    int(time.Since(currentTime).Milliseconds()),
	}

	if IsVerificationError(err) {
		info.Status = data.StatusInvalid
	} else if IsConnectionRefused(err) {
		info.Status = data.StatusOffline
	}

	resultch <- info
}

func handleIPPort(ip, port string, resultch chan data.DomainTrackingInfo) {
	start := time.Now()
	address := net.JoinHostPort(ip, port)
	conn, err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		resultch <- data.DomainTrackingInfo{
			ServerIP:   ip,
			Port:       port,
			Error:      err.Error(),
			Status:     determineStatusFromError(err),
			LastPollAt: start,
			Latency:    int(time.Since(start).Milliseconds()),
		}
		return
	}
	defer conn.Close()

	resultch <- data.DomainTrackingInfo{
		ServerIP:   ip,
		Port:       port,
		Status:     data.StatusHealthy,
		LastPollAt: start,
		Latency:    int(time.Since(start).Milliseconds()),
	}
}

func determineStatusFromError(err error) string {
	if IsVerificationError(err) {
		return data.StatusInvalid
	} else if IsConnectionRefused(err) {
		return data.StatusOffline
	}
	return data.StatusUnresponsive
}

func IsVerificationError(err error) bool {
	return strings.Contains(err.Error(), "tls: failed to verify")
}

func IsConnectionRefused(err error) bool {
	return strings.Contains(err.Error(), "connect: connection refused")
}

func extKeyUsageToString(usage x509.ExtKeyUsage) string {
	switch usage {
	case x509.ExtKeyUsageAny:
		return "any"
	case x509.ExtKeyUsageServerAuth:
		return "server auth"
	case x509.ExtKeyUsageClientAuth:
		return "client auth"
	case x509.ExtKeyUsageCodeSigning:
		return "code signing"
	case x509.ExtKeyUsageEmailProtection:
		return "email protection"
	case x509.ExtKeyUsageIPSECEndSystem:
		return "IPS SEC system"
	case x509.ExtKeyUsageIPSECTunnel:
		return "IPS SEC tunnel"
	case x509.ExtKeyUsageIPSECUser:
		return "IPS SEC user"
	case x509.ExtKeyUsageTimeStamping:
		return "time stamping"
	case x509.ExtKeyUsageOCSPSigning:
		return "OCSP signing"
	case x509.ExtKeyUsageMicrosoftServerGatedCrypto:
		return "Microsoft server gated crypto"
	case x509.ExtKeyUsageNetscapeServerGatedCrypto:
		return "Netscape server gated crypto"
	case x509.ExtKeyUsageMicrosoftCommercialCodeSigning:
		return "Microsoft commercial code signing"
	case x509.ExtKeyUsageMicrosoftKernelCodeSigning:
		return "Microsoft kernel code signing"
	default:
		return ""
	}
}
func encodedPEMFromCert(cert *x509.Certificate) string {
	b := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}
	buf := new(bytes.Buffer)
	err := pem.Encode(buf, b)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func publicKeyFromCert(cert *x509.Certificate) string {
	pubKeyBytes, _ := x509.MarshalPKIXPublicKey(cert.PublicKey)
	pubKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}
	return sha1Hex(pubKeyBlock.Bytes)
}

func IsDomainName(s string) bool {

	l := len(s)
	if l == 0 || l > 254 || l == 254 && s[l-1] != '.' {
		return false
	}

	last := byte('.')
	nonNumeric := false // true once we've seen a letter or hyphen
	partlen := 0
	parts := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		default:
			return false
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z':
			nonNumeric = true
			partlen++
		case '0' <= c && c <= '9':
			// fine
			partlen++
		case c == '-':
			// Byte before dash cannot be dot.
			if last == '.' {
				return false
			}
			partlen++
			nonNumeric = true
		case c == '.':
			// Byte before dot cannot be dot, dash.
			if last == '.' || last == '-' {
				return false
			}
			if partlen > 63 || partlen == 0 {
				return false
			}
			partlen = 0
			parts++
		}
		last = c
	}
	if last == '-' || partlen > 63 {
		return false
	}

	return nonNumeric && (parts > 1 || parts > 0 && last != '.')
}

func keyUsageToString(usage x509.KeyUsage) string {
	switch usage {
	case x509.KeyUsageDigitalSignature:
		return "digital signature"
	case x509.KeyUsageContentCommitment:
		return "content commitment"
	case x509.KeyUsageKeyEncipherment:
		return "key encipherment"
	case x509.KeyUsageDataEncipherment:
		return "data encipherment"
	case x509.KeyUsageKeyAgreement:
		return "key agreement"
	case x509.KeyUsageCertSign:
		return "certificate sign"
	case x509.KeyUsageCRLSign:
		return "CRL sign"
	case x509.KeyUsageEncipherOnly:
		return "encipher only"
	case x509.KeyUsageDecipherOnly:
		return "decipher only"
	default:
		return "digital signature"
	}
}

func expandCIDR(CIDR string) (chan string, error) {
	// parse CIDR
	_, ipnet, err := net.ParseCIDR(CIDR)
	if err != nil {
		return nil, err
	}

	// general check for unsupported cases
	mOnes, mBits := ipnet.Mask.Size()
	if mBits == 128 && mOnes < 64 {
		return nil, fmt.Errorf("%s: IPv6 mask is too wide, use one from range /[64-128]", CIDR)
	}

	// create channel to deliver output
	outputChan := make(chan string)

	// switch branch to IPv4 / IPv6
	switch mBits {
	case 32: // IPv4:
		go func() {
			// convert to uint32, for convenient bitwise operation
			ip32 := binary.BigEndian.Uint32(ipnet.IP)
			mask32 := binary.BigEndian.Uint32(ipnet.Mask)

			// create buffer
			buf := new(bytes.Buffer)
			for mask := uint32(0); mask <= ^mask32; mask++ {
				// build IP as byte slice
				buf.Reset()
				err := binary.Write(buf, binary.BigEndian, ip32^mask)
				if err != nil {
					panic(err)
				}
				// yield stringified IP
				outputChan <- net.IP(buf.Bytes()).String()
			}
			close(outputChan)
		}()

	case 128: // IPv6
		go func() {
			// convert lower halves to uint64, for convenient bitwise operation
			ip64 := binary.BigEndian.Uint64(ipnet.IP[8:])
			mask64 := binary.BigEndian.Uint64(ipnet.Mask[8:])

			buf := new(bytes.Buffer)

			// write portion of IP that will not change during expansion
			buf.Write(ipnet.IP[:8])
			for mask := uint64(0); mask <= ^mask64; mask++ {
				// build IP as byte slice
				buf.Truncate(8)
				err := binary.Write(buf, binary.BigEndian, ip64^mask)
				if err != nil {
					panic(err)
				}
				// yield stringified IP
				outputChan <- net.IP(buf.Bytes()).String()
			}
			close(outputChan)
		}()
	}
	return outputChan, nil
}

func isCIDR(value string) bool {
	return strings.Contains(value, `/`)
}

var portRegexp, bracketRegexp *regexp.Regexp

func init() {
	portRegexp = regexp.MustCompile(`^(.*?)(:(\d+))?$`)
	bracketRegexp = regexp.MustCompile(`^\[.*\]$`)
}

func splitHostPort(addr string) (host, port string) {
	// split host and port
	portMatch := portRegexp.FindStringSubmatch(addr)
	host = portMatch[1]
	port = portMatch[3]
	isIPv6 := strings.Contains(host, `:`)

	// skip further checks for bracketed IPv6
	if isIPv6 && bracketRegexp.MatchString(host) {
		host = strings.TrimPrefix(host, `[`)
		host = strings.TrimSuffix(host, `]`)
		return
	}

	// no port found, skip futher checks
	if port == "" {
		return
	}

	// skip futher checks for CIDR
	if isCIDR(host) {
		return
	}

	// check ambiguous cases for IPv6
	if isIPv6 {
		// if port is longer than 4 digits -> it is truly a port
		if len(port) > 4 {
			return
		}

		// cancel port if whole thing parses as valid IPv6
		hostPort := fmt.Sprintf(`%s:%s`, host, port)
		if net.ParseIP(hostPort) != nil {
			host, port = hostPort, ``
			return
		}
	}
	return
}

var loomingTreshold = time.Hour * 24 * 7 * 2 // 2 weeks
func getStatus(expires time.Time) string {
	if expires.Before(time.Now()) {
		return "expired"
	}
	timeLeft := time.Until(expires)
	if timeLeft < loomingTreshold {
		return "expires"
	}
	return "healthy"
}

func sha1Hex(b []byte) string {
	sha1Hash := sha1.Sum(b)
	return hex.EncodeToString(sha1Hash[:])
}
