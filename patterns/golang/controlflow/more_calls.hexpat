arch x86_64;
platform linux, darwin, windows, freebsd;

pattern go_gc_writebarrier_check {
    name: "go: GC write barrier check";
    library: "go-runtime";
    description: "Write barrier enabled check: cmp [runtime.writeBarrier], 0; jz";
    instr match {
        CMPQ ptr, 0x0
        JEQ @skip
    }
    gen { if gcWriteBarrier.enabled { }
    bind { ptr as "wb" }
}

pattern go_flag_Var {
    name: "go: flag.FlagSet.Var";
    library: "go-stdlib";
    description: "Register a flag value via flag.FlagSet.Var";
    instr match { CALL flag_FlagSet_Var }
    gen { flag.Var(value, name, usage) }
}

pattern go_net_http_request_formvalue {
    name: "go: net/http.Request.FormValue";
    library: "go-stdlib";
    description: "Get form value from HTTP request";
    instr match { CALL net_http_Request_FormValue }
    gen { r.FormValue(key) }
}

pattern go_net_http_redirect {
    name: "go: net/http.Redirect";
    library: "go-stdlib";
    description: "HTTP redirect response";
    instr match { CALL net_http_Redirect }
    gen { http.Redirect(w, r, url, code) }
}

pattern go_fmt_Appendf {
    name: "go: fmt.Appendf";
    library: "go-stdlib";
    description: "Formatted append via fmt.Appendf";
    instr match { CALL fmt_Appendf }
    gen { fmt.Appendf(buf, format, args...) }
}

pattern go_fmt_Sprintf {
    name: "go: fmt.Sprintf";
    library: "go-stdlib";
    description: "Formatted string via fmt.Sprintf";
    instr match { CALL fmt_Sprintf }
    gen { fmt.Sprintf(format, args...) }
}

pattern go_fmt_Errorf {
    name: "go: fmt.Errorf";
    library: "go-stdlib";
    description: "Formatted error via fmt.Errorf";
    instr match { CALL fmt_Errorf }
    gen { fmt.Errorf(format, args...) }
}

pattern go_fmt_Fprintln {
    name: "go: fmt.Fprintln";
    library: "go-stdlib";
    description: "Formatted println to writer via fmt.Fprintln";
    instr match { CALL fmt_Fprintln }
    gen { fmt.Fprintln(w, args...) }
}

pattern go_time_Time_Add {
    name: "go: time.Time.Add";
    library: "go-stdlib";
    description: "Add duration to time via time.Time.Add";
    instr match { CALL time_Time_Add }
    gen { t.Add(d) }
}

pattern go_math_rand_uint64n {
    name: "go: math/rand.rand.uint64n";
    library: "go-stdlib";
    description: "Random uint64 in range via math/rand.Rand.uint64n";
    instr match { CALL math_rand_v2_Rand_uint64n }
    gen { rand.Int64N(n) }
}

pattern go_net_http_Server_Serve {
    name: "go: net/http.Server.Serve";
    library: "go-stdlib";
    description: "HTTP server serve via net/http.Server.Serve";
    instr match { CALL net_http_Server_Serve }
    gen { server.Serve(listener) }
}

pattern go_sconf_Parse {
    name: "go: sconf.Parse";
    library: "go-stdlib";
    description: "Parse config via sconf.Parse";
    instr match { CALL sconf_Parse }
    gen { sconf.Parse(r, cfg) }
}

pattern go_sconf_Write {
    name: "go: sconf.Write";
    library: "go-stdlib";
    description: "Write config via sconf.Write";
    instr match { CALL sconf_Write }
    gen { sconf.Write(w, cfg) }
}

pattern go_sconf_Describe {
    name: "go: sconf.Describe";
    library: "go-stdlib";
    description: "Describe config via sconf.Describe";
    instr match { CALL sconf_Describe }
    gen { sconf.Describe(w, cfg) }
}

pattern go_promauto_NewCounter {
    name: "go: promauto.NewCounter";
    library: "go-stdlib";
    description: "Create prometheus counter via promauto.Factory.NewCounter";
    instr match { CALL Factory_NewCounter }
    gen { promauto.NewCounter(opts) }
}

pattern go_promauto_NewGauge {
    name: "go: promauto.NewGauge";
    library: "go-stdlib";
    description: "Create prometheus gauge via promauto.Factory.NewGauge";
    instr match { CALL Factory_NewGauge }
    gen { promauto.NewGauge(opts) }
}

pattern go_promauto_NewGaugeVec {
    name: "go: promauto.NewGaugeVec";
    library: "go-stdlib";
    description: "Create prometheus gauge vec via promauto.Factory.NewGaugeVec";
    instr match { CALL Factory_NewGaugeVec }
    gen { promauto.NewGaugeVec(opts, labels) }
}

pattern go_prom_GaugeVec_WithLabelValues {
    name: "go: prometheus.GaugeVec.WithLabelValues";
    library: "go-stdlib";
    description: "Get gauge with labels via prometheus.GaugeVec.WithLabelValues";
    instr match { CALL prometheus_GaugeVec_WithLabelValues }
    gen { gauge.WithLabelValues(labels...) }
}

pattern go_os_Process_wait {
    name: "go: os.Process.wait";
    library: "go-stdlib";
    description: "Wait for process via os.Process.wait";
    instr match { CALL os_Process_wait }
    gen { proc.Wait() }
}

pattern go_os_exec_Cmd_Start {
    name: "go: os/exec.Cmd.Start";
    library: "go-stdlib";
    description: "Start command via os/exec.Cmd.Start";
    instr match { CALL os_exec_Cmd_Start }
    gen { cmd.Start() }
}

pattern go_html_template_Template_Parse {
    name: "go: html/template.Template.Parse";
    library: "go-stdlib";
    description: "Parse HTML template via html/template.Template.Parse";
    instr match { CALL html_template_Template_Parse }
    gen { templ.Parse(text) }
}

pattern go_text_template_Template_Funcs {
    name: "go: text/template.Template.Funcs";
    library: "go-stdlib";
    description: "Register template functions via text/template.Template.Funcs";
    instr match { CALL text_template_Template_Funcs }
    gen { templ.Funcs(funcMap) }
}

pattern go_flag_FlagSet_PrintDefaults {
    name: "go: flag.FlagSet.PrintDefaults";
    library: "go-stdlib";
    description: "Print flag defaults via flag.FlagSet.PrintDefaults";
    instr match { CALL flag_FlagSet_PrintDefaults }
    gen { flag.PrintDefaults() }
}

pattern go_flag_Parse {
    name: "go: flag.Parse";
    library: "go-stdlib";
    description: "Parse flags via flag.Parse";
    instr match { CALL flag_Parse }
    gen { flag.Parse() }
}

pattern go_flag_FlagSet_Parse {
    name: "go: flag.FlagSet.Parse";
    library: "go-stdlib";
    description: "Parse flags via flag.FlagSet.Parse";
    instr match { CALL flag_FlagSet_Parse }
    gen { fs.Parse(arguments) }
}
