package render

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/strongdm/comply/internal/model"
)

func pdf(output string, live bool, errCh chan error, wg *sync.WaitGroup) {
	var pdfWG sync.WaitGroup

	errOutputCh := make(chan error)

	for {
		_, data, err := loadWithStats()
		if err != nil {
			errCh <- errors.Wrap(err, "unable to load data")
			return
		}

		policies, err := model.ReadPolicies()
		if err != nil {
			errCh <- errors.Wrap(err, "unable to read policies")
			return
		}
		for _, policy := range policies {
			renderToFilesystem(&pdfWG, errOutputCh, data, policy, live)
		}

		narratives, err := model.ReadNarratives()
		if err != nil {
			errCh <- errors.Wrap(err, "unable to read narratives")
			return
		}

		for _, narrative := range narratives {
			renderToFilesystem(&pdfWG, errOutputCh, data, narrative, live)
		}

		pdfWG.Wait()

		if !live {
			wg.Done()
			return
		}
		<-subscribe()
	}
}
