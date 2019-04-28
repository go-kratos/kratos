package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"k8s.io/apimachinery/pkg/util/sets"
)

type DirOptions struct {
	NoParentOwners bool `json:"no_parent_owners,omitempty"`
}

type Config struct {
	Approvers         []string    `json:"approvers,omitempty"`
	Reviewers         []string    `json:"reviewers,omitempty"`
	RequiredReviewers []string    `json:"required_reviewers,omitempty"`
	Labels            []string    `json:"labels,omitempty"`
	Options           *DirOptions `json:"options,omitempty"`
}

// Empty checks if a SimpleConfig could be considered empty
func (s *Config) Empty() bool {
	return len(s.Approvers) == 0 && len(s.Reviewers) == 0 && len(s.RequiredReviewers) == 0 && len(s.Labels) == 0
}

type Owner struct {
	Config `json:",inline"`
}

type contributor struct {
	Owner    []string
	Author   []string
	Reviewer []string
}

// FullConfig contains Filters which apply specific Config to files matching its regexp
type FullConfig struct {
	Options DirOptions        `json:"options,omitempty"`
	Filters map[string]Config `json:"filters,omitempty"`
}

func readContributor(content []byte) (c *contributor) {
	var (
		lines      []string
		lineStr    string
		curSection string
	)
	c = &contributor{}
	lines = strings.Split(string(content), "\n")
	for _, lineStr = range lines {
		if lineStr == "" {
			continue
		}
		if strings.Contains(strings.ToLower(lineStr), "owner") {
			curSection = "owner"
			continue
		}
		if strings.Contains(strings.ToLower(lineStr), "author") {
			curSection = "author"
			continue
		}
		if strings.Contains(strings.ToLower(lineStr), "reviewer") {
			curSection = "reviewer"
			continue
		}
		switch curSection {
		case "owner":
			c.Owner = append(c.Owner, strings.TrimSpace(lineStr))
		case "author":
			c.Author = append(c.Author, strings.TrimSpace(lineStr))
		case "reviewer":
			c.Reviewer = append(c.Reviewer, strings.TrimSpace(lineStr))
		}
	}
	return
}

func Label(path string) string {
	var result string
	if filepath.HasPrefix(path, "app") {
		path = strings.Replace(path, "/CONTRIBUTORS.md", "", 1)
		path = strings.Replace(path, "/OWNERS", "", 1)
		path = strings.Replace(path, "app/", "", 1)
		if len(strings.Split(path, "/")) == 3 {
			result = path
		} else {
			if len(strings.Split(path, "/")) == 2 && filepath.HasPrefix(path, "infra") {
				result = path
			}
		}
	}
	return result
}

func main() {
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if filepath.HasPrefix(path, "vendor") && filepath.HasPrefix(path, "build") {
			return filepath.SkipDir
		}
		if path == "CONTRIBUTORS.md" {
			return nil
		}
		if !info.IsDir() && info.Name() == "CONTRIBUTORS.md" || info.Name() == "OWNERS" {
			owner := Owner{}
			approves := sets.NewString()
			reviewers := sets.NewString()
			labels := sets.NewString()
			if info.Name() == "CONTRIBUTORS.md" {
				content, err := ioutil.ReadFile(path)
				if err != nil {
					fmt.Printf("fail to read contributor  %q: %v\n", path, err)
					return err
				}

				c := readContributor(content)

				approves.Insert(c.Owner...)
				reviewers.Insert(c.Author...)
				reviewers.Insert(c.Reviewer...)
			}

			if strings.Contains(path, "/main/") {
				labels.Insert("main")
			} else {
				if strings.Contains(path, "/ep/") {
					labels.Insert("ep")
				} else {
					if strings.Contains(path, "/live/") {
						labels.Insert("live")
					} else {
						if strings.Contains(path, "/openplatform/") {
							labels.Insert("openplatform")
						} else {
							if strings.Contains(path, "bbq/") {
								labels.Insert("bbq")
							}
						}
					}
				}
			}

			if strings.Contains(path, "admin/") {
				labels.Insert("admin")
			} else {
				if strings.Contains(path, "common/") {
					labels.Insert("common")
				} else {
					if strings.Contains(path, "interface/") {
						labels.Insert("interface")
					} else {
						if strings.Contains(path, "job/") {
							labels.Insert("job")
						} else {
							if strings.Contains(path, "service/") {
								labels.Insert("service")
							} else {
								if strings.Contains(path, "infra/") {
									labels.Insert("infra")
								} else {
									if strings.Contains(path, "tool/") {
										labels.Insert("tool")
									}
								}
							}
						}
					}
				}
			}
			oldyaml, err := ioutil.ReadFile(strings.Replace(path, "CONTRIBUTORS.md", "OWNERS", -1))
			if err == nil {
				var owner Owner
				err = yaml.Unmarshal(oldyaml, &owner)
				if err != nil || owner.Empty() {
					c, err := ParseFullConfig(oldyaml)
					if err != nil {
						return err
					}
					data, err := yaml.Marshal(c)
					if err != nil {
						fmt.Printf("fail to Marshal %q: %v\n", path, err)
						return nil
					}
					data = append([]byte("# See the OWNERS docs at https://go.k8s.io/owners\n\n"), data...)
					ownerpath := strings.Replace(path, "CONTRIBUTORS.md", "OWNERS", 1)
					err = ioutil.WriteFile(ownerpath, data, 0644)
					if err != nil {
						fmt.Printf("fail to write yaml %q: %v\n", path, err)
						return err
					}
					return nil
				}
				approves.Insert(owner.Approvers...)
				reviewers.Insert(owner.Reviewers...)
				labels.Insert(owner.Labels...)
			}
			labels.Insert(Label(path))
			approves.Delete("all", "")
			reviewers.Delete("all", "")
			labels.Delete("all", "")
			owner.Approvers = approves.List()
			owner.Reviewers = reviewers.List()
			owner.Labels = labels.List()
			if strings.Contains(path, "app") && len(strings.Split(path, "/")) > 4 {
				owner.Options = &DirOptions{}
				owner.Options.NoParentOwners = true
			}
			if strings.Contains(path, "library/ecode") || strings.Contains(path, "app/tool") || strings.Contains(path, "app/infra") && len(strings.Split(path, "/")) > 2 {
				owner.Options = &DirOptions{}
				owner.Options.NoParentOwners = true
			}
			data, err := yaml.Marshal(owner)
			if err != nil {
				fmt.Printf("fail to Marshal %q: %v\n", path, err)
				return nil
			}
			data = append([]byte("# See the OWNERS docs at https://go.k8s.io/owners\n\n"), data...)
			ownerpath := strings.Replace(path, "CONTRIBUTORS.md", "OWNERS", 1)
			err = ioutil.WriteFile(ownerpath, data, 0644)
			if err != nil {
				fmt.Printf("fail to write yaml %q: %v\n", path, err)
				return err
			}
			return nil
		}
		return nil
	})
}

// ParseFullConfig will unmarshal OWNERS file's content into a FullConfig
// Returns an error if the content cannot be unmarshalled
func ParseFullConfig(b []byte) (FullConfig, error) {
	full := new(FullConfig)
	err := yaml.Unmarshal(b, full)
	return *full, err
}
