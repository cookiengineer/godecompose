arch x86_64;
platform linux, darwin, windows, freebsd;

pattern context_WithCancel {
    name: "context.WithCancel";
    library: "go-stdlib";
    description: "Returns a copy of parent with a new Done channel";

    instr match {
        CALL context_WithCancel
    }

    gen {
        ctx, cancel := context.WithCancel($parent);
    }

    bind {
        parent as "parent";
    }
}

pattern context_WithTimeout {
    name: "context.WithTimeout";
    library: "go-stdlib";
    description: "Returns a copy of parent with timeout";

    instr match {
        CALL context_WithTimeout
    }

    gen {
        ctx, cancel := context.WithTimeout($parent, $timeout);
    }

    bind {
        parent as "parent";
        timeout as "timeout";
    }
}

pattern context_WithDeadline {
    name: "context.WithDeadline";
    library: "go-stdlib";
    description: "Returns a copy of parent with deadline";

    instr match {
        CALL context_WithDeadline
    }

    gen {
        ctx, cancel := context.WithDeadline($parent, $deadline);
    }

    bind {
        parent as "parent";
        deadline as "deadline";
    }
}

pattern context_Background {
    name: "context.Background";
    library: "go-stdlib";
    description: "Returns the background context";

    instr match {
        CALL context_Background
    }

    gen {
        ctx := context.Background();
    }
}

pattern context_TODO {
    name: "context.TODO";
    library: "go-stdlib";
    description: "Returns the TODO context";

    instr match {
        CALL context_TODO
    }

    gen {
        ctx := context.TODO();
    }
}

pattern context_WithValue {
    name: "context.WithValue";
    library: "go-stdlib";
    description: "Returns a copy of parent with associated key-value pair";

    instr match {
        CALL context_WithValue
    }

    gen {
        ctx := context.WithValue($parent, $key, $val);
    }

    bind {
        parent as "parent";
        key as "key";
        val as "val";
    }
}
