package file

import (
	"os"
	"path/filepath"
	"path"

	"go-common/library/log"
)

func (f *File) getFiles() map[string]os.FileInfo {
	paths := map[string]os.FileInfo{}

	for _, p := range f.c.Paths {
		// logs in docker overlay2
		if MergedDir := f.ctx.Value("MergedDir"); MergedDir != nil {
			p = path.Join(MergedDir.(string), p)
		}

		matches, err := filepath.Glob(p)
		if err != nil {
			log.Error("glob(%s) failed: %v", p, err)
			continue
		}

		// Check any matched files to see if we need to start a harvester
		for _, file := range matches {
			// check if the file is in the exclude_files list
			//if f.isFileExcluded(file) {
			//	log.Info("input", "Exclude file: %s", file)
			//	continue
			//}

			// Fetch Lstat File info to detected also symlinks
			fileInfo, err := os.Lstat(file)
			if err != nil {
				log.Warn("lstat(%s) failed: %s", file, err)
				continue
			}

			if fileInfo.IsDir() {
				log.Warn("Skipping directory: %s", file)
				continue
			}

			isSymlink := fileInfo.Mode()&os.ModeSymlink > 0
			if isSymlink && !f.c.Symlinks {
				log.Warn("File %s skipped as it is a symlink.", file)
				continue
			}

			// Fetch Stat file info which fetches the inode. In case of a symlink, the original inode is fetched
			fileInfo, err = os.Stat(file)
			if err != nil {
				log.Warn("stat(%s) failed: %s", file, err)
				continue
			}

			// If symlink is enabled, it is checked that original is not part of same input
			// It original is harvested by other input, states will potentially overwrite each other
			//if p.config.Symlinks {
			//	for _, finfo := range paths {
			//		if os.SameFile(finfo, fileInfo) {
			//			log.Info("Same file found as symlink and originap. Skipping file: %s", file)
			//			continue OUTER
			//		}
			//	}
			//}
			paths[file] = fileInfo
		}
	}

	return paths
}
