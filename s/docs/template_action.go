// author: wsfuyibing <websearch@163.com>
// date: 2021-12-16

package docs

var templateAction = `# {{TITLE}}

接口 ：·{{METHOD}} {{ROUTE}}· <br />
版本 ：·{{VERSION}}· <br />

> {{DESCRIPTION}}

### 【入参】

{{REQUEST}}

### 【出参】

{{RESPONSE}}

--------
入口 ：·{{CALLABLE_NAME}}.{{CALLABLE_FUNC}}()· <br />
源码 ：·{{SOURCE_FILE}}: {{SOURCE_LINE}}· <br />
更新 ：·{{UPDATED}}·
`
