package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/mohae/deepcopy"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/test-infra/prow/config"
)

// DefaultTriggerFor returns the default regexp string used to match comments
// that should trigger the job with this name.
func DefaultTriggerFor(name string) string {
	return fmt.Sprintf(`(?m)^\+test( | .* )%s,?($|\s.*)`, name)
}

// DefaultRerunCommandFor returns the default rerun command for the job with
// this name.
func DefaultRerunCommandFor(name string) string {
	return fmt.Sprintf("+test %s", name)
}

type Image struct {
	Image []struct {
		Name  string `yaml:"name"`
		Image string `yaml:"image"`
	} `yaml:"image"`
}

type Global struct {
	Template   *config.Config
	AppendTask *config.Config
	AlwaysRun  *config.Config
	Result     *config.Config

	Image         map[string]string
	Labels        sets.String
	DefaultLabels sets.String

	TemplateLabels Configuration
}

var GlobalStatue Global

type Owner struct {
	Approvers []string `yaml:"approvers"`
	Reviewers []string `yaml:"reviewers"`
	Labels    []string `yaml:"labels"`
}

// LabelTarget specifies the intent of the label (PR or issue)
type LabelTarget string

const (
	bothTarget = "both"
)

type Label struct {
	// Name is the current name of the label
	Name string `json:"name"`
	// Color is rrggbb or color
	Color string `json:"color"`
	// Description is brief text explaining its meaning, who can apply it
	Description string `json:"description"` // What does this label mean, who can apply it
	// Target specifies whether it targets PRs, issues or both
	Target LabelTarget `json:"target"`
	// ProwPlugin specifies which prow plugin add/removes this label
	ProwPlugin string `json:"prowPlugin,omitempty"`
	// AddedBy specifies whether human/munger/bot adds the label
	AddedBy string `json:"addedBy"`
	// Previously lists deprecated names for this label
	Previously []Label `json:"previously,omitempty"`
	// DeleteAfter specifies the label is retired and a safe date for deletion
	DeleteAfter *time.Time `json:"deleteAfter,omitempty"`
}

// RepoConfig contains only labels for the moment
type RepoConfig struct {
	Labels []Label `json:"labels"`
}

// Configuration is a list of Required Labels to sync in all kubernetes repos
type Configuration struct {
	Repos   map[string]RepoConfig `json:"repos,omitempty"`
	Default RepoConfig            `json:"default"`
}

func generate() {
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, "app") && len(strings.Split(path, "/")) > 5 {
			return filepath.SkipDir
		}
		if strings.HasPrefix(path, "vendor") || strings.HasPrefix(path, "build") || strings.HasPrefix(path, ".rider") || strings.HasPrefix(path, ".git") {
			return nil
		}
		if info.Name() == "OWNERS" && !info.IsDir() {
			var owner Owner
			yamlFile, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			err = yaml.Unmarshal(yamlFile, &owner)
			if err != nil {
				return err
			}
			if len(owner.Labels) == 0 {
				return nil
			}
			GlobalStatue.Labels.Insert(owner.Labels...)
			ts, ok := GlobalStatue.Template.JobConfig.Presubmits["platform/go-common"]
			if !ok {
				fmt.Println("wrong project name")
				return nil
			}
			labels := sets.NewString(owner.Labels...)
			isLib := labels.Has("library")
			labels.Delete("library", "admin", "interface", "infra", "common", "service", "job", "vendor", "tool")
			labels.Delete("bbq", "ep", "ops", "video", "openplatform", "main", "live")
			labels.Delete("new-project", "new-main-service-project", "new-main-job-project", "new-main-interface-project", "new-main-admin-project")
			owner.Labels = labels.List()
			if len(owner.Labels) == 0 {
				return nil
			}

			for _, t := range ts {
				if isLib {
					if t.Name == "__bazel_build_job_name__" || t.Name == "__bazel_test_job_name__" {
						continue
					}
				}
				v := (deepcopy.Copy(t)).(config.Presubmit)
				v.Name = JobName(v.Name, owner.Labels[0])
				v.Context = Trigger(v.Name, owner.Labels[0])
				//v.Spec.Containers[0].Image = JobImage(v.Spec.Containers[0].Image)
				v.Spec.Containers[0].Name = v.Name
				for index, arg := range v.Spec.Containers[0].Args {
					if strings.Contains(arg, "<<bazel_dir_param>>") {
						v.Spec.Containers[0].Args[index] = JobBazelPath(v.Spec.Containers[0].Args[index], owner.Labels[0])
					}
				}
				v.UntrustedLabels = []string{}
				v.Trigger = DefaultTriggerFor(v.Name)
				v.RerunCommand = DefaultRerunCommandFor(v.Name)
				v.RunPRPushed = true
				v.TrustedLabels = append(v.TrustedLabels, owner.Labels[0])
				v.UntrustedLabels = append(v.UntrustedLabels, t.UntrustedLabels...)
				GlobalStatue.Result.Presubmits["platform/go-common"] = append(GlobalStatue.Result.Presubmits["platform/go-common"], v)
			}
			return nil
		}
		return nil
	})
	//GlobalStatue.AlwaysRun.Presubmits["platform/go-common"][0].RunAfterSuccess = append(GlobalStatue.AlwaysRun.Presubmits["platform/go-common"][0].RunAfterSuccess, GlobalStatue.Result.Presubmits["platform/go-common"]...)
	//GlobalStatue.AlwaysRun.Presubmits["platform/go-common"][0].RunAfterSuccess = append(GlobalStatue.AlwaysRun.Presubmits["platform/go-common"][0].RunAfterSuccess, GlobalStatue.AppendTask.Presubmits["platform/go-common"]...)
	for _, v := range GlobalStatue.AppendTask.Presubmits["platform/go-common"] {
		v.Trigger = DefaultTriggerFor(v.Name)
		v.RerunCommand = DefaultRerunCommandFor(v.Name)
		GlobalStatue.Result.Presubmits["platform/go-common"] = append(GlobalStatue.Result.Presubmits["platform/go-common"], v)
	}
	d, err := yaml.Marshal(GlobalStatue.Result)
	if err != nil {
		fmt.Println("fail to Marshal")
	}
	ioutil.WriteFile("./build/root/go_common_job.yaml", d, 0644)
	generateLabel()
}

