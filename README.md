# go-kit-auto-code
根据参数自动生成go-kit框架的接口的重复代码


ps：有些服务已经设置好template1~9，不需要再设置，没有设置的一定要设置模板。



svc:


    不要改的：
    template1~9 是先定义在文件中的目标位置，在endpoint.go\handler.go\service.go中，可以参考staff-svc
    templateA~D 是创建文件时用的，不需任何操作

    要改的：
    service := 服务的名字,小写
    name :=  接口名字，首字母小写
    method := 方法名，首字母大写
    filename :=  要创建的文件名
    filePathPre := 本地文件绝对路径，创建文件用的
    
api:

    不要改的：
    template1~9 是先定义在文件中的目标位置，在endpoint.go\handler.go\service.go中，可以参考staff-api
    templateA~D 是创建文件时用的，不需任何操作
    
    要改的：
    service := 服务的名字,小写
    name :=  接口名字，首字母小写
    method := 方法名，首字母大写
    filename :=  要创建的文件名
    filePathPre := 本地文件绝对路径，创建文件用的
    
    
gw:

    不要改的：
    template1~4 是先定义在文件中的目标位置，在endpoint.go\handler.go\service.go中，已经都定义好了
    templateA~D 是创建文件时用的，不需任何操作


    要改的：
    service := 服务的名字,小写
    name :=  接口名字，首字母小写
    method := 方法名，首字母大写
    METHOD := restful方法名，全部大写
    router := 路由
    filename :=  要创建的文件名
    gw := 前台or后台
    filePathPre := 本地文件绝对路径，创建文件用的
    name_method := 打印的日志
	