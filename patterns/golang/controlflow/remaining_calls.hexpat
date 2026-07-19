arch x86_64;
platform linux, darwin, windows, freebsd;

pattern go_os_Remove {
    name: "go: os.Remove";
    library: "go-stdlib";
    description: "Remove file via os.Remove";
    instr match { CALL os_Remove }
    gen { os.Remove(name) }
}

pattern go_os_Lstat {
    name: "go: os.Lstat";
    library: "go-stdlib";
    description: "Stat symlink via os.Lstat";
    instr match { CALL os_Lstat }
    gen { os.Lstat(name) }
}

pattern go_net_Listen {
    name: "go: net.Listen";
    library: "go-stdlib";
    description: "Listen on network address via net.Listen";
    instr match { CALL net_Listen }
    gen { net.Listen(network, addr) }
}

pattern go_net_http_Request_BasicAuth {
    name: "go: net/http.Request.BasicAuth";
    library: "go-stdlib";
    description: "Extract basic auth via Request.BasicAuth";
    instr match { CALL net_http_Request_BasicAuth }
    gen { r.BasicAuth() }
}

pattern go_strings_TrimRight {
    name: "go: strings.TrimRight";
    library: "go-stdlib";
    description: "Trim trailing characters via strings.TrimRight";
    instr match { CALL strings_TrimRight }
    gen { strings.TrimRight(s, cut) }
}

pattern go_runtime_typedslicecopy {
    name: "go: runtime.typedslicecopy";
    library: "go-runtime";
    description: "Copy typed slice via runtime.typedslicecopy";
    instr match { CALL runtime_typedslicecopy }
    gen { copy(dst, src) }
}

pattern go_runtime_intstring {
    name: "go: runtime.intstring";
    library: "go-runtime";
    description: "Convert int to string via runtime.intstring";
    instr match { CALL runtime_intstring }
    gen { string(rune(n)) }
}

pattern go_runtime_moveSliceNoCapNoScan {
    name: "go: runtime.moveSliceNoCapNoScan";
    library: "go-runtime";
    description: "Move slice data";
    instr match { CALL runtime_moveSliceNoCapNoScan }
    gen { /* moveSliceNoCapNoScan */ }
}

pattern go_time_NewTimer {
    name: "go: time.NewTimer";
    library: "go-stdlib";
    description: "Create new timer via time.NewTimer";
    instr match { CALL time_NewTimer }
    gen { time.NewTimer(d) }
}

pattern go_syscall_symlinkat {
    name: "go: syscall.symlinkat";
    library: "go-runtime";
    description: "Create symlink via syscall.symlinkat";
    instr match { CALL syscall_symlinkat }
    gen { syscall.Symlinkat(new, dir, old) }
}

pattern go_syscall_rawSyscallNoError {
    name: "go: syscall.rawSyscallNoError.abi0";
    library: "go-runtime";
    description: "Raw syscall via syscall.rawSyscallNoError";
    instr match { CALL syscall_rawSyscallNoError_abi0 }
    gen { /* raw syscall */ }
}

pattern go_os_Process_signal {
    name: "go: os.Process.signal";
    library: "go-stdlib";
    description: "Send signal to process";
    instr match { CALL os_Process_signal }
    gen { proc.Signal(sig) }
}

pattern go_os_File_Sync {
    name: "go: os.File.Sync";
    library: "go-stdlib";
    description: "Sync file to disk";
    instr match { CALL os_File_Sync }
    gen { f.Sync() }
}

pattern go_net_Resolver_LookupTXT {
    name: "go: net.Resolver.LookupTXT";
    library: "go-stdlib";
    description: "DNS TXT record lookup";
    instr match { CALL net_Resolver_LookupTXT }
    gen { resolver.LookupTXT(ctx, name) }
}

pattern go_net_url_escape {
    name: "go: net/url.escape";
    library: "go-stdlib";
    description: "URL escape internal function";
    instr match { CALL net_url_escape }
    gen { url.PathEscape(s) }
}

pattern go_internal_filepathlite_Dir {
    name: "go: filepath.Dir";
    library: "go-stdlib";
    description: "Get directory part of path";
    instr match { CALL internal_filepathlite_Dir }
    gen { filepath.Dir(path) }
}

pattern go_log_output {
    name: "go: log.Logger.output";
    library: "go-stdlib";
    description: "Log output via log.Logger.output";
    instr match { CALL log_Logger_output }
    gen { log.Output(calldepth, msg) }
}

pattern go_semver_parse {
    name: "go: semver.Parse";
    library: "go-stdlib";
    description: "Parse semantic version";
    instr match { CALL semver_parse }
    gen { semver.Parse(v) }
}
