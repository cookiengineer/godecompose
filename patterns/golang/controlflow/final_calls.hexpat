arch x86_64;
platform linux, darwin, windows, freebsd;

pattern go_internal_filepathlite_Base {
    name: "go: filepath.Base";
    library: "go-stdlib";
    description: "Get base filename via internal/filepathlite.Base";
    instr match { CALL internal_filepathlite_Base }
    gen { filepath.Base(path) }
}

pattern go_time_Time_Year {
    name: "go: time.Time.Year";
    library: "go-stdlib";
    description: "Get year from time via time.Time.Year";
    instr match { CALL time_Time_Year }
    gen { t.Year() }
}

pattern go_time_Time_Month {
    name: "go: time.Time.Month";
    library: "go-stdlib";
    description: "Get month from time via time.Time.Month";
    instr match { CALL time_Time_Month }
    gen { t.Month() }
}

pattern go_time_Time_Day {
    name: "go: time.Time.Day";
    library: "go-stdlib";
    description: "Get day from time via time.Time.Day";
    instr match { CALL time_Time_Day }
    gen { t.Day() }
}

pattern go_time_Time_Date {
    name: "go: time.Date";
    library: "go-stdlib";
    description: "Construct time via time.Date";
    instr match { CALL time_Date }
    gen { time.Date(year, month, day, hour, min, sec, nsec, loc) }
}

pattern go_runtime_panicwrap {
    name: "go: runtime.panicwrap";
    library: "go-runtime";
    description: "Wrapper panic via runtime.panicwrap";
    instr match { CALL runtime_panicwrap }
    gen { panic(err) }
}

pattern go_os_MkdirAll {
    name: "go: os.MkdirAll";
    library: "go-stdlib";
    description: "Create directory tree via os.MkdirAll";
    instr match { CALL os_MkdirAll }
    gen { os.MkdirAll(path, perm) }
}

pattern go_time_Time_After {
    name: "go: time.Time.After";
    library: "go-stdlib";
    description: "Check if time is after another via time.Time.After";
    instr match { CALL time_Time_After }
    gen { t.After(u) }
}

pattern go_time_Time_Before {
    name: "go: time.Time.Before";
    library: "go-stdlib";
    description: "Check if time is before another via time.Time.Before";
    instr match { CALL time_Time_Before }
    gen { t.Before(u) }
}

pattern go_encoding_json_Marshal {
    name: "go: json.Marshal";
    library: "go-stdlib";
    description: "Marshal to JSON via encoding/json.Marshal";
    instr match { CALL encoding_json_Marshal }
    gen { json.Marshal(v) }
}
