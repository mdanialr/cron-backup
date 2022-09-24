package archive

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// ZipFile archive the given bytes as the body and filename as the archive header.
func ZipFile(filename string, body []byte) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	defer w.Close()

	f, err := w.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create filename %s in archive: %s", filename, err)
	}
	_, err = f.Write(body)
	if err != nil {
		return nil, fmt.Errorf("failed to write body bytes to filename %s in archive: %s", filename, err)
	}

	return buf, nil
}

// ZipDir archive the given directory. Also follow symlink.
func ZipDir(dir string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	defer w.Close()

	fnWalkDir := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		// If it's dir, just return. So this could be archived recursively.
		if info.IsDir() {
			return nil
		}

		headerName, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		f, err := w.Create(headerName)
		if err != nil {
			return err
		}

		fl, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fl.Close()

		if _, err = io.Copy(f, fl); err != nil {
			return err
		}

		return nil
	}

	return buf, zipWalkDir(dir, fnWalkDir)
}

// zipWalkDir extends filepath.WalkDir to also follow symlinks
func zipWalkDir(path string, walkFn fs.WalkDirFunc) error {
	return walkDir(path, path, walkFn)
}

// walkDir wrap the real fs.WalkDirFunc that do the zip
func walkDir(filename, linkDirname string, walkDirFn fs.WalkDirFunc) error {
	symWalkDirFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fName, err := filepath.Rel(filename, path)
		if err != nil {
			return err
		}
		path = filepath.Join(linkDirname, fName)

		info, err := d.Info()
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			finalPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				return err
			}

			inf, err := os.Stat(finalPath)
			if err != nil {
				return walkDirFn(path, d, err)
			}
			if inf.IsDir() {
				return walkDir(finalPath, path, walkDirFn)
			}
		}

		return walkDirFn(path, d, err)
	}
	return filepath.WalkDir(filename, symWalkDirFunc)
}
