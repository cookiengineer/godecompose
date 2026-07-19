arch x86_64;
platform linux, darwin, windows, freebsd;

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
