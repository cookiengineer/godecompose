arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bufio_Scanner_Text {
    name: "bufio.Scanner.Text"; library: "go-stdlib"; instr match { CALL bufio___Scanner__Text }
    gen { $s.Text(); } bind { s as "s"; }
}
