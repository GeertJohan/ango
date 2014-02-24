package {{.PackageName}}

//++ create interface that implements procedures

type {{.Service.Name}}Handler interface {
	{{range .Service.ServerProcedures}}
	{{.Name}}()
	{{end}}
}