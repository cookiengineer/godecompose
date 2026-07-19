arch x86_64;
platform linux, darwin, windows, freebsd;

pattern url_Values_Get {
    name: "url.Values.Get"; library: "go-stdlib";
    instr match { CALL net_url___Values__Get }
    gen { $v.Get($key); } bind { v as "v"; key as "key"; }
}
