# Go Standard Library Pattern Coverage Plan

## Pattern Structure

Each stdlib pattern follows the high-level CALL-matching approach:
- A single `CALL` instruction that matches by symbol name in the disassembled GoSyntax
- An optional `gen` block describing the recovered Go source
- An optional `bind` block for variable renaming

Symbol names in Go compiled binaries follow Go syntax:
- Package functions: `package.Function` → pattern: `CALL package_Function`
- Pointer receiver methods: `package.(*Type).Method` → pattern: `CALL package___Type__Method`
- Value receiver methods: `package.Type.Method` → pattern: `CALL package_Type_Method`

The fuzzy matcher normalizes underscores to periods, lowercases, and does prefix matching, so patterns like `sync_Mutex_Lock` match against `CALL sync.(*Mutex).lockSlow(SB)`.

## Package Coverage Status

### ✓ Already covered (existing patterns + newly implemented)
- `fmfmt`: Fprintln, Fprintf, Sprintf, Errorf, Fprint
- `errors`: New, Is, As, Unwrap **NEW**
- `log`: Println, Printf, Fatal, Fatalf, Panic **NEW**
- `context`: WithCancel, WithTimeout, WithDeadline, Background, TODO **NEW**
- `flag`: Parse, NewFlagSet, FlagSet.String/Int/Bool **NEW**
- `sync`: Mutex.Lock/Unlock, WaitGroup.Add/Done/Wait, Once.Do, Pool.Get/Put
- `os`: File.Write/Read/Close, Open, Create, Remove, Stat, Mkdir, Getenv
- `io`: ReadAll, Copy, WriteString
- `time`: Now, Sleep, After, Since, NewTicker
- `strconv`: Itoa, Atoi, FormatInt, FormatFloat, ParseInt, ParseFloat, Quote, Unquote **NEW**
- `strings`: Contains, HasPrefix, HasSuffix, Join, Split, ReplaceAll, ToLower, ToUpper, TrimSpace, TrimPrefix, TrimSuffix, Index, Count, Repeat, Fields, NewReader, NewReplacer, Builder.WriteString/WriteByte/String/Reset **NEW**
- `bytes`: Contains, Equal, Compare, Index, Join, Split, HasPrefix, NewReader, NewBuffer, Buffer.WriteString/WriteByte/Bytes/String/Reset **NEW**
- `bufio`: NewReader, NewWriter, NewScanner, Reader.ReadString/ReadByte/ReadBytes/ReadLine/Peek, Writer.WriteString/WriteByte/Flush/Reset, Scanner.Scan/Text/Bytes/Er **rNEW**
- `math`: Abs, Max, Min, Sqrt, Pow, Sin, Cos, Tan, Floor, Ceil, Round, Log, Exp, Mod **NEW**
- `math/rand`: Intn, Float64, Int, New **NEW**
- `sort`: Ints, Strings, Float64s, Slice, Search, Stable, IsSorted **NEW**
- `net`: Dial, DialTimeout, Listen, ResolveTCPAddr, SplitHostPort, JoinHostPort **NEW**
- `net/http`: Get, Post, NewRequest, ListenAndServe, NewServeMux, ServeMux.HandleFunc, Client.Do, Response.Body.Close **NEW**
- `net/url`: Parse, Values.Get/Set/Add/Encode **NEW**
- `encoding/json`: Marshal, Unmarshal, NewEncoder, NewDecoder, Encode, Decode
- `encoding/base64`: StdEncoding.EncodeToString/DecodeString, URLEncoding.EncodeToString **NEW**
- `encoding/hex`: EncodeToString, DecodeString, NewEncoder, NewDecoder **NEW**
- `encoding/xml`: Marshal, Unmarshal, NewEncoder, NewDecoder **NEW**
- `encoding/binary`: Read, Write **NEW**
- `regexp`: Compile, MustCompile, MatchString, Regexp.MatchString/FindString/FindAllString/ReplaceAllString/Split **NEW**
- `path/filepath`: Join, Base, Dir, Ext, Abs, Walk, Glob **NEW**
- `crypto/tls`: Dial, Listen, LoadX509KeyPair **NEW**
- `crypto/sha256`: Sum256, New
- `crypto/md5`: Sum
- `crypto/rand`: Read
- `crypto/tls`: Dial
- `crypto/aes`: NewCipher
- `crypto/hmac`: New
- `crypto/x509`: ParseCertificate
- All runtime: memory, goroutines, channels, maps, slices, interfaces, defer, locks, GC, syscalls, reflection

### To implement

