package render

import (
	"sync"

	"github.com/strongdm/comply/internal/model"
)

func pdf(output string, live bool, wg *sync.WaitGroup) {
	var pdfWG sync.WaitGroup

	for {
		data, err := loadWithStats()
		if err != nil {
			// TODO: errors channel or quit channel or panic?
			panic(err)
		}

		for _, policy := range model.ReadPolicies() {
			renderPolicyToDisk(pdfWG, data, policy)
		}

		for _, narrative := range model.ReadNarratives() {
			renderNarrativeToDisk(pdfWG, data, narrative)
		}

		pdfWG.Wait()

		if !live {
			wg.Done()
			return
		}
		<-subscribe()
	}
}
