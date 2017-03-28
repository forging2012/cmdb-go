package cmdb

import (
	"bytes"
	"html/template"

	"io/ioutil"

	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/russross/blackfriday"
)

type System struct {
	SystemCode          string  `json:"systemCode" yaml:"systemCode"`
	Name                string  `json:"name" yaml:"name"`
	Description         string  `json:"description" yaml:"description"`
	ServiceTier         string  `json:"serviceTier" yaml:"serviceTier"`
	LifecycleStage      string  `json:"lifecycleStage" yaml:"lifecycleStage"`
	ArchitectureDiagram string  `json:"architectureDiagram" yaml:"architectureDiagram"`
	Troubleshooting     string  `json:"troubleshooting" yaml:"troubleshooting"`
	MoreInformation     []Link  `json:"moreInformation" yaml:"moreInformation"`
	Monitoring          []Link  `json:"monitoring" yaml:"monitoring"`
	GitRepository       string  `json:"gitRepository" yaml:"gitRepository"`
	HostPlatform        string  `json:"hostPlatform" yaml:"hostPlatform"`
	PrimaryContact      Contact `json:"-" yaml:"primaryContact"`
	SecondaryContact    Contact `json:"-" yaml:"secondaryContact"`
	Programme           Contact `json:"-" yaml:"programme"`
	ProductOwner        Contact `json:"-" yaml:"productOwner"`
	TechnicalLead       Contact `json:"-" yaml:"technicalLead"`
}

func (s *System) ImportMarkdown(mdFile string, parentSystem string) string {
	if len(mdFile) <= 5 || mdFile[0:5] != "MD://" {
		return mdFile
	}

	truePath := filepath.Join(filepath.Dir(parentSystem), mdFile[5:])

	logrus.Infof("Importing markdown from %s", truePath)

	input, err := ioutil.ReadFile(truePath)
	if err != nil {
		logrus.Error(err)
		return mdFile
	}
	output := blackfriday.MarkdownCommon(input)
	return string(output)
}

func (s *System) ToSystemAttributes() SystemAttributes {

	sa := SystemAttributes{}
	sa.SystemCode = s.SystemCode
	sa.Name = s.Name
	sa.Description = s.Description
	sa.ServiceTier = s.ServiceTier
	sa.LifecycleStage = s.LifecycleStage
	sa.ArchitectureDiagram = s.ArchitectureDiagram
	sa.Troubleshooting = s.Troubleshooting
	sa.MoreInformation = linksToString(s.MoreInformation)
	sa.Monitoring = linksToString(s.Monitoring)
	sa.GitRepository = s.GitRepository
	sa.HostPlatform = s.HostPlatform

	return sa
}

var linkTemplate = `<table>
<tr>
	<th>Name</th>
	<th>Link</th>
</tr>
{{range .}}
<tr>
	<td>{{.Name}}</td>
	<td>{{.Link}}</td>
</tr>
{{end}}
</table>`

func linksToString(links []Link) string {
	var buffer bytes.Buffer
	t := template.Must(template.New("").Parse(linkTemplate))
	t.Execute(&buffer, links)
	return buffer.String()
}