| Package | Functions / Methods | Priority |
|---|---|---|
| `errors` | New, Is, As, Unwrap, Join | HIGH |
| `log` | Println, Printf, Fatal, Fatalf, Panic, Panicf | HIGH |
| `context` | WithCancel, WithTimeout, WithDeadline, WithValue, Background, TODO | HIGH |
| `flag` | Parse, NewFlagSet, FlagSet.String/Int/Bool, FlagSet.Parse | HIGH |
| `strconv` | Itoa, Atoi, FormatInt, FormatFloat, ParseInt, ParseFloat, Quote, Unquote | HIGH |
| `strings` | Contains, HasPrefix, HasSuffix, Join, Split, ReplaceAll, ToLower, ToUpper, TrimSpace, TrimPrefix, TrimSuffix, Index, LastIndex, Count, Repeat, Fields, NewReader, NewReplacer, Builder.WriteString/WriteRune/WriteByte/Reset/String | HIGH |
| `bytes` | Contains, Equal, Compare, Index, Join, Split, HasPrefix, Buffer.WriteString/WriteByte/Bytes/Reset/String, NewReader, NewBuffer | HIGH |
| `bufio` | NewReader, NewWriter, NewScanner, Reader.ReadString/ReadByte/ReadBytes/ReadLine/Peek, Writer.WriteString/WriteByte/Flush/Reset, Scanner.Scan/Text/Bytes/Err | HIGH |
| `math` | Abs, Max, Min, Sqrt, Pow, Sin, Cos, Tan, Floor, Ceil, Round, Log, Log2, Log10, Exp, Mod, Remainder, Hypot | MED |
| `math/rand` | Intn, Float64, Int, Int31, Int63, Perm, Shuffle, Seed, New, NewSource, Read | MED |
| `math/big` | Int.SetString, Int.String, Int.Add, Int.Mul, Int.Div, Int.Mod, Rat.SetString, Float.SetString | LOW |
| `sort` | Ints, Strings, Float64s, Slice, Search, SearchInts, SearchStrings, SearchFloat64s, Stable, Reverse, IsSorted | MED |
| `encoding/base64` | StdEncoding.EncodeToString/DecodeString, URLEncoding, RawStdEncoding | MED |
| `encoding/hex` | EncodeToString, DecodeString, NewEncoder, NewDecoder, Dump | MED |
| `encoding/xml` | Marshal, MarshalIndent, Unmarshal, NewEncoder, NewDecoder, Escape, EscapeText | MED |
| `encoding/binary` | Read, Write, Size, Varint, PutUvarint, PutVarint, LittleEndian.Uint32, BigEndian.Uint16 | MED |
| `encoding/csv` | NewReader, NewWriter, Reader.Read/ReadAll, Writer.Write/WriteAll/Flush | LOW |
| `encoding/gob` | NewEncoder, NewDecoder, Register, RegisterName | LOW |
| `regexp` | Compile, CompilePOSIX, Match, MatchString, MustCompile, QuoteMeta, Regexp.MatchString/FindString/FindAllString/ReplaceAllString/ReplaceAllStringFunc/Split/SubexpNames | MED |
| `path/filepath` | Join, Base, Dir, Ext, Abs, Clean, Rel, Split, Walk, WalkDir, Match, Glob, IsAbs, EvalSymlinks | MED |
| `net/http` | Get, Post, PostForm, Head, NewRequest, NewServeMux, ListenAndServe, ListenAndServeTLS, Server.ListenAndServe, ServeMux.Handle/HandleFunc, Client.Do/Get/Post, Response.Body.Close, Request.WithContext | HIGH |
| `net/url` | Parse, ParseRequestURI, Values.Get/Set/Add/Encode/Del, QueryEscape/QueryUnescape, PathEscape/PathUnescape | MED |
| `net` | Dial, DialTimeout, Listen, ListenPacket, Dialer.Dial/DialContext, Conn.SetDeadline, ResolveIPAddr, ResolveTCPAddr, ResolveUDPAddr, SplitHostPort, JoinHostPort, LookupHost, LookupAddr, InterfaceAddrs, InterfaceByName | MED |
| `crypto/tls` | Dial, DialWithDialer, Listen, NewListener, LoadX509KeyPair, X509KeyPair, Server.ListenAndServeTLS, Client, Config | MED |
| `crypto/rsa` | GenerateKey, EncryptOAEP, DecryptOAEP, SignPKCS1v15, VerifyPKCS1v15, EncryptPKCS1v15, DecryptPKCS1v15 | LOW |
| `html/template` | New, Must, ParseFiles, ParseGlob, Template.Execute/ExecuteTemplate | LOW |
| `text/template` | New, Must, ParseFiles, ParseGlob, Template.Execute/ExecuteTemplate | LOW |
| `mime` | TypeByExtension, ExtensionsByType, AddExtensionType, FormatMediaType, ParseMediaType | LOW |
| `mime/multipart` | NewReader, NewWriter, Reader.NextPart/ReadForm, Writer.CreateFormFile/CreatePart/FormDataContentType/Close | LOW |
| `compress/gzip` | NewReader, NewWriter, NewWriterLevel, Reader.Read/Close, Writer.Write/Close/Flush/Reset | LOW |
| `compress/zlib` | NewReader, NewWriter, NewWriterLevel | LOW |
| `compress/flate` | NewReader, NewWriter, NewWriterDict | LOW |
| `archive/tar` | NewReader, NewWriter, Reader.Next/Read, Writer.WriteHeader/Write/Close/Flush | LOW |
| `archive/zip` | OpenReader, NewReader, NewWriter, File.Open, Writer.Create/CreateHeader/Close | LOW |
| `reflect` | TypeOf, ValueOf, DeepEqual, Value.Interface/Int/String/Bool/Float/Len/Index/Field/MapIndex/MapKeys/Set/SetString/Elem/Kind/Slice/Slice3/Close/IsNil/IsValid, Type.Name/Kind/NumMethod/NumField/Field/Size/Align/String/Elem/Comparable | MED |
| `unsafe` | Sizeof, Offsetof, Alignof, Pointer, Add, Slice, SliceData, String, StringData | LOW |
| `container/list` | New, List.Init/Len/Front/Back/PushFront/PushBack/InsertBefore/InsertAfter/MoveToFront/MoveToBack/Remove | LOW |
| `container/heap` | Init, Push, Pop, Remove, Fix | LOW |
| `container/ring` | New, Ring.Len/Next/Prev/Link/Unlink/Move/Do | LOW |
| `sync/atomic` | LoadInt32, StoreInt32, AddInt32, SwapInt32, CompareAndSwapInt32, (and int64/uint32/uint64/uintptr/pointer variants), Bool, Value | LOW |
