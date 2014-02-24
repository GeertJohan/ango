package {{.PackageName}}

//++ create interface that implements procedures

// {{.Service.CapitalizedName}}Handler types all methods that can be called by the client
type {{.Service.CapitalizedName}}Handler interface {
	{{range .Service.ServerProcedures}}
		{{.CapitalizedName}} ( {{.GoArgs}} )( {{.GoRets}} )
	{{end}}
}

// New{{.Service.CapitalizedName}}Handler must return a new instance implementing {{.Service.CapitalizedName}}Handler
type New{{.Service.CapitalizedName}}Handler func()(handler {{.Service.CapitalizedName}}Handler)

func StartServer(fn New{{.Service.CapitalizedName}}Handler) {
	//++ do something with New{{.Service.CapitalizedName}}Handler
}