package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"
)

const (
	template1 = "\t{{name}}{{method}}Factory := create{{Service}}SvcFactory({{service}}endpoint.Make{{Name}}{{method}}Endpoint, {{service}}svcport)\n\t{{name}}{{method}}Endpointer := sd.NewEndpointer({{service}}Instancer, {{name}}{{method}}Factory, logger)\n\t{{name}}{{method}}Balancer := kitlb.NewRoundRobin({{name}}{{method}}Endpointer)\n\t{{name}}{{method}}Retry := kitlb.Retry(retryMax, retryTimeout, {{name}}{{method}}Balancer)\n\tendpoints.{{Name}}{{method}}Endpoint = {{name}}{{method}}Retry\n"
	template2 = "\t{{name}}{{method}}Endpoint := kitgrpc.NewClient(\n\t\tconn,\n\t\t\"xtech.{{service}}.v1.{{Service}}SvcService\",\n\t\t\"{{Name}}{{method}}\",\n\t\tendpoint.Encode{{Name}}{{method}}Request,\n\t\tendpoint.Decode{{Name}}{{method}}Response,\n\t\t{{service}}proto.{{Name}}{{method}}RespProto{},\n\t).Endpoint()\n"
	template3 = "\t\t{{Name}}{{method}}Endpoint: {{name}}{{method}}Endpoint,"
	template4 = "\t{{Name}}{{method}}Endpoint kitendpoint.Endpoint"
	template5 = "\t\t{{Name}}{{method}}Endpoint: Make{{Name}}{{method}}Endpoint(service),"
	template6 = "\t{{name}}{{method}}Handler kitgrpc.Handler"
	template7 = "\t{{name}}{{method}}Handler := kitgrpc.NewServer(\n\t\tendpoints.{{Name}}{{method}}Endpoint,\n\t\tendpoint.Decode{{Name}}{{method}}Request,\n\t\tendpoint.Encode{{Name}}{{method}}Response,\n\t\toptions...,\n\t)\n"
	template8 = "\t\t{{name}}{{method}}Handler: {{name}}{{method}}Handler,"
	template9 = "\t{{Name}}{{method}}(ctx context.Context, reqproto *{{service}}proto.{{Name}}{{method}}ReqProto) (*{{service}}proto.{{Name}}{{method}}RespProto, error)"

	templateA = "package grpc\n\nimport (\n\t\"xtech-kit/proto/{{service}}/v1\"\n)\n\nfunc New{{Name}}{{method}}ReqProto() *{{service}}proto.{{Name}}{{method}}ReqProto {\n\treturn &{{service}}proto.{{Name}}{{method}}ReqProto{\n\t}\n}\n"
	templateB = "package endpoint\n\nimport (\n\t\"context\"\n\n\t\"xtech-kit/proto/{{service}}/v1\"\n\t\"xtech-kit/svc/{{service}}-svc/pkg/{{service}}/service\"\n\n\t\"github.com/go-kit/kit/endpoint\"\n)\n\nfunc Encode{{Name}}{{method}}Request(ctx context.Context, req interface{}) (interface{}, error) {\n\treqProto := req.(*{{service}}proto.{{Name}}{{method}}ReqProto)\n\n\treturn reqProto, nil\n}\n\nfunc Decode{{Name}}{{method}}Request(ctx context.Context, req interface{}) (interface{}, error) {\n\treqProto := req.(*{{service}}proto.{{Name}}{{method}}ReqProto)\n\n\treturn reqProto, nil\n}\n\nfunc Encode{{Name}}{{method}}Response(ctx context.Context, resp interface{}) (interface{}, error) {\n\trespProto := resp.(*{{service}}proto.{{Name}}{{method}}RespProto)\n\n\treturn respProto, nil\n}\n\nfunc Decode{{Name}}{{method}}Response(ctx context.Context, resp interface{}) (interface{}, error) {\n\trespProto := resp.(*{{service}}proto.{{Name}}{{method}}RespProto)\n\n\treturn respProto, nil\n}\n\nfunc Make{{Name}}{{method}}Endpoint(service service.{{Service}}SvcService) endpoint.Endpoint {\n\treturn func(ctx context.Context, req interface{}) (interface{}, error) {\n\t\treqProto := req.(*{{service}}proto.{{Name}}{{method}}ReqProto)\n\n\t\trespProto, err := service.{{Name}}{{method}}(ctx, reqProto)\n\n\t\tif err != nil {\n\t\t\treturn nil, err\n\t\t}\n\n\t\treturn respProto, nil\n\t}\n}\n\n// for grpc client\nfunc (u {{Service}}SvcEndpoints) {{Name}}{{method}}(ctx context.Context, reqproto *{{service}}proto.{{Name}}{{method}}ReqProto) (*{{service}}proto.{{Name}}{{method}}RespProto, error) {\n\tresp, err := u.{{Name}}{{method}}Endpoint(ctx, reqproto)\n\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\n\treturn resp.(*{{service}}proto.{{Name}}{{method}}RespProto), nil\n}\n"
	templateC = "package grpc\n\nimport (\n\t\"context\"\n\n\t\"xtech-kit/proto/{{service}}/v1\"\n)\n\nfunc (p *grpcServer) {{Name}}{{method}}(ctx context.Context, req *{{service}}proto.{{Name}}{{method}}ReqProto) (*{{service}}proto.{{Name}}{{method}}RespProto, error) {\n\t_, resp, err := p.{{name}}{{method}}Handler.ServeGRPC(ctx, req)\n\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\n\trespProto := resp.(*{{service}}proto.{{Name}}{{method}}RespProto)\n\n\treturn respProto, nil\n}\n"
	templateD = "package service\n\nimport (\n\t\"context\"\n\n\t\"xtech-kit/const\"\n\t\"xtech-kit/error\"\n\t\"xtech-kit/proto/{{service}}/v1\"\n\tdbrepository \"xtech-kit/svc/{{service}}-svc/pkg/{{service}}/repository/database\"\n\n\t\"github.com/xianyu-tech/tlog\"\n\t\"github.com/xianyu-tech/util/log\"\n)\n\nfunc (p *basic{{Service}}SvcService) {{Name}}{{method}}(ctx context.Context, reqproto *{{service}}proto.{{Name}}{{method}}ReqProto) (*{{service}}proto.{{Name}}{{method}}RespProto, error) {\n\tdefer func() {\n\t\tif rec := recover(); rec != nil {\n\t\t\terrMsg := tlog.Error(\"{{service}} auth fetch (%d, %s) err (recover %v).\", reqproto.Get{{Service}}Id(), reqproto.GetIdentifier(), logutil.PrintPanic())\n\n\t\t\ttlog.AsyncSend(tlog.NewRawLogError(comconst.ErrPanicRecovery(\"{{service}} auth fetch\"), errMsg))\n\t\t}\n\t}()\n\n\n\t{{name}}{{method}}RespProto := &{{service}}proto.{{Name}}{{method}}RespProto{\n\t\tErrCode: errconst.COMMON_API_ERROR_OK,\n\t}\n\treturn {{name}}{{method}}RespProto, nil\n}\n"
)

