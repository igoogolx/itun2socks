package routes

import (
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func authRouter() http.Handler {
	r := chi.NewRouter()
	r.Post("/verify-password", verifyPassword)
	return r
}

type verifyPasswordRequest struct {
	Password string `json:"password"`
}

// verifyPassword checks the given password against the current OS user's credentials.
// On macOS: uses `dscl . -authonly <username> <password>`
// On Windows: uses PowerShell LogonUser via Win32 API
func verifyPassword(w http.ResponseWriter, r *http.Request) {
	var req verifyPasswordRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil || req.Password == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{"verified": false, "error": "password required"})
		return
	}

	verified := false
	var verifyErr string

	switch runtime.GOOS {
	case "darwin":
		verified, verifyErr = verifyMacOS(req.Password)
	case "windows":
		verified, verifyErr = verifyWindows(req.Password)
	default:
		verifyErr = "unsupported platform"
	}

	if verifyErr != "" {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{"verified": false, "error": verifyErr})
		return
	}

	if !verified {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, render.M{"verified": false})
		return
	}

	render.JSON(w, r, render.M{"verified": true})
}

func verifyMacOS(password string) (bool, string) {
	// Get current username
	u, err := user.Current()
	if err != nil {
		return false, "cannot determine current user"
	}
	username := u.Username
	// Strip domain prefix if present (e.g. "DOMAIN\user")
	if idx := strings.LastIndex(username, "\\"); idx >= 0 {
		username = username[idx+1:]
	}

	// dscl . -authonly <username> <password>
	// Returns exit code 0 if correct, non-zero if wrong
	cmd := exec.Command("/usr/bin/dscl", ".", "-authonly", username, password)
	// Suppress output
	cmd.Stdout = nil
	cmd.Stderr = nil
	err = cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			_ = exitErr
			return false, "" // wrong password, not an error
		}
		return false, err.Error()
	}
	return true, ""
}

func verifyWindows(password string) (bool, string) {
	username := os.Getenv("USERNAME")
	if username == "" {
		return false, "cannot determine current user"
	}

	script := `
Add-Type -TypeDefinition @"
using System;
using System.Runtime.InteropServices;
public class WinAuth {
    [DllImport("advapi32.dll", SetLastError=true)]
    public static extern bool LogonUser(string user, string domain, string pass,
        int logonType, int logonProvider, out IntPtr token);
    [DllImport("kernel32.dll")] public static extern bool CloseHandle(IntPtr h);
}
"@
$token = [IntPtr]::Zero
$ok = [WinAuth]::LogonUser("` + username + `", $env:USERDOMAIN, "` + password + `", 2, 0, [ref]$token)
if ($token -ne [IntPtr]::Zero) { [WinAuth]::CloseHandle($token) | Out-Null }
Write-Output $ok
`
	cmd := exec.Command("powershell.exe", "-noprofile", "-NonInteractive", "-command", script)
	out, err := cmd.Output()
	if err != nil {
		return false, err.Error()
	}
	return strings.TrimSpace(string(out)) == "True", ""
}
