package main

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/format/diff"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func test() {
	r, _ := git.Init(memory.NewStorage(), nil)
	testBranch := &config.Branch{
		Name:   "foo",
		Remote: "origin",
		Merge:  "refs/heads/foo",
	}
	r.CreateBranch(testBranch)

	testBranch = &config.Branch{
		Name:   "foo2",
		Remote: "origin",
		Merge:  "refs/heads/foo2",
	}
	r.CreateBranch(testBranch)

	cfg, _ := r.Config()
	for k := range cfg.Branches {
		println(k)
	}

}

func main() {

	r, e := git.PlainOpen("/Users/gao/Code/utils/")
	if e != nil {
		panic(e)
	}

	// cfg, e := r.Storer.Config()
	// if e != nil {
	// 	panic(e)
	// }

	// for n := range cfg.Branches {
	// 	println(n)
	// }

	t0, e := lookCommit(r, "master")
	if e != nil {
		panic(e)
	}
	t1, e := lookCommit(r, "f1")
	if e != nil {
		panic(e)
	}
	println(t0.Hash.String(), t1.Hash.String())
	changes, e := object.DiffTree(t0, t1)
	if e != nil {
		panic(e)
	}

	for _, d := range changes {
		p, e := d.Patch()
		if e != nil {
			panic(e)
		}
		// f, t, e := d.Files()
		// if e != nil {
		// 	panic(e)
		// }
		for _, fp := range p.FilePatches() {

			hks := diff.NewHunksGenerator(fp.Chunks(), diff.DefaultContextLines).Generate()
			for _, hk := range hks {
				s := hk.UnSafe()
				println(s.ToLine, s.ToCount)
			}

			for _, chk := range fp.Chunks() {
				if chk.Type() == 0 {
					continue
				}
				// println(chk.Type(), chk.Content())
			}
		}
	}

}

func lookCommit(repo *git.Repository, branch string) (*object.Tree, error) {
	bm, e := repo.Branch(branch)
	if e != nil {
		panic(e)
	}

	ref, e := repo.Reference(bm.Merge, true)
	if e != nil {
		panic(e)
	}

	c, e := repo.CommitObject(ref.Hash())
	if e != nil {
		panic(e)
	}

	t, e := c.Tree()
	if e != nil {
		panic(e)
	}
	return t, e

}
