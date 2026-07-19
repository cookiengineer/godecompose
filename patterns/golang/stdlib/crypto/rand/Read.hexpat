arch x86_64;
platform linux, darwin, windows, freebsd;

pattern crypto_rand_Read {
    name: "crypto/rand.Read";
    library: "go-stdlib";
    description: "Reads cryptographically secure random bytes";

    instr match {
        MOVQ buf, AX
        CALL crypto_rand_Read
    }

    gen {
        crypto_rand.Read($buf);
    }

    bind {
        buf as "buf";
    }
}
