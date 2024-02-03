package ip

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/paulmuenzner/powerplantmanager/config"
	"net"
	"net/http"
	"reflect"
	"strings"

	"github.com/seancfoley/ipaddress-go/ipaddr"
)

func RealIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remoteAddr := r.RemoteAddr
		if ip := r.Header.Get("X-Real-IP"); ip != "" {
			remoteAddr = ip
		} else if ips := r.Header.Get("X-Forwarded-For"); ips != "" {
			remoteAddr = strings.Split(ips, ",")[0]
		}
		r.Header.Set("X-Remote-IP", r.RemoteAddr)
		r.RemoteAddr = remoteAddr

		next.ServeHTTP(w, r)
	})
}

func ExtractIP(r *http.Request) (string, error) {
	// Get the client's IP address
	clientIP := r.RemoteAddr

	// Extract the IP address from the remote address if needed
	if strings.Contains(clientIP, ":") {
		ip, _, err := net.SplitHostPort(clientIP)
		if err != nil {
			return "", fmt.Errorf("error getting client IP: %w", err)
		}
		clientIP = ip
	}

	return clientIP, nil
}

// ////////////////////////////////////////////////////////
// Normalize / extend IPv6 addresses to enable comparisons
func NormalizeIP(ipStr string) (string, error) {
	// Implement your logic to normalize the IP address
	// This can involve converting to lowercase, removing whitespace, etc.
	parsedIP := net.ParseIP(ipStr)
	if parsedIP == nil {
		return "", fmt.Errorf("invalid IP address: %s", ipStr)
	}

	ipv4Regex := config.Regex.Ipv4
	ipv6Regex := config.Regex.Ipv6

	if ipv6Regex.MatchString(ipStr) {
		fullIP := ipaddr.NewIPAddressString(ipStr).GetAddress().ToFullString()
		if reflect.TypeOf(fullIP).Kind() == reflect.String {
			return fullIP, nil
		}

	} else if ipv4Regex.MatchString(ipStr) {
		return parsedIP.String(), nil
	}

	return "", errors.New("Cannot normalize IP address with 'NormalizeIP'. Provided argument: " + ipStr)
}

func FullIPv6(ip net.IP) string {
	dst := make([]byte, hex.EncodedLen(len(ip)))
	_ = hex.Encode(dst, ip)
	return string(dst[0:4]) + ":" +
		string(dst[4:8]) + ":" +
		string(dst[8:12]) + ":" +
		string(dst[12:16]) + ":" +
		string(dst[16:20]) + ":" +
		string(dst[20:24]) + ":" +
		string(dst[24:28]) + ":" +
		string(dst[28:])
}
