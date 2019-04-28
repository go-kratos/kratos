package service

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"go-common/app/admin/ep/melloi/conf"
	"go-common/app/admin/ep/melloi/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

// Upload upload file
func (s *Service) Upload(c context.Context, uploadParam *model.UploadParam, formFile multipart.File, header *multipart.FileHeader) (id int, scriptPath string, err error) {

	var destFile *os.File

	//创建脚本保存路径
	if uploadParam.ScriptPath == "" {
		if scriptPath, err = s.uniqueFolderPath(uploadParam.Path); err != nil {
			return
		}
	} else {
		scriptPath = uploadParam.ScriptPath
	}

	//创建保存文件
	if destFile, err = os.Create(scriptPath + header.Filename); err != nil {
		log.Error("Create failed,error(%v)", err)
		return
	}
	defer destFile.Close()

	// 读取表单文件，写入保存文件
	if _, err = io.Copy(destFile, formFile); err != nil {
		log.Error("write file error (%v)", err)
		return
	}

	// jmx文件
	if strings.Contains(header.Filename, ".jmx") {
		treePath := model.TreePath{Department: uploadParam.Department, Project: uploadParam.Project, App: uploadParam.APP}
		script := model.Script{
			ProjectName: uploadParam.TestName,
			//TestName:    uploadParam.TestName,
			SavePath:            destFile.Name(),
			UpdateBy:            uploadParam.UserName,
			TreePath:            treePath,
			Upload:              true,
			Active:              1,
			ResJtl:              scriptPath + "jtl",
			JmeterLog:           scriptPath + "jmg",
			ScriptPath:          scriptPath,
			ArgumentString:      "[{\"\":\"\"}]",
			TestType:            1,
			Domain:              uploadParam.Domains,
			Fusing:              uploadParam.Fusing,
			UseBusinessStop:     uploadParam.UseBusinessStop,
			BusinessStopPercent: uploadParam.BusinessStopPercent,
		}
		if id, _, _, err = s.dao.AddScript(&script); err != nil {
			log.Error("s.dao.AddScript err :(%v)", err)
			return
		}
	}

	return
}

// UploadAndGetProtoInfo upload and get proto info
func (s *Service) UploadAndGetProtoInfo(c context.Context, uploadParam *model.UploadParam, formFile multipart.File, header *multipart.FileHeader) (res map[string]interface{}, err error) {
	// 获取上传路径, Path 为root_path ,ScriptPath 为 root_path + import路径
	_, scriptPath, err := s.Upload(c, uploadParam, formFile, header)
	if err != nil {
		log.Error("Write file failed, error(%v)", err)
		return
	}

	// 解析proto文件
	path, fileName := filepath.Split(path.Join(scriptPath, header.Filename))
	if res, err = s.ProtoParsing(path, fileName); err != nil {
		log.Error("parser grpc error(%v)", err)
		return
	}

	// 生成import文件路径
	if err = s.CreateGRPCImportDir(res, uploadParam.Path); err != nil {
		return
	}

	return
}

// CreateGRPCImportDir create grpc import path
func (s *Service) CreateGRPCImportDir(res map[string]interface{}, rootPath string) (err error) {
	if ipts, ok := res["import"]; ok {
		for _, ipt := range ipts.([]string) {
			iptPath, _ := filepath.Split(ipt)
			protoPath := &model.ProtoPathModel{RootPath: rootPath, ExtraPath: iptPath}
			if err = s.CreateProtoImportDir(protoPath); err != nil {
				log.Error("create proto import dir error(%v)", err)
				return
			}
		}
	}
	return
}

// CompileProtoFile compile proto file get jar path
func (s *Service) CompileProtoFile(protoPath, filename string) (jarPath string, err error) {
	return s.createGrpcJar(protoPath, filename)
}

