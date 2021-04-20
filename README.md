# Sanakirja.org CLI
CLI-based client for sanakirja.org, dictionary site made by HS (Helsingin Sanomat)

### Usage
- `-i`: Interactive mode (WIP!)
- `--src`: Source language, language list available with command `langlist`
- `--target`: Source language, language list available with command `langlist`
- `--text`: Text or a word to search in sanakirja.org
- `--exp`: Shows examples with results
- `--syn`: Shows synonyms with results
- `--def`: Shows definitions with results
- `--pro`: Shows pronunciation with results
- `--alt`: Shows alternative spellings
- `-c`: Command to run. Available commands: `langlist` (lists all available languages)
- `-l`: Application language, available: fi, sv, en, fr
- `--json`: Outputs everything in JSON format
- `--prettify`: Prettify JSON
- `--vv`: Verbose level, Only errors: 0, Warnings: 1, Logs: 2, Debug: 3
- `-v`: Shows current version

### Compiling and installation

1. Install go to your system using this [guide](https://golang.org/doc/install)
2. Create directory `build` with command `mkdir build`
3. Build with command `go build -i -o build/sanakirja Main.go`
4. Compiled binary is located in folder `build`

For linux, you could copy that binary to /usr/bin, which would make this utility accessible from your cli.
Run `cp build/sanakirja /usr/bin/`
