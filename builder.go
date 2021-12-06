package golangasm

import (
	"fmt"

	"github.com/wdvxdr1123/golang-asm/asm/arch"
	"github.com/wdvxdr1123/golang-asm/obj"
	"github.com/wdvxdr1123/golang-asm/objabi"
)

// Builder allows you to assemble a series of instructions.
type Builder struct {
	ctxt *obj.Link
	arch *arch.Arch

	first *obj.Prog
	last  *obj.Prog

	// bulk allocator.
	block *[]obj.Prog
	used  int
}

// Root returns the first instruction.
func (b *Builder) Root() *obj.Prog {
	return b.first
}

// NewProg returns a new instruction structure.
func (b *Builder) NewProg() *obj.Prog {
	return b.progAlloc()
}

func (b *Builder) progAlloc() *obj.Prog {
	var p *obj.Prog

	if b.used >= len(*b.block) {
		p = b.ctxt.NewProg()
	} else {
		p = &(*b.block)[b.used]
		b.used++
	}

	p.Ctxt = b.ctxt
	return p
}

// AddInstruction adds an instruction to the list of instructions
// to be assembled.
func (b *Builder) AddInstruction(p *obj.Prog) {
	if b.first == nil {
		b.first = p
		b.last = p
	} else {
		b.last.Link = p
		b.last = p
	}
}

// Assemble generates the machine code from the given instructions.
func (b *Builder) Assemble() []byte {
	var fn interface{} = &obj.FuncInfo{
		Text: b.first,
	}
	s := &obj.LSym{
		Extra: &fn,
	}
	b.arch.Assemble(b.ctxt, s, b.progAlloc)
	return s.P
}

// NewBuilder constructs an assembler for the given architecture.
func NewBuilder(archStr string, cacheSize int) (*Builder, error) {
	a := arch.Set(archStr, false)
	ctxt := obj.Linknew(a.LinkArch)
	ctxt.Headtype = objabi.Hlinux
	ctxt.DiagFunc = func(in string, args ...interface{}) {
		fmt.Printf(in+"\n", args...)
	}
	a.Init(ctxt)

	block := make([]obj.Prog, cacheSize)

	return &Builder{
		ctxt:  ctxt,
		arch:  a,
		block: &block,
	}, nil
}
