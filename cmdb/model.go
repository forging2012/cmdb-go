package cmdb

type Relationship struct {
	SubjectType      string `json:"subjectType"`
	SubjectID        string `json:"subjectID"`
	RelationshipType string `json:"relationshipType"`
	ObjectType       string `json:"objectType"`
	ObjectID         string `json:"objectID"`
}



type Link struct {
	Name string `yaml:"name"`
	Link string `yaml:"link"`
}

type Contact struct {
	Entries []ContactEntry `json:"contact" yaml:"contact"`
}

type ContactEntry struct {
	DataItemID  string `json:"dataItemID" yaml:"dataItemID"`
	ContactPref string `json:"contactPref"`
	ContactType string `json:"contactType"`
	Email       string `json:"email"`
	Name        string `json:"name" yaml:"name"`
	Phone       string `json:"phone"`
	Programme   string `json:"programme"`
	Slack       string `json:"slack"`
	SupportRota string `json:"supportRota"`
}

type SystemAttributes struct {
	SystemCode          string `json:"systemCode"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	ServiceTier         string `json:"serviceTier"`
	LifecycleStage      string `json:"lifecycleStage"`
	ArchitectureDiagram string `json:"architectureDiagram"`
	MoreInformation     string `json:"moreInformation"`
	Troubleshooting     string `json:"troubleshooting"`
	Monitoring          string `json:"monitoring"`
	GitRepository       string `json:"gitRepository"`
	HostPlatform        string `json:"hostPlatform"`
}
