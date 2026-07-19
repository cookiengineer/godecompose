arch x86_64;
platform linux, darwin, windows, freebsd;

pattern url_Values_Add {
    name: "url.Values.Add"; library: "go-stdlib";
    instr match { CALL net_url___Values__Add }
    gen { $v.Add($key, $value); } bind { v as "v"; key as "key"; value as "value"; }
}
