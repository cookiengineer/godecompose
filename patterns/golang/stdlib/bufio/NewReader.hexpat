arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bufio_NewReader {
    name: "bufio.NewReader"; library: "go-stdlib"; instr match { CALL bufio_NewReader }
    gen { bufio.NewReader($rd); } bind { rd as "rd"; }
}
