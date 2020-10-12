// Package core contains the main functionality to crawl the directory and apply the appropriate patches to files
package core

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sergi/go-diff/diffmatchpatch"

	"io/ioutil"
	"os"
	"strings"
)

var touchedByFuse = " managed by fuse"

// WorkItem represents an intended pair of content to be patched. An original absolute destination pointing to the initial version
// and an absolute path pointing to the updated content to be patched in the original destination.
type WorkItem struct {
	OriginalAbsPath  string
	UpdateAbsPath    string
	CommonPath       string
	ID               string
	CommentDelimiter string
}

// WorkItemResult represents the result of a WorkItem.
// It contains the Diff and Patches between the original and the update.
// If an error occurred processing the associated work item, err will contain the error.
type WorkItemResult struct {
	WorkItemID      string
	OriginalAbsPath string
	UpdateAbsPath   string
	ResultText      string
	Err             error
	HasDiffs        bool
}

// Write serializes the WorkItemResult content at wr.OriginalAbsPath
func (wr *WorkItemResult) Write() (err error) {
	path := wr.OriginalAbsPath[:strings.LastIndex(wr.OriginalAbsPath, "/")]

	// check if the directory path of the file exists. Create it if it doesn't exist.
	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "Unable to create folder: "+path)
		}
	}

	// write only mode and truncate the file to 0 bytes before writing
	f, err := os.OpenFile(wr.OriginalAbsPath, os.O_TRUNC|os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)

	if err != nil {
		return errors.Wrap(err, "Write result error")
	}

	n, err := f.Write([]byte(wr.ResultText))

	if err != nil {
		return errors.Wrap(err, "Write result error")
	}

	// make sure the file gets written
	err = f.Sync()

	if err != nil {
		return errors.Wrap(err, "Flush result error")
	}

	err = f.Close()

	if err != nil {
		return errors.Wrap(err, "Close fd on result write")
	}

	log.Debug().
		Int("totalBytes", n).
		Str("file", wr.OriginalAbsPath).
		Msg("Successfully written patch.")

	return nil
}

// ComputeDiffPatch computes the diff and all Patches and return the appropriate result
func (w *WorkItem) ComputeDiffPatch() (result WorkItemResult) {
	//TODO: improve this approach to contemplate big files. Chunk by chunk read and compare

	// 1. read update content
	updateContent, err := ioutil.ReadFile(w.UpdateAbsPath)
	if err != nil {
		return errorResult(w, err)
	}
	decoratedContent := decorateUpdateContent(string(updateContent), w.CommentDelimiter)

	// 2. if the file doesn't exist in the target repo don't compute anything
	_, err = os.Stat(w.OriginalAbsPath)
	if os.IsNotExist(err) {
		return WorkItemResult{
			WorkItemID:      w.ID,
			OriginalAbsPath: w.OriginalAbsPath,
			UpdateAbsPath:   w.UpdateAbsPath,
			ResultText:      decoratedContent,
			Err:             nil,
			HasDiffs:        true,
		}
	}
	// 3. Read the original file
	originalBytes, err := ioutil.ReadFile(w.OriginalAbsPath)
	originalContent := string(originalBytes)

	if err != nil {
		return errorResult(w, err)
	}

	log.Info().
		Str("originalFile", w.OriginalAbsPath).
		Str("updateFile", w.UpdateAbsPath).
		Interface("originalContent", originalContent).
		Interface("decoratedContent", decoratedContent).
		Msg("About to diff and patch")

	// 4. compute diffs and patches
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(originalContent, decoratedContent, true)
	patches := dmp.PatchMake(diffs)
	patchText := dmp.PatchToText(patches)
	patchResult, patchesResult := dmp.PatchApply(patches, originalContent)

	log.Debug().
		Interface("results", patchesResult).
		Interface("diffs", diffs).
		Str("patchText", patchText).
		Str("patchResult", patchResult).
		Send()

	hasDiffs := false
	for _, diff := range diffs {
		if diff.Type != diffmatchpatch.DiffEqual {
			hasDiffs = true
			break
		}
	}

	return WorkItemResult{
		WorkItemID:      w.ID,
		OriginalAbsPath: w.OriginalAbsPath,
		UpdateAbsPath:   w.UpdateAbsPath,
		ResultText:      patchResult,
		Err:             nil,
		HasDiffs:        hasDiffs,
	}
}

func errorResult(w *WorkItem, err error) WorkItemResult {
	return WorkItemResult{
		Err:             err,
		WorkItemID:      w.ID,
		OriginalAbsPath: w.OriginalAbsPath,
		UpdateAbsPath:   w.UpdateAbsPath,
	}
}

func decorateUpdateContent(content, commentDelimiter string) string {
	return commentDelimiter + touchedByFuse + "\n" + content
}
