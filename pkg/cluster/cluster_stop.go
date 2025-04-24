package cluster

import (
	"fmt"
	"github.com/bitfield/script"
	"os"
	"path/filepath"
)

func (impl *clusterImpl) Stop() error {
	return impl.stopConsul()
}

func (impl *clusterImpl) stopConsul() error {
	output, err := script.Exec(`consul leave`).String()
	if err != nil {
		fmt.Println("leave consul failed, output: ", output)
	}
	if impl.needClean {
		if err = os.RemoveAll(filepath.Join(os.TempDir(), "consul")); err != nil {
			return err
		}
	}
	return nil
}
