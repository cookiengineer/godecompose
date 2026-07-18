arch x86_64;
platform linux, darwin, windows, freebsd;

pattern io_ReadAll {
    name: "io.ReadAll";
    library: "go-stdlib";
    description: "Reads everything from a reader";

    instr match {
        MOVQ r, AX
        CALL io_ReadAll
    }

    gen {
        io.ReadAll($r);
    }

    bind {
        r as "r";
    }
}

pattern io_Copy {
    name: "io.Copy";
    library: "go-stdlib";
    description: "Copies from reader to writer";

    instr match {
        MOVQ dst, AX
        MOVQ src, BX
        CALL io_Copy
    }

    gen {
        io.Copy($dst, $src);
    }

    bind {
        dst as "dst";
        src as "src";
    }
}

pattern io_WriteString {
    name: "io.WriteString";
    library: "go-stdlib";
    description: "Writes a string to a writer";

    instr match {
        MOVQ w, AX
        MOVQ s, BX
        CALL io_WriteString
    }

    gen {
        io.WriteString($w, $s);
    }

    bind {
        w as "w";
        s as "s";
    }
}

pattern time_Now {
    name: "time.Now";
    library: "go-stdlib";
    description: "Returns current time";

    instr match {
        CALL time_Now
    }

    gen {
        time.Now();
    }
}

pattern time_Sleep {
    name: "time.Sleep";
    library: "go-stdlib";
    description: "Sleeps for a duration";

    instr match {
        MOVQ d, AX
        CALL time_Sleep
    }

    gen {
        time.Sleep($d);
    }

    bind {
        d as "d";
    }
}

pattern time_After {
    name: "time.After";
    library: "go-stdlib";
    description: "Returns a channel that fires after a duration";

    instr match {
        MOVQ d, AX
        CALL time_After
    }

    gen {
        time.After($d);
    }

    bind {
        d as "d";
    }
}

pattern time_Since {
    name: "time.Since";
    library: "go-stdlib";
    description: "Returns time since a given time";

    instr match {
        MOVQ t, AX
        CALL time_Since
    }

    gen {
        time.Since($t);
    }

    bind {
        t as "t";
    }
}

pattern time_Tick {
    name: "time.NewTicker";
    library: "go-stdlib";
    description: "Creates a new ticker";

    instr match {
        MOVQ d, AX
        CALL time_NewTicker
    }

    gen {
        time.NewTicker($d);
    }

    bind {
        d as "d";
    }
}
