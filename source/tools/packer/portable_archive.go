package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Include the application directory, not the staging directory or sibling archives.
func portableArchive(ctx context.Context, app, path string, windows bool) (err error) {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, out.Close()) }()
	var zw *zip.Writer
	var tw *tar.Writer
	if windows {
		zw = zip.NewWriter(out)
		defer func() { err = errors.Join(err, zw.Close()) }()
	} else {
		gz := gzip.NewWriter(out)
		defer func() { err = errors.Join(err, gz.Close()) }()
		tw = tar.NewWriter(gz)
		defer func() { err = errors.Join(err, tw.Close()) }()
	}
	return filepath.WalkDir(app, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.IsDir() && !info.Mode().IsRegular() {
			return fmt.Errorf("cannot archive non-regular file %s", path)
		}
		rel, err := filepath.Rel(filepath.Dir(app), path)
		if err != nil {
			return err
		}
		name := filepath.ToSlash(rel)
		if info.IsDir() {
			name += "/"
		}
		var w io.Writer
		if windows {
			h, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}
			h.Name, h.Method = name, zip.Deflate
			w, err = zw.CreateHeader(h)
			if err != nil {
				return err
			}
		} else {
			h, err := tar.FileInfoHeader(info, "")
			if err != nil {
				return err
			}
			h.Name, h.Uid, h.Gid, h.Uname, h.Gname = name, 0, 0, "", ""
			if err := tw.WriteHeader(h); err != nil {
				return err
			}
			w = tw
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, file)
		return errors.Join(err, file.Close())
	})
}
