# 文档定义

**语法**

```shell
go run main.go docs
```

**目录**

> 命令执行时扫描控制器目录 (`app/controllers`) 下的文件 (`*controller.go`)

**包名**

> 命令执行时, 需要项目下有包文件（`go.mod`）定义, 从中读取包名.

**全局**

> 全局选项定义在文档文件（`app/controllers/doc.go`）中, 接受以下注解

1. @Host(string) - 部署IP.
   ```text
    // @Host(0.0.0.0)
    package controllers
    ```
2. @Port(int) - 端口号.
   ```text
    // @Port(8080)
    package controllers
    ```
3. @Version(string) - 版本号.
   ```text
    // @Version(1.2.3)
    package controllers
    ```
4. @Domain(string) - 域名.
   ```text
    // @Domain(example.com)
    package controllers
    ```
5. @DomainPrefix(string) - 域名前缀.
   ```text
    // @DomainPrefix(www)
    package controllers
    ```

**控制器**

1. @RoutePrefix(string) - 路由前缀.
   ```text
    // @RoutePrefix(/manage)
    type Controller struct{
    }
    ```

**接口**

1. @Ignore(bool) - 是否忽略导出文档, 默认: false.
   ```text
    // @Ignore()
    // @Ignore(false)
    func (o *Controller) GetList(ctx iris.Context) interface{} {
    }
    ```
2. @Request(Struct) - 入参结构体路径.
   ```text
    // @Request(app/logics/topic.CreateRequest)
    func (o *Controller) PostCreate(ctx iris.Context) interface{} {
    }
    ```
3. @Response(Struct) - 出参结构体路径.
   ```text
    // @Response(app/logics/topic.CreateResponse)
    func (o *Controller) PostCreate(ctx iris.Context) interface{} {
    }
    ```
4. @Sdk(func) - 导出SDK名称.
   ```text
    // @Sdk(CreateTopic)
    func (o *Controller) PostCreate(ctx iris.Context) interface{} {
    }
    ```
5. @Version(string) - 接口版本号.
   ```text
    // @Version(1.2.3)
    func (o *Controller) PostCreate(ctx iris.Context) interface{} {
    }
    ```

**标签**

1. json - 输入/输出JSON格式时字段名.
   ```text
   type Example struct {
       Id int `json:"id"`
   }
   ```
2. url - 通过URL传递参数时输入字段.
   ```text
   type Example struct {
       Id int `url:"id"`
   }
   ```
3. form - 通过FORM表单传递入参时字段名.
   ```text
   type Example struct {
       Id int `form:"id"`
   }
   ```
4. desc - 导出文档时字段描述.
   ```text
   type Example struct {
       Id int `desc:"在数据表Table中的主键值"`
   }
   ```
5. label - 导出文档时字段标题(名称).
   ```text
   type Example struct {
       Id int `label:"ID号"`
   }
   ```
6. mock - 指定Mock数据.
   ```text
   type Example struct {
       Id int `mock:"10"`
   }
   ```
7. validate - 指定校验规则.
   ```text
   type Example struct {
       Id int `validate:"required,min=5"`
   }
   ```

