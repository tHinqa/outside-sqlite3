// Copyright (c) 2013 Tony Wilson. All rights reserved.
// See LICENCE file for permissions and restrictions.

//Package sqlite3 provides API definitions for accessing the sqlite3 dll.
//Vendor documentation can be found at http://www.sqlite.org/c3ref/
package sqlite3

import (
	"github.com/tHinqa/outside"
	. "github.com/tHinqa/outside/types"
)

func init() {
	outside.AddDllApis(dll, false, apiList)
}

func UTF16() {
	outside.AddDllApis(dll, true, api16)
}

func UTF8() {
	outside.AddDllApis(dll, false, api8)
}

type (
	Char         int8
	Sqlite3      struct{}
	Backup       struct{}
	Blob         struct{}
	Context      struct{}
	Mutex        struct{}
	Stmt         struct{}
	Value        struct{}
	UnsignedChar uint8
	VaList       uintptr
	Private      struct{}
	Void         struct{}

	SyscallPtr func()

	// Callback func(*Void, int, **Char, **Char) int
	Callback func(*Void, int, *[1000]uintptr, *[1000]uintptr) int

	Vtab struct {
		Module *Module
		Ref    int
		ErrMsg *Char
	}

	VtabCursor struct {
		Vtab *Vtab
	}

	IndexConstraint struct {
		Column     int
		Op         UnsignedChar
		Usable     UnsignedChar
		TermOffset int
	}

	IndexOrderby struct {
		Column int
		Desc   UnsignedChar
	}

	IndexConstraintUsage struct {
		argvIndex int
		Omit      UnsignedChar
	}

	IndexInfo struct {
		ConstraintCount  int
		Constraint       *IndexConstraint
		OrderByCount     int
		OrderBy          *IndexOrderby
		ConstraintUsage  *IndexConstraintUsage
		IdxNum           int
		IdxStr           *Char
		NeedToFreeIdxStr int
		OrderByConsumed  int
		EstimatedCost    float64
	}

	Module struct {
		Version int
		Create  func(_ *Sqlite3, aux *Void, argc int,
			argv **Char, ppVTab **Vtab, _ **Char) int
		Connect func(_ *Sqlite3, aux *Void, argc int,
			argv **Char, vTab **Vtab, _ **Char) int
		BestIndex  func(v *Vtab, _ *IndexInfo) int
		Disconnect func(v *Vtab) int
		Destroy    func(v *Vtab) int
		Open       func(v *Vtab, c **VtabCursor) int
		Close      func(*VtabCursor) int
		Filter     func(_ *VtabCursor, num int, str *Char,
			argc int, argv **Value) int
		Next         func(*VtabCursor) int
		Eof          func(*VtabCursor) int
		Column       func(*VtabCursor, *Context, int) int
		Rowid        func(_ *VtabCursor, rowid *int64) int
		Update       func(*Vtab, int, **Value, *int64) int
		Begin        func(pVTab *Vtab) int
		Sync         func(pVTab *Vtab) int
		Commit       func(pVTab *Vtab) int
		Rollback     func(pVTab *Vtab) int
		FindFunction func(vtab *Vtab, args int, name *Char,
			Func *func( //???
				*Context, int, **Value), argsPtr **Void) int
		Rename     func(vtab *Vtab, new *Char) int
		Savepoint  func(vTab *Vtab, _ int) int
		Release    func(vTab *Vtab, _ int) int
		RollbackTo func(vTab *Vtab, _ int) int
	}

	Vfs struct {
		Version  int
		OsFile   int
		Pathname int
		Next     *Vfs
		Name     *Char
		AppData  *Void
		Open     func(_ *Vfs, name *Char,
			_ *File, flags int, outFlags *int) int
		Delete func(_ *Vfs, name *Char, syncDir int) int
		Access func(_ *Vfs, name *Char,
			flags int, pResOut *int) int
		FullPathname func(_ *Vfs, name *Char,
			nOut int, out *Char) int
		DlOpen  func(_ *Vfs, filename *Char)
		DlError func(_ *Vfs, nByte int, errMsg *Char)
		//TODO(t): This is why we use Go
		//void (*(*xDlSym)(*Vfs,*Void, symbol  *Char))(void);
		DlClose          func(*Vfs, *Void)
		Randomness       func(_ *Vfs, bytes int, out *Char) int
		Sleep            func(_ *Vfs, microseconds int) int
		CurrentTime      func(*Vfs, *float64) int
		GetLastError     func(*Vfs, int, *Char) int
		CurrentTimeInt64 func(*Vfs, *int64) int
		SetSystemCall    func(
			_ *Vfs, name *Char, _ SyscallPtr) int
		GetSystemCall  func(_ *Vfs, name *Char) SyscallPtr
		NextSystemCall func(_ *Vfs, name *Char) *Char
	}

	File struct {
		Methods *IOMethods
	}

	IOMethods struct {
		Version int
		Close   func(_ *File) int
		Read    func(_ *File, _ *Void,
			amt int, ofst int64) int
		Write func(_ *File, _ *Void,
			amt int, ofst int64) int
		Truncate func(_ *File, size int64) int
		Sync     func(_ *File, flags int) int
		FileSize func(
			_ *File, size *int64) int
		Lock              func(*File, int) int
		Unlock            func(*File, int) int
		CheckReservedLock func(_ *File, resOut *int) int
		FileControl       func(
			_ *File, op int, arg *Void) int
		SectorSize            func(*File) int
		DeviceCharacteristics func(*File) int
		ShmMap                func(
			_ *File, pg int, pgsz int, _ int, _ **Void) int
		ShmLock func(
			_ *File, offset int, n int, flags int) int
		ShmBarrier func(*File)
		ShmUnmap   func(_ *File, deleteFlag int) int
		Fetch      func(
			_ *File, ofst int64, amt int, pp **Void) int
		Unfetch func(
			_ *File, ofst int64, p *Void) int
	}

	RtreeGeometry struct {
		Context    *Void
		ParamCount int
		Param      *float64
		User       *Void
		DelUser    func(*Void)
	}
)

