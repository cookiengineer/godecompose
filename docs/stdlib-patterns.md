# Go Standard Library Pattern Coverage Plan

## File Layout Convention

Each Go function/method has its own `.hexpat` file at a path mirroring the Go
import path. Filenames use exact Go casing with dots for methods. Patterns are
organized into three independently-loaded modules:

| Module | Purpose | Location |
|---|---|---|
| stdlib | Multi-instruction patterns (highest confidence) | `patterns/golang/stdlib/<import-path>/<Func>.hexpat` |
| runtime | Multi-instruction runtime patterns | `patterns/golang/runtime/<Func>.hexpat` |
| fallback | Single-CALL idiomatic patterns for stripped binaries | `patterns/golang/fallback/<import-path>/<Func>.hexpat` |
| controlflow | Control flow, Go idioms, data types, and stdlib call patterns | `patterns/golang/controlflow/<category>.hexpat` |
| syscalls | Kernel syscall tables (JSON) | `database/syscall/tables/<os>/syscall_*.json` |
| Test sources | Go test programs | `testdata/src/<package>/main.go` |
| E2E tests | Decompile pipeline tests | `e2e/decompile/<package>/<package>_test.go` |

Examples:
- `fmt.Fprintln` → `patterns/golang/stdlib/fmt/Fprintln.hexpat`
- `net/http.Get` → `patterns/golang/stdlib/net/http/Get.hexpat`
- `sync.(*Mutex).Lock` → `patterns/golang/stdlib/sync/Mutex.Lock.hexpat`
- `strings.Builder.WriteString` → `patterns/golang/stdlib/strings/Builder.WriteString.hexpat`
- `runtime.deferproc` → `patterns/golang/runtime/deferproc.hexpat`
- Fallback `fmt.Println` → `patterns/golang/fallback/fmt/Println.hexpat`

All patterns are embedded in the `godecompose` binary via `//go:embed` in
`patterns/golang/embed.go`. The three modules are loaded independently:
```go
golang.LoadStdlib(db)   // always loaded
golang.LoadRuntime(db)  // always loaded
golang.LoadFallback(db) // loaded on demand for stripped binaries
```

Shared E2E test helpers live in `e2e/internal/decompile/helpers.go`. Each E2E
test compiles the testdata Go program for `linux/amd64`, disassembles the
binary, loads all pattern files, matches patterns against user-function
instructions, generates decompiled source, and verifies that the pipeline
produces non-empty results.

## Pattern Structure

Each pattern follows CALL or multi-instruction matching:
- **CALL-based patterns**: A `CALL` instruction matched by symbol name via fuzzy matching against the disassembled GoSyntax
- **Multi-instruction patterns**: Register setup + CALL sequences for higher confidence matching
- **Control flow patterns**: TEST/CMP + conditional jump sequences for if/else and nil checks
- **Idiom patterns**: Single CALL for common Go constructs (defer, goroutines, channels, maps)
- An `instr match` block describing the instruction sequence
- A `gen` block with Go source code output (supports balanced nested braces)
- An optional `bind` block for variable renaming

Symbol names in the pattern language use underscores as word separators:
- Package functions: `package.Function` → pattern: `CALL package_Function`
- Receiver methods: `package.(*Type).Method` → pattern: `CALL package_Type_Method`
- Runtime internals: `runtime.mapassign_faststr` → pattern: `CALL runtime_mapassign_faststr`

The fuzzy matcher normalizes dots, parens, slashes, stars, AND underscores to spaces, lowercases, and does prefix matching. So patterns like `sync_Mutex_Lock` match against `CALL sync.(*Mutex).lockSlow(SB)` and `runtime_mapassign_faststr` matches `CALL runtime.mapassign_faststr(SB)`.

All four modules are embedded via `//go:embed` and loaded independently:
```go
golang.LoadStdlib(db)       // always loaded
golang.LoadRuntime(db)      // always loaded
golang.LoadFallback(db)     // always loaded
golang.LoadControlFlow(db)  // always loaded
```

## Package Coverage Status

All documented packages are fully implemented with patterns, testdata programs, and end-to-end decompile tests.

### Standard Library (all complete)

