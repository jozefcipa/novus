package paths

import (
	"github.com/jozefcipa/novus/internal/logger"
)

func Resolve() {
	resolveNovusDirs()
	resolveSSLCertDirs()
	resolveSudoDirs()

	logger.Debugf("All paths have been resolved.")
}
