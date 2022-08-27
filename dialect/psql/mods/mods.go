package mods

import (
	"io"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/clause"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/expr"
	"github.com/stephenafamo/bob/mods"
)

//nolint:gochecknoglobals
var bmod = expr.Builder[psql.Expression, psql.Expression]{}

type Distinct struct {
	On []any
}

func (di Distinct) WriteSQL(w io.Writer, d bob.Dialect, start int) ([]any, error) {
	w.Write([]byte("DISTINCT"))
	return bob.ExpressSlice(w, d, start, di.On, " ON (", ", ", ")")
}

func With[Q interface{ AppendWith(clause.CTE) }](name string, columns ...string) CteChain[Q] {
	return CteChain[Q](func() clause.CTE {
		return clause.CTE{
			Name:    name,
			Columns: columns,
		}
	})
}

type FromItemMod[Q interface {
	SetTableAlias(alias string, columns ...string)
	SetOnly(bool)
	SetLateral(bool)
	SetWithOrdinality(bool)
}] struct{}

func As[Q interface{ SetTableAlias(string, ...string) }](alias string, columns ...string) bob.Mod[Q] {
	return mods.QueryModFunc[Q](func(q Q) {
		q.SetTableAlias(alias, columns...)
	})
}

func (FromItemMod[Q]) Only() bob.Mod[Q] {
	return mods.QueryModFunc[Q](func(q Q) {
		q.SetOnly(true)
	})
}

func (FromItemMod[Q]) Lateral() bob.Mod[Q] {
	return mods.QueryModFunc[Q](func(q Q) {
		q.SetLateral(true)
	})
}

func (FromItemMod[Q]) WithOrdinality() bob.Mod[Q] {
	return mods.QueryModFunc[Q](func(q Q) {
		q.SetWithOrdinality(true)
	})
}

type JoinChain[Q interface{ AppendJoin(clause.Join) }] func() clause.Join

func (j JoinChain[Q]) Apply(q Q) {
	q.AppendJoin(j())
}

func (j JoinChain[Q]) As(alias string) JoinChain[Q] {
	jo := j()
	jo.Alias = alias

	return JoinChain[Q](func() clause.Join {
		return jo
	})
}

func (j JoinChain[Q]) Natural() bob.Mod[Q] {
	jo := j()
	jo.Natural = true

	return mods.Join[Q](jo)
}

func (j JoinChain[Q]) On(on ...any) bob.Mod[Q] {
	jo := j()
	jo.On = append(jo.On, on...)

	return mods.Join[Q](jo)
}

func (j JoinChain[Q]) OnEQ(a, b any) bob.Mod[Q] {
	jo := j()
	jo.On = append(jo.On, bmod.X(a).EQ(b))

	return mods.Join[Q](jo)
}

func (j JoinChain[Q]) Using(using ...any) bob.Mod[Q] {
	jo := j()
	jo.Using = using

	return mods.Join[Q](jo)
}

type joinable interface{ AppendJoin(clause.Join) }

func InnerJoin[Q joinable](e any) JoinChain[Q] {
	return JoinChain[Q](func() clause.Join {
		return clause.Join{
			Type: clause.InnerJoin,
			To:   e,
		}
	})
}

func LeftJoin[Q joinable](e any) JoinChain[Q] {
	return JoinChain[Q](func() clause.Join {
		return clause.Join{
			Type: clause.LeftJoin,
			To:   e,
		}
	})
}

func RightJoin[Q joinable](e any) JoinChain[Q] {
	return JoinChain[Q](func() clause.Join {
		return clause.Join{
			Type: clause.RightJoin,
			To:   e,
		}
	})
}

func FullJoin[Q joinable](e any) JoinChain[Q] {
	return JoinChain[Q](func() clause.Join {
		return clause.Join{
			Type: clause.FullJoin,
			To:   e,
		}
	})
}

func CrossJoin[Q joinable](e any) bob.Mod[Q] {
	return mods.Join[Q]{
		Type: clause.CrossJoin,
		To:   e,
	}
}

