// This is the yacc file for ango
// To build it:
// go tool yacc -p "ango" ango.y (produces y.go)

%{

package parser

import (
	//++ yacc imports
)

%}

%union {
	txt string
}

%type	<txt>	expr expr1 expr2 expr3

%token	<txt>	IDENT SERVER CLIENT 

%%

IDENT

%%

//++ yacc helper functions