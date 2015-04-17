package bindings

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

var _assets_pwlogin_html = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\x52\x4d\x6e\xf2\x30\x10\x5d\xc3\x29\x86\x88\xe5\x17\xb2\xe7\x33\x91\xda\xd2\x45\xa5\x56\xb0\xe8\x8f\xba\xaa\x1c\x3c\x25\x11\x89\x9d\xda\x4e\x69\x84\xb8\x4b\xcf\xd2\x93\xd5\x1e\x43\x00\x09\xb5\x5d\xf9\x65\xf2\x7e\x46\xa3\xc7\x06\xd3\xd9\xd5\xfd\xf3\xfc\x1a\x72\x5b\x95\x69\x9f\x85\xa7\xc7\x72\xe4\xc2\xbd\x3d\x66\x0b\x5b\x62\x7a\xab\x96\x85\x84\x75\x61\x73\x98\xd5\x28\x6f\xa6\x2c\x09\x3f\x1c\x35\xd9\x71\x59\xa6\x44\x4b\x1a\x63\xb5\x92\x4b\x30\xb6\x2d\x71\x12\x2d\x54\xa9\xf4\x58\xa3\xf8\x0f\x51\xba\xd9\xc0\xe8\x89\x6b\x09\xdb\x2d\x4b\x02\x8f\x24\xaf\x4a\x57\xc0\x17\xb6\x50\x72\x12\x79\xd2\x25\x37\xd8\xe8\xd2\xf1\x22\xa8\xd0\xe6\x4a\x84\xf9\x1d\x61\x3f\xf6\xba\x1e\x1b\xc4\x31\xac\x10\x6b\xd0\xf8\xd6\xa0\xb1\x50\x73\xcd\x9d\x00\xb5\x81\x38\x26\x8e\x93\x69\x2e\x97\x08\xc3\x15\xb6\xff\x60\xf8\xce\xcb\x06\x61\x3c\x81\xd1\xa3\x47\xc6\x99\x91\x55\x21\xeb\xc6\x82\x6d\x6b\xb7\x74\x5e\x08\x81\x32\x02\xe9\xbc\x28\xd8\x6b\x69\x19\x52\x87\x51\x30\x72\xc3\xaf\xcf\x64\x9f\x84\xd2\x2f\xd7\x27\x43\xcb\x33\x3a\x11\x61\x1d\x80\x43\x22\x7d\x30\xa8\xdd\x05\xc5\xd1\xe8\x24\xdd\xe2\x87\xdd\x67\xbf\x34\x8e\xec\x61\x94\xa4\x07\x8d\x43\xfa\x9c\xf3\x9c\x1b\xb3\x56\x5a\xfc\xe4\x5e\xef\x38\x5d\x42\x37\xf8\x4b\xc2\x89\x95\x69\xb2\xaa\x38\xac\xca\x33\xa5\x6d\x77\xa2\x0b\xfa\x3a\xf6\xfc\xd5\xa0\xf4\x3d\xeb\x0c\xa8\x75\xe7\x97\x72\x60\x7f\x5c\x96\xf8\xf2\x50\x13\x43\x03\x5d\x23\xa9\xc5\xdf\x01\x00\x00\xff\xff\x62\xcb\x23\x5a\xdd\x02\x00\x00")

func assets_pwlogin_html_bytes() ([]byte, error) {
	return bindata_read(
		_assets_pwlogin_html,
		"assets/pwlogin.html",
	)
}

func assets_pwlogin_html() (*asset, error) {
	bytes, err := assets_pwlogin_html_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "assets/pwlogin.html", size: 733, mode: os.FileMode(420), modTime: time.Unix(1429209851, 0)}
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
	"assets/pwlogin.html": assets_pwlogin_html,
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
	"assets": &_bintree_t{nil, map[string]*_bintree_t{
		"pwlogin.html": &_bintree_t{assets_pwlogin_html, map[string]*_bintree_t{
		}},
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

