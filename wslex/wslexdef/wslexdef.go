package wslexdef

import "encoding/xml"

const (
	// DefaultDistribution ...
	DefaultDistribution = "DefaultDistribution"
	// BasePath ...
	BasePath         = "BasePath"
	DistributionName = "DistributionName"
	DefaultUID       = "DefaultUid"
	Flags            = "Flags"
	State            = "State"
	Version          = "Version"

	// AppxManifestXML comment para o go-lint ficar feliz
	AppxManifestXML = "AppxManifest.xml"

	// WslexLxssKey key
	WslexLxssKey = "Software\\Microsoft\\Windows\\CurrentVersion\\Lxss"
)

// WslDistro struct
type WslDistro struct {
	Default    bool
	BasePath   string
	Name       string
	DefaultUID uint64
	State      uint64
	Flags      uint64
	Version    uint64

	// Those are related to filesystem
	Executable string
}

// WslXMLRoot struct
type WslXMLRoot struct {
	XMLName      xml.Name `xml:"Package"`
	Applications WslXMLApplications
}

// WslXMLApplications struct
type WslXMLApplications struct {
	XMLName xml.Name `xml:"Applications"`
	App     WslXMLApplication
}

// WslXMLApplication stuct
type WslXMLApplication struct {
	XMLName    xml.Name `xml:"Application"`
	Executable string   `xml:"Executable,attr"`
	Id         string   `xml:"Id,attr"`
}
