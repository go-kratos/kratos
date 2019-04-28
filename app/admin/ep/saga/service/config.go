package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/ep/saga/conf"
	"go-common/app/admin/ep/saga/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/BurntSushi/toml"
)

const (
	_configFlag     = "\r\n"
	_sagaConfigFlag = "[[property.repos]]"
)

const (
	_svenConfigAppName   = "app_name"
	_svenConfigEnv       = "env"
	_svenConfigZone      = "zone"
	_svenConfigTreeID    = "tree_id"
	_svenConfigToken     = "token"
	_svenConfigBuild     = "build"
	_svenConfigUser      = "user"
	_svenConfigData      = "data"
	_svenConfigNames     = "names"
	_svenConfigMark      = "mark"
	_svenConfigConfigIDs = "config_ids"
	_svenConfigForce     = "force"
	_svenConfigIncrement = "increment"
)

const (
	_formatStr    = `        %s=%s`
	_formatStrQuo = `        %s="%s"`
	_formatValue  = `        %s=%v`
	_formatInt    = `        %s=%d`
)

const (
	_defaultBranch      = "master"
	_defaultLockTimeout = 600
)

const (
	_repoURL            = "URL"
	_repoGroup          = "Group"
	_repoName           = "Name"
	_repoLanguage       = "Language"
	_repoLockTimeout    = "LockTimeout"
	_repoAuthBranches   = "AuthBranches"
	_repoTargetBranches = "TargetBranches"
)

var (
	sagaConfigCnName = []string{
		"仓库地址",
		"仓库组名",
		"仓库名称",
		"仓库别名",
		"开发语言",
		"权限分支",
		"目标分支",
		"MR锁定超时时间(s)",
		"最少review人数",
		"是否关联pipeline",
		"自动合并",
		"权限限制",
		"准入标签",
		"超级权限用户",
	}

	sagaConfigMark = []string{
		"仓库地址",
		"仓库组名",
		"仓库名称",
		"仓库别名",
		"仓库使用语言",
		"saga的权限管控将以此分支配置的CONTRIBUTORS.md为准，即使CONTRIBUTORS.md在其他分支上更改了，也都会以此分支的鉴权信息为准。",
		"配置的分支可以触发saga行为，如配置 targetBranches为master、release分支，则MR的目标分支为master或release时，都能触发saga行为。支持通配",
		"每个仓库在每个时间点只能允许一个MR在合并。MR合并时会取得一个锁并在其合并结束后将其释放；如果获取的锁MR在lockTimeout时间内都未能结束，则会认为超时并将锁自动释放，这样其他MR才能有机会获取到独享锁及合并的机会。",
		"最终合并前除了需要owner点赞外，还需通过权限文件配置的Reviewer中minReviewer数量的人点赞后可合并。",
		"配置后saga将会检查pipeline执行结果，并会在最后merge前再次retry pipeline。不配置的话saga不会对pipeline执行结果进行判断。",
		"此配置以 relatePipeline 为基础，打开后saga在MR最终合并前将不再retry pipeline，并且 +mr 时如果pipeline还在运行中，待pipeline运行通过后MR将会自动合并。",
		"打开后owner的权限将只限定在当前目录，即如果子目录配置了owner等信息，根目录的owner等将不能再管控子目录。",
		"如果配置了标签，saga只合入打了此label的MR。",
		"如果有配置，原来的鉴权文件CONTRIBUTORS.md将会失效，需要合并的MR都必须通过super users的review。super users本身也拥有+mr直接合并的权利。",
	}
)

// SagaUserList ...
func (s *Service) SagaUserList(c context.Context) (resp []string, err error) {
	resp = conf.Conf.Property.Sven.SagaConfigsParam.UserList
	return
}

