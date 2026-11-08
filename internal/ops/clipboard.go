package ops
type Mode int

const(
	Copy Mode= iota
	Cut
)
type Clipboard struct{
	Source string
	Mode   Mode
}
func(c *Clipboard)Set(source string,mode Mode){
c.Source=source
c.Mode =mode
}
func (c *Clipboard) Clear(){
	c.Source=""
}