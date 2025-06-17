package cluster

import (
	"fmt"
	"github.com/bitfield/script"
	"os"
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
		if err = os.RemoveAll(impl.getConsulDataDir()); err != nil {
			return err
		}
	}
	return nil
}