type OrderBy[Q interface{ AppendOrder(clause.OrderDef) }] func() clause.OrderDef

func (s OrderBy[Q]) Apply(q Q) {
	q.AppendOrder(s())
}

func (o OrderBy[Q]) Asc() OrderBy[Q] {
	order := o()
	order.Direction = "ASC"

	return OrderBy[Q](func() clause.OrderDef {
		return order
	})
}

func (o OrderBy[Q]) Desc() OrderBy[Q] {
	order := o()
	order.Direction = "DESC"

	return OrderBy[Q](func() clause.OrderDef {
		return order
	})
}

func (o OrderBy[Q]) Using(operator string) OrderBy[Q] {
	order := o()
	order.Direction = "USING " + operator

	return OrderBy[Q](func() clause.OrderDef {
		return order
	})
}

func (o OrderBy[Q]) NullsFirst() OrderBy[Q] {
	order := o()
	order.Nulls = "FIRST"

	return OrderBy[Q](func() clause.OrderDef {
		return order
	})
}

func (o OrderBy[Q]) NullsLast() OrderBy[Q] {
	order := o()
	order.Nulls = "LAST"

	return OrderBy[Q](func() clause.OrderDef {
		return order
	})
}

func (o OrderBy[Q]) Collate(collation string) OrderBy[Q] {
	order := o()
	order.CollationName = collation

	return OrderBy[Q](func() clause.OrderDef {
		return order
	})
}

type CteChain[Q interface{ AppendWith(clause.CTE) }] func() clause.CTE

func (c CteChain[Q]) Apply(q Q) {
	q.AppendWith(c())
}

func (c CteChain[Q]) As(q bob.Query) CteChain[Q] {
	cte := c()
	cte.Query = q
	return CteChain[Q](func() clause.CTE {
		return cte
	})
}

func (c CteChain[Q]) NotMaterialized() CteChain[Q] {
	b := false
	cte := c()
	cte.Materialized = &b
	return CteChain[Q](func() clause.CTE {
		return cte
	})
}

func (c CteChain[Q]) Materialized() CteChain[Q] {
	b := true
	cte := c()
	cte.Materialized = &b
	return CteChain[Q](func() clause.CTE {
		return cte
	})
}

func (c CteChain[Q]) SearchBreadth(setCol string, searchCols ...string) CteChain[Q] {
	cte := c()
	cte.Search = clause.CTESearch{
		Order:   clause.SearchDepth,
		Columns: searchCols,
		Set:     setCol,
	}
	return CteChain[Q](func() clause.CTE {
		return cte
	})
}

func (c CteChain[Q]) SearchDepth(setCol string, searchCols ...string) CteChain[Q] {
	cte := c()
	cte.Search = clause.CTESearch{
		Order:   clause.SearchDepth,
		Columns: searchCols,
		Set:     setCol,
	}
	return CteChain[Q](func() clause.CTE {
		return cte
	})
}

func (c CteChain[Q]) Cycle(set, using string, cols ...string) CteChain[Q] {
	cte := c()
	cte.Cycle.Set = set
	cte.Cycle.Using = using
	cte.Cycle.Columns = cols
	return CteChain[Q](func() clause.CTE {
		return cte
	})
}

func (c CteChain[Q]) CycleValue(value, defaultVal any) CteChain[Q] {
	cte := c()
	cte.Cycle.SetVal = value
	cte.Cycle.DefaultVal = defaultVal
	return CteChain[Q](func() clause.CTE {
		return cte
	})
}

type LockChain[Q interface{ SetFor(clause.For) }] func() clause.For

func (l LockChain[Q]) Apply(q Q) {
	q.SetFor(l())
}

func (l LockChain[Q]) NoWait() LockChain[Q] {
	lock := l()
	lock.Wait = clause.LockWaitNoWait
	return LockChain[Q](func() clause.For {
		return lock
	})
}

