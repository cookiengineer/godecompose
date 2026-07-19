arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bufio_Writer_WriteString {
    name: "bufio.Writer.WriteString"; library: "go-stdlib"; instr match { CALL bufio___Writer__WriteString }
    gen { $w.WriteString($s); } bind { w as "w"; s as "s"; }
}
