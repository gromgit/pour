// +build linux,amd64

package config

const JSON_URL = "https://formulae.brew.sh/api/formula-linux.json"
const DEFAULT_PREFIX = "/home/linuxbrew/.linuxbrew"

func GetOS() (os string) {
os = "Linux64"
return
}