// QueryAllConfigFile ...
func (s *Service) QueryAllConfigFile(c context.Context, sessionID string, isSaga bool) (resp *model.ConfigData, err error) {
	var (
		url          = conf.Conf.Property.Sven.Configs + "?app_name=%s&tree_id=%s&env=%s&zone=%s&build_id=%s"
		sagaConfig   = conf.Conf.Property.Sven.SagaConfigsParam
		runnerConfig = conf.Conf.Property.Sven.ConfigsParam
	)

	if isSaga {
		url = fmt.Sprintf(url, sagaConfig.AppName, strconv.Itoa(sagaConfig.TreeID), sagaConfig.Env, sagaConfig.Zone, strconv.Itoa(sagaConfig.BuildId))
	} else {
		url = fmt.Sprintf(url, runnerConfig.AppName, strconv.Itoa(runnerConfig.TreeID), runnerConfig.Env, runnerConfig.Zone, strconv.Itoa(runnerConfig.BuildId))
	}

	return s.dao.QueryAllConfigFile(c, sessionID, url)
}

// QueryConfigFileContent ...
func (s *Service) QueryConfigFileContent(c context.Context, sessionID string) (content string, err error) {
	var (
		url      = conf.Conf.Property.Sven.ConfigValue
		fileName = conf.Conf.Property.Sven.SagaConfigsParam.FileName
		configs  *model.ConfigData
	)

	if configs, err = s.QueryAllConfigFile(c, sessionID, true); err != nil {
		return
	}
	for _, confValue := range configs.BuildFiles {
		if confValue.Name == fileName {
			log.Info("QueryConfigFileContent get config name: %s", fileName)
			id := strconv.Itoa(confValue.ID)
			url = fmt.Sprintf(url+"?config_id=%s", id)
			if content, err = s.dao.QueryConfigFileContent(c, sessionID, url); err != nil {
				return
			}
		}
	}
	return
}

// QueryProjectSagaConfig ...
func (s *Service) QueryProjectSagaConfig(c context.Context, sessionID string, projectID int) (sagaConfig *model.RepoConfig, err error) {
	var (
		projectInfo *model.ProjectInfo
		content     string
	)

	if content, err = s.QueryConfigFileContent(c, sessionID); err != nil {
		return
	}

	index := strings.Index(content, _sagaConfigFlag)
	if index < 0 {
		return
	}

	content = content[index:]
	log.Info("QueryProjectSagaConfig content: %s", content)

	Conf := &model.Config{}
	if _, err = toml.Decode(content, &Conf); err != nil {
		log.Error("QueryProjectSagaConfig Decode err(%+v)", err)
		return
	}
	if projectInfo, err = s.dao.ProjectInfoByID(projectID); err != nil {
		log.Error("QueryProjectSagaConfig ProjectInfoByID err(%+v)", err)
		return
	}
	projectUrl := strings.Replace(projectInfo.Repo, "git-test", "git", 1)
	for _, r := range Conf.Property.Repos {
		if r.URL == projectUrl {
			sagaConfig = r
			return
		}
	}
	return
}

// UpdateConfig ...
func (s *Service) UpdateConfig(c context.Context, sessionID, user, configFileName, configContent, mark string, isSaga bool) (resp *model.CommonResp, err error) {
	var (
		reqUrl = conf.Conf.Property.Sven.ConfigUpdate
		params = url.Values{}
	)

	if isSaga {
		params.Set(_svenConfigAppName, conf.Conf.Property.Sven.SagaConfigsParam.AppName)
		params.Set(_svenConfigEnv, conf.Conf.Property.Sven.SagaConfigsParam.Env)
		params.Set(_svenConfigZone, conf.Conf.Property.Sven.SagaConfigsParam.Zone)
		params.Set(_svenConfigTreeID, strconv.Itoa(conf.Conf.Property.Sven.SagaConfigsParam.TreeID))
		params.Set(_svenConfigToken, conf.Conf.Property.Sven.SagaConfigsParam.Token)
	} else {
		params.Set(_svenConfigAppName, conf.Conf.Property.Sven.ConfigsParam.AppName)
		params.Set(_svenConfigEnv, conf.Conf.Property.Sven.ConfigsParam.Env)
		params.Set(_svenConfigZone, conf.Conf.Property.Sven.ConfigsParam.Zone)
		params.Set(_svenConfigTreeID, strconv.Itoa(conf.Conf.Property.Sven.ConfigsParam.TreeID))
		params.Set(_svenConfigToken, conf.Conf.Property.Sven.ConfigsParam.Token)
	}

	data := `[{"name":"%s","comment":"%s","mark":"%s"}]`
	data = fmt.Sprintf(data, configFileName, configContent, mark)
	params.Set(_svenConfigData, data)
	params.Set(_svenConfigUser, user)

	log.Info("UpdateConfig params:%v", params)

	if resp, err = s.dao.RequestConfig(c, sessionID, reqUrl, params); err != nil {
		return
	}
	if resp.Code == ecode.OK.Code() && resp.Message == "0" {
		log.Info("RequestConfig success")
		resp = nil
	}
	return
}

