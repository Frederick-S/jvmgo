package heap

import (
	"strings"

	"github.com/Frederick-S/jvmgo/classfile"
)

type Class struct {
	accessFlags             uint16
	name                    string
	superClassName          string
	interfaceNames          []string
	constantPool            *ConstantPool
	fields                  []*Field
	methods                 []*Method
	sourceFileName          string
	classLoader             *ClassLoader
	superClass              *Class
	interfaces              []*Class
	instanceVariablesCount  uint
	staticVariablesCount    uint
	staticVariables         Variables
	isInitializationStarted bool
	javaClass               *Object
}

func newClass(classFile *classfile.ClassFile) *Class {
	class := &Class{}
	class.accessFlags = classFile.GetAccessFlags()
	class.name = classFile.GetClassName()
	class.superClassName = classFile.GetSuperClassName()
	class.interfaceNames = classFile.GetInterfaceNames()
	class.constantPool = newConstantPool(class, classFile.GetConstantPool())
	class.fields = newFields(class, classFile.GetFields())
	class.methods = newMethods(class, classFile.GetMethods())
	class.sourceFileName = getSourceFileName(classFile)

	return class
}

func getSourceFileName(classFile *classfile.ClassFile) string {
	sourceFileAttribute := classFile.GetSourceFileAttribute()

	if sourceFileAttribute != nil {
		return sourceFileAttribute.GetFileName()
	}

	return "Unknown"
}

func (class *Class) GetName() string {
	return class.name
}

func (class *Class) GetConstantPool() *ConstantPool {
	return class.constantPool
}

func (class *Class) GetSuperClass() *Class {
	return class.superClass
}

func (class *Class) GetSourceFileName() string {
	return class.sourceFileName
}

func (class *Class) GetClassLoader() *ClassLoader {
	return class.classLoader
}

func (class *Class) GetStaticVariables() Variables {
	return class.staticVariables
}

func (class *Class) IsInitializationStarted() bool {
	return class.isInitializationStarted
}

func (class *Class) GetJavaClass() *Object {
	return class.javaClass
}

func (class *Class) GetJavaName() string {
	return strings.Replace(class.name, "/", ".", -1)
}

func (class *Class) IsPublic() bool {
	return class.accessFlags&ACC_PUBLIC != 0
}

func (class *Class) IsFinal() bool {
	return class.accessFlags&ACC_FINAL != 0
}

func (class *Class) IsSuper() bool {
	return class.accessFlags&ACC_SUPER != 0
}

func (class *Class) IsInterface() bool {
	return class.accessFlags&ACC_INTERFACE != 0
}

func (class *Class) IsAbstract() bool {
	return class.accessFlags&ACC_ABSTRACT != 0
}

func (class *Class) IsSynthetic() bool {
	return class.accessFlags&ACC_SYNTHETIC != 0
}

func (class *Class) IsAnnotation() bool {
	return class.accessFlags&ACC_ANNOTATION != 0
}

func (class *Class) IsEnum() bool {
	return class.accessFlags&ACC_ENUM != 0
}

func (class *Class) IsAccessibleTo(otherClass *Class) bool {
	return class.IsPublic() || class.GetPackageName() == otherClass.GetPackageName()
}

func (class *Class) IsAssignableFrom(otherClass *Class) bool {
	if otherClass == class {
		return true
	}

	if !otherClass.IsArray() {
		if !otherClass.IsInterface() {
			if !class.IsInterface() {
				return otherClass.IsSubClassOf(class)
			}

			return otherClass.IsImplementsFrom(class)
		}

		if !class.IsInterface() {
			return class.IsJavaObjectClass()
		}

		return class.IsSuperInterfaceOf(otherClass)
	}

	if !class.IsArray() {
		if !class.IsInterface() {
			return class.IsJavaObjectClass()
		}

		return class.IsJavaCloneableClass() || class.IsJavaSerializableClass()
	}

	return otherClass.GetArrayElementClass() == class.GetArrayElementClass() ||
		class.GetArrayElementClass().IsAssignableFrom(otherClass.GetArrayElementClass())
}

