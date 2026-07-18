arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bytes_Contains {
    name: "bytes.Contains"; library: "go-stdlib"; instr match { CALL bytes_Contains }
    gen { bytes.Contains($b, $subslice); } bind { b as "b"; subslice as "subslice"; }
}
pattern bytes_Equal {
    name: "bytes.Equal"; library: "go-stdlib"; instr match { CALL bytes_Equal }
    gen { bytes.Equal($a, $b); } bind { a as "a"; b as "b"; }
}
pattern bytes_Compare {
    name: "bytes.Compare"; library: "go-stdlib"; instr match { CALL bytes_Compare }
    gen { bytes.Compare($a, $b); } bind { a as "a"; b as "b"; }
}
pattern bytes_Index {
    name: "bytes.Index"; library: "go-stdlib"; instr match { CALL bytes_Index }
    gen { bytes.Index($s, $sep); } bind { s as "s"; sep as "sep"; }
}
pattern bytes_LastIndex {
    name: "bytes.LastIndex"; library: "go-stdlib"; instr match { CALL bytes_LastIndex }
    gen { bytes.LastIndex($s, $sep); } bind { s as "s"; sep as "sep"; }
}
pattern bytes_Join {
    name: "bytes.Join"; library: "go-stdlib"; instr match { CALL bytes_Join }
    gen { bytes.Join($s, $sep); } bind { s as "s"; sep as "sep"; }
}
pattern bytes_Split {
    name: "bytes.Split"; library: "go-stdlib"; instr match { CALL bytes_Split }
    gen { bytes.Split($s, $sep); } bind { s as "s"; sep as "sep"; }
}
pattern bytes_HasPrefix {
    name: "bytes.HasPrefix"; library: "go-stdlib"; instr match { CALL bytes_HasPrefix }
    gen { bytes.HasPrefix($s, $prefix); } bind { s as "s"; prefix as "prefix"; }
}
pattern bytes_NewReader {
    name: "bytes.NewReader"; library: "go-stdlib"; instr match { CALL bytes_NewReader }
    gen { bytes.NewReader($b); } bind { b as "b"; }
}
pattern bytes_NewBuffer {
    name: "bytes.NewBuffer"; library: "go-stdlib"; instr match { CALL bytes_NewBuffer }
    gen { bytes.NewBuffer($buf); } bind { buf as "buf"; }
}
pattern bytes_Buffer_WriteString {
    name: "bytes.Buffer.WriteString"; library: "go-stdlib"; instr match { CALL bytes___Buffer__WriteString }
    gen { $b.WriteString($s); } bind { b as "b"; s as "s"; }
}
pattern bytes_Buffer_WriteByte {
    name: "bytes.Buffer.WriteByte"; library: "go-stdlib"; instr match { CALL bytes___Buffer__WriteByte }
    gen { $b.WriteByte($c); } bind { b as "b"; c as "c"; }
}
pattern bytes_Buffer_Bytes {
    name: "bytes.Buffer.Bytes"; library: "go-stdlib"; instr match { CALL bytes___Buffer__Bytes }
    gen { $b.Bytes(); } bind { b as "b"; }
}
pattern bytes_Buffer_String {
    name: "bytes.Buffer.String"; library: "go-stdlib"; instr match { CALL bytes___Buffer__String }
    gen { $b.String(); } bind { b as "b"; }
}
pattern bytes_Buffer_Reset {
    name: "bytes.Buffer.Reset"; library: "go-stdlib"; instr match { CALL bytes___Buffer__Reset }
    gen { $b.Reset(); } bind { b as "b"; }
}

pattern bufio_NewReader {
    name: "bufio.NewReader"; library: "go-stdlib"; instr match { CALL bufio_NewReader }
    gen { bufio.NewReader($rd); } bind { rd as "rd"; }
}
pattern bufio_NewWriter {
    name: "bufio.NewWriter"; library: "go-stdlib"; instr match { CALL bufio_NewWriter }
    gen { bufio.NewWriter($wr); } bind { wr as "wr"; }
}
pattern bufio_NewScanner {
    name: "bufio.NewScanner"; library: "go-stdlib"; instr match { CALL bufio_NewScanner }
    gen { bufio.NewScanner($r); } bind { r as "r"; }
}
pattern bufio_Reader_ReadString {
    name: "bufio.Reader.ReadString"; library: "go-stdlib"; instr match { CALL bufio___Reader__ReadString }
    gen { $r.ReadString($delim); } bind { r as "r"; delim as "delim"; }
}
pattern bufio_Reader_ReadByte {
    name: "bufio.Reader.ReadByte"; library: "go-stdlib"; instr match { CALL bufio___Reader__ReadByte }
    gen { $r.ReadByte(); } bind { r as "r"; }
}
pattern bufio_Reader_ReadBytes {
    name: "bufio.Reader.ReadBytes"; library: "go-stdlib"; instr match { CALL bufio___Reader__ReadBytes }
    gen { $r.ReadBytes($delim); } bind { r as "r"; delim as "delim"; }
}
pattern bufio_Reader_ReadLine {
    name: "bufio.Reader.ReadLine"; library: "go-stdlib"; instr match { CALL bufio___Reader__ReadLine }
    gen { $r.ReadLine(); } bind { r as "r"; }
}
pattern bufio_Reader_Peek {
    name: "bufio.Reader.Peek"; library: "go-stdlib"; instr match { CALL bufio___Reader__Peek }
    gen { $r.Peek($n); } bind { r as "r"; n as "n"; }
}
pattern bufio_Writer_WriteString {
    name: "bufio.Writer.WriteString"; library: "go-stdlib"; instr match { CALL bufio___Writer__WriteString }
    gen { $w.WriteString($s); } bind { w as "w"; s as "s"; }
}
pattern bufio_Writer_WriteByte {
    name: "bufio.Writer.WriteByte"; library: "go-stdlib"; instr match { CALL bufio___Writer__WriteByte }
    gen { $w.WriteByte($c); } bind { w as "w"; c as "c"; }
}
pattern bufio_Writer_Flush {
    name: "bufio.Writer.Flush"; library: "go-stdlib"; instr match { CALL bufio___Writer__Flush }
    gen { $w.Flush(); } bind { w as "w"; }
}
pattern bufio_Writer_Reset {
    name: "bufio.Writer.Reset"; library: "go-stdlib"; instr match { CALL bufio___Writer__Reset }
    gen { $w.Reset($wr); } bind { w as "w"; wr as "wr"; }
}
pattern bufio_Scanner_Scan {
    name: "bufio.Scanner.Scan"; library: "go-stdlib"; instr match { CALL bufio___Scanner__Scan }
    gen { $s.Scan(); } bind { s as "s"; }
}
pattern bufio_Scanner_Text {
    name: "bufio.Scanner.Text"; library: "go-stdlib"; instr match { CALL bufio___Scanner__Text }
    gen { $s.Text(); } bind { s as "s"; }
}
pattern bufio_Scanner_Bytes {
    name: "bufio.Scanner.Bytes"; library: "go-stdlib"; instr match { CALL bufio___Scanner__Bytes }
    gen { $s.Bytes(); } bind { s as "s"; }
}
pattern bufio_Scanner_Err {
    name: "bufio.Scanner.Err"; library: "go-stdlib"; instr match { CALL bufio___Scanner__Err }
    gen { $s.Err(); } bind { s as "s"; }
}
