// +build linux

package ploop

import (
	"fmt"
	"github.com/docker/docker/daemon/graphdriver"
	"github.com/docker/docker/pkg/parsers"
	ploop "github.com/kolyshkin/goploop"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

const (
	mountPoint  = "ploop"
	snapshots   = "snapshots"
	currentData = "data"
)

type Driver struct {
	base       string
	sync.Mutex                   // synchronizes access to idmap
	idmap      map[string]string // docker id -> ploop uuid for all snapshots
	ploop      ploop.Ploop
}

type options struct {
	mode ploop.ImageMode
	size uint64
}

func init() {
	graphdriver.Register("ploop", Init)
}

func Init(base string, opt []string) (graphdriver.Driver, error) {
	var err error
	options, err := parseOptions(opt)

	if err != nil {
		return nil, err
	}

	// Call Create, call Open
	createParam := ploop.CreateParam{
		Size: options.size,
		Mode: options.mode,
		File: base,
	}

	if err = ploop.Create(&createParam); err != nil {
		return nil, err
	}

	p, err := ploop.Open(path.Join(base, "DiskDescriptor.xml"))

	if err != nil {
		return nil, err
	}

	if err = os.MkdirAll(path.Join(base, mountPoint), 0755); err != nil {
		return nil, err
	}

	mountParam := ploop.MountParam{Target: path.Join(base, mountPoint)}

	device, err := p.Mount(&mountParam)

	if err != nil {
		return nil, err
	}

	return Driver{base, sync.Mutex{}, make(map[string]string), p}, nil
}

func parseOptions(opt []string) (options, error) {
	var options options

	for _, option := range opt {
		key, value, err := parsers.ParseKeyValueOpt(option)
		if err != nil {
			return options, err
		}

		key = strings.ToLower(key)

		switch key {
		case "ploop.size":
			size, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return options, err
			}
			options.size = size
		case "ploop.mode":
			var mode ploop.ImageMode
			switch value {
			case "expanded":
				mode = ploop.Expanded
			case "preallocated":
				mode = ploop.Preallocated
			case "raw":
				mode = ploop.Raw
			default:
				return options, fmt.Errorf("Unknown ploop.mode value %s", value)
			}
			options.mode = mode
		default:
			return options, fmt.Errorf("Unknown option %s", key)
		}
	}

	return options, nil
}

func (d *Driver) String() string {
	return "Ploop Driver"
}

func (d *Driver) Create(id, parent string) error {
	var err error

	if parent == "" {

	} else {
		if !d.Exists(parent) {
			return fmt.Errorf("No such ploop snapshot with id %s", parent)
		}
	}

	return err
}

func (d *Driver) Remove(id string) error {

}

func (d *Driver) Get(id, mountLabel string) (dir string, err error) {
	// Use MountParam with uuid field set to snapshot id.
}

func (d *Driver) Put(id string) error {

}

func (d *Driver) Exists(id string) bool {

}

func (d *Driver) Status() [][2]string {

}

func (d *Driver) GetMetadata(id string) (map[string]string, error) {

}

func (d *Driver) Cleanup() error {

}

func (d *Driver) Diff(id, parent string) (archive.Archive, error) {

}

func (d *Driver) Changes(id, parent string) ([]archive.Change, error) {

}

func (d *Driver) ApplyDiff(id, parent string, diff archive.ArchiveReader) (size int64, err error) {

}

func (d *Driver) DiffSize(id, parent string) (size int64, err error) {

}
