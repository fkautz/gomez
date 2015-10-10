package gomez

import (
	"errors"
)

type SymbolTable struct {
	frames []*frame
}

type frame struct {
	name      string
	variables map[string]value
}

type value struct {
	name         string
	valueType    []string
	internalName string
}

func (vt *SymbolTable) PushFrame() {
	newFrame := &frame{}
	newFrame.variables = make(map[string]value)
	vt.frames = append(vt.frames, newFrame)
}

func (vt *SymbolTable) PopFrame() {
	if len(vt.frames) > 0 {
		vt.frames = vt.frames[0 : len(vt.frames)-1]
	}
}

func (vt *SymbolTable) AddSymbol(name string, valueType []string, internalName string) string {
	curFrame := vt.frames[len(vt.frames)-1]
	newValueType := make([]string, len(valueType))
	copy(newValueType, valueType)
	v := value{
		name:      name,
		valueType: newValueType,
		internalName: internalName,
	}
	curFrame.variables[name] = v
	return name
}

func (vt *SymbolTable) FindVariable(name string) (string, []string, string, error) {
	for i := len(vt.frames) - 1; i >= 0; i-- {
		if variableFrame, ok := vt.frames[i].variables[name]; ok {
			vfName := variableFrame.name
			vfValueType := variableFrame.valueType
			vfInternalName := variableFrame.internalName
			return vfName, vfValueType, vfInternalName, nil
		}
	}
	return "", nil, "", errors.New("Unable to find variable: " + name)
}