//protocFile protoc file
func (s *Service) protocFile(scriptPath, protoPath string) (err error) {

	var (
		protoCMD      string
		protoJava     string
		protoFileList []string
		line          string
		protoFiles    = fmt.Sprintf("cd %s && find .  -name '*.proto' |  awk -F '\\.\\/' '{print $2}'", scriptPath)
		stdout        io.ReadCloser
	)
	// 获取scriptPath 目录下所有 proto文件
	cmd := exec.Command("/bin/bash", "-c", protoFiles)
	if stdout, err = cmd.StdoutPipe(); err != nil {
		log.Error("run protoFiles (%s) error(%v)", protoFiles, err)
		return
	}

	cmd.Start()

	reader := bufio.NewReader(stdout)
	for {
		if line, err = reader.ReadString('\n'); err != nil || io.EOF == err {
			break
		}
		if line != "" {
			protoFileList = append(protoFileList, line)
		}
	}

	// 编译所有proto文件
	for _, line := range protoFileList {
		protoCMD = fmt.Sprintf("protoc -I %s --java_out=%s %s", scriptPath, protoPath, path.Join(scriptPath, line))
		protoJava = fmt.Sprintf("protoc --plugin=protoc-gen-grpc-java=%s --grpc-java_out=%s -I %s %s",
			s.c.Grpc.ProtoJavaPluginPath, protoPath, scriptPath, path.Join(scriptPath, line))

		if err = exec.Command("/bin/bash", "-c", protoCMD).Run(); err != nil {
			log.Error("protoc --proto_path (%s) error (%v)", protoCMD, err)
			return ecode.MelloiProtocError
		}
		// 编译所有proto文件
		for _, line := range protoFileList {
			protoCMD = fmt.Sprintf("protoc -I %s --java_out=%s %s", scriptPath, protoPath, path.Join(scriptPath, line))
			protoJava = fmt.Sprintf("protoc --plugin=protoc-gen-grpc-java=%s --grpc-java_out=%s -I %s %s",
				s.c.Grpc.ProtoJavaPluginPath, protoPath, scriptPath, path.Join(scriptPath, line))

			if err = exec.Command("/bin/bash", "-c", protoCMD).Run(); err != nil {
				log.Error("protoc --proto_path (%s) error (%v)", protoCMD, err)
				return ecode.MelloiProtocError
			}

			if err = exec.Command("/bin/bash", "-c", protoJava).Run(); err != nil {
				log.Error("protoc java (%s) error (%v)", protoCMD, err)
				return ecode.MelloiProtocError
			}
		}
		if err = exec.Command("/bin/bash", "-c", protoJava).Run(); err != nil {
			log.Error("protoc java (%s) error (%v)", protoCMD, err)
			return ecode.MelloiProtocError
		}
	}
	return
}

//createGrpcJar crate grpc jar
func (s *Service) createGrpcJar(scriptPath, filename string) (jarPath string, err error) {
	protoPath := path.Join(scriptPath, "proto")
	if err = os.MkdirAll(protoPath, pathPerm); err != nil {
		log.Error("create proto path err ... (%v)", err)
		return
	}

	var (
		cmdJava     = fmt.Sprintf(" cd %s && find . -name '*.java' |grep -v '^\\./\\.' > sources.txt  ", protoPath)
		cmdJavac    = fmt.Sprintf("cd %s && javac -Djava.ext.dirs=%s -encoding utf-8 @sources.txt ", protoPath, s.c.Jmeter.JmeterExtLibPath)
		jarFilename = strings.Replace(filename, ".proto", "", -1) + ".jar"
		cmdJar      = fmt.Sprintf("cd %s && jar -cvf %s .", protoPath, jarFilename)
	)

	if err = s.protocFile(scriptPath, protoPath); err != nil {
		log.Error("protoc error (%v)", err)
		return
	}

	// 编译java文件
	if err = exec.Command("/bin/bash", "-c", cmdJava).Run(); err != nil {
		log.Error("find *.java (%s) error(%+v)", cmdJava, errors.WithStack(err))
		return
	}
	if err = exec.Command("/bin/bash", "-c", cmdJavac).Run(); err != nil {
		log.Error("javac protoc javac (%s) error(%v)", cmdJavac, err)
		err = ecode.MelloiJavacCompileError
		return
	}
	// 生成jar文件
	if err = exec.Command("/bin/bash", "-c", cmdJar).Run(); err != nil {
		log.Error("protoc jar (%s) error(%v)", cmdJar, err)
		err = ecode.MelloiJarError
		return
	}
	jarPath = path.Join(protoPath, jarFilename)
	return
}

//DownloadFile Download File
func (s *Service) DownloadFile(c *bm.Context, filePath string, w http.ResponseWriter) (err error) {
	var (
		Fi   *os.File
		info os.FileInfo
	)
	if info, err = os.Stat(filePath); err != nil {
		return
	}
	//如果文件太大，就禁止下载
	if info.Size() > conf.Conf.Melloi.MaxDowloadSize {
		err = ecode.MelloiBeyondFileSize
		return
	}
	if Fi, err = os.Open(filePath); err != nil {
		log.Error("open file err :(%v)", err)
		return
	}
	defer Fi.Close()
	w.Header().Set("Content-Disposition", "attachment; filename="+Fi.Name())
	if _, err = io.Copy(w, io.Reader(Fi)); err != nil {
		log.Error("copy file err :(%v)", err)
		err = ecode.MelloiCopyFileErr
		return
	}
	return
}

//ReadFile read file
func (s *Service) ReadFile(c *bm.Context, file, limit string) (data string, err error) {
	var (
		info os.FileInfo
		buff []byte
		cmd  *exec.Cmd
	)
	log.Info("开始读取文件--------- : (%s)", file)
	if info, err = os.Stat(file); err != nil {
		return
	}
	if info.Size() > conf.Conf.Melloi.MaxFileSize {
		cmd = exec.Command("/bin/bash", "-c", "tail -1500 "+file)
	} else {
		cmd = exec.Command("/bin/bash", "-c", "cat "+file)
	}
	buff, err = cmd.Output()
	if err != nil {
		return
	}
	data = string(buff)
	return
}

//IsFileExists is file exists
func (s *Service) IsFileExists(c *bm.Context, fileName string) (bool, error) {
	return exists(fileName)
}

// exists exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
