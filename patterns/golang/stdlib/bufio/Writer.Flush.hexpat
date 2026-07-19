arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bufio_Writer_Flush {
    name: "bufio.Writer.Flush"; library: "go-stdlib"; instr match { CALL bufio___Writer__Flush }
    gen { $w.Flush(); } bind { w as "w"; }
}