var (
	Initialize func() int

	Shutdown func() int

	OSInit func() int

	OSEnd func() int

	Config func(int, ...VArg) int

	DbConfig func(_ *Sqlite3, op int, _ ...VArg) int

	ExtendedResultCodes func(_ *Sqlite3, onoff int) int

	LastInsertRowid func(*Sqlite3) int64

	Changes func(*Sqlite3) int

	TotalChanges func(*Sqlite3) int

	Interrupt func(*Sqlite3)

	Complete func(sql VString) int

	BusyHandler func(
		*Sqlite3, func(*Void, int) int, *Void) int

	BusyTimeout func(_ *Sqlite3, ms int) int

	GetTable func(db *Sqlite3, sql string,
		result ***Char, row *int, column *int,
		errmsg **Char) int

	FreeTable func(result **Char)

	Mprintf func(*Char, ...VArg) string

	Vmprintf func(*Char, VaList) string

	Snprintf func(int, string, string, ...VArg) string

	Vsnprintf func(int, string, string, VaList) string

	Malloc func(int) *Void

	Realloc func(*Void, int) *Void

	Free func(*Void)

	MemoryUsed func() int64

	MemoryHighwater func(resetFlag int) int64

	Randomness func(N int, P *Void)

	SetAuthorizer func(_ *Sqlite3,
		auth func(*Void, int, string, string, string, string) int,
		userData *Void) int

	Trace func(_ *Sqlite3,
		trace func(*Void, string),
		_ *Void) *Void

	Profile func(_ *Sqlite3,
		profile func(*Void, string, uint64),
		_ *Void) *Void

	ProgressHandler func(
		*Sqlite3, int, func(*Void) int, *Void)

	Open func(filename VString, db **Sqlite3) int

	OpenV2 func(filename string, db **Sqlite3,
		flags OpenFlags, vfs string) int

	UriParameter func(
		filename string, param string) string

	UriBoolean func(
		file string, param string, deflt int) int

	UriInt64 func(
		string, string, int64) int64

	Errcode func(db *Sqlite3) int

	ExtendedErrcode func(db *Sqlite3) int

	Errmsg func(*Sqlite3) VString

	Errstr func(int) string

	Limit func(_ *Sqlite3, id int, newVal int) int

	Prepare func(db *Sqlite3, sql string,
		bytes int, stmt **Stmt, tail **Char) int

	PrepareV2 func(db *Sqlite3, sql string,
		bytes int, stmt **Stmt, tail **Char) int

	Prepare16 func(db *Sqlite3, sql *Void,
		bytes int, stmt **Stmt, tail **Void) int

	Prepare16V2 func(db *Sqlite3, sql *Void,
		bytes int, stmt **Stmt, tail **Void) int

	Sql func(pStmt *Stmt) string

	StmtReadonly func(*Stmt) int

	StmtBusy func(*Stmt) int

	BindBlob func(_ *Stmt,
		_ int, _ *byte, n int, _ func(*Void)) int

	BindDouble func(
		*Stmt, int, float64) int

	BindInt func(
		*Stmt, int, int) int

	BindInt64 func(*Stmt, int, int64) int

	BindNull func(*Stmt, int) int

	BindText func(_ *Stmt,
		_ int, _ VString, n int, _ func(*Void)) int

	BindValue func(
		*Stmt, int, *Value) int

	BindZeroblob func(_ *Stmt, _ int, n int) int

	BindParameterCount func(*Stmt) int

	BindParameterName func(*Stmt, int) string

	BindParameterIndex func(
		_ *Stmt, name string) int

	ClearBindings func(*Stmt) int

	ColumnCount func(pStmt *Stmt) int

	ColumnName func(_ *Stmt, N int) VString

	ColumnDatabaseName func(*Stmt, int) VString

	ColumnTableName func(*Stmt, int) VString

	ColumnOriginName func(*Stmt, int) VString

	ColumnDecltype func(*Stmt, int) VString

	Step func(*Stmt) int

	DataCount func(pStmt *Stmt) int

	ColumnBlob func(_ *Stmt, col int) *Void

	ColumnBytes func(_ *Stmt, col int) int

	ColumnDouble func(_ *Stmt, col int) float64

	ColumnInt func(_ *Stmt, col int) int

	ColumnInt64 func(_ *Stmt, col int) int64

	ColumnText func(_ *Stmt, col int) VString

	ColumnType func(_ *Stmt, col int) int

	ColumnValue func(
		_ *Stmt, col int) *Value

	Finalize func(stmt *Stmt) int

	Reset func(stmt *Stmt) int

	CreateFunction func(db *Sqlite3,
		functionName VString, args, textRep int, app *Void,
		fnc func(*Context, int, **Value),
		step func(*Context, int, **Value),
		final func(*Context)) int

	CreateFunctionV2 func(
		db *Sqlite3,
		functionName string,
		args int,
		textRep int,
		app *Void,
		fnc func(
			*Context, int, **Value),
		step func(*Context, int, **Value),
		final func(*Context),
		destroy func(*Void)) int

	AggregateCount func(*Context) int

	Expired func(*Stmt) int

	TransferBindings func(*Stmt, *Stmt) int

	GlobalRecover func() int

	ThreadCleanup func()

	MemoryAlarm func(
		func(*Void, int64, int),
		*Void,
		int64) int

	ValueBlob func(*Value) *Void

	ValueBytes func(*Value) int

	ValueDouble func(*Value) float64

	ValueInt func(*Value) int

	ValueInt64 func(*Value) int64

	ValueText func(*Value) VString

	ValueText16le func(*Value) *Void

	ValueText16be func(*Value) *Void

	ValueType func(*Value) int

	ValueNumericType func(*Value) int

	AggregateContext func(
		_ *Context, nBytes int) *Void

	UserData func(*Context) *Void

	ContextDbHandle func(*Context) *Sqlite3

	GetAuxdata func(_ *Context, N int) *Void

	SetAuxdata func(
		_ *Context, N int, _ *Void, _ func(*Void))

	ResultBlob func(
		*Context, *Void, int, func(*Void))

	ResultDouble func(*Context, float64)

	ResultError func(*Context, VString, int)

	ResultErrorToobig func(*Context)

	ResultErrorNomem func(*Context)

	ResultErrorCode func(*Context, int)

	ResultInt func(*Context, int)

	ResultInt64 func(*Context, int64)

	ResultNull func(*Context)

	ResultText func(
		*Context, VString, int, func(*Void))

	ResultText16le func(
		*Context, *Void, int, func(*Void))

	ResultText16be func(
		*Context, *Void, int, func(*Void))

	ResultValue func(
		*Context, *Value)

	ResultZeroblob func(_ *Context, n int)

	CreateCollation func(
		_ *Sqlite3, name VString, textRep int, arg *Void,
		compare func(*Void, int, *Void, int, *Void) int) int

	CreateCollationV2 func(
		_ *Sqlite3, name string, textRep int, arg *Void,
		compare func(*Void, int, *Void, int, *Void) int,
		destroy func(*Void)) int

	CollationNeeded func(*Sqlite3, *Void,
		func(_ *Void, _ *Sqlite3, textRep int, _ VString)) int

	Sleep func(int) int

	GetAutocommit func(*Sqlite3) int

	DbHandle func(*Stmt) *Sqlite3

	DbFilename func(db *Sqlite3, dbName string) string

	DbReadonly func(db *Sqlite3, dbName string) int

	NextStmt func(
		db *Sqlite3, stmt *Stmt) *Stmt

	CommitHook func(
		*Sqlite3, func(*Void) int, *Void) *Void

	RollbackHook func(*Sqlite3, func(*Void), *Void) *Void

	UpdateHook func(*Sqlite3,
		func(*Void, int, string, string, int64),
		*Void) *Void

	EnableSharedCache func(int) int

	ReleaseMemory func(int) int

	DbReleaseMemory func(*Sqlite3) int

	SoftHeapLimit64 func(int64) int64

	SoftHeapLimit func(int)

	TableColumnMetadata func(db *Sqlite3,
		dbName, tableName, columnName string,
		dataType, collSeq **Char,
		notNull, primaryKey, autoinc *int) int

	LoadExtension func(db *Sqlite3,
		file, proc string, errMsg **Char) int

	EnableLoadExtension func(db *Sqlite3, onoff int) int

	AutoExtension func(entryPoint func()) int

	CancelAutoExtension func(entryPoint func()) int

	ResetAutoExtension func()

	CreateModule func(db *Sqlite3,
		name string, p *Module, clientData *Void) int

	CreateModuleV2 func(db *Sqlite3,
		name string, p *Module, clientData *Void,
		destroy func(*Void)) int

	DeclareVtab func(_ *Sqlite3, sql string) int

	OverloadFunction func(
		_ *Sqlite3, funcName string, args int) int

	BlobReopen func(*Blob, int64) int

	BlobClose func(*Blob) int

	BlobBytes func(*Blob) int

	BlobRead func(
		_ *Blob, buffer *Void, bytes int, offset int) int

	BlobWrite func(
		_ *Blob, buffer *Void, bytes int, offset int) int

	VfsFind func(vfsName string) *Vfs

	VfsRegister func(_ *Vfs, makeDflt int) int

	VfsUnregister func(*Vfs) int

	MutexAlloc func(int) *Mutex

	MutexFree func(*Mutex)

	MutexEnter func(*Mutex)

	MutexTry func(*Mutex) int

	MutexLeave func(*Mutex)

	MutexHeld func(*Mutex) int

	MutexNotheld func(*Mutex) int

	DbMutex func(*Sqlite3) *Mutex

	FileControl func(
		_ *Sqlite3, dbName string, op int, _ *Void) int

	TestControl func(op int, _ ...VArg) int

	Status func(
		op int, current, highwater *int, resetFlag int) int

	DbStatus func(
		_ *Sqlite3, op int, cur, hiwtr *int, resetFlg int) int

	StmtStatus func(
		_ *Stmt, op int, resetFlg int) int

	Close func(*Sqlite3) int

	CloseV2 func(*Sqlite3) int

	Exec func(_ *Sqlite3, sql string,
		callback Callback, _ *Void, errmsg **Char) int

	Threadsafe func() int

	Libversion func() string

	Sourceid func() string

	LibversionNumber func() int

	CompileoptionUsed func(optName string) int

	CompileoptionGet func(n int) string

	BackupStep func(
		p *Backup,
		nPage int) int

	BackupFinish func(
		p *Backup) int

	BackupRemaining func(
		p *Backup) int

	BackupPagecount func(
		p *Backup) int

	UnlockNotify func(
		pBlocked *Sqlite3,
		notify func(arg **Void, nArg int), notifyArg *Void) int

	Stricmp func(string, string) int

	Strnicmp func(string, string, int) int

	Strglob func(glob string, str string) int

	Log func(errCode int, format string, v ...VArg)

	WalHook func(
		*Sqlite3, func(
			*Void, *Sqlite3, string, int) int, *Void) *Void

	WalAutocheckpoint func(db *Sqlite3, N int) int

	WalCheckpoint func(db *Sqlite3, dbName string) int

	WalCheckpointV2 func(db *Sqlite3,
		dbName string, mode int, nLog, nCkpt *int) int

	VtabConfig func(_ *Sqlite3, op int, v ...VArg) int

	VtabOnConflict func(*Sqlite3) int

	RtreeGeometryCallback func(
		db *Sqlite3,
		geom string,
		geomfunc func(
			_ *RtreeGeometry,
			n int,
			a *float64,
			res *int) int,
		context *Void) int

	BackupInit func(
		dest *Sqlite3, destName string,
		source *Sqlite3, sourceName string) *Backup

	BlobOpen func(
		_ *Sqlite3,
		db, table, column string,
		row int64,
		flags int,
		blob **Blob) int
)

