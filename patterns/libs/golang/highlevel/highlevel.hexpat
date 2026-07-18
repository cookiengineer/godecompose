arch x86_64;
platform linux, darwin, windows, freebsd;

pattern fmt_Println_string {
    name: "fmt.Println";
    library: "go-stdlib";
    description: "fmt.Println with a string argument — recovers to fmt.Println(str)";

    instr match {
        CALL fmt_Fprintln
    }

    gen {
        fmt.Println($str);
    }

    bind {
        str as "str";
    }
}

pattern fmt_Printf_format {
    name: "fmt.Printf";
    library: "go-stdlib";
    description: "fmt.Printf with format string and args — recovers to fmt.Printf(fmt, args...)";

    instr match {
        CALL fmt_Fprintf
    }

    gen {
        fmt.Printf($format, $args...);
    }

    bind {
        format as "fmt";
        args as "args";
    }
}

pattern fmt_Sprintf_format {
    name: "fmt.Sprintf";
    library: "go-stdlib";
    description: "fmt.Sprintf with format string and args — recovers to fmt.Sprintf(fmt, args...)";

    instr match {
        CALL fmt_Sprintf
    }

    gen {
        s := fmt.Sprintf($format, $args...);
    }

    bind {
        format as "fmt";
        args as "args";
    }
}

pattern fmt_Errorf_format {
    name: "fmt.Errorf";
    library: "go-stdlib";
    description: "fmt.Errorf with format string and args — recovers to fmt.Errorf(fmt, args...)";

    instr match {
        CALL fmt_Errorf
    }

    gen {
        err := fmt.Errorf($format, $args...);
    }

    bind {
        format as "fmt";
        args as "args";
    }
}

pattern os_File_Write_call {
    name: "os.File.Write";
    library: "go-stdlib";
    description: "os.File.Write method call";

    instr match {
        CALL os_File_Write
    }

    gen {
        os.File.Write($buf);
    }

    bind {
        buf as "buf";
    }
}

pattern os_File_Read_call {
    name: "os.File.Read";
    library: "go-stdlib";
    description: "os.File.Read method call";

    instr match {
        CALL os_File_Read
    }

    gen {
        os.File.Read($buf);
    }

    bind {
        buf as "buf";
    }
}

pattern os_File_Close_call {
    name: "os.File.Close";
    library: "go-stdlib";
    description: "os.File.Close method call";

    instr match {
        CALL os_File_Close
    }

    gen {
        os.File.Close();
    }
}

pattern time_Sleep_call {
    name: "time.Sleep";
    library: "go-stdlib";
    description: "time.Sleep call";

    instr match {
        CALL time_Sleep
    }

    gen {
        time.Sleep($d);
    }

    bind {
        d as "d";
    }
}

pattern time_Now_call {
    name: "time.Now";
    library: "go-stdlib";
    description: "time.Now call";

    instr match {
        CALL time_Now
    }

    gen {
        now := time.Now();
    }
}

pattern sync_Mutex_Lock_call {
    name: "sync.Mutex.Lock";
    library: "go-stdlib";
    description: "sync.Mutex.Lock call in ABIInternal";

    instr match {
        CALL sync_Mutex_Lock
    }

    gen {
        $mu.Lock();
    }

    bind {
        mu as "mu";
    }
}

pattern sync_Mutex_Unlock_call {
    name: "sync.Mutex.Unlock";
    library: "go-stdlib";
    description: "sync.Mutex.Unlock call in ABIInternal";

    instr match {
        CALL sync_Mutex_Unlock
    }

    gen {
        $mu.Unlock();
    }

    bind {
        mu as "mu";
    }
}

pattern sync_WaitGroup_Done_call {
    name: "sync.WaitGroup.Done";
    library: "go-stdlib";
    description: "sync.WaitGroup.Done via Add(-1)";

    instr match {
        CALL sync_WaitGroup_Add
    }

    gen {
        $wg.Done();
    }

    bind {
        wg as "wg";
    }
}

pattern json_Marshal_call {
    name: "encoding/json.Marshal";
    library: "go-stdlib";
    description: "json.Marshal call";

    instr match {
        CALL json_Marshal
    }

    gen {
        data, err := json.Marshal($v);
    }

    bind {
        v as "v";
    }
}

