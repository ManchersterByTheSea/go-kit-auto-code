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
	template1 = "\t{{Name}}{{method}}Endpoint kitendpoint.Endpoint"
	template2 = "\t\t{{Name}}{{method}}Endpoint: Make{{Name}}{{method}}Endpoint(service),"
	template3 = "\t{{name}}{{method}}Handler := kithttp.NewServer(\n\t\tendpoints.{{Name}}{{method}}Endpoint,\n\t\tendpoint.Decode{{Name}}{{method}}Request,\n\t\tendpoint.Encode{{Name}}{{method}}Response,\n\t\toptions...)\n\trouter.Handle(\"{{router}}\", {{name}}{{method}}Handler).Methods(\"{{METHOD}}\")\n"
	template4 = "\t{{Name}}{{method}}(ctx context.Context, reqdata *data.{{Name}}{{method}}Request) (*data.{{Name}}{{method}}Response, error)"

	templateA      = "package data\n\nimport (\n\t\"net/http\"\n)\n\ntype {{Name}}{{method}}Request struct {\n\tRequest *http.Request `json:\"-\"`\n\n\t// 用户id\n\t// required: true\n\tUserId int64 `json:\"user_id\"`\n\t// 防跨域secret\n\t// required: true\n\tCsrfSecret string `json:\"-\"`\n\t// accessToken\n\t// required: true\n\tAccessToken string `json:\"-\"`\n}\n\ntype {{Name}}{{method}}Response struct {\n\t// 错误码\n\tErrCode string `json:\"error\"`\n\t// 错误信息\n\tErrMsg string `json:\"message\"`\n\n\tData *{{Name}}Data\n}\n\ntype {{Name}}Data struct {\n\tVoteId    int64 `json:\"vote_id\"`\n\tResultId  int64 `json:\"result_id\"`\n\tUserId    int64 `json:\"user_id\"`\n\tTeamId    int64 `json:\"team_id\"`\n\tCreatedAt int64 `json:\"created_at\"`\n}\n"
	templateB_get  = "package endpoint\n\nimport (\n\t\"context\"\n\t\"net/http\"\n\t\"strconv\"\n\t\"strings\"\n\n\t\"xtech-kit/const\"\n\t\"xtech-kit/error\"\n\t\"xtech-kit/gw/data\"\n\t\"xtech-kit/gw/{{gw}}-gw/pkg/{{gw}}/data\"\n\t\"xtech-kit/gw/{{gw}}-gw/pkg/{{gw}}/service\"\n\n\t\"github.com/go-kit/kit/endpoint\"\n\t\"github.com/json-iterator/go\"\n\t\"github.com/xianyu-tech/tlog\"\n)\n\nfunc Decode{{Name}}{{Method}}Request(ctx context.Context, r *http.Request) (interface{}, error) {\n\t{{name}}{{Method}}Request := &data.{{Name}}{{Method}}Request{}\n\n\tparams := r.URL.Query()\n\n\tvalStaffIds := params[\"staff_id\"]\n\tvalStaffId := \"\"\n\n\tif len(valStaffIds) > 0 {\n\t\tvalStaffId = valStaffIds[0]\n\t}\n\n\tif valStaffId != \"\" {\n\t\tsrcStaffId, err := strconv.ParseInt(valStaffId, 10, 64)\n\n\t\tif err != nil {\n\t\t\terrMsg := tlog.Error(\"decode {{name_method}} request (%v) err (parse int %v).\", r, err)\n\n\t\t\ttlog.AsyncSend(tlog.NewRawLogError(err, errMsg).AttachRequest(r))\n\n\t\t\treturn nil, comconst.ErrParamIllegal(\"staff id\")\n\t\t}\n\n\t\tif srcStaffId <= 0 {\n\t\t\terrMsg := tlog.Error(\"decode {{name_method}} request (%v) err (staff id illegal).\", r)\n\n\t\t\ttlog.AsyncSend(tlog.NewRawLogError(comconst.ErrParamIllegal(\"staff id\"), errMsg).AttachRequest(r))\n\n\t\t\treturn nil, comconst.ErrParamIllegal(\"staff id\")\n\t\t}\n\n\t\t{{name}}{{Method}}Request.StaffId = srcStaffId\n\t} else {\n\t\terrMsg := tlog.Error(\"decode {{name_method}} request (%v) err (staff id illegal).\", r)\n\n\t\ttlog.AsyncSend(tlog.NewRawLogError(comconst.ErrParamIllegal(\"staff id\"), errMsg).AttachRequest(r))\n\n\t\treturn nil, comconst.ErrParamIllegal(\"staff id\")\n\t}\n\n\tcsrfSecret := strings.TrimSpace(r.Header.Get(\"X-CSRF-Token\"))\n\n\tif csrfSecret == \"\" {\n\t\terrMsg := tlog.Error(\"decode {{name_method}} request (%v) err (csrf secret illegal).\", r)\n\n\t\ttlog.AsyncSend(tlog.NewRawLogError(comconst.ErrParamIllegal(\"csrf secret\"), errMsg).AttachRequest(r))\n\n\t\treturn nil, comconst.ErrParamIllegal(\"csrf secret\")\n\t}\n\n\t{{name}}{{Method}}Request.CsrfSecret = csrfSecret\n\n\taccessToken, err := r.Cookie(\"AccessToken\")\n\n\tif err != nil || accessToken == nil || accessToken.Value == \"\" {\n\t\tsrcAccessToken := strings.TrimSpace(r.Header.Get(\"AccessToken\"))\n\n\t\tif srcAccessToken == \"\" {\n\t\t\terrMsg := tlog.Error(\"decode {{name_method}} request (%v) err (access token %v).\", r, err)\n\n\t\t\ttlog.AsyncSend(tlog.NewRawLogError(err, errMsg).AttachRequest(r))\n\n\t\t\treturn nil, comconst.ErrParamIllegal(\"access token\")\n\t\t} else {\n\t\t\t{{name}}{{Method}}Request.AccessToken = srcAccessToken\n\t\t}\n\t} else {\n\t\t{{name}}{{Method}}Request.AccessToken = accessToken.Value\n\t}\n\n\t{{name}}{{Method}}Request.Request = r\n\n\treturn {{name}}{{Method}}Request, nil\n}\n\nfunc Encode{{Name}}{{Method}}Response(ctx context.Context, w http.ResponseWriter, resp interface{}) error {\n\tw.Header().Set(\"Content-Type\", \"application/json; charset=utf-8\")\n\n\trespData := resp.(*data.{{Name}}{{Method}}Response)\n\n\tretData := &respdata.Response{\n\t\tError: errconst.COMMON_API_CODE_OK,\n\t}\n\n\tif respData.ErrCode != errconst.COMMON_API_ERROR_OK {\n\t\tretData.Error = errconst.RESP_API_ERROR_DESC[respData.ErrCode]\n\t\tretData.Message = respData.ErrCode\n\n\t\tw.WriteHeader(errconst.RESP_API_STATUS_CODE_DESC[respData.ErrCode])\n\t} else {\n\t\tretData.Data = respData.Data\n\t}\n\n\terr := jsoniter.NewEncoder(w).Encode(retData)\n\n\treturn err\n}\n\nfunc Make{{Name}}{{Method}}Endpoint(service service.{{Gw}}Service) endpoint.Endpoint {\n\treturn func(ctx context.Context, req interface{}) (interface{}, error) {\n\t\treqData := req.(*data.{{Name}}{{Method}}Request)\n\n\t\trespData, _ := service.{{Name}}{{Method}}(ctx, reqData)\n\n\t\treturn respData, nil\n\t}\n}\n"
	templateB_post = "package endpoint\n\nimport (\n\t\"context\"\n\t\"net/http\"\n\t\"strings\"\n\n\t\"xtech-kit/const\"\n\t\"xtech-kit/error\"\n\t\"xtech-kit/gw/data\"\n\t\"xtech-kit/gw/{{gw}}-gw/pkg/{{gw}}/data\"\n\t\"xtech-kit/gw/{{gw}}-gw/pkg/{{gw}}/service\"\n\n\t\"github.com/go-kit/kit/endpoint\"\n\t\"github.com/json-iterator/go\"\n\t\"github.com/xianyu-tech/tlog\"\n)\n\n\nfunc Decode{{Name}}{{method}}Request(ctx context.Context, r *http.Request) (interface{}, error) {\n\t{{name}}{{method}}Request := &data.{{Name}}{{method}}Request{}\n\n\terr := jsoniter.NewDecoder(r.Body).Decode({{name}}{{method}}Request)\n\n\tif err != nil {\n\t\terrMsg := tlog.Error(\"decode {{name_method}} request (%v) err (unmarshal %v).\", r, err)\n\n\t\ttlog.AsyncSend(tlog.NewRawLogError(comconst.ErrParamIllegal(\"unmarshal\"), errMsg).AttachRequest(r))\n\n\t\treturn nil, comconst.ErrParamIllegal(\"unmarshal\")\n\t}\n\t\n\n\tstaffId := {{name}}{{method}}Request.StaffId\n\n\tif staffId == 0 {\n\t\terrMsg := tlog.Error(\"decode {{name_method}} request (%v) err (staff id illegal).\", r)\n\n\t\ttlog.AsyncSend(tlog.NewRawLogError(comconst.ErrParamIllegal(\"staff id\"), errMsg).AttachRequest(r))\n\n\t\treturn nil, comconst.ErrParamIllegal(\"staff id\")\n\t}\n\n\tcsrfSecret := strings.TrimSpace(r.Header.Get(\"X-CSRF-Token\"))\n\n\tif csrfSecret == \"\" {\n\t\terrMsg := tlog.Error(\"decode {{name_method}} request (%v) err (csrf secret illegal).\", r)\n\n\t\ttlog.AsyncSend(tlog.NewRawLogError(comconst.ErrParamIllegal(\"csrf secret\"), errMsg).AttachRequest(r))\n\n\t\treturn nil, comconst.ErrParamIllegal(\"csrf secret\")\n\t}\n\n\t{{name}}{{method}}Request.CsrfSecret = csrfSecret\n\n\taccessToken, err := r.Cookie(\"AccessToken\")\n\n\tif err != nil || accessToken == nil || accessToken.Value == \"\" {\n\t\tsrcAccessToken := strings.TrimSpace(r.Header.Get(\"AccessToken\"))\n\n\t\tif srcAccessToken == \"\" {\n\t\t\terrMsg := tlog.Error(\"decode {{name_method}} request (%v) err (access token illegal).\", r)\n\n\t\t\ttlog.AsyncSend(tlog.NewRawLogError(comconst.ErrParamIllegal(\"access token\"), errMsg).AttachRequest(r))\n\n\t\t\treturn nil, comconst.ErrParamIllegal(\"access token\")\n\t\t} else {\n\t\t\t{{name}}{{method}}Request.AccessToken = srcAccessToken\n\t\t}\n\t} else {\n\t\t{{name}}{{method}}Request.AccessToken = accessToken.Value\n\t}\n\n\t{{name}}{{method}}Request.Request = r\n\n\treturn {{name}}{{method}}Request, nil\n}\n\nfunc Encode{{Name}}{{method}}Response(ctx context.Context, w http.ResponseWriter, resp interface{}) error {\n\tw.Header().Set(\"Content-Type\", \"application/json; charset=utf-8\")\n\n\trespData := resp.(*data.{{Name}}{{method}}Response)\n\n\tretData := &respdata.Response{\n\t\tError: errconst.COMMON_API_CODE_OK,\n\t}\n\tif respData.ErrCode != errconst.COMMON_API_ERROR_OK {\n\t\tretData.Error = errconst.RESP_API_ERROR_DESC[respData.ErrCode]\n\t\tretData.Message = respData.ErrCode\n\t\tw.WriteHeader(errconst.RESP_API_STATUS_CODE_DESC[respData.ErrCode])\n\t} else {\n\t\tretData.Data = respData.Data\n\t}\n\n\terr := jsoniter.NewEncoder(w).Encode(retData)\n\n\treturn err\n}\n\nfunc Make{{Name}}{{method}}Endpoint(service service.{{Gw}}Service) endpoint.Endpoint {\n\treturn func(ctx context.Context, req interface{}) (interface{}, error) {\n\t\treqData := req.(*data.{{Name}}{{method}}Request)\n\n\t\trespData, _ := service.{{Name}}{{method}}(ctx, reqData)\n\n\t\treturn respData, nil\n\t}\n}\n"
	templateC      = "package service\n\nimport (\n\t\"context\"\n\t\n\t\"xtech-kit/const\"\n\t\"xtech-kit/error\"\n\t\"xtech-kit/gw/{{gw}}-gw/pkg/{{gw}}/data\"\n\t{{service}}cligrpc \"xtech-kit/api/{{service}}-api/client/grpc\"\n\n\t\"github.com/xianyu-tech/tlog\"\n\t\"github.com/xianyu-tech/util/log\"\n)\n\nfunc (p *basic{{Gw}}Service) {{Name}}{{method}}(ctx context.Context, reqdata *data.{{Name}}{{method}}Request) (*data.{{Name}}{{method}}Response, error) {\n\tdefer func() {\n\t\tif rec := recover(); rec != nil {\n\t\t\terrMsg := tlog.Error(\"official vote del (%d, %d) err (recover %v).\",\n\t\t\t\treqdata.UserId, reqdata.ResultId, logutil.PrintPanic())\n\n\t\t\ttlog.AsyncSend(tlog.NewRawLogError(comconst.ErrPanicRecovery(\"official vote del\"), errMsg))\n\t\t}\n\t}()\n\n\trequest := reqdata.Request\n\n\tuserId := reqdata.UserId\n\tresultId := reqdata.ResultId\n\n\t{{name}}{{method}}ReqProto := {{service}}cligrpc.New{{Name}}{{method}}ReqProto(userId, resultId)\n\n\t{{name}}{{method}}RespProto, err := p.{{service}}Endpoints.{{Name}}{{method}}(ctx, {{name}}{{method}}ReqProto)\n\n\tif err != nil {\n\t\terrMsg := tlog.Error(\"official vote del (%d, %d) err ({{service}} api rpc %v).\", userId, resultId, err)\n\n\t\ttlog.AsyncSend(tlog.NewRawLogError(err, errMsg).AttachRequest(request))\n\n\t\t{{name}}{{method}}Response := &data.{{Name}}{{method}}Response{\n\t\t\tErrCode: errconst.APP_GW_ERROR_{{service}}_API_SERVER_ABNORMAL,\n\t\t\tErrMsg:  errMsg,\n\t\t}\n\n\t\treturn {{name}}{{method}}Response, nil\n\t}\n\n\tif {{name}}{{method}}RespProto.{{method}}ErrCode() != errconst.COMMON_API_ERROR_OK {\n\t\terrMsg := tlog.Error(\"official vote del (%d, %d) err (official vote del %s).\",\n\t\t\tuserId, resultId, {{name}}{{method}}RespProto.{{method}}ErrMsg())\n\n\t\ttlog.AsyncSend(tlog.NewRawLogError(comconst.ErrSvcExecute(\"official vote del\"), errMsg).AttachRequest(request))\n\n\t\t{{name}}{{method}}Response := &data.{{Name}}{{method}}Response{\n\t\t\tErrCode: {{name}}{{method}}RespProto.{{method}}ErrCode(),\n\t\t\tErrMsg:  errMsg,\n\t\t}\n\n\t\treturn {{name}}{{method}}Response, nil\n\t}\n\n\t{{name}}{{method}}Response := &data.{{Name}}{{method}}Response{\n\t\tErrCode: errconst.COMMON_API_ERROR_OK,\n\t}\n\n\tvoteProto := {{name}}{{method}}RespProto.{{method}}Vote()\n\tif voteProto != nil {\n\t\t{{name}}{{method}}Response.Data = &data.{{Name}}Data{\n\t\t\tVoteId:    voteProto.GetVoteId(),\n\t\t\tUserId:    voteProto.GetUserId(),\n\t\t\tResultId:  voteProto.GetResultId(),\n\t\t\tTeamId:    voteProto.GetTeamId(),\n\t\t\tCreatedAt: voteProto.GetCreatedAt(),\n\t\t}\n\t}\n\n\treturn {{name}}{{method}}Response, nil\n}\n\nfunc (l permissionMiddleware) {{Name}}{{method}}(ctx context.Context, reqdata *data.{{Name}}{{method}}Request) (*data.{{Name}}{{method}}Response, error) {\n\tpermit, logErr := l.checkPermission(ctx, reqdata.Request, reqdata.serId, reqdata.CsrfSecret, reqdata.AccessToken)\n\n\tif logErr != nil {\n\t\terrMsg := tlog.Error(\"check permission err (check permission %s).\", logErr.ErrMsg())\n\n\t\tresponse := &data.{{Name}}{{method}}Response{\n\t\t\tErrCode: logErr.ErrCode(),\n\t\t\tErrMsg:  errMsg,\n\t\t}\n\n\t\treturn response, nil\n\t}\n\n\tif permit == false {\n\t\terrMsg := tlog.Error(\"check permission err (no permission).\")\n\n\t\tresponse := &data.{{Name}}{{method}}Response{\n\t\t\tErrCode: errconst.APP_GW_ERROR_PRIVILEGE_API_NO_PERMISSION,\n\t\t\tErrMsg:  errMsg,\n\t\t}\n\n\t\treturn response, nil\n\t}\n\n\treturn l.next.{{Name}}{{method}}(ctx, reqdata)\n}\n"
)

