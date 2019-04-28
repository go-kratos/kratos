package dump

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/juju/errors"
	. "github.com/siddontang/go-mysql/mysql"
)

// Unlick mysqldump, Dumper is designed for parsing and syning data easily.
type Dumper struct {
	// mysqldump execution path, like mysqldump or /usr/bin/mysqldump, etc...
	ExecutionPath string

	Addr     string
	User     string
	Password string

	// Will override Databases
	Tables  []string
	TableDB string

	Databases []string

	Charset string

	IgnoreTables map[string][]string

	ErrOut io.Writer

	masterDataSkipped bool
	maxAllowedPacket  int
}

func NewDumper(executionPath string, addr string, user string, password string) (*Dumper, error) {
	if len(executionPath) == 0 {
		return nil, nil
	}

	path, err := exec.LookPath(executionPath)
	if err != nil {
		return nil, errors.Trace(err)
	}

	d := new(Dumper)
	d.ExecutionPath = path
	d.Addr = addr
	d.User = user
	d.Password = password
	d.Tables = make([]string, 0, 16)
	d.Databases = make([]string, 0, 16)
	d.Charset = DEFAULT_CHARSET
	d.IgnoreTables = make(map[string][]string)
	d.masterDataSkipped = false

	d.ErrOut = os.Stderr

	return d, nil
}

func (d *Dumper) SetCharset(charset string) {
	d.Charset = charset
}

func (d *Dumper) SetErrOut(o io.Writer) {
	d.ErrOut = o
}

// In some cloud MySQL, we have no privilege to use `--master-data`.
func (d *Dumper) SkipMasterData(v bool) {
	d.masterDataSkipped = v
}

func (d *Dumper) SetMaxAllowedPacket(i int) {
	d.maxAllowedPacket = i
}

func (d *Dumper) AddDatabases(dbs ...string) {
	d.Databases = append(d.Databases, dbs...)
}

func (d *Dumper) AddTables(db string, tables ...string) {
	if d.TableDB != db {
		d.TableDB = db
		d.Tables = d.Tables[0:0]
	}

	d.Tables = append(d.Tables, tables...)
}

func (d *Dumper) AddIgnoreTables(db string, tables ...string) {
	t, _ := d.IgnoreTables[db]
	t = append(t, tables...)
	d.IgnoreTables[db] = t
}

func (d *Dumper) Reset() {
	d.Tables = d.Tables[0:0]
	d.TableDB = ""
	d.IgnoreTables = make(map[string][]string)
	d.Databases = d.Databases[0:0]
}

func (d *Dumper) Dump(w io.Writer) error {
	args := make([]string, 0, 16)

	// Common args
	seps := strings.Split(d.Addr, ":")
	args = append(args, fmt.Sprintf("--host=%s", seps[0]))
	if len(seps) > 1 {
		args = append(args, fmt.Sprintf("--port=%s", seps[1]))
	}

	args = append(args, fmt.Sprintf("--user=%s", d.User))
	args = append(args, fmt.Sprintf("--password=%s", d.Password))

	if !d.masterDataSkipped {
		args = append(args, "--master-data")
	}

	if d.maxAllowedPacket > 0 {
		// mysqldump param should be --max-allowed-packet=%dM not be --max_allowed_packet=%dM
		args = append(args, fmt.Sprintf("--max-allowed-packet=%dM", d.maxAllowedPacket))
	}

	args = append(args, "--single-transaction")
	args = append(args, "--skip-lock-tables")

	// Disable uncessary data
	args = append(args, "--compact")
	args = append(args, "--skip-opt")
	args = append(args, "--quick")

	// We only care about data
	args = append(args, "--no-create-info")

	// Multi row is easy for us to parse the data
	args = append(args, "--skip-extended-insert")

	for db, tables := range d.IgnoreTables {
		for _, table := range tables {
			args = append(args, fmt.Sprintf("--ignore-table=%s.%s", db, table))
		}
	}

	if len(d.Tables) == 0 && len(d.Databases) == 0 {
		args = append(args, "--all-databases")
	} else if len(d.Tables) == 0 {
		args = append(args, "--databases")
		args = append(args, d.Databases...)
	} else {
		args = append(args, d.TableDB)
		args = append(args, d.Tables...)

		// If we only dump some tables, the dump data will not have database name
		// which makes us hard to parse, so here we add it manually.

		w.Write([]byte(fmt.Sprintf("USE `%s`;\n", d.TableDB)))
	}

	if len(d.Charset) != 0 {
		args = append(args, fmt.Sprintf("--default-character-set=%s", d.Charset))
	}

	cmd := exec.Command(d.ExecutionPath, args...)

	cmd.Stderr = d.ErrOut
	cmd.Stdout = w

	return cmd.Run()
}

// Dump MySQL and parse immediately
func (d *Dumper) DumpAndParse(h ParseHandler) error {
	r, w := io.Pipe()

	done := make(chan error, 1)
	go func() {
		err := Parse(r, h, !d.masterDataSkipped)
		r.CloseWithError(err)
		done <- err
	}()

	err := d.Dump(w)
	w.CloseWithError(err)

	err = <-done

	return errors.Trace(err)
}
