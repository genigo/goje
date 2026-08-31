package goje

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"sync"
	"testing"
)

// testDriverName is the name the fake driver below is registered under.
const testDriverName = "goje-test-driver"

// testDriver is registered exactly once: sql.Register panics on duplicate names.
var testDriver = &fakeDriver{}

func init() {
	sql.Register(testDriverName, testDriver)
}

// fakeDriver records every DSN sql.Open hands over. It never talks to a real
// database: OpenConnector runs eagerly inside sql.Open (that is how we capture
// the DSN without opening a connection) and Connect always fails.
type fakeDriver struct {
	mu   sync.Mutex
	dsns []string
}

func (d *fakeDriver) Open(name string) (driver.Conn, error) {
	d.record(name)
	return nil, errors.New("goje test driver can't open connections")
}

func (d *fakeDriver) OpenConnector(name string) (driver.Connector, error) {
	d.record(name)
	return &fakeConnector{}, nil
}

func (d *fakeDriver) record(name string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.dsns = append(d.dsns, name)
}

// lastDSN returns the most recent DSN the driver saw
func (d *fakeDriver) lastDSN() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	if len(d.dsns) == 0 {
		return ""
	}
	return d.dsns[len(d.dsns)-1]
}

// fakeConnector satisfies driver.Connector without touching the network
type fakeConnector struct{}

func (c *fakeConnector) Connect(context.Context) (driver.Conn, error) {
	return nil, errors.New("goje test driver can't connect")
}

func (c *fakeConnector) Driver() driver.Driver {
	return testDriver
}

func (c *fakeConnector) Close() error {
	return nil
}

func TestNewDBConnectionWithDriverPassesDSN(t *testing.T) {
	prevDB, prevDialect := DefatultDB, dialect
	defer func() {
		DefatultDB = prevDB
		dialect = prevDialect
	}()

	conn := &DBConfig{
		Driver:   DriverMysql,
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "u",
		Password: "p",
		Schema:   "s",
	}

	db, err := NewDBConnectionWithDriver(testDriverName, conn)
	if err != nil {
		t.Fatalf("NewDBConnectionWithDriver(%q): %+v", testDriverName, err)
	}
	defer db.Close()

	if got, want := testDriver.lastDSN(), conn.String(); got != want {
		t.Errorf("dsn = %q, want %q", got, want)
	}
}

func TestNewDBConnectionWithDriverKeepsDialect(t *testing.T) {
	prevDB, prevDialect := DefatultDB, dialect
	defer func() {
		DefatultDB = prevDB
		dialect = prevDialect
	}()

	// the dialect follows conn.Driver, not the custom driver name
	cases := []struct {
		driver string
		pg     bool
	}{
		{DriverMysql, false},
		{DriverPostgres, true},
	}

	for _, c := range cases {
		_, err := NewDBConnectionWithDriver(testDriverName, &DBConfig{
			Driver:   c.driver,
			Host:     "127.0.0.1",
			Port:     3306,
			User:     "u",
			Password: "p",
			Schema:   "s",
		})
		if err != nil {
			t.Fatalf("driver %q: %+v", c.driver, err)
		}
		if IsPostgres() != c.pg {
			t.Errorf("driver %q: IsPostgres() = %v, want %v", c.driver, IsPostgres(), c.pg)
		}
	}
}

func TestNewDBConnectionWithDriverUnknownDriver(t *testing.T) {
	prevDB, prevDialect := DefatultDB, dialect
	defer func() {
		DefatultDB = prevDB
		dialect = prevDialect
	}()

	// the custom driver name is not a way around the driver whitelist
	_, err := NewDBConnectionWithDriver(testDriverName, &DBConfig{Driver: "oracle"})
	if !errors.Is(err, ErrUnknownDBDriver) {
		t.Errorf("err = %+v, want ErrUnknownDBDriver", err)
	}
}

func TestInitDBWithDriver(t *testing.T) {
	prevDB, prevDialect := DefatultDB, dialect
	defer func() {
		DefatultDB = prevDB
		dialect = prevDialect
	}()

	err := InitDBWithDriver(testDriverName, &DBConfig{
		Driver:   DriverMysql,
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "u",
		Password: "p",
		Schema:   "s",
	})
	if err != nil {
		t.Fatalf("InitDBWithDriver(%q): %+v", testDriverName, err)
	}

	if DefatultDB == nil {
		t.Error("InitDBWithDriver did not set DefatultDB")
	}
}
