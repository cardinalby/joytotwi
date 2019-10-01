package fill

import log "github.com/sirupsen/logrus"

type fillStats struct {
	postedCount, errorsCount, existentCount int
}

func (stats *fillStats) print() {
	log.Infof("Posted: %d", stats.postedCount)
	log.Infof("Errors: %d", stats.errorsCount)
	log.Infof("Exists: %d", stats.existentCount)
}