func replaceimage() {
	ts, ok := GlobalStatue.Template.JobConfig.Presubmits["platform/go-common"]
	if !ok {
		fmt.Println("wrong project name")
		return
	}
	for _, t := range ts {
		t.Spec.Containers[0].Image = JobImage(t.Spec.Containers[0].Image)
	}
	at, ok := GlobalStatue.AppendTask.JobConfig.Presubmits["platform/go-common"]
	if !ok {
		fmt.Println("wrong project name")
		return
	}
	for _, t := range at {
		t.Spec.Containers[0].Image = JobImage(t.Spec.Containers[0].Image)
	}

	ar, ok := GlobalStatue.AlwaysRun.JobConfig.Presubmits["platform/go-common"]
	if !ok {
		fmt.Println("wrong project name")
		return
	}
	for _, a := range ar {
		a.Spec.Containers[0].Image = JobImage(a.Spec.Containers[0].Image)
	}

}

func generateLabel() (err error) {
	var repo RepoConfig
	repo.Labels = []Label{}
	for _, label := range GlobalStatue.Labels.List() {
		if GlobalStatue.DefaultLabels.Has(label) {
			continue
		}
		repo.Labels = append(repo.Labels, Label{
			Name:        label,
			Color:       "0052cc",
			Description: "Categorizes an issue or PR as relevant to " + label,
			Target:      bothTarget,
			AddedBy:     "anyone",
			ProwPlugin:  "label",
		})
	}
	GlobalStatue.TemplateLabels.Repos = make(map[string]RepoConfig)
	GlobalStatue.TemplateLabels.Repos["platform/go-common"] = repo
	d, err := yaml.Marshal(GlobalStatue.TemplateLabels)
	if err != nil {
		fmt.Println("fail to Marshal")
	}
	ioutil.WriteFile("./build/root/labels.yaml", d, 0644)
	return nil
}

func ReadTemplate() (err error) {
	GlobalStatue.Template, err = config.Load("./build/config.yaml", "./build/template/task")
	if err != nil {
		fmt.Println(err)
		return err
	}
	GlobalStatue.AlwaysRun, err = config.Load("./build/config.yaml", "./build/template/always_run.yaml")
	if err != nil {
		fmt.Println(err)
		return err
	}
	GlobalStatue.AppendTask, err = config.Load("./build/config.yaml", "./build/template/append_task")
	if err != nil {
		fmt.Println(err)
		return err
	}

	yamlFile, err := ioutil.ReadFile("./build/template/image.yaml")
	if err != nil {
		fmt.Println("yamlFile.Get err ", err)
	}
	var i Image
	err = yaml.Unmarshal(yamlFile, &i)
	if err != nil {
		fmt.Println("Unmarshal: ", err)
		return err
	}
	for _, im := range i.Image {
		GlobalStatue.Image[im.Name] = im.Image
	}

	labelTemp, err := ioutil.ReadFile("./build/labels-temp.yaml")
	if err != nil {
		fmt.Println("labels-temp err ", err)
	}
	var labels Configuration
	err = yaml.Unmarshal(labelTemp, &labels)
	if err != nil {
		fmt.Println("Unmarshal: ", err)
		return err
	}
	for _, n := range labels.Default.Labels {
		GlobalStatue.DefaultLabels.Insert(n.Name)
	}
	GlobalStatue.TemplateLabels = labels
	replaceimage()
	return nil
}

func init() {
	GlobalStatue.TemplateLabels = Configuration{}
	GlobalStatue.Labels = sets.NewString()
	GlobalStatue.DefaultLabels = sets.NewString()
	GlobalStatue.Result = &config.Config{}
	GlobalStatue.Result.Presubmits = make(map[string][]config.Presubmit)
	GlobalStatue.Image = make(map[string]string)
}

func main() {
	ReadTemplate()
	generate()

}
