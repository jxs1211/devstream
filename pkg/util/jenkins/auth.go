package jenkins

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/bndr/gojenkins"

	"github.com/devstream-io/devstream/pkg/util/log"
)

type GitlabCredentials struct {
	XMLName     xml.Name `xml:"com.dabsquared.gitlabjenkins.connection.GitLabApiTokenImpl"`
	ID          string   `xml:"id"`
	Scope       string   `xml:"scope"`
	Description string   `xml:"description"`
	APIToken    string   `xml:"apiToken"`
}

func (j *jenkins) CreateGiltabCredential(id, gitlabToken string) error {
	cred := GitlabCredentials{
		ID:          id,
		Scope:       credentialScope,
		APIToken:    gitlabToken,
		Description: id,
	}
	return j.createCredential(id, cred)
}

func (j *jenkins) CreateSSHKeyCredential(id, userName, privateKey string) error {
	cred := gojenkins.SSHCredentials{
		ID:       id,
		Scope:    credentialScope,
		Username: userName,
		PrivateKeySource: &gojenkins.PrivateKey{
			Value: privateKey,
			Class: gojenkins.KeySourceDirectEntryType,
		},
		Description: id,
	}
	return j.createCredential(id, cred)
}

func (j *jenkins) CreatePasswordCredential(id, userName, password string) error {
	cred := gojenkins.UsernameCredentials{
		ID:          id,
		Scope:       credentialScope,
		Username:    userName,
		Password:    password,
		Description: id,
	}
	return j.createCredential(id, cred)
}

func (j *jenkins) CreateSecretCredential(id, secretText string) error {
	cred := gojenkins.StringCredentials{
		ID:          id,
		Scope:       credentialScope,
		Secret:      secretText,
		Description: id,
	}
	return j.createCredential(id, cred)

}

func (j *jenkins) createCredential(id string, cred interface{}) error {
	cm := &gojenkins.CredentialsManager{
		J: &j.Jenkins,
	}
	err := cm.Add(j.ctx, domain, cred)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Debugf("jenkins credential %s exist, try to update it", id)
			return cm.Update(j.ctx, domain, id, cred)
		}
		return fmt.Errorf("could not create credential: %v", err)
	}

	// get credential to validate creation
	getCred := map[string]string{}
	if err = cm.GetSingle(j.ctx, domain, id, getCred); err != nil {
		return fmt.Errorf("could not get credential: %v", err)
	}
	return nil
}
