package ci

import (
	"errors"
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/server"
	"github.com/devstream-io/devstream/pkg/util/downloader"
	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/template"
)

type ciConfigMap map[string]string

type CIConfig struct {
	Type server.CIServerType `validate:"oneof=jenkins github gitlab" mapstructure:"type"`
	// ConfigLocation represent location of ci config, it can be a remote location or local location
	ConfigLocation string `validate:"required_without=ConfigContents" mapstructure:"configLocation"`
	// Contents respent map of ci fileName to fileContent
	ConfigContentMap ciConfigMap            `validate:"required_without=ConfigLocation" mapstructure:"configContents"`
	Vars             map[string]interface{} `mapstructure:"vars"`
}

// SetContent is used to config ConfigContentMap for ci
func (c *CIConfig) SetContent(content string) {
	ciFileName := c.newCIServerClient().CIFilePath()
	if c.ConfigContentMap == nil {
		c.ConfigContentMap = ciConfigMap{}
	}
	c.ConfigContentMap[ciFileName] = content
}

func (c *CIConfig) SetContentMap(contentMap map[string]string) {
	c.ConfigContentMap = contentMap
}

func (c *CIConfig) getGitfileMap() (gitFileMap git.GitFileContentMap, err error) {
	if len(c.ConfigContentMap) == 0 {
		// 1. if ConfigContentMap is empty, get GitFileContentMap from ConfigLocation
		gitFileMap, err = c.getConfigContentFromLocation()
	} else {
		// 2. else render CIConfig.ConfigContentMap values
		gitFileMap = make(git.GitFileContentMap)
		ciServerClient := c.newCIServerClient()
		for filePath, content := range c.ConfigContentMap {
			scmFilePath := ciServerClient.GetGitNameFunc()(filePath, "")
			scmFileContent, err := c.renderContent(content)
			if err != nil {
				log.Debugf("ci render git files failed: %+v", err)
				return nil, err
			}
			gitFileMap[scmFilePath] = []byte(scmFileContent)
		}
	}
	if len(gitFileMap) == 0 {
		return nil, errors.New("ci can't get valid ci files")
	}
	return gitFileMap, err
}

func (c *CIConfig) newCIServerClient() (ciClient server.CIServerOptions) {
	return server.NewCIServer(c.Type)
}

func (c *CIConfig) renderContent(ciFileContent string) (string, error) {
	needRenderContent := len(c.Vars) > 0
	if needRenderContent {
		return template.New().FromContent(ciFileContent).SetDefaultRender(ciTemplateName, c.Vars).Render()
	}
	return ciFileContent, nil

}
func (c *CIConfig) getConfigContentFromLocation() (git.GitFileContentMap, error) {
	// 1. get resource
	getClient := downloader.ResourceClient{
		Source: c.ConfigLocation,
	}
	ciConfigPath, err := getClient.GetWithGoGetter()
	if err != nil {
		return nil, fmt.Errorf("ci get files by %s failed: %w", c.ConfigLocation, err)
	}
	defer getClient.CleanUp()
	// 2. get ci content map from ciConfigPath
	ciClient := c.newCIServerClient()
	return file.GetFileMap(
		ciConfigPath, ciClient.FilterCIFilesFunc(),
		ciClient.GetGitNameFunc(), processCIFilesFunc(c.Vars),
	)
}