// PublicConfig ...
func (s *Service) PublicConfig(c context.Context, sessionID, user, configFileName, mark string, isSaga bool) (resp *model.CommonResp, err error) {
	var (
		reqUrl = conf.Conf.Property.Sven.TagUpdate
		params = url.Values{}
	)

	if isSaga {
		params.Set(_svenConfigAppName, conf.Conf.Property.Sven.SagaConfigsParam.AppName)
		params.Set(_svenConfigEnv, conf.Conf.Property.Sven.SagaConfigsParam.Env)
		params.Set(_svenConfigZone, conf.Conf.Property.Sven.SagaConfigsParam.Zone)
		params.Set(_svenConfigTreeID, strconv.Itoa(conf.Conf.Property.Sven.SagaConfigsParam.TreeID))
		params.Set(_svenConfigToken, conf.Conf.Property.Sven.SagaConfigsParam.Token)
		params.Set(_svenConfigBuild, conf.Conf.Property.Sven.SagaConfigsParam.Build)

	} else {
		params.Set(_svenConfigAppName, conf.Conf.Property.Sven.ConfigsParam.AppName)
		params.Set(_svenConfigEnv, conf.Conf.Property.Sven.ConfigsParam.Env)
		params.Set(_svenConfigZone, conf.Conf.Property.Sven.ConfigsParam.Zone)
		params.Set(_svenConfigTreeID, strconv.Itoa(conf.Conf.Property.Sven.ConfigsParam.TreeID))
		params.Set(_svenConfigToken, conf.Conf.Property.Sven.ConfigsParam.Token)
		params.Set(_svenConfigBuild, conf.Conf.Property.Sven.ConfigsParam.Build)
	}
	params.Set(_svenConfigForce, strconv.Itoa(conf.Conf.Property.Sven.SagaConfigsParam.Force))
	params.Set(_svenConfigIncrement, strconv.Itoa(conf.Conf.Property.Sven.SagaConfigsParam.Increment))
	params.Set(_svenConfigMark, mark)
	params.Set(_svenConfigUser, user)
	params.Set(_svenConfigNames, configFileName)
	params.Set(_svenConfigConfigIDs, "")

	if resp, err = s.dao.RequestConfig(c, sessionID, reqUrl, params); err != nil {
		return
	}
	if resp.Code == ecode.OK.Code() && resp.Message == "0" {
		log.Info("RequestConfig success")
		resp = nil
	}
	return
}

// ParseRequestConfig ...
func (s *Service) ParseRequestConfig(projectInfo *model.ProjectInfo, configs []model.ConfigSagaItem) (requestConfig *model.RepoConfig, requestConfigStr string, err error) {
	var (
		configStr string
		content   []byte
	)

	configStr = doBasicDefault(projectInfo)

	for _, config := range configs {

		rv := reflect.ValueOf(config.Value)
		if !rv.IsValid() {
			configStr, _ = doDefault(config.Name, configStr)
			continue
		}

		if rv.Kind() == reflect.Slice {
			log.Info("%s is slice", config.Name)
			if rv.IsNil() {
				configStr, _ = doDefault(config.Name, configStr)
				continue
			}
			if content, err = json.Marshal(config.Value); err != nil {
				log.Error("ParseRequestConfig err(%+v)", err)
				return
			}

			configTr := fmt.Sprintf(_formatStr, config.Name, string(content))
			configStr = configStr + configTr + _configFlag

		} else {
			configTr := fmt.Sprintf(_formatValue, config.Name, config.Value)
			configStr = configStr + configTr + _configFlag
		}
	}
	log.Info("ParseRequestConfig: %s", configStr)
	requestConfigStr = configStr

	requestConfig = &model.RepoConfig{}
	if _, err = toml.Decode(configStr, &requestConfig); err != nil {
		log.Error("ParseRequestConfig toml decode err(%+v)", err)
		return
	}
	return
}

