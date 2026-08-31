package goje

import (
	"strings"
	"testing"
)

// withDialect activates a dialect for the test and restores the previous one
func withDialect(t *testing.T, d Dialect) {
	t.Helper()
	prev := dialect
	dialect = d
	t.Cleanup(func() { dialect = prev })
}

func TestDialectBindAndPlaceholders(t *testing.T) {
	if got := MysqlDialect.Bind(3); got != "?" {
		t.Errorf("mysql bind = %q, want ?", got)
	}
	if got := PostgresDialect.Bind(3); got != "$3" {
		t.Errorf("postgres bind = %q, want $3", got)
	}
	if got := PostgresDialect.Placeholders(4, 3); got != "$4,$5,$6" {
		t.Errorf("postgres placeholders = %q, want $4,$5,$6", got)
	}
	if got := MysqlDialect.Placeholders(1, 3); got != "?,?,?" {
		t.Errorf("mysql placeholders = %q, want ?,?,?", got)
	}
}

func TestRebind(t *testing.T) {
	withDialect(t, MysqlDialect)
	if got := rebind("WHERE a = ? AND b IN(?,?)"); got != "WHERE a = ? AND b IN(?,?)" {
		t.Errorf("mysql rebind changed the query: %s", got)
	}

	withDialect(t, PostgresDialect)
	if got := rebind("WHERE a = ? AND b IN(?,?)"); got != "WHERE a = $1 AND b IN($2,$3)" {
		t.Errorf("postgres rebind = %q, want $1..$3 numbering", got)
	}
}

func TestQouteColumnPerDialect(t *testing.T) {
	withDialect(t, MysqlDialect)
	if got := qouteColumn("user.name"); got != "`user`.`name`" {
		t.Errorf("mysql quote = %q", got)
	}

	withDialect(t, PostgresDialect)
	if got := qouteColumn("user.name"); got != `"user"."name"` {
		t.Errorf("postgres quote = %q", got)
	}
	// raw fragments pass through untouched on both dialects
	if got := qouteColumn("COUNT(*) AS total"); got != "COUNT(*) AS total" {
		t.Errorf("raw fragment quoted: %q", got)
	}
}

func TestSQLConditionBuilderPostgres(t *testing.T) {
	withDialect(t, PostgresDialect)

	query, args, err := SQLConditionBuilder([]QueryInterface{
		Where("age > ?", 18),
		WhereIn("role", "admin", "user"),
		Order("name"),
		Limit(10),
		Offset(5),
	})
	if err != nil {
		t.Fatal(err)
	}

	want := "  WHERE (age > $1) AND (role IN($2,$3)) ORDER BY name LIMIT $4 OFFSET $5"
	if query != want {
		t.Errorf("query = %q\nwant       %q", query, want)
	}
	if len(args) != 5 {
		t.Errorf("args = %v, want 5 entries", args)
	}
}

func TestSQLConditionBuilderMysqlUnchanged(t *testing.T) {
	withDialect(t, MysqlDialect)

	query, args, err := SQLConditionBuilder([]QueryInterface{
		Where("age > ?", 18),
		WhereIn("role", "admin", "user"),
		Limit(10),
	})
	if err != nil {
		t.Fatal(err)
	}

	want := "  WHERE (age > ?) AND (role IN(?,?)) LIMIT ?"
	if query != want {
		t.Errorf("query = %q\nwant       %q", query, want)
	}
	if len(args) != 4 {
		t.Errorf("args = %v, want 4 entries", args)
	}
}

func TestFindInSetPerDialect(t *testing.T) {
	withDialect(t, MysqlDialect)
	if got := FindInSet("roles", "admin").query; !strings.Contains(got, "FIND_IN_SET(?, `roles`)") {
		t.Errorf("mysql FindInSet = %q", got)
	}

	withDialect(t, PostgresDialect)
	if got := FindInSet("roles", "admin").query; got != "? = ANY(string_to_array(\"roles\", ','))" {
		t.Errorf("postgres FindInSet = %q", got)
	}
}

func TestInsertHelpers(t *testing.T) {
	withDialect(t, MysqlDialect)
	if got := insertValues(2, 2); got != "(?,?),(?,?)" {
		t.Errorf("mysql insertValues = %q", got)
	}
	if got := insertTable("users"); got != "users" {
		t.Errorf("mysql insertTable = %q", got)
	}

	withDialect(t, PostgresDialect)
	if got := insertValues(2, 2); got != "($1,$2),($3,$4)" {
		t.Errorf("postgres insertValues = %q", got)
	}
	if got := insertTable("users"); got != `"users"` {
		t.Errorf("postgres insertTable = %q", got)
	}
	cols := insertColumns([]string{"id", "name"})
	if cols[0] != `"id"` || cols[1] != `"name"` {
		t.Errorf("postgres insertColumns = %v", cols)
	}
	if got := updateColExpr("users", "name"); got != `"name" = ?` {
		t.Errorf("postgres updateColExpr = %q", got)
	}
}

func TestDBConfigStringPerDriver(t *testing.T) {
	mysql := DBConfig{
		Driver: "mysql", User: "root", Password: "p", Host: "localhost", Port: 3306, Schema: "app",
	}
	if got, want := mysql.String(), "root:p@tcp(localhost:3306)/app?parseTime=True"; got != want {
		t.Errorf("mysql dsn = %q, want %q", got, want)
	}

	pg := DBConfig{
		Driver: "postgres", User: "postgres", Password: "secret", Host: "db.local", Port: 5432, Schema: "app",
		Flags: map[string]string{"sslmode": "require"},
	}
	want := "postgres://postgres:secret@db.local:5432/app?sslmode=require"
	if got := pg.String(); got != want {
		t.Errorf("postgres dsn = %q, want %q", got, want)
	}
}

func TestResolveDriver(t *testing.T) {
	withDialect(t, MysqlDialect)

	if _, err := resolveDriver("mysql"); err != nil {
		t.Errorf("mysql rejected: %+v", err)
	}
	if _, err := resolveDriver(""); err != nil {
		t.Errorf("empty driver (mysql default) rejected: %+v", err)
	}
	if _, err := resolveDriver("postgres"); err != nil {
		t.Errorf("postgres rejected: %+v", err)
	}
	if IsPostgres() == false && CurrentDialect().Name != "mysql" {
		t.Error("dialect not tracked")
	}
	if _, err := resolveDriver("oracle"); err != ErrUnknownDBDriver {
		t.Errorf("unknown driver error = %+v, want ErrUnknownDBDriver", err)
	}

	// postgres activation
	if name, err := resolveDriver("postgres"); err != nil || name != "pgx" {
		t.Errorf("postgres driver name = %q err=%+v, want pgx", name, err)
	}
	if !IsPostgres() {
		t.Error("resolveDriver(postgres) did not activate the postgres dialect")
	}
}
