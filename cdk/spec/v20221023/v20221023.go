package spec

type File struct {
	Version string  `field:"required" yaml:"version"`
	Project Project `field:"required" yaml:"project"`
}

type Project struct {
	Kind       Kind       `field:"required" yaml:"kind"`
	Name       string     `field:"required" yaml:"name"`
	Properties Properties `field:"required" yaml:"properties"`
}

type Kind string

const (
	Kind_App    Kind = "app"
	Kind_Domain Kind = "domain"
)

type Properties struct {
	App      App       `field:"optional" yaml:"app"`
	Api      Api       `field:"optional" yaml:"api"`
	Handlers []Handler `field:"optional" yaml:"handlers"`
	Domains  []Domain  `field:"optional" yaml:"domains"`
}

type App struct {
	Name string `field:"required" yaml:"name"`
}

type Domain struct {
	Name    string   `field:"required" yaml:"name"`
	Indexes []string `field:"optional" yaml:"indexes"`
}

type Api struct {
	DataSources []DataSource `field:"optional" yaml:"datasources"`
	Functions   []Function   `field:"optional" yaml:"functions"`
	Resolvers   []Resolver   `field:"optional" yaml:"resolvers"`
}

type DataSource struct {
	Name string `field:"required" yaml:"name"`
	Type DSType `field:"required" yaml:"type"`
}

type DSType string

const (
	DSType_Lambda DSType = "lambda"
	DSType_Table  DSType = "table"
	DSType_Mock   DSType = "mock"
)

type Function struct {
	Name       string `field:"required" yaml:"name"`
	DataSource string `field:"required" yaml:"datasource"`
}

type Resolver struct {
	Name      string   `field:"required" yaml:"name"`
	Operation string   `field:"required" yaml:"operation"`
	Functions []string `field:"required" yaml:"functions"`
}

type Handler struct {
	Name     string      `field:"required" yaml:"name"`
	Type     HandlerType `field:"required" yaml:"type"`
	Events   []string    `field:"optional" yaml:"events"`
	Commands []string    `field:"optional" yaml:"commands"`
}

type HandlerType string

const (
	HandlerType_Basic    HandlerType = "basic"
	HandlerType_Standard HandlerType = "standard"
	HandlerType_Advanced HandlerType = "advanced"
)