func (l LockChain[Q]) SkipLocked() LockChain[Q] {
	lock := l()
	lock.Wait = clause.LockWaitSkipLocked
	return LockChain[Q](func() clause.For {
		return lock
	})
}

type WindowMod[Q interface{ AppendWindow(clause.NamedWindow) }] struct {
	Name string
	*WindowChain[*WindowMod[Q]]
}

func (w WindowMod[Q]) Apply(q Q) {
	q.AppendWindow(clause.NamedWindow{
		Name:       w.Name,
		Definition: w.def,
	})
}

type WindowChain[T any] struct {
	def  clause.WindowDef
	Wrap T
}

func (w *WindowChain[T]) From(name string) T {
	w.def.SetFrom(name)
	return w.Wrap
}

func (w *WindowChain[T]) PartitionBy(condition ...any) T {
	w.def.AddPartitionBy(condition...)
	return w.Wrap
}

func (w *WindowChain[T]) OrderBy(order ...any) T {
	w.def.AddOrderBy(order...)
	return w.Wrap
}

func (w *WindowChain[T]) Range() T {
	w.def.SetMode("RANGE")
	return w.Wrap
}

func (w *WindowChain[T]) Rows() T {
	w.def.SetMode("ROWS")
	return w.Wrap
}

func (w *WindowChain[T]) Groups() T {
	w.def.SetMode("GROUPS")
	return w.Wrap
}

func (w *WindowChain[T]) FromUnboundedPreceding() T {
	w.def.SetStart("UNBOUNDED PRECEDING")
	return w.Wrap
}

func (w *WindowChain[T]) FromPreceding(exp any) T {
	w.def.SetStart(bob.ExpressionFunc(
		func(w io.Writer, d bob.Dialect, start int) ([]any, error) {
			return bob.ExpressIf(w, d, start, exp, true, "", " PRECEDING")
		}),
	)
	return w.Wrap
}

func (w *WindowChain[T]) FromCurrentRow() T {
	w.def.SetStart("CURRENT ROW")
	return w.Wrap
}

func (w *WindowChain[T]) FromFollowing(exp any) T {
	w.def.SetStart(bob.ExpressionFunc(
		func(w io.Writer, d bob.Dialect, start int) ([]any, error) {
			return bob.ExpressIf(w, d, start, exp, true, "", " FOLLOWING")
		}),
	)
	return w.Wrap
}

func (w *WindowChain[T]) ToPreceding(exp any) T {
	w.def.SetEnd(bob.ExpressionFunc(
		func(w io.Writer, d bob.Dialect, start int) ([]any, error) {
			return bob.ExpressIf(w, d, start, exp, true, "", " PRECEDING")
		}),
	)
	return w.Wrap
}

func (w *WindowChain[T]) ToCurrentRow(count int) T {
	w.def.SetEnd("CURRENT ROW")
	return w.Wrap
}

func (w *WindowChain[T]) ToFollowing(exp any) T {
	w.def.SetEnd(bob.ExpressionFunc(
		func(w io.Writer, d bob.Dialect, start int) ([]any, error) {
			return bob.ExpressIf(w, d, start, exp, true, "", " FOLLOWING")
		}),
	)
	return w.Wrap
}

func (w *WindowChain[T]) ToUnboundedFollowing() T {
	w.def.SetEnd("UNBOUNDED FOLLOWING")
	return w.Wrap
}

func (w *WindowChain[T]) ExcludeNoOthers() T {
	w.def.SetExclusion("NO OTHERS")
	return w.Wrap
}

func (w *WindowChain[T]) ExcludeCurrentRow() T {
	w.def.SetExclusion("CURRENT ROW")
	return w.Wrap
}

func (w *WindowChain[T]) ExcludeGroup() T {
	w.def.SetExclusion("GROUP")
	return w.Wrap
}

func (w *WindowChain[T]) ExcludeTies() T {
	w.def.SetExclusion("TIES")
	return w.Wrap
}
