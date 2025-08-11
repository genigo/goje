package goje

/*
*

	Limit Query

*
*/
type QueryLimit int

func (l QueryLimit) GetType() string {
	return QueryTypeLimit
}

func (l QueryLimit) GetQuery() string {
	return "LIMIT ?"
}

func (l QueryLimit) GetArgs() []any {
	return []any{l}
}

func Limit(limit int) QueryLimit {
	return QueryLimit(limit)
}

/**
	Offset Query
**/

type QueryOffset int

func (l QueryOffset) GetType() string {
	return QueryTypeOffset
}

func (l QueryOffset) GetQuery() string {
	return "OFFSET ?"
}

func (l QueryOffset) GetArgs() []any {
	return []any{l}
}

func Offset(offset int) QueryOffset {
	return QueryOffset(offset)
}

/**
	order Query
**/

type QueryOrder struct {
	query string
	args  []any
}

func (o QueryOrder) GetType() string {
	return QueryTypeOrder
}

func (o QueryOrder) GetQuery() string {
	return o.query
}

func (o QueryOrder) GetArgs() []any {
	return o.args
}

func Order(query string, args ...any) QueryOrder {
	return QueryOrder{
		query: query,
		args:  args,
	}
}

/**
	Group Query
**/

type QueryGroup struct {
	query string
	args  []any
}

func (h QueryGroup) GetType() string {
	return QueryTypeGroup
}

func (h QueryGroup) GetQuery() string {
	return h.query
}

func (h QueryGroup) GetArgs() []any {
	return h.args
}

func GroupBy(query string, args ...any) QueryGroup {
	return QueryGroup{
		query: query,
		args:  args,
	}
}

/**
	having Query
**/

type QueryHaving struct {
	query string
	args  []any
}

func (h QueryHaving) GetType() string {
	return QueryTypeHaving
}

func (h QueryHaving) GetQuery() string {
	return h.query
}

func (h QueryHaving) GetArgs() []any {
	return h.args
}

func Having(query string, args ...any) QueryHaving {
	return QueryHaving{
		query: query,
		args:  args,
	}
}

/**
	Where Query
**/

type QueryWhere struct {
	query string
	args  []any
}

func (q QueryWhere) GetType() string {
	return QueryTypeWhere
}

func (q QueryWhere) GetQuery() string {
	return q.query
}

func (q QueryWhere) GetArgs() []any {
	return q.args
}

func Where(query string, args ...any) QueryWhere {
	return QueryWhere{
		query: query,
		args:  args,
	}
}

/**
	Where in Query
**/

type QueryWhereIn struct {
	query string
	args  []any
}

func (q QueryWhereIn) GetType() string {
	return QueryTypeWhereIn
}

func (q QueryWhereIn) GetQuery() string {
	return q.query
}

func (q QueryWhereIn) GetArgs() []any {
	return q.args
}

func WhereIn(columnName string, args ...any) QueryWhereIn {
	return QueryWhereIn{
		query: columnName,
		args:  args,
	}
}

/**
	Where in Query
**/

type QueryWhereNotIn struct {
	query string
	args  []any
}

func (q QueryWhereNotIn) GetType() string {
	return QueryTypeWhereNotIn
}

func (q QueryWhereNotIn) GetQuery() string {
	return q.query
}

func (q QueryWhereNotIn) GetArgs() []any {
	return q.args
}

func WhereNotIn(columnName string, args ...any) QueryWhereNotIn {
	return QueryWhereNotIn{
		query: columnName,
		args:  args,
	}
}

/**
	Join Query Builder
**/

type QueryJoin struct {
	joinType string
	table    string
	on       string
	args     []any
}

func (q QueryJoin) GetType() string {
	return QueryTypeJoin
}

func (q QueryJoin) GetQuery() string {
	var on string
	if q.on != "" {
		on = " ON " + q.on
	}
	return " " + q.joinType + " JOIN " + q.table + on + " "
}

func (q QueryJoin) GetArgs() []any {
	return q.args
}

func InnerJoin(table string, on string, args ...any) QueryJoin {
	return QueryJoin{
		on:       on,
		joinType: Inner,
		table:    table,
		args:     args,
	}
}

func OuterJoin(table string, on string, args ...any) QueryJoin {
	return QueryJoin{
		on:       on,
		joinType: Outer,
		table:    table,
		args:     args,
	}
}

func NaturalJoin(table string, on string, args ...any) QueryJoin {
	return QueryJoin{
		on:       on,
		joinType: Natural,
		table:    table,
		args:     args,
	}
}

func RightJoin(table string, on string, args ...any) QueryJoin {
	return QueryJoin{
		on:       on,
		joinType: Right,
		table:    table,
		args:     args,
	}
}

func LeftJoin(table string, on string, args ...any) QueryJoin {
	return QueryJoin{
		on:       on,
		joinType: Left,
		table:    table,
		args:     args,
	}
}
