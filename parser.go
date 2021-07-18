// MIT License
//
// Copyright (c) 2019 Aloïs Micard
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
//	copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package mvnparser

import (
	"encoding/xml"
	"io"
)

// Represent a POM file
type MavenProject struct {
	XMLName              xml.Name             `xml:"project,omitempty"`
	ModelVersion         string               `xml:"modelVersion,omitempty"`
	Parent               Parent               `xml:"parent,omitempty"`
	GroupId              string               `xml:"groupId,omitempty"`
	ArtifactId           string               `xml:"artifactId,omitempty"`
	Version              string               `xml:"version,omitempty"`
	Packaging            string               `xml:"packaging,omitempty"`
	Name                 string               `xml:"name,omitempty"`
	Description          string               `xml:"description,omitempty"`
	Repositories         []Repository         `xml:"repositories>repository,omitempty"`
	Properties           Properties           `xml:"properties,omitempty"`
	DependencyManagement DependencyManagement `xml:"dependencyManagement,omitempty"`
	Dependencies         []Dependency         `xml:"dependencies>dependency,omitempty"`
	Profiles             []Profile            `xml:"profiles,omitempty,omitempty"`
	Build                Build                `xml:"build,omitempty"`
	PluginRepositories   []PluginRepository   `xml:"pluginRepositories>pluginRepository,omitempty"`
	Modules              []string             `xml:"modules>module,omitempty"`
}

// Represent the parent of the project
type Parent struct {
	GroupId    string `xml:"groupId,omitempty"`
	ArtifactId string `xml:"artifactId,omitempty"`
	Version    string `xml:"version,omitempty"`
}

// Represent a dependency of the project
type Dependency struct {
	XMLName    xml.Name    `xml:"dependency,omitempty"`
	GroupId    string      `xml:"groupId,omitempty"`
	ArtifactId string      `xml:"artifactId,omitempty"`
	Version    string      `xml:"version,omitempty"`
	Classifier string      `xml:"classifier,omitempty"`
	Type       string      `xml:"type,omitempty"`
	Scope      string      `xml:"scope,omitempty"`
	Exclusions []Exclusion `xml:"exclusions>exclusion,omitempty"`
}

// Represent an exclusion
type Exclusion struct {
	XMLName    xml.Name `xml:"exclusion,omitempty"`
	GroupId    string   `xml:"groupId,omitempty"`
	ArtifactId string   `xml:"artifactId,omitempty"`
}

type DependencyManagement struct {
	Dependencies []Dependency `xml:"dependencies>dependency,omitempty"`
}

// Represent a repository
type Repository struct {
	Id   string `xml:"id,omitempty"`
	Name string `xml:"name,omitempty"`
	Url  string `xml:"url,omitempty"`
}

type Profile struct {
	Id    string `xml:"id,omitempty"`
	Build Build  `xml:"build,omitempty"`
}

type Build struct {
	// todo: final name ?
	Plugins []Plugin `xml:"plugins>plugin,omitempty"`
}

type Plugin struct {
	XMLName    xml.Name `xml:"plugin,omitempty"`
	GroupId    string   `xml:"groupId,omitempty"`
	ArtifactId string   `xml:"artifactId,omitempty"`
	Version    string   `xml:"version,omitempty"`
	//todo something like: Configuration map[string]string `xml:"configuration"`
	// todo executions
}

// Represent a pluginRepository
type PluginRepository struct {
	Id   string `xml:"id,omitempty"`
	Name string `xml:"name,omitempty"`
	Url  string `xml:"url,omitempty"`
}

// Represent Properties
type Properties map[string]string

func (p *Properties) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*p = map[string]string{}
	for {
		key := ""
		value := ""
		token, err := d.Token()
		if err == io.EOF {
			break
		}
		switch tokenType := token.(type) {
		case xml.StartElement:
			key = tokenType.Name.Local
			err := d.DecodeElement(&value, &start)
			if err != nil {
				return err
			}
			(*p)[key] = value
		}
	}
	return nil
}

func (s Properties) MarshalXML(e *xml.Encoder, start xml.StartElement) error {

	tokens := []xml.Token{start}

	for key, value := range s {
		t := xml.StartElement{Name: xml.Name{"", key}}
		tokens = append(tokens, t, xml.CharData(value), xml.EndElement{t.Name})
	}

	tokens = append(tokens, xml.EndElement{start.Name})

	for _, t := range tokens {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}

	// flush to ensure tokens are written
	return e.Flush()
}
