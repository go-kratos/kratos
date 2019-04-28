package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"go-common/app/job/main/appstatic/model"
	"go-common/library/log"
)

const (
	errFormat = "Func:[%s] - Step:[%s] - Error:[%v]"
	_zipType  = "application/zip"
	_tmpSlash = "*"
	// file type
	_fullPackage           = 0
	_diffPackge            = 1
	_calculationInProgress = 2
)

// Sync modified season data to the license owner
func (s *Service) calDiffproc() {
	defer s.waiter.Done()
	var (
		file *model.ResourceFile
		err  error
	)
	for {
		if s.daoClosed {
			log.Info("DB closed!")
			return
		}
		time.Sleep(time.Duration(s.c.Cfg.Diff.FreDiff))
		// check recently created files
		if file, err = s.dao.DiffNew(ctx); err != nil {
			log.Error(errFormat, "callDiffproc", "diffNew", err)
			continue
		}
		// if recently created file not found, we retry failed files
		if file == nil {
			if file, err = s.dao.DiffRetry(ctx); err != nil {
				log.Error(errFormat, "callDiffproc", "DiffRetry", err)
				continue
			}
			if file == nil {
				log.Error(errFormat, "callDiffproc", "DiffRetry", "No Diff To Calculate")
				continue
			}
		}
		// begin calculate
		log.Info("Calculate Diff for (fileID:%d, ResID: %d, FileName:%s)", file.ID, file.ResourceID, file.Name)
		if err := s.calDiff(file); err != nil {
			log.Error(errFormat, "calDiffproc", "calDiff", err)
			continue
		}
	}
}

// get the file info from DB and download the file to local
func (s *Service) catchFile(file *model.ResourceFile, version int) (newPath string, err error) {
	var (
		size   int64
		newPkg *model.ResourceFile
	)
	// get the new package info
	if newPkg, err = s.getFile(file.ResourceID, _fullPackage, version); err != nil {
		log.Error(errFormat, "calDiff", fmt.Sprintf("getFileNew %v", file), err)
		s.chgFStatus(file.ID, _diffPackge)
		return
	}
	// save the new package
	newPkg.Name = strings.Replace(newPkg.Name, "/", _tmpSlash, 1)    // avoid open file error
	newPath = fmt.Sprintf("%s/%s", s.c.Cfg.Diff.Folder, newPkg.Name) // combine the local path
	if size, err = s.dao.DownloadFile(ctx, newPkg.URL, newPath); err != nil {
		log.Error(errFormat, "calDiff", fmt.Sprintf("saveFile %v", file), err)
		s.chgFStatus(file.ID, _diffPackge)
		return
	}
	// check result file
	if _, err = os.Stat(newPath); os.IsNotExist(err) {
		log.Error(errFormat, "diffCmd", "IsNotExist", newPath+" - File Not Exist")
		return
	}
	log.Info("Save File From URL [%s] in [%s], Size: %d", newPkg.URL, newPath, size)
	return
}

// calculate the diff for one file struct and upload the result to BFS and fill the URL
func (s *Service) calDiff(file *model.ResourceFile) (err error) {
	var (
		newPath    string
		oldPath    string
		patchPath  string
		patchURL   string
		patchFInfo *model.FileInfo
	)
	// update the status of the file, to avoid being picked by another
	if err = s.chgFStatus(file.ID, _calculationInProgress); err != nil {
		log.Error(errFormat, "calDiff", "chgFStatus", err)
		return
	}
	// save the new file
	if newPath, err = s.catchFile(file, 0); err != nil {
		log.Error(errFormat, "callDiff", "catchFile", err)
		return
	}
	// save the old file
	if oldPath, err = s.catchFile(file, int(file.FromVer)); err != nil {
		log.Error(errFormat, "callDiff", "catchFile", err)
		return
	}
	// exec bsdiff to get the patch file and upload it
	file.Name = strings.Replace(file.Name, "/", _tmpSlash, 1)
	patchPath = s.c.Cfg.Diff.Folder + "/" + file.Name
	if patchFInfo, patchURL, err = s.diffCmd(file.Name, patchPath, newPath, oldPath); err != nil {
		log.Error(errFormat, "calDiff", "diffCmd", err)
		s.chgFStatus(file.ID, _diffPackge)
		return
	}
	log.Info("Upload Path File From [%s] to [%s], Size: %d", patchPath, patchURL)
	// save the url to the file
	if err = s.dao.SaveFile(ctx, file.ID, &model.FileInfo{
		Name: patchFInfo.Name,
		Size: patchFInfo.Size,
		Md5:  patchFInfo.Md5,
		URL:  patchURL}); err != nil {
		log.Error(errFormat, "calDiff", "updateURL", err)
		return
	}
	// delete all the packages used
	if err = delPkgs(newPath, oldPath, patchPath); err != nil {
		log.Error(errFormat, "calDiff", "delPkgs", err)
	}
	return
}