// doBasicDefault ...
func doBasicDefault(projectInfo *model.ProjectInfo) (configStr string) {
	var configTr string

	configStr = _configFlag

	configTr = fmt.Sprintf(_formatStrQuo, _repoURL, projectInfo.Repo)
	configStr = configStr + configTr + _configFlag

	configTr = fmt.Sprintf(_formatStrQuo, _repoGroup, projectInfo.SpaceName)
	configStr = configStr + configTr + _configFlag

	configTr = fmt.Sprintf(_formatStrQuo, _repoName, projectInfo.Name)
	configStr = configStr + configTr + _configFlag

	return configStr
}

// doDefault ...
func doDefault(name, config string) (configStr string, err error) {
	var (
		content   []byte
		defaultBr = []string{_defaultBranch}
	)
	configStr = config

	if strings.ToLower(name) == strings.ToLower(_repoLockTimeout) {

		configTr := fmt.Sprintf(_formatInt, name, _defaultLockTimeout)
		configStr = configStr + configTr + _configFlag
	}

	if strings.ToLower(name) == strings.ToLower(_repoAuthBranches) || strings.ToLower(name) == strings.ToLower(_repoTargetBranches) {

		if content, err = json.Marshal(defaultBr); err != nil {
			log.Error("Marshal err(%+v)", err)
			return
		}
		configTr := fmt.Sprintf(_formatStr, name, string(content))
		configStr = configStr + configTr + _configFlag
	}
	return
}

// ParseSvenConfig ...
func (s *Service) ParseSvenConfig(c context.Context, sessionID, projectUrl string) (fileContent, svenConfig string, err error) {
	var (
		content        string
		projectConfigs []string
	)

	if fileContent, err = s.QueryConfigFileContent(c, sessionID); err != nil {
		return
	}
	log.Info("ParseSvenConfig fileContent : %s", fileContent)

	index := strings.Index(fileContent, _sagaConfigFlag)
	if index < 0 {
		log.Warn("ParseSvenConfig not found any config flag: %s", projectUrl)
		return
	}
	content = fileContent[index:]

	projectConfigs = strings.Split(content, _sagaConfigFlag)
	for i := 0; i < len(projectConfigs); i++ {

		if strings.Contains(projectConfigs[i], projectUrl) {
			svenConfig = projectConfigs[i]
			return
		}
	}
	return
}

// ReplaceConfig ...
func (s *Service) ReplaceConfig(c context.Context, username, sessionID string, projectInfo *model.ProjectInfo, req *model.ConfigList) (newConfig string, err error) {
	var (
		requestConfig    *model.RepoConfig
		requestConfigStr string
		fileContent      string
		svenConfig       string
	)

	if requestConfig, requestConfigStr, err = s.ParseRequestConfig(projectInfo, req.Configs); err != nil {
		return
	}

	if fileContent, svenConfig, err = s.ParseSvenConfig(c, sessionID, projectInfo.Repo); err != nil {
		return
	}
	if len(svenConfig) <= 0 {
		return
	}

	index := strings.Index(svenConfig, "#")
	if index > 0 {
		annotate := svenConfig[index:]
		requestConfigStr = requestConfigStr + _configFlag + "    " + annotate
	}

	log.Info("ReplaceConfig requestConfig: %v", requestConfig)
	log.Info("ReplaceConfig requestConfigStr: %s", requestConfigStr)
	log.Info("ReplaceConfig svenConfig: %s", svenConfig)

	newConfig = strings.Replace(fileContent, svenConfig, requestConfigStr, 1)
	log.Info("ReplaceConfig newConfig: %s", newConfig)

	return
}

