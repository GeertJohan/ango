// This is the yacc file for ango
// To build it:
// go tool yacc -p "ango" ango.y (produces y.go)

%{

package parser

import (
	//++
)

%}

%union {
	num *big.Rat
}

%type	<num>	expr expr1 expr2 expr3

%token	<num>	NUM

%%



%%
