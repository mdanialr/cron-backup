package arch

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Zip zipping exactly only one file
func Zip(src string, zipName string) error {
	zipBuffer, err := os.Create(zipName)
	if err != nil {
		return err
	}
	defer zipBuffer.Close()

	zipWriter := zip.NewWriter(zipBuffer)
	defer zipWriter.Close()

	fl, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fl.Close()

	info, err := fl.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Method = zip.Deflate

	headerWriter, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(headerWriter, fl)
	if err != nil {
		return err
	}

	return nil
}

// ZipDir zipping entire dir contents uses go stdlib
func ZipDir(src string, zipName string) error {
	zipBuffer, err := os.Create(zipName)
	if err != nil {
		return err
	}
	defer zipBuffer.Close()

	zipWriter := zip.NewWriter(zipBuffer)
	defer zipWriter.Close()

	nzWalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Method = zip.Deflate

		// If its dir, just return. So this could be archived recursively.
		if info.IsDir() {
			return nil
		}

		// Preserve dir structure in the zip archive.
		header.Name, err = filepath.Rel(src, path)
		if err != nil {
			return err
		}

		headerWriter, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		fl, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fl.Close()

		_, err = io.Copy(headerWriter, fl)
		if err != nil {
			return err
		}

		return nil
	})
	return nil
}

// nzWalkDir extends filepath.WalkDir to also follow symlinks
func nzWalkDir(path string, walkFn fs.WalkDirFunc) error {
	return walkDir(path, path, walkFn)
}

// walkDir wrap the real fs.WalkDirFunc that do the zip
func walkDir(filename string, linkDirname string, walkDirFn fs.WalkDirFunc) error {
	symWalkDirFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fname, err := filepath.Rel(filename, path)
		if err != nil {
			return err
		}
		path = filepath.Join(linkDirname, fname)

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
