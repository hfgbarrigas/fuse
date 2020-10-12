// Package core contains the main functionality to crawl the directory and apply the appropriate patches to files
package core

import (
	"fuse/internal/util"

	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// CrawlResult stores the final result of the crawling process
type CrawlResult struct {
	WithDiffs uint32
	Error     uint32
}

// Crawl traverses the provided directory tree structure looking for text files. It will not follow sym links.
func Crawl(contentDir, targetDir string, concurrency int8) (chan *CrawlResult, error) {
	var queue = make(chan WorkItem, concurrency)

	contentAbs, targetAbs, err := validatePaths(contentDir, targetDir)

	if err != nil {
		return nil, errors.Wrap(err, "Crawling error")
	}

	// start queue worker to consume work items
	result := initWorker(queue)

	// crawl directory tree. This is the queue producer
	err = crawlDirectory(contentAbs, targetAbs, queue)

	if err != nil {
		return nil, errors.Wrap(err, "Crawling error")
	}

	return result, nil
}

// traverse the directory tree and for each valid file put it in the queue to be processed. Once done, close the queue channel.
func crawlDirectory(contentAbs, targetAbs string, queue chan WorkItem) error {
	log.Debug().
		Msg("Crawling: " + contentAbs)

	// the walk is sync, once it's done tell the worker that no more that will be sent
	defer close(queue)

	// walk the directory tree
	err := filepath.Walk(contentAbs, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "Crawling error")
		}

		// skip directories
		if info.IsDir() {
			return nil
		}

		// skip empty files
		if info.Size() == 0 {
			return nil
		}

		isText, err := util.IsTextFile(path)

		if err != nil {
			return errors.Wrap(err, "Crawling error")
		}

		// skip files that are not text based
		if !isText {
			log.Debug().
				Str("filepath", path).
				Msg("Skipping because its not text.")
			return nil
		}

		commonPath := strings.Replace(path, contentAbs, "", -1)
		wiID, err := uuid.NewRandom()

		if err != nil {
			return errors.Wrap(err, "Crawling error")
		}

		queue <- WorkItem{
			OriginalAbsPath: targetAbs + commonPath,
			UpdateAbsPath:   path,
			ID:              wiID.String(),
			CommonPath:      commonPath,
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "Crawling error")
	}

	return nil
}

func initWorker(queue chan WorkItem) chan *CrawlResult {
	results := make(chan WorkItemResult)
	done := make(chan *CrawlResult)
	wg := sync.WaitGroup{}
	toReturn := &CrawlResult{
		WithDiffs: 0,
		Error:     0,
	}

	go func() {
		for {
			wi, hasItems := <-queue
			if hasItems {
				wg.Add(1)
				// process each WI async
				go func(w WorkItem) {
					log.Debug().
						Interface("workItem", w).
						Msg("Processing working item.")

					result := w.ComputeDiffPatch()

					// write the result to the original destination
					if result.Err == nil {
						if err := result.Write(); err != nil {
							result.Err = err // update the err in case the write fails
						}
					}

					results <- result
				}(wi)
			} else {
				// when there are no more WI we need to wait for all results to be computed
				wg.Wait()
				done <- toReturn
			}
		}
	}()

	go func() {
		for res := range results {
			if res.Err != nil {
				toReturn.Error++
				log.Error().
					Stack().
					Err(res.Err).
					Str("WorkItemID", res.WorkItemID).
					Str("OriginalAbsPath", res.OriginalAbsPath).
					Str("UpdateAbsPath", res.UpdateAbsPath).
					Msg("Error processing WI.")
			} else if res.HasDiffs {
				toReturn.WithDiffs++
				log.Info().
					Str("WorkItemID", res.WorkItemID).
					Str("OriginalAbsPath", res.OriginalAbsPath).
					Str("UpdateAbsPath", res.UpdateAbsPath).
					Msg("Successful WI DIFF.")
			} else {
				log.Info().
					Str("WorkItemID", res.WorkItemID).
					Str("OriginalAbsPath", res.OriginalAbsPath).
					Str("UpdateAbsPath", res.UpdateAbsPath).
					Msg("Successful WI.")
			}
			wg.Done()
		}
	}()

	return done
}

func validatePaths(startDir, targetDir string) (startAbs, targetDirAbs string, err error) {
	startAbs, err = filepath.Abs(startDir)

	if err != nil {
		return "", "", err
	}

	targetDirAbs, err = filepath.Abs(targetDir)

	if err != nil {
		return "", "", err
	}

	startDirAStat, err := os.Stat(startAbs)

	if err != nil {
		return "", "", err
	}

	targetDirAStat, err := os.Stat(targetDirAbs)

	if err != nil {
		return "", "", err
	}

	if !startDirAStat.IsDir() {
		return "", "", errors.New("start path is not a directory")
	}

	if !targetDirAStat.IsDir() {
		return "", "", errors.New("target path is not a directory")
	}

	return startAbs, targetDirAbs, nil
}
