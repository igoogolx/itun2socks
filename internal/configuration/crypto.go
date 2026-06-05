package configuration

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const encPrefix = "enc:"

// getMachineKey derives a 32-byte AES key from a machine-specific identifier.
// On Windows it uses the MachineGuid registry value.
// On macOS it uses the hardware UUID from IOPlatformExpertDevice.
// Falls back to hostname on other platforms.
func getMachineKey() ([]byte, error) {
	var machineID string

	if runtime.GOOS == "windows" {
		out, err := exec.Command("reg", "query",
			`HKLM\SOFTWARE\Microsoft\Cryptography`, "/v", "MachineGuid").Output()
		if err == nil {
			for _, line := range strings.Split(string(out), "\n") {
				if strings.Contains(line, "MachineGuid") {
					parts := strings.Fields(line)
					if len(parts) >= 3 {
						machineID = strings.TrimSpace(parts[len(parts)-1])
						break
					}
				}
			}
		}
	} else if runtime.GOOS == "darwin" {
		// Use macOS hardware UUID which is stable across reboots and hostname changes
		out, err := exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice").Output()
		if err == nil {
			for _, line := range strings.Split(string(out), "\n") {
				if strings.Contains(line, "IOPlatformUUID") {
					parts := strings.Split(line, "=")
					if len(parts) >= 2 {
						machineID = strings.Trim(strings.TrimSpace(parts[1]), "\"")
						break
					}
				}
			}
		}
	}

	if machineID == "" {
		h, err := os.Hostname()
		if err != nil {
			machineID = "lux-default-key"
		} else {
			machineID = h
		}
	}

	salt := "lux_proxy_encryption_v1"
	hash := sha256.Sum256([]byte(salt + ":" + machineID))
	return hash[:], nil
}

// encryptPassword encrypts a plaintext password using AES-GCM.
// Returns "enc:<base64>" or the original string on error.
func encryptPassword(plaintext string) string {
	if plaintext == "" || strings.HasPrefix(plaintext, encPrefix) {
		return plaintext
	}

	key, err := getMachineKey()
	if err != nil {
		return plaintext
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return plaintext
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return plaintext
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return plaintext
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return encPrefix + base64.StdEncoding.EncodeToString(ciphertext)
}

// decryptPassword decrypts an "enc:<base64>" password.
// Returns the original string if it's not encrypted or on error.
func decryptPassword(encrypted string) string {
	if !strings.HasPrefix(encrypted, encPrefix) {
		return encrypted
	}

	key, err := getMachineKey()
	if err != nil {
		return encrypted
	}

	data, err := base64.StdEncoding.DecodeString(encrypted[len(encPrefix):])
	if err != nil {
		return encrypted
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return encrypted
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return encrypted
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return encrypted
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return encrypted
	}

	return string(plaintext)
}

// isEncrypted returns true if the value has the enc: prefix.
func isEncrypted(value string) bool {
	return strings.HasPrefix(value, encPrefix)
}

// passwordFields are the proxy map keys that hold sensitive credentials.
var passwordFields = []string{"password", "passwd", "auth_str"}

// encryptProxyPasswords returns a copy of the proxy map with passwords encrypted.
func encryptProxyPasswords(proxy map[string]any) map[string]any {
	result := make(map[string]any, len(proxy))
	for k, v := range proxy {
		result[k] = v
	}
	for _, field := range passwordFields {
		if val, ok := result[field]; ok {
			if str, ok := val.(string); ok && str != "" && !isEncrypted(str) {
				result[field] = encryptPassword(str)
			}
		}
	}
	return result
}

// decryptProxyPasswords returns a copy of the proxy map with passwords decrypted.
func decryptProxyPasswords(proxy map[string]any) map[string]any {
	result := make(map[string]any, len(proxy))
	for k, v := range proxy {
		result[k] = v
	}
	for _, field := range passwordFields {
		if val, ok := result[field]; ok {
			if str, ok := val.(string); ok && isEncrypted(str) {
				result[field] = decryptPassword(str)
			}
		}
	}
	return result
}

// stripProxyPasswords returns a copy of the proxy map with password fields removed.
func stripProxyPasswords(proxy map[string]any) map[string]any {
	result := make(map[string]any, len(proxy))
	for k, v := range proxy {
		result[k] = v
	}
	for _, field := range passwordFields {
		delete(result, field)
	}
	return result
}

// StripProxyPasswords is the exported version for use by other packages.
func StripProxyPasswords(proxy map[string]any) map[string]any {
	return stripProxyPasswords(proxy)
}

// LockProxyPassword sets the passwordLocked flag on a proxy.
// The password is kept intact (encrypted on disk) so the proxy still works,
// but the API will refuse to reveal it until unlocked.
func LockProxyPassword(proxyId string) error {
	config, err := Read()
	if err != nil {
		return err
	}

	found := false
	for i, proxy := range config.Proxy {
		if id, ok := proxy["id"].(string); ok && id == proxyId {
			config.Proxy[i]["passwordLocked"] = true
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("proxy not found: %s", proxyId)
	}

	return Write(config)
}

// ResetProxyPassword removes the lock AND wipes the stored password.
// The user must re-enter the password via the edit form afterward.
func ResetProxyPassword(proxyId string) error {
	config, err := Read()
	if err != nil {
		return err
	}

	found := false
	for i, proxy := range config.Proxy {
		if id, ok := proxy["id"].(string); ok && id == proxyId {
			config.Proxy[i]["passwordLocked"] = false
			for _, field := range passwordFields {
				config.Proxy[i][field] = ""
			}
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("proxy not found: %s", proxyId)
	}

	return Write(config)
}

// MigrateEncryptPasswords checks if any proxy passwords are stored in plaintext
// and encrypts them. Called once at startup to transparently upgrade old configs.
func MigrateEncryptPasswords() error {
	config, err := Read()
	if err != nil {
		return err
	}

	needsWrite := false
	for _, proxy := range config.Proxy {
		for _, field := range passwordFields {
			if val, ok := proxy[field]; ok {
				if str, ok := val.(string); ok && str != "" && !isEncrypted(str) {
					needsWrite = true
					break
				}
			}
		}
		if needsWrite {
			break
		}
	}

	if needsWrite {
		// Write() will encrypt all plaintext passwords
		return Write(config)
	}
	return nil
}