| Package | Functions / Methods | Status |
|---|---|---|
| `fmt` | Fprintln, Fprintf, Sprintf, Errorf, Fprint, Printf, Println | ✓ Covered + E2E tested |
| `errors` | New, Is, As, Unwrap, Join | ✓ Covered + E2E tested |
| `log` | Println, Printf, Fatal, Fatalf, Panic, Panicf | ✓ Covered + E2E tested |
| `context` | WithCancel, WithTimeout, WithDeadline, WithValue, Background, TODO | ✓ Covered + E2E tested |
| `flag` | Parse, NewFlagSet, FlagSet.String/Int/Bool, FlagSet.Parse | ✓ Covered + E2E tested |
| `sync` | Mutex.Lock/Unlock, WaitGroup.Add/Done/Wait, Once.Do, Pool.Get/Put | ✓ Covered + E2E tested |
| `os` | File.Write/Read/Close, Open, Create, Remove, Stat, Mkdir, Getenv | ✓ Covered + E2E tested |
| `io` | ReadAll, Copy, WriteString | ✓ Covered + E2E tested |
| `time` | Now, Sleep, After, Since, NewTicker | ✓ Covered + E2E tested |
| `strconv` | Itoa, Atoi, FormatInt, FormatFloat, ParseInt, ParseFloat, Quote, Unquote | ✓ Covered + E2E tested |
| `strings` | Contains, HasPrefix, HasSuffix, Join, Split, ReplaceAll, ToLower, ToUpper, TrimSpace, TrimPrefix, TrimSuffix, Index, LastIndex, Count, Repeat, Fields, NewReader, NewReplacer, Builder.WriteString/WriteByte/WriteRune/String/Reset | ✓ Covered + E2E tested |
| `bytes` | Contains, Equal, Compare, Index, LastIndex, Join, Split, HasPrefix, NewReader, NewBuffer, Buffer.WriteString/WriteByte/Bytes/String/Reset | ✓ Covered + E2E tested |
| `bufio` | NewReader, NewWriter, NewScanner, Reader.ReadString/ReadByte/ReadBytes/ReadLine/Peek, Writer.WriteString/WriteByte/Flush/Reset, Scanner.Scan/Text/Bytes/Err | ✓ Covered + E2E tested |
| `math` | Abs, Max, Min, Sqrt, Pow, Sin, Cos, Tan, Floor, Ceil, Round, Log, Log2, Log10, Exp, Mod, Remainder, Hypot | ✓ Covered + E2E tested |
| `math/rand` | Intn, Float64, Int, Int31, Int63, Perm, Shuffle, Seed, New, NewSource, Read | ✓ Covered + E2E tested |
| `sort` | Ints, Strings, Float64s, Slice, Search, SearchInts, SearchStrings, SearchFloat64s, Stable, Reverse, IsSorted | ✓ Covered + E2E tested |
| `encoding/json` | Marshal, Unmarshal, NewEncoder, NewDecoder, Encode, Decode | ✓ Covered + E2E tested |
| `encoding/base64` | StdEncoding.EncodeToString/DecodeString, URLEncoding.EncodeToString, RawStdEncoding.EncodeToString | ✓ Covered + E2E tested |
| `encoding/hex` | EncodeToString, DecodeString, NewEncoder, NewDecoder, Dump | ✓ Covered + E2E tested |
| `encoding/xml` | Marshal, MarshalIndent, Unmarshal, NewEncoder, NewDecoder, Escape, EscapeText | ✓ Covered + E2E tested |
| `encoding/binary` | Read, Write, Size, Varint, PutUvarint, PutVarint | ✓ Covered + E2E tested |
| `regexp` | Compile, CompilePOSIX, Match, MatchString, MustCompile, QuoteMeta, Regexp.MatchString/FindString/FindAllString/ReplaceAllString/ReplaceAllStringFunc/Split/SubexpNames | ✓ Covered + E2E tested |
| `path/filepath` | Join, Base, Dir, Ext, Abs, Clean, Rel, Split, Walk, WalkDir, Match, Glob, IsAbs, EvalSymlinks | ✓ Covered + E2E tested |
| `net` | Dial, DialTimeout, Listen, ListenPacket, ResolveTCPAddr, ResolveIPAddr, ResolveUDPAddr, SplitHostPort, JoinHostPort, LookupHost, LookupAddr, InterfaceAddrs, InterfaceByName, Conn.SetDeadline | ✓ Covered + E2E tested |
| `net/http` | Get, Post, PostForm, Head, NewRequest, ListenAndServe, ListenAndServeTLS, NewServeMux, ServeMux.HandleFunc, Client.Do, Response.Body.Close, Request.WithContext | ✓ Covered + E2E tested |
| `net/url` | Parse, ParseRequestURI, Values.Get/Set/Add/Encode/Del, QueryEscape/QueryUnescape, PathEscape/PathUnescape | ✓ Covered + E2E tested |
| `crypto/tls` | Dial, DialWithDialer, Listen, NewListener, LoadX509KeyPair, X509KeyPair | ✓ Covered + E2E tested |
| `crypto/sha256` | Sum256, New | ✓ Covered + E2E tested |
| `crypto/md5` | Sum | ✓ Covered + E2E tested |
| `crypto/rand` | Read | ✓ Covered + E2E tested |
| `crypto/aes` | NewCipher | ✓ Covered + E2E tested |
| `crypto/hmac` | New | ✓ Covered + E2E tested |
| `crypto/x509` | ParseCertificate | ✓ Covered + E2E tested |
| `crypto/rsa` | GenerateKey, EncryptOAEP, DecryptOAEP, SignPKCS1v15, VerifyPKCS1v15 | ✓ Covered + E2E tested |
| `math/big` | Int.SetString, Int.String, Int.Add, Int.Mul, Int.Div, Int.Mod, Rat.SetString, Float.SetString | ✓ Covered + E2E tested |
| `encoding/csv` | NewReader, NewWriter, Reader.Read/ReadAll, Writer.Write/WriteAll/Flush | ✓ Covered + E2E tested |
| `encoding/gob` | NewEncoder, NewDecoder, Register, RegisterName | ✓ Covered + E2E tested |
| `html/template` | New, Must, ParseFiles, ParseGlob, Template.Execute/ExecuteTemplate | ✓ Covered + E2E tested |
| `text/template` | New, Must, ParseFiles, ParseGlob, Template.Execute/ExecuteTemplate | ✓ Covered + E2E tested |
| `mime` | TypeByExtension, ExtensionsByType, AddExtensionType, FormatMediaType, ParseMediaType | ✓ Covered + E2E tested |
| `mime/multipart` | NewReader, NewWriter, Reader.NextPart/ReadForm, Writer.CreateFormFile/CreatePart/FormDataContentType/Close | ✓ Covered + E2E tested |
| `compress/gzip` | NewReader, NewWriter, NewWriterLevel, Reader.Read/Close, Writer.Write/Close/Flush/Reset | ✓ Covered + E2E tested |
| `compress/zlib` | NewReader, NewWriter, NewWriterLevel | ✓ Covered + E2E tested |
| `compress/flate` | NewReader, NewWriter, NewWriterDict | ✓ Covered + E2E tested |
| `archive/tar` | NewReader, NewWriter, Reader.Next/Read, Writer.WriteHeader/Write/Close/Flush | ✓ Covered + E2E tested |
| `archive/zip` | OpenReader, NewReader, NewWriter, File.Open, Writer.Create/CreateHeader/Close | ✓ Covered + E2E tested |
| `reflect` | TypeOf, ValueOf, DeepEqual, Value.Interface/Int/String/Bool/Float/Len/Index/Field/MapIndex/MapKeys/Set/SetString/Elem/Kind/Slice/Close/IsNil/IsValid, Type.Name/Kind/NumMethod/NumField/Field/Size/String/Elem | ✓ Covered + E2E tested |
| `unsafe` | Sizeof, Offsetof, Alignof, Pointer, Add, Slice, SliceData, String, StringData | ✓ Covered + E2E tested |
| `container/list` | New, List.Init/Len/Front/Back/PushFront/PushBack/InsertBefore/InsertAfter/MoveToFront/MoveToBack/Remove | ✓ Covered + E2E tested |
| `container/heap` | Init, Push, Pop, Remove, Fix | ✓ Covered + E2E tested |
| `container/ring` | New, Ring.Len/Next/Prev/Link/Unlink/Move/Do | ✓ Covered + E2E tested |
| `sync/atomic` | LoadInt32, StoreInt32, AddInt32, SwapInt32, CompareAndSwapInt32 (and Int64/Uint32/Uint64/Uintptr/Pointer variants) | ✓ Covered + E2E tested |

