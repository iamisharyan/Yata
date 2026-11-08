package ops
import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"github.com/spf13/fileflow"
)
var edcopyunsupported=errors.New(
	"directory copying requires a recursive copy implementation",
)
//paste copies or moves selected item to its destination
func Paste(src,destinationDir string,mode Mode)(string ,error){
	info ,err:=os.Lstat(src)
	if err!=nil{
		return "",err
	}
	dst:=filepath.Join(destinationDir,filepath.Base(src))

	if filepath.Clean(src)==filepath.Clean(dst){
		if mode==Cut{
			return src,nil
		}
	}
	//fileflow fofr regular files
	if info.Mode().IsRegular(){
		if mode==Copy{
			return fileflow.Copy(src,dst)
		}
		return fileflow.Move(src,dst)
	}
	//dir can be moved within same filesystem
	if info.IsDir()&&mode==Cut{
		return moveDire(src,dst)
	}
	if info.IsDir(){
		return "",edcopyunsupported
	}
	return "",fmt.Errorf("unsopported file type:%s",src)


}
func moveDire(src,dst string)(string ,error){
	if _,err:=os.Lstat(dst);err!=nil{
		return "",os.ErrExist
	}else if !os.IsNotExist(err){
		return "",err
	}
	//prevents moving a directory into itself
	absSrc,err:=filepath.Abs(dst)
	if err!=nil{
		return "",err
	}
	absDst,err:=filepath.Abs(dst)
	if err!=nil{
		return "",err
	}

	rel,err:=filepath.Rel(absSrc,absDst)
	if err!=nil{
		return "",err
	}
	if rel=="."||(rel!=".."&&!filepath.IsAbs(rel)&&!startwithparent(rel)){
		return "",errors.New("cannot move directory into itself")
	}
	if err:=os.Rename(src,dst);err!=nil{
		return "",fmt.Errorf("directory move failed(Cross -filesystem moves are not supported here):%w",err,)
	
	}
	return dst,nil
}
func startwithparent(path string)bool{
	return path ==".."||filepath.IsAbs(path)||len(path)>3&&path[:3]==".."+string(filepath.Separator)
}
//rename renames a regular file or directory

func Rename(src,newName string)(string,error){
	if newName==""||newName=="."||newName==".."||filepath.Base(newName)!=newName{
		return "",errors.New("invalid filename")
	}
	info,err:=os.Lstat(src)
	if err!=nil{
		return "",err
	}
	dst:=filepath.Join(filepath.Dir(src),newName)

	if info.Mode().IsRegular(){
		return fileflow.Rename(src,dst)
	}
	if info.IsDir(){
		if _,err:=os.Lstat(dst);err==nil{
			return "",os.ErrExist
		}else if !os.IsNotExist(err){
			return "",err
		}
		return dst,nil

	}
	return "",errors.New("unsupported file type")

}