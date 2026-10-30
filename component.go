package thsdk

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// componentFileNames maps GOOS to the runtime component's file name.
//
// The Python SDK keeps the component under thsdk/libs/<os>/<arch>/; this port
// mirrors that layout so a copied libs/ tree just works.
var componentFileNames = map[string]string{
	"windows": "hq.dll",
	"darwin":  "hq.dylib",
	"linux":   "hq.so",
}

// EnvRuntimePath overrides the component location, like the Python SDK's
// THSDK_RUNTIME_PATH.
const EnvRuntimePath = "THSDK_RUNTIME_PATH"

// resolveComponent finds the native runtime component.
//
// Order: explicit option, THSDK_RUNTIME_PATH, the file next to the executable,
// libs/<os>/<arch> next to the executable, libs/<os>/<arch> in the working
// directory, and finally whatever the OS loader can find by name.
func resolveComponent(explicit string) (string, error) {
	if explicit != "" {
		return expandHome(explicit)
	}
	if fromEnv := os.Getenv(EnvRuntimePath); fromEnv != "" {
		return expandHome(fromEnv)
	}

	name, ok := componentFileNames[runtime.GOOS]
	if !ok {
		return "", fmt.Errorf("thsdk: 当前平台 %s 暂不受支持", runtime.GOOS)
	}
	arch := runtime.GOARCH
	switch arch {
	case "amd64", "arm64":
	default:
		return "", fmt.Errorf("thsdk: 当前架构 %s 暂不受支持", arch)
	}

	var candidates []string
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(dir, name),
			filepath.Join(dir, "libs", familyForGOOS(), arch, name),
		)
	}
	cwd, err := os.Getwd()
	if err == nil {
		candidates = append(candidates,
			filepath.Join(cwd, "libs", familyForGOOS(), arch, name),
			filepath.Join(cwd, name),
		)
	}
	for _, candidate := range candidates {
		if info, statErr := os.Stat(candidate); statErr == nil && !info.IsDir() {
			return candidate, nil
		}
	}
	// Let the platform loader take over (PATH, LD_LIBRARY_PATH, rpath).
	return name, nil
}

func familyForGOOS() string {
	switch runtime.GOOS {
	case "darwin":
		return "darwin"
	case "linux":
		return "linux"
	default:
		return "windows"
	}
}

func expandHome(path string) (string, error) {
	if path != "~" && !strings.HasPrefix(path, "~/") && !strings.HasPrefix(path, "~\\") {
		return filepath.Abs(path)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Abs(filepath.Join(home, strings.TrimPrefix(strings.TrimPrefix(path, "~"), "/")))
}
