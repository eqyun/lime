package sublime

import (
	"code.google.com/p/log4go"
	"fmt"
	"lime/3rdparty/libs/gopy/lib"
	"lime/backend"
	"lime/backend/primitives"
)

var (
	_ = backend.View{}
	_ = primitives.Region{}
)

var (
	_windowCommandGlueClass = py.Class{
		Name:    "sublime.WindowCommandGlue",
		Pointer: (*WindowCommandGlue)(nil),
	}
	_textCommandGlueClass = py.Class{
		Name:    "sublime.TextCommandGlue",
		Pointer: (*TextCommandGlue)(nil),
	}
	_applicationCommandGlueClass = py.Class{
		Name:    "sublime.ApplicationCommandGlue",
		Pointer: (*ApplicationCommandGlue)(nil),
	}
)

type (
	CommandGlue struct {
		py.BaseObject
		inner py.Object
	}
	WindowCommandGlue struct {
		py.BaseObject
		CommandGlue
	}
	TextCommandGlue struct {
		py.BaseObject
		CommandGlue
	}
	ApplicationCommandGlue struct {
		py.BaseObject
		CommandGlue
	}
)

func (c *CommandGlue) PyInit(args *py.Tuple, kwds *py.Dict) error {
	if args.Size() != 1 {
		return fmt.Errorf("Expected only 1 argument not %d", args.Size())
	}
	if v, err := args.GetItem(0); err != nil {
		return err
	} else {
		c.inner = v
	}
	// TODO: look into ref counting convention
	c.inner.Incref()
	return nil
}

func (c *CommandGlue) CreatePyArgs(args backend.Args) (ret *py.Dict, err error) {
	if r, err := toPython(args); err != nil {
		return nil, err
	} else {
		return r.(*py.Dict), nil
	}
}

func (c *CommandGlue) callBool(name string, args backend.Args) bool {
	if pyargs, err := c.CreatePyArgs(args); err != nil {
		log4go.Error(err)
	} else if r, err := c.CallMethodObjArgs(name, pyargs); err != nil {
		log4go.Error(err)
	} else if r, ok := r.(*py.Bool); ok {
		return r.Bool()
	}
	return true
}

func (c *CommandGlue) IsEnabled(args backend.Args) bool {
	return c.callBool("is_enabled", args)
}

func (c *CommandGlue) IsVisible(args backend.Args) bool {
	return c.callBool("is_visible", args)
}

func (c *CommandGlue) Description(args backend.Args) string {
	if pyargs, err := c.CreatePyArgs(args); err != nil {
		log4go.Error(err)
	} else if r, err := c.CallMethodObjArgs("description", pyargs); err != nil {
		log4go.Error(err)
	} else if r, ok := r.(*py.String); ok {
		return r.String()
	}
	return ""
}

func (c *TextCommandGlue) Run(v *backend.View, e *backend.Edit, args backend.Args) error {
	if pyv, err := toPython(v); err != nil {
		return err
	} else if pye, err := toPython(e); err != nil {
		return err
	} else if pyargs, err := c.CreatePyArgs(args); err != nil {
		return err
	} else if obj, err := c.inner.Base().CallFunctionObjArgs(pyv); err != nil {
		return err
	} else {
		if _, err := obj.Base().CallMethodObjArgs("run_", pye, pyargs); err != nil {
			return err
		}
	}
	return nil
}

func (c *WindowCommandGlue) Run(w *backend.Window, args backend.Args) error {
	if pyw, err := toPython(w); err != nil {
		return err
	} else if pyargs, err := c.CreatePyArgs(args); err != nil {
		return err
	} else if obj, err := c.inner.Base().CallFunctionObjArgs(pyw); err != nil {
		return err
	} else if _, err := obj.Base().CallMethodObjArgs("run", pyargs); err != nil {
		return err
	}
	return nil
}

func (c *ApplicationCommandGlue) Run(args backend.Args) error {
	if pyargs, err := c.CreatePyArgs(args); err != nil {
		return err
	} else if obj, err := c.inner.Base().CallFunctionObjArgs(); err != nil {
		return err
	} else if _, err := obj.Base().CallMethodObjArgs("run", pyargs); err != nil {
		return err
	}
	return nil
}

func (c *ApplicationCommandGlue) IsChecked(args backend.Args) bool {
	return c.callBool("is_checked", args)
}