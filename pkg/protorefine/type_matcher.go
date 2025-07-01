package protorefine

import (
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/types/descriptorpb"
)

type protobufTypeMatcher struct {
	nestedNames map[string]struct{}
	founds      map[string]desc.Descriptor
}

func newProtobufTypeMatcher() *protobufTypeMatcher {
	return &protobufTypeMatcher{
		nestedNames: make(map[string]struct{}),
		founds:      make(map[string]desc.Descriptor),
	}
}

func (m *protobufTypeMatcher) traverse(pbPkgName string, currentMsgDescriptor *desc.MessageDescriptor) {
	// if curren msg descriptor exists in nested then ignore it
	// because protoc_gen will compile it as expected
	for _, t := range currentMsgDescriptor.GetNestedMessageTypes() {
		m.nestedNames[t.GetFullyQualifiedName()] = struct{}{}
	}

	current := m.getToBeCheckedMessageType(pbPkgName, currentMsgDescriptor)
	if current == nil {
		return
	}

	if _, exists1 := m.nestedNames[current.GetFullyQualifiedName()]; !exists1 {
		if _, exists2 := m.founds[current.GetFullyQualifiedName()]; !exists2 {
			m.founds[current.GetFullyQualifiedName()] = current
		}
	}

	for _, field := range current.GetFields() {
		switch field.GetType() {
		case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
			next := m.getToBeCheckedMessageType(pbPkgName, field.GetMessageType())
			if next == nil {
				continue
			}

			// avoid recursive definition
			// message AreaNode{
			//  ...
			//  AreaNode children = 1;
			//  ...
			// }
			if _, exists := m.founds[next.GetFullyQualifiedName()]; !exists && current.GetFullyQualifiedName() != next.GetFullyQualifiedName() {
				m.traverse(pbPkgName, next)
			}
		case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
			v := field.GetEnumType()
			if _, exists := m.founds[v.GetFullyQualifiedName()]; !exists {
				m.founds[v.GetFullyQualifiedName()] = v
			}
		}

	}
}

func (m *protobufTypeMatcher) getToBeCheckedMessageType(pbPkgName string, d *desc.MessageDescriptor) *desc.MessageDescriptor {
	// if it is not the same package, ignore it
	if d.GetFile().GetPackage() != pbPkgName {
		return nil
	}

	// if it is not map, return it
	if !d.IsMapEntry() {
		return d
	}

	// if it is map, return the value message type
	mapValue := d.GetFields()[1]
	if d.GetFields()[1].GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
		return mapValue.GetMessageType()
	}
	return nil
}