// delete all the packages used to generate the diff pkg
func delPkgs(newPath string, oldPath string, patchPath string) (err error) {
	if err = deleteFile(newPath); err != nil {
		log.Error(errFormat, "delPkgs", "NewPath", err)
		return
	}
	if err = deleteFile(oldPath); err != nil {
		log.Error(errFormat, "delPkgs", "oldPath", err)
		return
	}
	if err = deleteFile(patchPath); err != nil {
		log.Error(errFormat, "delPkgs", "patchPath", err)
	}
	return
}

// delete one file
func deleteFile(path string) (err error) {
	// check file
	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = fmt.Errorf("File %s Not exist", path)
		log.Error(errFormat, "deleteFile", "IsNotExist", "File Not Exist")
		return
	}
	if err = os.Remove(path); err != nil {
		log.Error(errFormat, "deleteFile", "Remove", err)
		return
	}
	log.Info("Delete %s Succesfully", path)
	return
}

// execute the bsdiff command to have the result
func (s *Service) diffCmd(patchName string, patchPath string, newPath string, oldPath string) (patchFInfo *model.FileInfo, location string, err error) {
	var (
		content []byte
		fInfo   os.FileInfo
	)
	begin := time.Now()
	cmd := exec.Command("bsdiff", oldPath, newPath, patchPath)
	// exec Command
	if err = cmd.Run(); err != nil {
		log.Error(errFormat, "diffCmd", "cmdRun", err)
		return
	}
	timeCost := time.Since(begin)
	log.Info("BSDiff Command Finished. Time Cost: %v, Cmd:%s %s %s %s", timeCost, "bsdiff", oldPath, newPath, patchPath)
	// check patch file
	if fInfo, err = os.Stat(patchPath); os.IsNotExist(err) {
		log.Error(errFormat, "diffCmd", "patchFileCheck", "File Not Exist")
		return
	}
	log.Info("Patch File Generated, Name: %s, Size: %s.", fInfo.Name(), fInfo.Size())
	// read patch file and upload
	if content, err = ioutil.ReadFile(patchPath); err != nil {
		log.Error(errFormat, "diffCmd", "ReadFile_Patch", err)
		return
	}
	if patchFInfo, err = s.ParseFile(content); err != nil {
		log.Error(errFormat, "diffCmd", "ParsePatchFile", err)
		return
	}
	patchFInfo.Name = rename(patchName, patchFInfo.Md5)
	// upload patch file to bfs
	location, err = s.Upload(context.Background(), patchFInfo.Name, _zipType, time.Now().Unix(), content)
	if err != nil {
		log.Error(errFormat, "diffCmd", "UploadPatch", err)
	}
	return
}

// split the patchName, pick the ModID (res[0]), VersionInfo (res[1]), Insert the md5 inside
func rename(patchName string, md5 string) (newName string) {
	res := strings.Split(patchName, _tmpSlash)
	if len(res) != 2 {
		log.Error("patchName %s can't split", patchName)
		return patchName
	}
	return res[0] + "_" + md5 + "/" + res[1]
}

// get file object ( struct )
func (s *Service) getFile(resID int, fileType int, version int) (file *model.ResourceFile, err error) {
	file = &model.ResourceFile{}
	var (
		res    *model.Resource // current version
		resHis *model.Resource // history version
		poolID int
	)
	if res, err = s.dao.ParseResID(ctx, resID); err != nil {
		log.Error("[getFile]-[findPool %d]-Error(%v)", resID, err)
		return
	}
	poolID = int(res.PoolID)
	if version != 0 { // full pkg of the history version
		if resHis, err = s.dao.ParseResVer(ctx, poolID, version); err != nil {
			log.Error("[getFile]-[findVersion]-Error(%v)", err)
			return
		}
		resID = int(resHis.ID)
	}
	if file, err = s.dao.ReadyFile(ctx, resID, fileType); err != nil {
		log.Error(errFormat, "getUrl", "First", err)
	}
	return
}

// change the file's status
func (s *Service) chgFStatus(fileID int, status int) (err error) {
	if err = s.dao.UpdateStatus(ctx, status, fileID); err != nil {
		log.Error(errFormat, "chgFStatus", "update", err)
	}
	return
}
