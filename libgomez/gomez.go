package libgomez

type Compiler struct{}

func (c *Compiler) BuildFile(file string) {
	// convert file to ll
}

func (c *Compiler) BuildDir(directory string) {
}

type Template interface {
    Emit()
}

type Module struct {}

type Function struct {}