func (class *Class) IsSubClassOf(otherClass *Class) bool {
	for currentClass := class.superClass; currentClass != nil; currentClass = currentClass.superClass {
		if currentClass == otherClass {
			return true
		}
	}

	return false
}

func (class *Class) IsImplementsFrom(otherInterface *Class) bool {
	for currentClass := class; currentClass != nil; currentClass = currentClass.superClass {
		for _, interfaceMember := range currentClass.interfaces {
			if interfaceMember == otherInterface || interfaceMember.IsSubInterfaceOf(otherInterface) {
				return true
			}
		}
	}

	return false
}

func (class *Class) IsSubInterfaceOf(otherInterface *Class) bool {
	for _, superInterface := range class.interfaces {
		if superInterface == otherInterface || superInterface.IsSubInterfaceOf(otherInterface) {
			return true
		}
	}

	return false
}

func (class *Class) IsSuperClassOf(otherClass *Class) bool {
	return otherClass.IsSubClassOf(class)
}

func (class *Class) IsSuperInterfaceOf(otherInterface *Class) bool {
	return otherInterface.IsSubInterfaceOf(class)
}

func (class *Class) IsJavaObjectClass() bool {
	return class.name == "java/lang/Object"
}

func (class *Class) IsJavaCloneableClass() bool {
	return class.name == "java/lang/Cloneable"
}

func (class *Class) IsJavaSerializableClass() bool {
	return class.name == "java/io/Serializable"
}

func (class *Class) IsPrimitive() bool {
	_, ok := primitiveTypes[class.name]

	return ok
}

func (class *Class) GetPackageName() string {
	index := strings.LastIndex(class.name, "/")

	if index >= 0 {
		return class.name[:index]
	}

	return ""
}

func (class *Class) GetMainMethod() *Method {
	return class.GetStaticMethod("main", "([Ljava/lang/String;)V")
}

func (class *Class) GetStaticMethod(methodName, descriptor string) *Method {
	for _, method := range class.methods {
		if method.IsStatic() && method.name == methodName && method.descriptor == descriptor {
			return method
		}
	}

	return nil
}

func (class *Class) GetInstanceMethod(methodName, methodDescriptor string) *Method {
	return class.GetMethod(methodName, methodDescriptor, false)
}

func (class *Class) GetMethod(name, descriptor string, isStatic bool) *Method {
	for currentClass := class; currentClass != nil; currentClass = currentClass.superClass {
		for _, method := range currentClass.methods {
			if method.IsStatic() == isStatic && method.name == name && method.descriptor == descriptor {
				return method
			}
		}
	}

	return nil
}

func (class *Class) GetField(name, descriptor string, isStatic bool) *Field {
	for currentClass := class; currentClass != nil; currentClass = currentClass.superClass {
		for _, field := range currentClass.fields {
			if field.IsStatic() == isStatic && field.name == name && field.descriptor == descriptor {
				return field
			}
		}
	}

	return nil
}

func (class *Class) StartInitialization() {
	class.isInitializationStarted = true
}

func (class *Class) GetClassInitializationMethod() *Method {
	return class.GetStaticMethod("<clinit>", "()V")
}

func (class *Class) NewObject() *Object {
	return newObject(class)
}

func (class *Class) GetArrayClass() *Class {
	arrayClassName := getArrayClassName(class.name)

	return class.classLoader.LoadClass(arrayClassName)
}

func (class *Class) GetReferenceVariable(fieldName, fieldDescriptor string) *Object {
	field := class.GetField(fieldName, fieldDescriptor, true)

	return class.staticVariables.GetReferenceValue(field.variableIndex)
}

func (class *Class) SetReferenceVariable(fieldName, fieldDescriptor string, referenceVariable *Object) {
	field := class.GetField(fieldName, fieldDescriptor, true)

	class.staticVariables.SetReferenceValue(field.variableIndex, referenceVariable)
}
