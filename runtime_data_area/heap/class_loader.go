package heap

import (
	"fmt"

	"github.com/Frederick-S/jvmgo/classfile"
	"github.com/Frederick-S/jvmgo/classpath"
)

type ClassLoader struct {
	classFinder   *classpath.ClassFinder
	loadedClasses map[string]*Class
}

func NewClassLoader(classFinder *classpath.ClassFinder) *ClassLoader {
	classLoader := &ClassLoader{
		classFinder:   classFinder,
		loadedClasses: make(map[string]*Class),
	}

	classLoader.loadBasicClasses()
	classLoader.loadPrimitiveTypeClasses()

	return classLoader
}

func (classLoader *ClassLoader) loadBasicClasses() {
	javaClassClass := classLoader.LoadClass("java/lang/Class")

	for _, class := range classLoader.loadedClasses {
		if class.javaClass == nil {
			class.javaClass = javaClassClass.NewObject()
			class.javaClass.extraData = class
		}
	}
}

func (classLoader *ClassLoader) loadPrimitiveTypeClasses() {
	for primitiveType, _ := range primitiveTypes {
		classLoader.loadPrimitiveTypeClass(primitiveType)
	}
}

func (classLoader *ClassLoader) loadPrimitiveTypeClass(className string) {
	class := &Class{
		accessFlags:             ACC_PUBLIC,
		name:                    className,
		classLoader:             classLoader,
		isInitializationStarted: true,
	}

	class.javaClass = classLoader.loadedClasses["java/lang/Class"].NewObject()
	class.javaClass.extraData = class

	classLoader.loadedClasses[className] = class
}

func (classLoader *ClassLoader) LoadClass(className string) *Class {
	class, ok := classLoader.loadedClasses[className]

	if ok {
		return class
	}

	if className[0] == '[' {
		class = classLoader.LoadArrayClass(className)
	} else {
		class = classLoader.LoadNonArrayClass(className)
	}

	javaClassClass, ok := classLoader.loadedClasses["java/lang/Class"]

	if ok {
		class.javaClass = javaClassClass.NewObject()
		class.javaClass.extraData = class
	}

	return class
}

func (classLoader *ClassLoader) LoadArrayClass(className string) *Class {
	class := &Class{
		accessFlags:             ACC_PUBLIC,
		name:                    className,
		classLoader:             classLoader,
		isInitializationStarted: true,
		superClass:              classLoader.LoadClass("java/lang/Object"),
		interfaces: []*Class{
			classLoader.LoadClass("java/lang/Cloneable"),
			classLoader.LoadClass("java/io/Serializable"),
		},
	}

	classLoader.loadedClasses[className] = class

	return class
}

func (classLoader *ClassLoader) LoadNonArrayClass(className string) *Class {
	classData, classpathEntry := classLoader.ReadClass(className)
	class := classLoader.DefineClass(classData)

	linkClass(class)

	fmt.Printf("[Loaded class %s from %s]\n", className, classpathEntry)

	return class
}

func (classLoader *ClassLoader) ReadClass(className string) ([]byte, classpath.ClasspathEntry) {
	classData, classpathEntry, err := classLoader.classFinder.ReadClass(className)

	if err != nil {
		panic("java.lang.ClassNotFoundException: " + className)
	}

	return classData, classpathEntry
}

func (classLoader *ClassLoader) DefineClass(classData []byte) *Class {
	class := parseClassData(classData)
	class.classLoader = classLoader

	resolveSuperClass(class)
	resolveInterfaces(class)

	classLoader.loadedClasses[class.name] = class

	return class
}

func parseClassData(classData []byte) *Class {
	classFile, err := classfile.Parse(classData)

	if err != nil {
		panic(err)
	}

	return newClass(classFile)
}

func resolveSuperClass(class *Class) {
	if class.name != "java/lang/Object" {
		class.superClass = class.classLoader.LoadClass(class.superClassName)
	}
}

func resolveInterfaces(class *Class) {
	interfacesCount := len(class.interfaceNames)

	if interfacesCount > 0 {
		class.interfaces = make([]*Class, interfacesCount)

		for i, interfaceName := range class.interfaceNames {
			class.interfaces[i] = class.classLoader.LoadClass(interfaceName)
		}
	}
}

func linkClass(class *Class) {
	verifyClass(class)
	prepareClass(class)
}

func verifyClass(class *Class) {
}

func prepareClass(class *Class) {
	assignInstanceFieldsVariableIndices(class)
	assignStaticFieldsVariableIndices(class)
	initializeStaticFinalVariables(class)
}

func assignInstanceFieldsVariableIndices(class *Class) {
	index := uint(0)

	if class.superClass != nil {
		index = class.superClass.instanceVariablesCount
	}

	for _, field := range class.fields {
		if !field.IsStatic() {
			field.variableIndex = index

			if field.IsLongOrDouble() {
				index += 2
			} else {
				index++
			}
		}
	}

	class.instanceVariablesCount = index
}

func assignStaticFieldsVariableIndices(class *Class) {
	index := uint(0)

	for _, field := range class.fields {
		if field.IsStatic() {
			field.variableIndex = index

			if field.IsLongOrDouble() {
				index += 2
			} else {
				index++
			}
		}
	}

	class.staticVariablesCount = index
}

func initializeStaticFinalVariables(class *Class) {
	class.staticVariables = newVariables(class.staticVariablesCount)

	for _, field := range class.fields {
		if field.IsStatic() && field.IsFinal() {
			initializeStaticFinalVariable(class, field)
		}
	}
}

func initializeStaticFinalVariable(class *Class, field *Field) {
	staticVariables := class.staticVariables
	constantPool := class.constantPool
	constantValueIndex := field.GetConstantValueIndex()
	variableIndex := field.GetVariableIndex()

	if constantValueIndex > 0 {
		switch field.GetDescriptor() {
		case "Z", "B", "C", "S", "I":
			value := constantPool.GetConstant(constantValueIndex).(int32)
			staticVariables.SetIntegerValue(variableIndex, value)
		case "J":
			value := constantPool.GetConstant(constantValueIndex).(int64)
			staticVariables.SetLongValue(variableIndex, value)
		case "F":
			value := constantPool.GetConstant(constantValueIndex).(float32)
			staticVariables.SetFloatValue(variableIndex, value)
		case "D":
			value := constantPool.GetConstant(constantValueIndex).(float64)
			staticVariables.SetDoubleValue(variableIndex, value)
		case "Ljava/lang/String;":
			goString := constantPool.GetConstant(constantValueIndex).(string)
			javaString := ConvertGoStringToJavaString(class.GetClassLoader(), goString)

			staticVariables.SetReferenceValue(variableIndex, javaString)
		}
	}
}