// ReleaseSagaConfig ...
func (s *Service) ReleaseSagaConfig(c context.Context, username, sessionID string, req *model.ConfigList) (resp *model.CommonResp, err error) {
	var (
		configFileName   = conf.Conf.Property.Sven.SagaConfigsParam.FileName
		sagaConfig       *model.RepoConfig
		projectInfo      *model.ProjectInfo
		projectID        = req.ProjectID
		newConfigContent string
	)

	if sagaConfig, err = s.QueryProjectSagaConfig(c, sessionID, req.ProjectID); err != nil || sagaConfig == nil {
		log.Error("ReleaseSagaConfig  err(%+v)", err)
		return
	}

	if projectInfo, err = s.dao.ProjectInfoByID(projectID); err != nil {
		log.Error("ProjectInfoByID err(%+v)", err)
		return
	}
	projectInfo.Repo = strings.Replace(projectInfo.Repo, "git-test", "git", 1)

	log.Info("ReleaseSagaConfig query project: %s, sagaConfig: %v", projectInfo.Name, sagaConfig)

	if sagaConfig.Name == projectInfo.Name {

		if newConfigContent, err = s.ReplaceConfig(c, username, sessionID, projectInfo, req); err != nil {
			return
		}

		year, month, day := time.Now().Date()
		monthInt := int(month)
		hour := time.Now().Hour()
		updateMark := fmt.Sprintf("%s-%d-%d-%d-%d | from saga-admin", username, year, monthInt, day, hour)

		newConfigContent = strconv.Quote(newConfigContent)[1 : len(strconv.Quote(newConfigContent))-1]
		log.Info("ReleaseSagaConfig newConfig: %s", newConfigContent)

		if _, err = s.UpdateConfig(c, sessionID, username, configFileName, newConfigContent, updateMark, true); err != nil {
			log.Error("UpdateConfig err(%+v)", err)
			return
		}
		if _, err = s.PublicConfig(c, sessionID, username, configFileName, updateMark, true); err != nil {
			log.Error("PublicConfig err(%+v)", err)
			return
		}
	}

	return
}

// OptionSaga ...
func (s *Service) OptionSaga(c context.Context, projectID, sessionID string) (resp []*model.OptionSagaItem, err error) {

	var sagaConfig *model.RepoConfig

	projectIDInt, _ := strconv.Atoi(projectID)
	if sagaConfig, err = s.QueryProjectSagaConfig(c, sessionID, projectIDInt); err != nil || sagaConfig == nil {
		log.Error("QueryProjectSagaConfig  err(%+v)", err)
		return
	}

	t := reflect.TypeOf(sagaConfig)
	if t.Kind() != reflect.Ptr {
		log.Info("OptionSaga the object is not a Ptr, but it is : %v", t.Kind())
		return
	}

	t = reflect.TypeOf(sagaConfig).Elem()
	if t.Kind() != reflect.Struct {
		log.Info("OptionSaga the object is not a struct, but it is : %v", t.Kind())
		return
	}

	v := reflect.ValueOf(sagaConfig).Elem()
	for i := 0; i < t.NumField(); i++ {

		f := t.Field(i)
		if f.Name == _repoURL || f.Name == _repoGroup || f.Name == _repoName || f.Name == _repoLanguage {
			continue
		}

		val := v.Field(i).Interface()
		log.Info("OptionSaga === %s: %v = %v", f.Name, f.Type, val)

		sagaItem := &model.OptionSagaItem{}

		sagaItem.Name = f.Name
		sagaItem.Value = val
		sagaItem.CNName = sagaConfigCnName[i]
		sagaItem.Remark = sagaConfigMark[i]

		configTr := fmt.Sprintf(`%v`, f.Type)
		sagaItem.Type = configTr

		resp = append(resp, sagaItem)
	}

	return
}
