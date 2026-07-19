arch x86_64;
platform linux, darwin, windows, freebsd;

pattern os_Getenv {
    name: "os.Getenv";
    library: "go-stdlib";
    description: "Gets an environment variable";

    instr match {
        MOVQ key, AX
        CALL os_Getenv
    }

    gen {
        os.Getenv($key);
    }

    bind {
        key as "key";
    }
}
