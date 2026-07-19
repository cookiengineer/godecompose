arch x86_64;
platform linux, darwin, windows, freebsd;

pattern go_runtime_newobject {
    name: "go: new(T) (runtime.newobject)";
    library: "go-datatype";
    description: "Allocate new zeroed object via runtime.newobject";
    instr match { CALL runtime_newobject }
    gen { x := new(T) }
}

pattern go_runtime_makeslice {
    name: "go: make slice (runtime.makeslice)";
    library: "go-datatype";
    description: "Slice allocation via runtime.makeslice";
    instr match { CALL runtime_makeslice }
    gen { make([]T, n) }
}

pattern go_runtime_makemap {
    name: "go: make map (runtime.makemap)";
    library: "go-datatype";
    description: "Map allocation via runtime.makemap";
    instr match { CALL runtime_makemap }
    gen { make(map[K]V) }
}

pattern go_runtime_convTstring {
    name: "go: string -> interface (runtime.convTstring)";
    library: "go-datatype";
    description: "Convert string to interface{}";
    instr match { CALL runtime_convTstring }
    gen { interface{}(s) }
}

pattern go_runtime_convT64 {
    name: "go: int64 -> interface (runtime.convT64)";
    library: "go-datatype";
    description: "Convert int64 to interface{}";
    instr match { CALL runtime_convT64 }
    gen { interface{}(v) }
}

pattern go_runtime_convT32 {
    name: "go: int32 -> interface (runtime.convT32)";
    library: "go-datatype";
    description: "Convert int32 to interface{}";
    instr match { CALL runtime_convT32 }
    gen { interface{}(v) }
}

pattern go_runtime_growslice {
    name: "go: grow slice (runtime.growslice)";
    library: "go-datatype";
    description: "Slice growth via runtime.growslice";
    instr match { CALL runtime_growslice }
    gen { append(s, elem...) }
}

pattern go_runtime_slicebytetostring {
    name: "go: []byte -> string (runtime.slicebytetostring)";
    library: "go-datatype";
    description: "Convert byte slice to string";
    instr match { CALL runtime_slicebytetostring }
    gen { string(b) }
}

pattern go_runtime_stringtoslicebyte {
    name: "go: string -> []byte (runtime.stringtoslicebyte)";
    library: "go-datatype";
    description: "Convert string to byte slice";
    instr match { CALL runtime_stringtoslicebyte }
    gen { []byte(s) }
}

pattern go_os_exec_command {
    name: "go: exec.Command";
    library: "go-datatype";
    description: "Execute command via os/exec.Command";
    instr match { CALL os_exec_Command }
    gen { exec.Command(name, args...) }
}

pattern go_encoding_json_encode {
    name: "go: json.Encode";
    library: "go-datatype";
    description: "JSON encode via encoding/json.Encoder.Encode";
    instr match { CALL encoding_json_Encoder_Encode }
    gen { json.NewEncoder(w).Encode(v) }
}

pattern go_net_http_serve {
    name: "go: http.Serve";
    library: "go-datatype";
    description: "HTTP server via net/http.Server.Serve";
    instr match { CALL net_http_Server_Serve }
    gen { server.Serve(listener) }
}
