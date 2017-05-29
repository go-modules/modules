// Package file provides an inject.Injector for file input.
// The tag value is used as the filename.
// Values of type *os.File will be set normally.
// Other values will be depending on the file type.
// File type is derived from the file extension, or optionally overridden by a tag option.
// Supported types include: txt, json, xml, gob
package file

import (
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"github.com/go-modules/modules/inject"
	"github.com/go-modules/modules/tags"
)

// Injector is an inject.Injector for file input.
var Injector = inject.InjectorFunc(Inject)

// Inject opens opens a file and sets the value via literal.Injector.
func Inject(value reflect.Value, fileName string) (bool, error) {
	if typeOfFile.AssignableTo(value.Type()) || typeOfFile.ConvertibleTo(value.Type()) {
		file, err := os.Open(fileName)
		if err != nil {
			return false, err
		}
		value.Elem().Set(reflect.ValueOf(file).Elem())
		return true, nil
	} else {
		tag, optionalType := tags.ParseTag(fileName)
		file, err := os.Open(tag)
		if err != nil {
			return false, err
		}
		fileType := string(optionalType)
		if fileType == "" {
			fileType = filepath.Ext(tag)[1:]
		}
		if fileType == "" {
			return false, fmt.Errorf("no extension or type option given for file: %s", tag)
		}
		switch fileType {
		case "txt":
			if bytes, err := ioutil.ReadAll(file); err != nil {
				return false, err
			} else {
				value.SetString(string(bytes))
			}
		case "json":
			if err := json.NewDecoder(file).Decode(value.Interface()); err != nil {
				return false, err
			}
		case "xml":
			if err := xml.NewDecoder(file).Decode(value.Interface()); err != nil {
				return false, err
			}
		case "gob":
			if err := gob.NewDecoder(file).Decode(value.Interface()); err != nil {
				return false, err
			}
		default:
			return false, fmt.Errorf("unable to read file %s, unrecognized file type %s", tag, fileType)
		}
		return true, nil
	}
}

var typeOfFile = reflect.TypeOf(new(os.File))
