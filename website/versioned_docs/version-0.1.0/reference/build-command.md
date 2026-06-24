---
title: build Command
---

`build` compiles a `.pgpackage` project into a `.pgpkg` archive.

```bash
pgpackage build --project <file.pgpackage> --output <dir-or-file>
```

## Required flags

- `--project`: path to the `.pgpackage` project file
- `--output`: output directory or a direct `.pgpkg` file path

## Output

The command prints the resolved package path to stdout.

When the output points to a directory, the file name is `<PackageId>.pgpkg`.