const (
	OK = iota
	ERROR
	INTERNAL
	PERM
	ABORT
	BUSY
	LOCKED
	NOMEM
	READONLY
	INTERRUPT
	IOERR
	CORRUPT
	NOTFOUND
	FULL
	CANTOPEN
	PROTOCOL
	EMPTY
	SCHEMA
	TOOBIG
	CONSTRAINT
	MISMATCH
	MISUSE
	NOLFS
	AUTH
	FORMAT
	RANGE
	NOTADB
	NOTICE
	WARNING
)
const (
	ROW = iota + 100
	DONE
)

const (
	INTEGER = iota + 1
	FLOAT
	TEXT
	BLOB
	NULL
)

type OpenFlags int

const (
	OPEN_READONLY OpenFlags = 1 << iota
	OPEN_READWRITE
	OPEN_CREATE
	OPEN_DELETEONCLOSE
	OPEN_EXCLUSIVE
	OPEN_AUTOPROXY
	OPEN_URI
	OPEN_MEMORY
	OPEN_MAIN_DB
	OPEN_TEMP_DB
	OPEN_TRANSIENT_DB
	OPEN_MAIN_JOURNAL
	OPEN_TEMP_JOURNAL
	OPEN_SUBJOURNAL
	OPEN_MASTER_JOURNAL
	OPEN_NOMUTEX
	OPEN_FULLMUTEX
	OPEN_SHAREDCACHE
	OPEN_PRIVATECACHE
	OPEN_WAL
)

