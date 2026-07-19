arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bufio_Writer_WriteByte {
    name: "bufio.Writer.WriteByte"; library: "go-stdlib"; instr match { CALL bufio___Writer__WriteByte }
    gen { $w.WriteByte($c); } bind { w as "w"; c as "c"; }
}
