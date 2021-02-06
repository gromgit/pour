// +build linux,amd64

package config

const OS_NAME = "Linux"
const JSON_URL = "https://formulae.brew.sh/api/formula-linux.json"
const DEFAULT_PREFIX = "/home/linuxbrew/.linuxbrew"

func GetDeps() []string {
	return []string{"patchelf"}
}

func GetOS() (os string) {
	os = "Linux64"
	return
}