var (
	templateMap = map[string]string{
		"template1": template1,
		"template2": template2,
		"template3": template3,
		"template4": template4,
	}
)

func main() {
	name := "informationCommentsByUser"
	service := "information"
	method := "Get"
	METHOD := "GET"
	router := "/sys-information/v1/comments-by-user"
	filename := "information_comments_by_user_get.go"
	gw := "sys"
	name_method := "moment posts by user get"
	filePathPre := "/Users/haikuotiankong/Desktop/go-code/src/xtech-kit/gw/" + gw + "-gw/pkg/" + gw

	Service := Ucfirst(service)
	Name := Ucfirst(name)
	Gw := Ucfirst(gw)

	endpointPath := filePathPre + "/endpoint/endpoint.go"
	httpHandlerPath := filePathPre + "/http/handler.go"
	servicePath := filePathPre + "/service/service.go"

	var pathArry = []string{endpointPath, httpHandlerPath, servicePath}

	for _, path := range pathArry {
		err := handlerFile(path, name, Name, method, router, METHOD)
		if err != nil {
			panic(err)
		}
	}

	newDataPath := filePathPre + "/data/" + filename
	createFile(newDataPath, templateA, name, Name, method, gw, Gw, service, Service, name_method)

	newEndpointPath := filePathPre + "/endpoint/" + filename
	if method == "Get" {
		createFile(newEndpointPath, templateB_get, name, Name, method, gw, Gw, service, Service, name_method)
	} else {
		createFile(newEndpointPath, templateB_post, name, Name, method, gw, Gw, service, Service, name_method)
	}

	newServicePath := filePathPre + "/service/" + filename
	createFile(newServicePath, templateC, name, Name, method, gw, Gw, service, Service, name_method)

}
func createFile(filename, template, name, Name, method, gw, Gw, service, Service, name_method string) {
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
	s = strings.ReplaceAll(s, "{{gw}}", gw)
	s = strings.ReplaceAll(s, "{{Gw}}", Gw)
	s = strings.ReplaceAll(s, "{{name_method}}", name_method)

	_, err1 = io.WriteString(f, s)
	if err1 != nil {
		log.Printf("Cannot write text file: %s, err: [%v]", filename, err1)
		return
	}
}

func handlerFile(filepath, name, Name, method, router, METHOD string) error {
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
				s = strings.ReplaceAll(s, "{{router}}", router)
				s = strings.ReplaceAll(s, "{{METHOD}}", METHOD)

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