var (
	templateMap = map[string]string{
		"template1": template1,
		"template2": template2,
		"template3": template3,
		"template4": template4,
		"template5": template5,
		"template6": template6,
		"template7": template7,
		"template8": template8,
		"template9": template9,
	}
)

func main() {

	service := "social"
	name := "momentPostsByUser"
	method := "Fetch"
	filename := "moment_posts_by_user_fetch.go"

	filePathPre := "/Users/haikuotiankong/Desktop/go-code/src/xtech-kit/svc/" + service + "-svc/"

	Service := Ucfirst(service)
	Name := Ucfirst(name)
	clientEndpoint := filePathPre + "client/endpoint/endpoint.go"
	clientGrpcHandler := filePathPre + "client/grpc/handler.go"
	pkgEndpoint := filePathPre + "pkg/" + service + "/endpoint/endpoint.go"
	pkgGrpcHandler := filePathPre + "pkg/" + service + "/grpc/handler.go"
	pkgService := filePathPre + "pkg/" + service + "/service/service.go"

	var pathArry = []string{clientEndpoint, clientGrpcHandler, pkgEndpoint, pkgGrpcHandler, pkgService}

	for _, path := range pathArry {
		err := handlerFile(path, name, Name, method, service, Service)
		if err != nil {
			panic(err)
		}
	}

	newClientGrpc := filePathPre + "client/grpc/" + filename
	createFile(newClientGrpc, templateA, name, Name, method, service, Service)

	newPkgEndpoint := filePathPre + "pkg/" + service + "/endpoint/" + filename
	createFile(newPkgEndpoint, templateB, name, Name, method, service, Service)

	newPkgGrpcHandler := filePathPre + "pkg/" + service + "/grpc/" + filename
	createFile(newPkgGrpcHandler, templateC, name, Name, method, service, Service)

	newPkgService := filePathPre + "pkg/" + service + "/service/" + filename
	createFile(newPkgService, templateD, name, Name, method, service, Service)

}
func createFile(filename, template, name, Name, method, service, Service string) {
	f, err1 := os.Create(filename)
	if err1 != nil {
		log.Printf("Cannot create text file: %s, err: [%v]", filename, err1)
		return
	}

	s := strings.ReplaceAll(template, "{{name}}", name)
	s = strings.ReplaceAll(s, "{{Name}}", Name)
	s = strings.ReplaceAll(s, "{{method}}", method)
	s = strings.ReplaceAll(s, "{{service}}", service)
	s = strings.ReplaceAll(s, "{{Service}}", Service)

	_, err1 = io.WriteString(f, s)
	if err1 != nil {
		log.Printf("Cannot write text file: %s, err: [%v]", filename, err1)
		return
	}
}

func handlerFile(filepath, name, Name, method, service, Service string) error {
	file, err := os.Open(filepath)
	if err != nil {
		log.Printf("Cannot open text file: %s, err: [%v]", filepath, err)
		return err
	}
	defer file.Close()

	var content string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text() // or

		strArry := regexp.MustCompile(`{{(template\d+)}}`).FindStringSubmatch(line)

		if len(strArry) == 2 {
			if template, ok := templateMap[strArry[1]]; ok {

				s := strings.ReplaceAll(template, "{{name}}", name)
				s = strings.ReplaceAll(s, "{{Name}}", Name)
				s = strings.ReplaceAll(s, "{{method}}", method)
				s = strings.ReplaceAll(s, "{{service}}", service)
				s = strings.ReplaceAll(s, "{{Service}}", Service)

				content += s + "\n"
			}
		}

		content += line + "\n"
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Cannot scanner text file: %s, err: [%v]", filepath, err)
		return err
	}

	writeFile(filepath, []byte(content))

	return nil
}
func writeFile(filepath string, content []byte) {
	err := ioutil.WriteFile(filepath, content, 0666)
	if err != nil {
		log.Fatal(err)
	}
}

func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}
