package main

import (
	"os"
	"text/template"
	"time"

	"gopkg.in/yaml.v3"
)

type YearMonth time.Time

type Person struct {
	Name         string        `yaml:"name"`
	Phone        string        `yaml:"phone"`
	Email        string        `yaml:"email"`
	PhotoURL     string        `yaml:"photoUrl"`
	Subtitle     string        `yaml:"subtitle"`
	Description  string        `yaml:"description"`
	Website      string        `yaml:"website"`
	Github       string        `yaml:"github"`
	Location     string        `yaml:"location"`
	LocationURL  string        `yaml:"locationUrl"`
	KeyLanguages string        `yaml:"keyLanguages"`
	KeySystems   string        `yaml:"keySystems"`
	Experience   []Company     `yaml:"experience"`
	Publications []Publication `yaml:"publications"`
	Projects     []Project     `yaml:"projects"`
}

type Company struct {
	CompanyName     string `yaml:"company"`
	Website         string `yaml:"website"`
	LogoURL         string `yaml:"logoUrl"`
	Description     string `yaml:"description"`
	Accomplishments []Item `yaml:"accomplishments"`
	Jobs            []Job  `yaml:"jobs"`
}

type Item struct {
	Thing string `yaml:"item"`
}

type Job struct {
	JobTitle string     `yaml:"role"`
	Start    *YearMonth `yaml:"start"`
	End      *YearMonth `yaml:"end,omitempty"`
}

type Publication struct {
	Title       string `yaml:"title"`
	URL         string `yaml:"url"`
	Description string `yaml:"description"`
	Image       string `yaml:"image"`
}

type Project struct {
	Title       string `yaml:"title"`
	URL         string `yaml:"url"`
	Description string `yaml:"description"`
}

func main() {
	yamlData, err := os.ReadFile("resume.yaml")
	if err != nil {
		panic(err)
	}

	var person Person
	err = yaml.Unmarshal(yamlData, &person)
	if err != nil {
		panic(err)
	}
	RenderTemplate(person, "static/template.html", "resume.html")
	RenderTemplate(person, "static/template.txt", "resume.txt")

	//spew.Dump(person)
}

func RenderTemplate(p Person, path, outputFile string) {
	// Define the Go template
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	tmpl, err := template.New("person").Funcs(template.FuncMap{"df": DateFormat}).Parse(string(data))
	if err != nil {
		panic(err)
	}
	f, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = tmpl.Execute(f, p)
	if err != nil {
		panic(err)
	}
}

// UnmarshalYAML takes the $MONTH-$YEAR value syntax for a start or end date
// and parses it into a proper time.Time.
func (ym *YearMonth) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var dateStr string
	if err := unmarshal(&dateStr); err != nil {
		return err
	}

	date, err := time.Parse("2006-01", dateStr)
	if err != nil {
		return err
	}

	*ym = YearMonth(date)
	return nil
}

func DateFormat(in YearMonth) string {
	t := time.Time(in)
	return t.Format("January 2006")
}
