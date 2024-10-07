# go-bcdedit

```text
$ go-bcdedit --help
  -create
        /create <id> --object-type <object type(e.g. 0x10200002)> [/d <description>]
        This command creates a new entry in the boot configuration data store.
  -createstore
        /createstore <bcd_file>
        Creates a new and empty boot configuration data store.
  -enum
        /enum all
        This command lists entries in a store.
  -json
        Output result as JSON
  -set
        /set <id> --value-type <ValueType(e.g. RegSz)> --value-raw "BASE64"
        /set <id> --value-type <ValueType(e.g. RegMultiSz)> --value "first" --value "second"
        This command sets an entry option value in the boot configuration data store.
  -store string
        Used to specify a BCD store.
```

# License

[GNU LESSER GENERAL PUBLIC LICENSE 2.1](./LICENSE)

We uses [hivex](https://github.com/libguestfs/hivex) which is licensed under LGPL2.1.