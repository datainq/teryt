package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/datainq/rwmc"
	"github.com/datainq/teryt"
)

func PrintLocations(loc *teryt.Location) int {
	for i, p := range loc.Parts {
		if i > 0 {
			fmt.Fprint(os.Stdout, ", ")
		}
		fmt.Fprint(os.Stdout, p)
	}
	fmt.Println()
	n := 1
	for _, v := range loc.Children {
		n += PrintLocations(v)
	}
	return n
}

// OpenFile opens a file to read. If file is *.zip it uses a unzip wrapper
// and opens the only file or the file with the name same as zip with
// extension replaced by .csv
func OpenFile(filePath string) (io.ReadCloser, error) {
	var dataReader io.ReadCloser
	if strings.HasSuffix(filePath, ".zip") {
		r, err := zip.OpenReader(filePath)
		if err != nil {
			return nil, fmt.Errorf("cannot open: %w", err)
		}
		fileName := path.Base(filePath[:len(filePath)-4]) + ".csv"
		for _, f := range r.File {
			if strings.HasPrefix(f.Name, fileName) || len(r.File) == 1 {
				dataReader, err = f.Open()
				if err != nil {
					_ = r.Close()
					return nil, fmt.Errorf("cannot open file: %filePath", err)
				}
				return rwmc.NewReadMultiCloser(dataReader, r), nil
			}
		}
		_ = r.Close()
		return nil, fmt.Errorf("file %q not found in ZIP: %q", fileName, filePath)
	}
	return os.Open(filePath)
}

func LoadData(tercFile, simcFile string) ([]teryt.SetTERC, []teryt.SetSIMC, error) {
	tercReader, err := OpenFile(tercFile)
	if err != nil {
		return nil, nil, fmt.Errorf("open TERC file: %w", err)
	}
	defer tercReader.Close()
	tercData, err := teryt.TERCFromReader(tercReader)
	if err != nil {
		return nil, nil, fmt.Errorf("reading TERC: %w", err)
	}

	simcReader, err := OpenFile(simcFile)
	if err != nil {
		return nil, nil, fmt.Errorf("open SIMC file: %w", err)
	}

	simcData, err := teryt.SIMCFromReader(simcReader)
	if err != nil {
		return nil, nil, fmt.Errorf("reading SIMC: %w", err)
	}
	return tercData, simcData, nil
}

func EnlistLocations(root *teryt.Location, cfg teryt.Config) []*teryt.Location {
	var ret []*teryt.Location
	if root.Name != "" {
		root.Build(cfg)
		ret = append(ret, root)
	}
	if len(root.Children) > 1 {
		for _, v := range root.Children {
			ret = append(ret, EnlistLocations(v, cfg)...)
		}
	}
	return ret
}
