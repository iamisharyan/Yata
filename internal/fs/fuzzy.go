package fs

import (
	"os"
	"path/filepath"
	"strings"
	"github.com/sahilm/fuzzy"
)
type FuzzyItem struct{
	Relpath string
	IsDir   bool

}
func CollectFiles(root string,showhidden bool)[]FuzzyItem{
	var items[]FuzzyItem
	_=filepath.WalkDir(root,func(path string,d os.DirEntry,err error)error{
		if err!=nil{
			return nil
		}
		if path==root{
			return nil
		}
		name:=d.Name()
		if !showhidden&&strings.HasPrefix(name,"."){
		if d.IsDir(){
			return filepath.SkipDir
		}
		return nil
	}
rel,err:=filepath.Rel(root,path)
if err!=nil{
	rel=path
}
items=append(items,FuzzyItem{
	Relpath:rel,
	IsDir:d.IsDir(),
})

	return nil
})
return items	
}
//filtering on query
func FuzzySearch(query string,items[]FuzzyItem)[]FuzzyItem{
	if strings.TrimSpace(query)==""{
		return items
	}
	target:=make([]string,len(items))
	for i,it:=range items{
		target[i]=it.Relpath
	}
	matches:=fuzzy.Find(query,target)
	var results []FuzzyItem
	for _, m:=range matches{
		results=append(results,items[m.Index])
	}
	return results
}

