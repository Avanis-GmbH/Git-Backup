package config

import (
	"os"
	"path/filepath"

	logr "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var instance *Config
var configPath = "./config"

type Config struct {
	Organizations map[string]Account `yaml:"Organizations"`
	OrgaName      string             `yaml:"OrgaName"`
	OrgaToken     string             `yaml:"OrgaToken"`
	OrgaRepoType  string             `yaml:"OrgaRepoType"`

	// CloneUserRepos bool               `yaml:"CloneUserRepos"`
	Users map[string]Account `yaml:"Users"`

	OutputPath     string `yaml:"OutputPath"`
	UpdateInterval string `yaml:"UpdateInterval"`

	ListReferences bool `yaml:"ListReferences"`
	LogCommits     bool `yaml:"LogCommits"`
	LogLevel       int  `yaml:"LogLevel"`
}

type Account struct {
	Name         string `yaml:"Name"`
	Token        string `yaml:"Token"`
	Option       string `yaml:"Option"`
	ValidateName bool   `yaml:"ValidateName"` // Whether the User-/OrgaName has to be contained in the "full_name" of the repository
	BackupRepos  bool   `yaml:"BackupRepos"`
}

func GetConfig() *Config {
	if instance == nil {
		err := initConfig()
		if err != nil {
			logr.Fatalf("[config] Error initializing the config: %s", err.Error())
		}
	}

	return instance
}

func initConfig() error {
	instance = &Config{}

	if _, err := os.Stat(filepath.Join(configPath, "config.yml")); err != nil {
		err = createConfig()
		if err != nil {
			return err
		}
	}

	file, err := os.Open(filepath.Join(configPath, "config.yml"))
	if err != nil {
		return err
	}
	defer file.Close()

	err = yaml.NewDecoder(file).Decode(instance)
	if err != nil {
		return err
	}

	return nil
}

func createConfig() error {
	config := &Config{
		OrgaName:     "Default Orga",
		OrgaToken:    "",
		OrgaRepoType: "all",

		Organizations: map[string]Account{
			"1st": {
				Name:         "1st Orga",
				Token:        "",
				Option:       "all",
				ValidateName: false,
				BackupRepos:  true,
			},
		},

		// CloneUserRepos: true,
		Users: map[string]Account{
			"1st": {
				Name:         "1st User",
				Token:        "",
				Option:       "owner",
				ValidateName: false,
				BackupRepos:  true,
			},
		},
		OutputPath:     "../Repo-Backups/",
		UpdateInterval: "0 */12 * * * *",

		ListReferences: true,
		LogCommits:     false,
		LogLevel:       6,
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(configPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(filepath.Join(configPath, "config.yml"), data, 0600)
	if err != nil {
		return err
	}
	logr.Info("[config] created default configuration, exiting...")
	os.Exit(0)
	return nil
}
