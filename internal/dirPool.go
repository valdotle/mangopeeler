package internal

import (
	"github.com/vitali-fedulov/images4"
)

type DirPool struct {
	threads   chan any
	completed chan *DirEntryResponse
}

func NewDirPool(threads, entries int) DirPool {
	p := DirPool{make(chan any, threads), make(chan *DirEntryResponse, entries)}
	for i := 0; i < threads; i++ {
		p.threads <- nil
	}

	return p
}

func (d DirPool) Add() {
	<-d.threads
}

func (d DirPool) RemoveWithResult(filename string, icon images4.IconT) {
	d.threads <- nil
	d.completed <- &DirEntryResponse{filename, icon}
}

func (d DirPool) Remove() {
	d.threads <- nil
	d.completed <- nil
}

func (d DirPool) Wait() {
	for i := 1; i < cap(d.completed); i++ {
		<-d.completed
	}
}

func (d DirPool) WaitForResults() []DirEntryResponse {
	var results []DirEntryResponse
	for i := 1; i < cap(d.completed); i++ {
		if r := <-d.completed; r != nil {
			results = append(results, *r)
		}
	}

	return results
}

type DirEntryResponse struct {
	Filename string
	Icon     images4.IconT
}
