package resume

type Link struct {
	Icon string `yaml:"icon"`
	URL  string `yaml:"url"`
}

type Header struct {
	Name     string `yaml:"name"`
	Summary  string `yaml:"summary"`
	Location string `yaml:"location"`
	Links    []Link `yaml:"links"`
}

type SkillAreas struct {
	Area   string   `yaml:"area"`
	Skills []string `yaml:"skills"`
}

type Experience struct {
	Company    string   `yaml:"company"`
	Title      string   `yaml:"title"`
	StartDate  string   `yaml:"start_date"`
	EndDate    string   `yaml:"end_date"`
	Summary    string   `yaml:"summary"`
	Highlights []string `yaml:"highlights"`
}

type Education struct {
	Institution  string   `yaml:"institution"`
	Achievements []string `yaml:"achievements"`
}

type Certification struct {
	Authority string `yaml:"authority"`
	Name      string `yaml:"name"`
	Link      Link   `yaml:"link"`
	Issued    string `yaml:"issued"`
}

type Resume struct {
	Header         Header          `yaml:"header"`
	Skills         []SkillAreas    `yaml:"skills"`
	Experience     []Experience    `yaml:"experience"`
	Education      []Education     `yaml:"education"`
	Certifications []Certification `yaml:"certifications"`
}
