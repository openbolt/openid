package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"os"
	"time"
	"io/ioutil"
	"path"
	"path/filepath"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindata_file_info struct {
	name string
	size int64
	mode os.FileMode
	modTime time.Time
}

func (fi bindata_file_info) Name() string {
	return fi.name
}
func (fi bindata_file_info) Size() int64 {
	return fi.size
}
func (fi bindata_file_info) Mode() os.FileMode {
	return fi.mode
}
func (fi bindata_file_info) ModTime() time.Time {
	return fi.modTime
}
func (fi bindata_file_info) IsDir() bool {
	return false
}
func (fi bindata_file_info) Sys() interface{} {
	return nil
}

var _hello_html = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x7c\x8f\x31\x4f\xc3\x30\x10\x85\x67\xfa\x2b\x4e\x1d\x3a\x62\x52\x3a\xa0\xca\x09\x42\x54\x48\x48\x48\xb0\x31\x46\xae\x7d\xd4\x27\xb9\xbe\xc8\xbe\x80\x9a\x5f\xcf\x85\x4c\x2c\x1d\xbc\x9c\xdf\xf7\xe9\x3d\x1b\xd1\x85\x6e\x75\x63\x85\x24\x61\xf7\x3e\x60\x7e\x3d\xc0\x33\xe7\x8c\x5e\xa0\xb9\xbd\x83\x8f\xc2\xdf\x14\xb0\x58\xb3\x44\x56\xd6\x2c\x8c\x3d\x72\xb8\xcc\x68\x6c\xba\x4f\x4c\x9e\xcf\x08\xc2\x30\x2b\x8e\x9c\xa4\xc2\x15\x19\x1c\xf0\xcc\x2a\x6a\xfe\xf8\x6d\xf7\xc6\x27\xca\x30\x56\xca\x27\x78\x1a\x25\x4e\x4a\x05\x84\x97\xc4\x3f\x9a\xda\xce\x29\x07\xb1\xe0\x57\xbb\x36\x4e\xff\xb9\xd0\x84\x8f\x3e\x11\x66\xe9\x29\xb4\x3e\x49\x62\xef\x52\xe4\x2a\xfb\x87\xdd\xee\x7e\x53\x3d\x0f\xd8\xea\xcb\x14\x36\x05\xeb\xc0\xb9\x62\x2f\x17\x3d\x7a\x55\xeb\x29\x50\xd1\x5a\xfd\x58\xa8\x8d\x22\x43\xdd\x1b\xf3\xdf\x61\xd6\x4b\x2f\x6b\xdc\xbc\x7a\x99\xfb\x1b\x00\x00\xff\xff\x88\x8d\xb3\xda\x31\x01\x00\x00")

func hello_html_bytes() ([]byte, error) {
	return bindata_read(
		_hello_html,
		"hello.html",
	)
}

func hello_html() (*asset, error) {
	bytes, err := hello_html_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "hello.html", size: 305, mode: os.FileMode(420), modTime: time.Unix(1429199371, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _finger_json = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\xce\x41\x0e\x82\x30\x10\x85\xe1\x3d\xa7\x20\x5d\x1b\x9a\x08\x0b\xc3\x65\x48\x53\x47\xa8\x60\x1f\xe9\x4c\x31\xd1\x78\x77\xa1\xa6\x4b\x0c\xfb\xf7\xcd\xfc\xef\xa2\x54\x8e\x39\x52\x50\x6d\xa9\x06\x91\x99\x5b\xad\x27\x58\x33\x0d\x60\x69\x2f\x4d\x53\xab\xd3\x3a\x32\x51\x06\x04\xf7\x32\xe2\xe0\x3b\xf2\xd7\x19\xce\xcb\x3e\xd2\x19\x50\xe2\x82\x91\x0e\xb1\x34\x4c\x24\x32\x05\xe7\x6f\x38\xa2\xb0\x7d\x3b\xeb\xa5\xd6\x59\xa5\x13\x81\x96\x75\xb7\x5b\x6c\xac\x45\xf4\xc2\x55\x0f\xf4\x13\x55\x16\x0f\x8d\x7c\x6b\xb3\xe3\x2f\xfe\xfe\x1c\xb9\x8b\xc1\xfd\x09\xb0\x14\x84\x55\xf1\x29\xbe\x01\x00\x00\xff\xff\x52\x82\x43\x96\x51\x01\x00\x00")

func finger_json_bytes() ([]byte, error) {
	return bindata_read(
		_finger_json,
		"finger.json",
	)
}

func finger_json() (*asset, error) {
	bytes, err := finger_json_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "finger.json", size: 337, mode: os.FileMode(420), modTime: time.Unix(1439732108, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if (err != nil) {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"hello.html": hello_html,
	"finger.json": finger_json,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() (*asset, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"finger.json": &_bintree_t{finger_json, map[string]*_bintree_t{
	}},
	"hello.html": &_bintree_t{hello_html, map[string]*_bintree_t{
	}},
}}

// Restore an asset under the given directory
func RestoreAsset(dir, name string) error {
        data, err := Asset(name)
        if err != nil {
                return err
        }
        info, err := AssetInfo(name)
        if err != nil {
                return err
        }
        err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
        if err != nil {
                return err
        }
        err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
        if err != nil {
                return err
        }
        err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
        if err != nil {
                return err
        }
        return nil
}

// Restore assets under the given directory recursively
func RestoreAssets(dir, name string) error {
        children, err := AssetDir(name)
        if err != nil { // File
                return RestoreAsset(dir, name)
        } else { // Dir
                for _, child := range children {
                        err = RestoreAssets(dir, path.Join(name, child))
                        if err != nil {
                                return err
                        }
                }
        }
        return nil
}

func _filePath(dir, name string) string {
        cannonicalName := strings.Replace(name, "\\", "/", -1)
        return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

