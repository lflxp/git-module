package git

import (
	"fmt"
	"strings"
)

type Entry struct {
	mode       string `json:"mode"`
	typ        string `json:"type"`
	id         string `json:"id"`
	name       string `json:"name"`
	size       string `json:"size"`
	commitid   string `json:"commitid"`
	commitinfo string `json:"commitinfo"`
	ctime      string `json:"ctime"`
	time       string `json:"time"`
	author     string `json:"author"`
}

// git ls-tree -l --full-tree HEAD cmd/|while read filename;do echo "$filename $(git log -1 --format="|%h|%s|%cr|%ai|%an" -- `echo $filename|awk '{print $5}'`)" ;done
// 100644 blob 206c042ada4cb14d8f068c364b4b623a45038425    1799	cmd/bolt.go |c25d9ea|boltapi|10 months ago|2020-03-20 19:17:37 +0800|yckj0834
// 100644|blob|8923a8f720912df0810262bfbe74a0ee65f07d1d|3496|cmd/tty.go |c9e64fd|remove local tty|9 months ago|2020-04-08 03:30:17 +0800|yckj0834
func NewEntry(data string) *Entry {
	tmp := strings.Split(data, "|")
	fmt.Println("data", data, tmp)
	return &Entry{
		mode:       tmp[0],
		typ:        tmp[1],
		id:         tmp[2],
		name:       tmp[3],
		size:       tmp[4],
		commitid:   tmp[5],
		commitinfo: tmp[6],
		ctime:      tmp[7],
		time:       tmp[8],
		author:     tmp[9],
	}
}

type Tree struct {
	dirPath string
	path    string
	branch  string
	entries []*Entry
}

func NewTree(dir, path, branch string) *Tree {
	return &Tree{
		dirPath: dir,
		path:    path,
		branch:  branch,
	}
}

// git ls-tree -l --full-tree HEAD cmd/|while read filename;do echo "$(echo $filename|awk '{print $1"|"$2"|"$3"|"$4"|"$5}') $(git log -1 --format="|%h|%s|%cr|%ai|%an" -- `echo $filename|awk '{print $5}'`)" ;done
func (t *Tree) LsTree() error {
	parseString := "|while read filename;do echo \"$(echo $filename|awk '{print $1\"|\"$2\"|\"$3\"|\"$4\"|\"$5}') $(git log -1 --format=\"|%h|%s|%cr|%ai|%an\" -- `echo $filename|awk '{print $5}'`)\" ;done"
	cmd := NewCommand("ls-tree")
	cmd.AddArgs("-l").AddArgs("--full-tree").AddArgs(t.branch).AddArgs(fmt.Sprintf("%s/", t.path))
	cmd.AddArgs(parseString)

	result, err := cmd.RunCmd(t.dirPath)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if len(result) > 0 {
		if t.entries == nil {
			t.entries = []*Entry{}
		}
	}

	for _, info := range strings.Split(result, "\t") {
		t.entries = append(t.entries, NewEntry(info))
	}

	return nil
}
