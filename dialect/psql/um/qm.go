package um

import (
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/clause"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/expr"
	"github.com/stephenafamo/bob/mods"
)

func With(name string, columns ...string) dialect.CTEChain[*dialect.UpdateQuery] {
	return dialect.With[*dialect.UpdateQuery](name, columns...)
}

func Recursive(r bool) bob.Mod[*dialect.UpdateQuery] {
	return mods.Recursive[*dialect.UpdateQuery](r)
}

func Only() bob.Mod[*dialect.UpdateQuery] {
	return mods.QueryModFunc[*dialect.UpdateQuery](func(u *dialect.UpdateQuery) {
		u.Only = true
	})
}

func Table(name any) bob.Mod[*dialect.UpdateQuery] {
	return mods.QueryModFunc[*dialect.UpdateQuery](func(u *dialect.UpdateQuery) {
		u.Table = clause.Table{
			Expression: name,
		}
	})
}

func TableAs(name any, alias string) bob.Mod[*dialect.UpdateQuery] {
	return mods.QueryModFunc[*dialect.UpdateQuery](func(u *dialect.UpdateQuery) {
		u.Table = clause.Table{
			Expression: name,
			Alias:      alias,
		}
	})
}

func Set(a string, b any) bob.Mod[*dialect.UpdateQuery] {
	return mods.Set[*dialect.UpdateQuery]{expr.OP("=", expr.Quote(a), b)}
}

func SetArg(a string, b any) bob.Mod[*dialect.UpdateQuery] {
	return mods.Set[*dialect.UpdateQuery]{expr.OP("=", expr.Quote(a), expr.Arg(b))}
}

func From(table any) dialect.FromChain[*dialect.UpdateQuery] {
	return dialect.From[*dialect.UpdateQuery](table)
}

func FromFunction(funcs ...*dialect.Function) dialect.FromChain[*dialect.UpdateQuery] {
	var table any

	if len(funcs) == 1 {
		table = funcs[0]
	}

	if len(funcs) > 1 {
		table = dialect.Functions(funcs)
	}

	return dialect.From[*dialect.UpdateQuery](table)
}

func InnerJoin(e any) dialect.JoinChain[*dialect.UpdateQuery] {
	return dialect.InnerJoin[*dialect.UpdateQuery](e)
}

func LeftJoin(e any) dialect.JoinChain[*dialect.UpdateQuery] {
	return dialect.LeftJoin[*dialect.UpdateQuery](e)
}

func RightJoin(e any) dialect.JoinChain[*dialect.UpdateQuery] {
	return dialect.RightJoin[*dialect.UpdateQuery](e)
}

func FullJoin(e any) dialect.JoinChain[*dialect.UpdateQuery] {
	return dialect.FullJoin[*dialect.UpdateQuery](e)
}

func CrossJoin(e any) bob.Mod[*dialect.UpdateQuery] {
	return dialect.CrossJoin[*dialect.UpdateQuery](e)
}

func Where(e bob.Expression) bob.Mod[*dialect.UpdateQuery] {
	return mods.Where[*dialect.UpdateQuery]{e}
}

func WhereClause(clause string, args ...any) bob.Mod[*dialect.UpdateQuery] {
	return mods.Where[*dialect.UpdateQuery]{expr.RawQuery(dialect.Dialect, clause, args...)}
}

func Returning(clauses ...any) bob.Mod[*dialect.UpdateQuery] {
	return mods.Returning[*dialect.UpdateQuery](clauses)
}