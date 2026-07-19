arch x86_64;
platform linux, darwin, windows, freebsd;

pattern http_ServeMux_HandleFunc {
    name: "http.ServeMux.HandleFunc"; library: "go-stdlib";
    instr match { CALL net_http___ServeMux__HandleFunc }
    gen { $mux.HandleFunc($pattern, $handler); } bind { mux as "mux"; pattern as "pattern"; handler as "handler"; }
}
