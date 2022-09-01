package callecr

import (
		"os"
		"path/filepath"
		"runtime"
		"testing"
		"strings"
)

var (
		testConfigFilename = filepath.Join("testdata", "shared_config")
)

func TestCreateSession(t *testing.T) {
	restoreEnvFn := initSessionTestEnv()
	defer restoreEnvFn()

	os.Setenv("AWS_SDK_LOAD_CONFIG", "0")
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", testConfigFilename)

	sess, err := CreateSession("full_profile", "full_profile_region")

	if err != nil {
		t.Errorf("expect nil, %v", err)
	}

	if e, a := "full_profile_region", *sess.Config.Region; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

	creds, err := sess.Config.Credentials.Get()
	if err != nil {
		t.Errorf("expect nil, %v", err)
	}
	if e, a := "full_profile_akid", creds.AccessKeyID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "full_profile_secret", creds.SecretAccessKey; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if v := creds.SessionToken; len(v) != 0 {
		t.Errorf("expect empty, got %v", v)
	}
	if e, a := "SharedConfigCredentials", creds.ProviderName; !strings.Contains(a, e) {
		t.Errorf("expect %v, to be in %v", e, a)
	}

}

// support_functions:

func getEnvs(envs []string) map[string]string {
	extraEnvs := make(map[string]string)
	for _, env := range envs {
		if val, ok := os.LookupEnv(env); ok && len(val) > 0 {
			extraEnvs[env] = val
		}
	}
	return extraEnvs
}

func popEnv(env []string) {
	os.Clearenv()

	for _, e := range env {
		p := strings.SplitN(e, "=", 2)
		k, v := p[0], ""
		if len(p) > 1 {
			v = p[1]
		}
		os.Setenv(k, v)
	}
}

func stashEnv(envToKeep ...string) func() {
	if runtime.GOOS == "windows" {
		envToKeep = append(envToKeep, "ComSpec")
		envToKeep = append(envToKeep, "SYSTEM32")
		envToKeep = append(envToKeep, "SYSTEMROOT")
	}
	envToKeep = append(envToKeep, "PATH")
	extraEnv := getEnvs(envToKeep)
	originalEnv := os.Environ()
	os.Clearenv() // clear env
	for key, val := range extraEnv {
		os.Setenv(key, val)
	}
	return func() {
		popEnv(originalEnv)
	}
}

func initSessionTestEnv() (oldEnv func()) {
	oldEnv = stashEnv()
	os.Setenv("AWS_CONFIG_FILE", "file_not_exists")
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "file_not_exists")

	return oldEnv
}