var dll = "Sqlite3.dll"

var apiList = outside.Apis{
	{"sqlite3_aggregate_context", &AggregateContext},
	{"sqlite3_aggregate_count", &AggregateCount},
	{"sqlite3_auto_extension", &AutoExtension},
	{"sqlite3_backup_finish", &BackupFinish},
	{"sqlite3_backup_init", &BackupInit},
	{"sqlite3_backup_pagecount", &BackupPagecount},
	{"sqlite3_backup_remaining", &BackupRemaining},
	{"sqlite3_backup_step", &BackupStep},
	{"sqlite3_bind_blob", &BindBlob},
	{"sqlite3_bind_double", &BindDouble},
	{"sqlite3_bind_int", &BindInt},
	{"sqlite3_bind_int64", &BindInt64},
	{"sqlite3_bind_null", &BindNull},
	{"sqlite3_bind_parameter_count", &BindParameterCount},
	{"sqlite3_bind_parameter_index", &BindParameterIndex},
	{"sqlite3_bind_parameter_name", &BindParameterName},
	{"sqlite3_bind_value", &BindValue},
	{"sqlite3_bind_zeroblob", &BindZeroblob},
	{"sqlite3_blob_bytes", &BlobBytes},
	{"sqlite3_blob_close", &BlobClose},
	{"sqlite3_blob_open", &BlobOpen},
	{"sqlite3_blob_read", &BlobRead},
	{"sqlite3_blob_reopen", &BlobReopen},
	{"sqlite3_blob_write", &BlobWrite},
	{"sqlite3_busy_handler", &BusyHandler},
	{"sqlite3_busy_timeout", &BusyTimeout},
	{"sqlite3_cancel_auto_extension", &CancelAutoExtension},
	{"sqlite3_changes", &Changes},
	{"sqlite3_clear_bindings", &ClearBindings},
	{"sqlite3_close", &Close},
	{"sqlite3_close_v2", &CloseV2},
	{"sqlite3_column_blob", &ColumnBlob},
	{"sqlite3_column_count", &ColumnCount},
	{"sqlite3_column_double", &ColumnDouble},
	{"sqlite3_column_int", &ColumnInt},
	{"sqlite3_column_int64", &ColumnInt64},
	{"sqlite3_column_type", &ColumnType},
	{"sqlite3_column_value", &ColumnValue},
	{"sqlite3_commit_hook", &CommitHook},
	{"sqlite3_compileoption_get", &CompileoptionGet},
	{"sqlite3_compileoption_used", &CompileoptionUsed},
	{"sqlite3_config", &Config},
	{"sqlite3_context_db_handle", &ContextDbHandle},
	{"sqlite3_create_collation_v2", &CreateCollationV2},
	{"sqlite3_create_function_v2", &CreateFunctionV2},
	{"sqlite3_create_module", &CreateModule},
	{"sqlite3_create_module_v2", &CreateModuleV2},
	{"sqlite3_data_count", &DataCount},
	{"sqlite3_db_config", &DbConfig},
	{"sqlite3_db_filename", &DbFilename},
	{"sqlite3_db_handle", &DbHandle},
	{"sqlite3_db_mutex", &DbMutex},
	{"sqlite3_db_readonly", &DbReadonly},
	{"sqlite3_db_release_memory", &DbReleaseMemory},
	{"sqlite3_db_status", &DbStatus},
	{"sqlite3_declare_vtab", &DeclareVtab},
	{"sqlite3_enable_load_extension", &EnableLoadExtension},
	{"sqlite3_enable_shared_cache", &EnableSharedCache},
	{"sqlite3_errcode", &Errcode},
	{"sqlite3_errstr", &Errstr},
	{"sqlite3_exec", &Exec},
	{"sqlite3_expired", &Expired},
	{"sqlite3_extended_errcode", &ExtendedErrcode},
	{"sqlite3_extended_result_codes", &ExtendedResultCodes},
	{"sqlite3_file_control", &FileControl},
	{"sqlite3_finalize", &Finalize},
	{"sqlite3_free", &Free},
	{"sqlite3_free_table", &FreeTable},
	{"sqlite3_get_autocommit", &GetAutocommit},
	{"sqlite3_get_auxdata", &GetAuxdata},
	{"sqlite3_get_table", &GetTable},
	{"sqlite3_global_recover", &GlobalRecover},
	{"sqlite3_initialize", &Initialize},
	{"sqlite3_interrupt", &Interrupt},
	{"sqlite3_last_insert_rowid", &LastInsertRowid},
	{"sqlite3_libversion", &Libversion},
	{"sqlite3_libversion_number", &LibversionNumber},
	{"sqlite3_limit", &Limit},
	{"sqlite3_load_extension", &LoadExtension},
	{"sqlite3_log", &Log},
	{"sqlite3_malloc", &Malloc},
	{"sqlite3_memory_alarm", &MemoryAlarm},
	{"sqlite3_memory_highwater", &MemoryHighwater},
	{"sqlite3_memory_used", &MemoryUsed},
	{"sqlite3_mprintf", &Mprintf},
	{"sqlite3_mutex_alloc", &MutexAlloc},
	{"sqlite3_mutex_enter", &MutexEnter},
	{"sqlite3_mutex_free", &MutexFree},
	{"sqlite3_mutex_leave", &MutexLeave},
	{"sqlite3_mutex_try", &MutexTry},
	{"sqlite3_next_stmt", &NextStmt},
	{"sqlite3_open_v2", &OpenV2},
	{"sqlite3_os_end", &OSEnd},
	{"sqlite3_os_init", &OSInit},
	{"sqlite3_overload_function", &OverloadFunction},
	{"sqlite3_prepare", &Prepare},
	{"sqlite3_prepare16", &Prepare16},
	{"sqlite3_prepare16_v2", &Prepare16V2},
	{"sqlite3_prepare_v2", &PrepareV2},
	{"sqlite3_profile", &Profile},
	{"sqlite3_progress_handler", &ProgressHandler},
	{"sqlite3_randomness", &Randomness},
	{"sqlite3_realloc", &Realloc},
	{"sqlite3_release_memory", &ReleaseMemory},
	{"sqlite3_reset", &Reset},
	{"sqlite3_reset_auto_extension", &ResetAutoExtension},
	{"sqlite3_result_blob", &ResultBlob},
	{"sqlite3_result_double", &ResultDouble},
	{"sqlite3_result_error_code", &ResultErrorCode},
	{"sqlite3_result_error_nomem", &ResultErrorNomem},
	{"sqlite3_result_error_toobig", &ResultErrorToobig},
	{"sqlite3_result_int", &ResultInt},
	{"sqlite3_result_int64", &ResultInt64},
	{"sqlite3_result_null", &ResultNull},
	{"sqlite3_result_text16be", &ResultText16be},
	{"sqlite3_result_text16le", &ResultText16le},
	{"sqlite3_result_value", &ResultValue},
	{"sqlite3_result_zeroblob", &ResultZeroblob},
	{"sqlite3_rollback_hook", &RollbackHook},
	{"sqlite3_rtree_geometry_callback", &RtreeGeometryCallback},
	{"sqlite3_set_authorizer", &SetAuthorizer},
	{"sqlite3_set_auxdata", &SetAuxdata},
	{"sqlite3_shutdown", &Shutdown},
	{"sqlite3_sleep", &Sleep},
	{"sqlite3_snprintf", &Snprintf},
	{"sqlite3_soft_heap_limit", &SoftHeapLimit},
	{"sqlite3_soft_heap_limit64", &SoftHeapLimit64},
	{"sqlite3_sourceid", &Sourceid},
	{"sqlite3_sql", &Sql},
	{"sqlite3_status", &Status},
	{"sqlite3_step", &Step},
	{"sqlite3_stmt_busy", &StmtBusy},
	{"sqlite3_stmt_readonly", &StmtReadonly},
	{"sqlite3_stmt_status", &StmtStatus},
	{"sqlite3_strglob", &Strglob},
	{"sqlite3_stricmp", &Stricmp},
	{"sqlite3_strnicmp", &Strnicmp},
	{"sqlite3_table_column_metadata", &TableColumnMetadata},
	{"sqlite3_test_control", &TestControl},
	{"sqlite3_thread_cleanup", &ThreadCleanup},
	{"sqlite3_threadsafe", &Threadsafe},
	{"sqlite3_total_changes", &TotalChanges},
	{"sqlite3_trace", &Trace},
	{"sqlite3_transfer_bindings", &TransferBindings},
	{"sqlite3_update_hook", &UpdateHook},
	{"sqlite3_uri_boolean", &UriBoolean},
	{"sqlite3_uri_int64", &UriInt64},
	{"sqlite3_uri_parameter", &UriParameter},
	{"sqlite3_user_data", &UserData},
	{"sqlite3_value_blob", &ValueBlob},
	{"sqlite3_value_double", &ValueDouble},
	{"sqlite3_value_int", &ValueInt},
	{"sqlite3_value_int64", &ValueInt64},
	{"sqlite3_value_numeric_type", &ValueNumericType},
	{"sqlite3_value_text16be", &ValueText16be},
	{"sqlite3_value_text16le", &ValueText16le},
	{"sqlite3_value_type", &ValueType},
	{"sqlite3_vfs_find", &VfsFind},
	{"sqlite3_vfs_register", &VfsRegister},
	{"sqlite3_vfs_unregister", &VfsUnregister},
	{"sqlite3_vmprintf", &Vmprintf},
	{"sqlite3_vsnprintf", &Vsnprintf},
	{"sqlite3_vtab_config", &VtabConfig},
	{"sqlite3_vtab_on_conflict", &VtabOnConflict},
	{"sqlite3_wal_autocheckpoint", &WalAutocheckpoint},
	{"sqlite3_wal_checkpoint", &WalCheckpoint},
	{"sqlite3_wal_checkpoint_v2", &WalCheckpointV2},
	{"sqlite3_wal_hook", &WalHook},
	// Undocumented {"sqlite3_win32_mbcs_to_utf8", &Win32MbcsToUtf8},
	// Undocumented {"sqlite3_win32_set_directory", &Win32SetDirectory},
	// Undocumented {"sqlite3_win32_sleep", &Win32Sleep},
	// Undocumented {"sqlite3_win32_utf8_to_mbcs", &Win32Utf8ToMbcs},
	// Undocumented {"sqlite3_win32_write_debug", &Win32WriteDebug},
}