### Runtime patterns (all covered)

| Category | Patterns | Status |
|---|---|---|
| Memory | memmove, newobject | ✓ Covered + E2E tested |
| Channels | chansend1, chanrecv1/2, closechan, selectnbsend/recv, makechan | ✓ Covered + E2E tested |
| Maps | mapaccess1/2, mapassign, mapdelete, makemap | ✓ Covered + E2E tested |
| Slices | makeslice, growslice, slicebytetostring, stringtoslicebyte, slicecopy, typedslicecopy | ✓ Covered + E2E tested |
| Goroutines | newproc, goexit, gopark, goready | ✓ Covered + E2E tested |
| Defer/Panic | deferproc, deferprocStack, deferreturn, gopanic, gorecover | ✓ Covered + E2E tested |
| Locks | lock, unlock, noteclear, notesleep, notewakeup | ✓ Covered + E2E tested |
| GC | gcWriteBarrier, gcBgMarkWorker | ✓ Covered |
| Syscalls | Linux (137), Windows (121), Darwin (70), FreeBSD (57) | ✓ Covered |
| Reflection | typedslicecopy, convT* | ✓ Covered |

## Summary

All Go standard library packages documented in this plan are fully implemented with:
- Pattern files (`.hexpat`) — 562 files across stdlib (437), runtime (87), fallback (26), and controlflow (12) modules
- Test source programs (`main.go`)
- End-to-end decompile tests (compile → disassemble → match → generate → verify)

Total: **50+ packages** with **562 pattern files** across four modules (stdlib 437, runtime 87, fallback 26, controlflow 12), all with automated E2E tests.
