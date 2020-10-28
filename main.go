package main

import (
	"fmt"
	"os"

	"github.com/emicklei/proto"
)

func main() {
	reader, _ := os.Open("/Users/haikuotiankong/Desktop/go-code/src/newproject/接口代码自动生成工具/test.proto")
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, _ := parser.Parse()

	//proto.Walk(definition, proto.WithService(handleService))
	proto.Walk(definition, proto.WithMessage(handleMessage))
}

func handleService(s *proto.Service) {
	fmt.Println(s.Name)
	fmt.Println(s.Elements)
	for _, v := range s.Elements {

		switch t := v.(type) {
		case *proto.NormalField:
			fmt.Println(t.Field.Name)
			fmt.Println(t.Field.Type)
		case *proto.MapField:
			fmt.Println(t.Name)

			// todo: other field

		}
	}

}

func handleMessage(s *proto.Message) {
	//fmt.Println(s.Position)
	//fmt.Println(s.Comment)
	fmt.Println(s.Name)
	fmt.Println(s.Elements)
	for _, v := range s.Elements {

		switch t := v.(type) {
		case *proto.NormalField:
			fmt.Println(t.Field.Name, t.Field.Type, t.Repeated, t.Options)

			for _, vv := range t.Options {
				fmt.Println(vv.Name)
			}

		case *proto.MapField:
			fmt.Println(t.Name, t.KeyType, t.Type)

			// todo: other field
		case *proto.EnumField:
			fmt.Println(t.Name, t.Elements)

		case *proto.OneOfField:
			fmt.Println(t.Name)
		}
	}
	//fmt.Println(s.Parent  )
}

//func (p *PParser) handleMessage(m *protoparser.Message) {
//	mes := &Message{
//		Name:   m.Name,
//		Fields: make([]Field, 0),
//	}
//
//	for _, v := range m.Elements {
//		switch t := v.(type) {
//		case *protoparser.NormalField:
//			f := getNormalField(t)
//			mes.Fields = append(mes.Fields, f)
//		case *protoparser.MapField:
//			f := getMapField(t)
//			mes.Fields = append(mes.Fields, f)
//
//			// todo: other field
//
//		}
//	}
//
//	p.protos.Messages[mes.Name] = mes
//
//	return
//}
