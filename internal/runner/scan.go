package runner

import (
	"io/fs"
	"path/filepath"
)

func (r *Runner) ReadScans() {
	var scans []string
	r.Logger.Debugf("searching path: %v", r.Options.Path)
	err := filepath.WalkDir(r.Options.Path, func(path string, d fs.DirEntry, err error) error {
		r.Logger.Debugf("found file: %v", filepath.Base(path))
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".xml" {
			scans = append(scans, path)
			r.Logger.Debugf("%v appears to be a valid scan", filepath.Base(path))
		}
		return nil
	})
	if err != nil {
		r.Logger.Errorf("error parsing path: %v", err)
	}
	r.Logger.Debugf("found scans: %v", scans)
}
