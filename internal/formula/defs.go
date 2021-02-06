package formula

type Formula struct {
	Name              string   `json:"name"`
	FullName          string   `json:"full_name"`
	Oldname           string   `json:"oldname"`
	Aliases           []string `json:"aliases"`
	VersionedFormulae []string `json:"versioned_formulae"`
	Desc              string   `json:"desc"`
	Homepage          string   `json:"homepage"`
	Versions          struct {
		Stable string `json:"stable"`
		Bottle bool   `json:"bottle"`
	} `json:"versions"`
	Revision      int `json:"revision"`
	VersionScheme int `json:"version_scheme"`
	Bottle        struct {
		Stable struct {
			Rebuild int    `json:"rebuild"`
			Cellar  string `json:"cellar"`
			Prefix  string `json:"prefix"`
			RootURL string `json:"root_url"`
			URL     string `json:"-"`
			Sha256  string `json:"-"`
			Files   struct {
				Catalina struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"catalina,omitempty"`
				Mojave struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"mojave,omitempty"`
				HighSierra struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"high_sierra,omitempty"`
				Sierra struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"sierra,omitempty"`
				ElCapitan struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"el_capitan,omitempty"`
				Yosemite struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"yosemite,omitempty"`
				Mavericks struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"mavericks,omitempty"`
				Linux64 struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"x86_64_linux,omitempty"`
			} `json:"files"`
		} `json:"stable"`
	} `json:"bottle,omitempty"`
	KegOnly                 bool     `json:"keg_only"`
	BottleDisabled          bool     `json:"bottle_disabled"`
	Options                 []string `json:"options"`
	BuildDependencies       []string `json:"build_dependencies"`
	Dependencies            []string `json:"dependencies"`
	RecommendedDependencies []string `json:"recommended_dependencies"`
	OptionalDependencies    []string `json:"optional_dependencies"`
	UsesFromMacos           []string `json:"uses_from_macos"`
	Requirements            []string `json:"requirements"`
	ConflictsWith           []string `json:"conflicts_with"`
	Caveats                 string   `json:"caveats"`
	Status                  int      `json:"-"`
	InstallDir              string   `json:"-"`
	InstallTime             string   `json:"-"`
	Pinned                  bool     `json:"-"`
}

type Formulas map[string]Formula

var StatusMap = map[int]string{
	INSTALLED: " ✓",
	OUTDATED:  " ✗",
}

// Various formula-related enumerations
const (
	RUN = iota
	BUILD
	RECOMMENDED
	OPTIONAL
	INSTALLED
	OUTDATED
	MISSING
)
