package tools

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/pkg/errors"
)

type windowsPlatformImpl struct {
}

func WindowsPlatform() *windowsPlatformImpl {
	return &windowsPlatformImpl{}
}

func (*windowsPlatformImpl) MoveFile(src, dst string) error {
	cmd := fmt.Sprintf(`powershell -Command Get-ChildItem "%s" | Move-Item -Destination "%s"" -Force`, src, dst)
	_, err := script.Exec(cmd).Stdout()
	if err != nil {
		return errors.Wrapf(err, "move file failed, src: %s, dst: %s", src, dst)
	}
	return nil
}

func (*windowsPlatformImpl) Unzip(zipFile, dst string) error {
	cmd := fmt.Sprintf(`powershell -Command Expand-Archive -Path "%s" -DestinationPath "%s" -Force`, zipFile, dst)
	err := script.Exec(cmd).Wait()
	if err != nil {
		return errors.Wrapf(err, "uncompress failed, file: %s", zipFile)
	}
	return err
}

func (*windowsPlatformImpl) UnzipFile(zipFile, sourceDir, targetDir string) error {
	cmd := fmt.Sprintf(`$tempDir = Join-Path $env:TEMP (New-Guid).Guid
		New-Item -ItemType Directory -Path $tempDir | Out-Null
		Expand-Archive -Path "%s" -DestinationPath $tempDir -Force
		$sourcePath = Join-Path $tempDir "%s"
		if (Test-Path $sourcePath) {
			Copy-Item -Path "$sourcePath\*" -Destination "%s" -Recurse -Force
		} else {
			Write-Host "错误: ZIP文件中不存在目录 %s"
			exit 1
		}
		Remove-Item -Path $tempDir -Recurse -Force`,
		zipFile, sourceDir, targetDir, sourceDir)
	err := script.Exec(cmd).Wait()
	if err != nil {
		return errors.Wrapf(err, "解压失败, 文件: %s", zipFile)
	}
	return err
}
