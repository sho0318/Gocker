package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

type OverlayManager struct {
	LowerDir string
	UpperDir string
	WorkDir  string
	Merged   string
	BaseDir  string
}

func NewOverlayManager(rootfs, runtimePath, id string) *OverlayManager {
	base := filepath.Join(runtimePath, id)
	return &OverlayManager{
		LowerDir: rootfs,
		UpperDir: filepath.Join(base, "upper"),
		WorkDir:  filepath.Join(base, "work"),
		Merged:   filepath.Join(base, "merged"),
		BaseDir:  base,
	}
}

func (m *OverlayManager) MountOverlay() (string, error) {
	os.MkdirAll(m.UpperDir, 0755)
	os.MkdirAll(m.WorkDir, 0755)
	os.MkdirAll(m.Merged, 0755)

	opts := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", m.LowerDir, m.UpperDir, m.WorkDir)
	if err := syscall.Mount("overlay", m.Merged, "overlay", 0, opts); err != nil {
		return "", err
	}
	return m.Merged, nil
}

func (m *OverlayManager) Unmount() error {
    return syscall.Unmount(m.Merged, syscall.MNT_DETACH)
}

func (m *OverlayManager) UnmountAndCleanup() error {
	return os.RemoveAll(m.BaseDir)
}