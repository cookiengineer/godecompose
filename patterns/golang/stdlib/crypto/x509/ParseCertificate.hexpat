arch x86_64;
platform linux, darwin, windows, freebsd;

pattern crypto_x509_ParseCertificate {
    name: "crypto/x509.ParseCertificate";
    library: "go-stdlib";
    description: "Parses an X.509 certificate";

    instr match {
        MOVQ asn1, AX
        CALL x509_ParseCertificate
    }

    gen {
        x509.ParseCertificate($asn1);
    }

    bind {
        asn1 as "data";
    }
}