pattern json_Unmarshal_call {
    name: "encoding/json.Unmarshal";
    library: "go-stdlib";
    description: "json.Unmarshal call";

    instr match {
        CALL json_Unmarshal
    }

    gen {
        err := json.Unmarshal($data, $v);
    }

    bind {
        data as "data";
        v as "v";
    }
}

pattern net_http_Get_call {
    name: "net/http.Get";
    library: "go-stdlib";
    description: "http.Get call";

    instr match {
        CALL http_Get
    }

    gen {
        resp, err := http.Get($url);
    }

    bind {
        url as "url";
    }
}

pattern crypto_sha256_Sum256_call {
    name: "crypto/sha256.Sum256";
    library: "go-stdlib";
    description: "sha256.Sum256 call";

    instr match {
        CALL sha256_Sum256
    }

    gen {
        h := sha256.Sum256($data);
    }

    bind {
        data as "data";
    }
}

pattern crypto_rand_Read_call {
    name: "crypto/rand.Read";
    library: "go-stdlib";
    description: "crypto/rand.Read call";

    instr match {
        CALL crypto_rand_Read
    }

    gen {
        crypto_rand.Read($buf);
    }

    bind {
        buf as "buf";
    }
}

pattern runtime_makeslice_call {
    name: "runtime.makeslice";
    library: "go-runtime";
    description: "runtime.makeslice call — slice allocation";

    instr match {
        CALL runtime_makeslice
    }

    gen {
        make([]T, $len, $cap);
    }

    bind {
        len as "len";
        cap as "cap";
    }
}

pattern runtime_growslice_call {
    name: "runtime.growslice";
    library: "go-runtime";
    description: "runtime.growslice call — append operation";

    instr match {
        CALL runtime_growslice
    }

    gen {
        append($s, ...);
    }

    bind {
        s as "s";
    }
}

pattern runtime_chan_recv {
    name: "runtime.chanrecv";
    library: "go-runtime";
    description: "Channel receive operation";

    instr match {
        CALL runtime_chanrecv1
        | CALL runtime_chanrecv2
    }

    gen {
        <-$ch;
    }

    bind {
        ch as "ch";
    }
}

pattern runtime_chan_send {
    name: "runtime.chansend";
    library: "go-runtime";
    description: "Channel send operation";

    instr match {
        CALL runtime_chansend1
    }

    gen {
        // ch <- val
    }

    bind {
        ch as "ch";
        val as "val";
    }
}

pattern runtime_map_access {
    name: "runtime.mapaccess";
    library: "go-runtime";
    description: "Map access operation";

    instr match {
        CALL runtime_mapaccess1
        | CALL runtime_mapaccess2
    }

    gen {
        // map access: $m[$key]
    }

    bind {
        m as "m";
        key as "key";
    }
}

pattern runtime_convT {
    name: "runtime.convT";
    library: "go-runtime";
    description: "Type conversion to interface";

    instr match {
        CALL runtime_convT64
        | CALL runtime_convTstring
        | CALL runtime_convTslice
        | CALL runtime_convT2E
        | CALL runtime_convT2I
    }

    gen {
        // convT: convert to interface
    }
}

pattern runtime_defer {
    name: "runtime.deferproc";
    library: "go-runtime";
    description: "Defer statement setup";

    instr match {
        CALL runtime_deferproc
        | CALL runtime_deferprocStack
    }

    gen {
        defer func();
    }
}

pattern runtime_panic {
    name: "runtime.gopanic";
    library: "go-runtime";
    description: "Panic call";

    instr match {
        CALL runtime_gopanic
    }

    gen {
        panic($val);
    }

    bind {
        val as "val";
    }
}

pattern runtime_type_assert {
    name: "runtime.assert";
    library: "go-runtime";
    description: "Type assertion operation";

    instr match {
        CALL runtime_assertE2I
        | CALL runtime_assertI2I
        | CALL runtime_assertE2I2
    }

    gen {
        // type assertion: val.(T)
    }

    bind {
        val as "val";
    }
}