var api16 = outside.Apis{
	{"sqlite3_bind_text16", &BindText},
	{"sqlite3_collation_needed16", &CollationNeeded},
	{"sqlite3_column_bytes16", &ColumnBytes},
	{"sqlite3_column_database_name16", &ColumnDatabaseName},
	{"sqlite3_column_decltype16", &ColumnDecltype},
	{"sqlite3_column_name16", &ColumnName},
	{"sqlite3_column_origin_name16", &ColumnOriginName},
	{"sqlite3_column_table_name16", &ColumnTableName},
	{"sqlite3_column_text16", &ColumnText},
	{"sqlite3_complete16", &Complete},
	{"sqlite3_create_collation16", &CreateCollation},
	{"sqlite3_create_function16", &CreateFunction},
	{"sqlite3_errmsg16", &Errmsg},
	{"sqlite3_open16", &Open},
	{"sqlite3_result_error16", &ResultError},
	{"sqlite3_result_text16", &ResultText},
	{"sqlite3_value_bytes16", &ValueBytes},
	{"sqlite3_value_text16", &ValueText},
}

var api8 = outside.Apis{
	{"sqlite3_bind_text", &BindText},
	{"sqlite3_collation_needed", &CollationNeeded},
	{"sqlite3_column_bytes", &ColumnBytes},
	{"sqlite3_column_database_name", &ColumnDatabaseName},
	{"sqlite3_column_decltype", &ColumnDecltype},
	{"sqlite3_column_name", &ColumnName},
	{"sqlite3_column_origin_name", &ColumnOriginName},
	{"sqlite3_column_table_name", &ColumnTableName},
	{"sqlite3_column_text", &ColumnText},
	{"sqlite3_complete", &Complete},
	{"sqlite3_create_collation", &CreateCollation},
	{"sqlite3_create_function", &CreateFunction},
	{"sqlite3_errmsg", &Errmsg},
	{"sqlite3_open", &Open},
	{"sqlite3_result_error", &ResultError},
	{"sqlite3_result_text", &ResultText},
	{"sqlite3_value_bytes", &ValueBytes},
	{"sqlite3_value_text", &ValueText},
}
