arch x86_64;
platform linux, darwin, windows, freebsd;

pattern go_slog_logger_log {
    name: "go: slog.Logger.log";
    library: "go-controlflow";
    description: "Structured logging via log/slog.Logger.log";
    instr match {
        CALL log_slog_Logger_log
    }
    gen { slog.Log(ctx, level, msg, attrs...) }
}

pattern go_time_sleep {
    name: "go: time.Sleep";
    library: "go-runtime";
    description: "Sleep via time.Sleep";
    instr match {
        CALL time_Sleep
    }
    gen { time.Sleep(d) }
}

pattern go_time_now {
    name: "go: time.Now";
    library: "go-runtime";
    description: "Get current time via time.Now";
    instr match {
        CALL time_Now
    }
    gen { time.Now() }
}

pattern go_time_since {
    name: "go: time.Since";
    library: "go-runtime";
    description: "Time since a point via time.Since";
    instr match {
        CALL time_Since
    }
    gen { time.Since(t) }
}

pattern go_time_until {
    name: "go: time.Until";
    library: "go-runtime";
    description: "Time until a point via time.Until";
    instr match {
        CALL time_Until
    }
    gen { time.Until(t) }
}

pattern go_time_timer_reset {
    name: "go: time.Timer.Reset";
    library: "go-runtime";
    description: "Reset timer via time.Timer.Reset";
    instr match {
        CALL time_Timer_Reset
    }
    gen { timer.Reset(d) }
}

pattern go_time_stop_timer {
    name: "go: time.stopTimer";
    library: "go-runtime";
    description: "Stop timer via time.stopTimer";
    instr match {
        CALL time_stopTimer
    }
    gen { stopTimer(t) }
}

pattern go_runtime_rand {
    name: "go: runtime.rand";
    library: "go-runtime";
    description: "Random number via runtime.rand or math/rand";
    instr match {
        CALL runtime_rand
    }
    gen { rand.Intn(n) }
}

pattern go_runtime_memequal {
    name: "go: runtime.memequal";
    library: "go-runtime";
    description: "Memory equality via runtime.memequal";
    instr match {
        CALL runtime_memequal
    }
    gen { memequal(a, b) }
}

pattern go_runtime_gcWriteBarrier1 {
    name: "go: gcWriteBarrier1";
    library: "go-runtime";
    description: "Write barrier for GC via gcWriteBarrier1";
    instr match { CALL runtime_gcWriteBarrier1 }
    gen { /* GC wb1 */ }
}

pattern go_runtime_gcWriteBarrier2 {
    name: "go: gcWriteBarrier2";
    library: "go-runtime";
    description: "Write barrier for GC via gcWriteBarrier2";
    instr match { CALL runtime_gcWriteBarrier2 }
    gen { /* GC wb2 */ }
}

pattern go_runtime_gcWriteBarrier3 {
    name: "go: gcWriteBarrier3";
    library: "go-runtime";
    description: "Write barrier for GC via gcWriteBarrier3";
    instr match { CALL runtime_gcWriteBarrier3 }
    gen { /* GC wb3 */ }
}

pattern go_runtime_gcWriteBarrier4 {
    name: "go: gcWriteBarrier4";
    library: "go-runtime";
    description: "Write barrier for GC via gcWriteBarrier4";
    instr match { CALL runtime_gcWriteBarrier4 }
    gen { /* GC wb4 */ }
}

pattern go_os_create_temp {
    name: "go: os.CreateTemp";
    library: "go-controlflow";
    description: "Create temporary file via os.CreateTemp";
    instr match {
        CALL os_CreateTemp
    }
    gen { os.CreateTemp(dir, pattern) }
}

pattern go_os_rename {
    name: "go: os.Rename";
    library: "go-stdlib";
    description: "Rename file via os.Rename";
    instr match {
        CALL os_Rename
    }
    gen { os.Rename(old, new) }
}

pattern go_os_exit {
    name: "go: os.Exit";
    library: "go-stdlib";
    description: "Exit process via os.Exit";
    instr match {
        CALL os_Exit
    }
    gen { os.Exit(code) }
}

pattern go_os_stat {
    name: "go: os.Stat";
    library: "go-stdlib";
    description: "File stat via os.Stat";
    instr match {
        CALL os_Stat
    }
    gen { os.Stat(path) }
}

pattern go_os_symlink {
    name: "go: os.Symlink";
    library: "go-stdlib";
    description: "Create symlink via os.Symlink";
    instr match {
        CALL os_Symlink
    }
    gen { os.Symlink(old, new) }
}

pattern go_syscall_mkdirat {
    name: "go: syscall.Mkdirat";
    library: "go-runtime";
    description: "Make directory via syscall.Mkdirat";
    instr match {
        CALL syscall_Mkdirat
    }
    gen { syscall.Mkdirat(dir, path, mode) }
}

pattern go_internal_sync_mutex_lock_slow {
    name: "go: sync.Mutex.lockSlow";
    library: "go-runtime";
    description: "Mutex slow-path lock via internal/sync.Mutex.lockSlow";
    instr match { CALL internal_sync_Mutex_lockSlow }
    gen { m.Lock() }
}

pattern go_internal_sync_mutex_unlock_slow {
    name: "go: sync.Mutex.unlockSlow";
    library: "go-runtime";
    description: "Mutex slow-path unlock via internal/sync.Mutex.unlockSlow";
    instr match { CALL internal_sync_Mutex_unlockSlow }
    gen { m.Unlock() }
}
