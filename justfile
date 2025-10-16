# set shell := ["pwsh", "", "-CommandWithArgs"]
# set positional-arguments
shebang := 'pwsh.exe'
# Variables
dist := ".dist"
exe_name := "korrectkm"
mod_name := "korrectkm"
css_input := "tailwind.css"
css_output := "embeded/root/assets/css/tailwind.css"
ld_flags :="-H=windowsgui -s -w -X korrectkm/config.Mode=production"

default:
  just --list

# Build production CSS
build:
    tailwindcss -i {{css_input}} -o {{css_output}} --minify

win64:
    #!{{shebang}}
    $env:Path = "C:\Go\go.125\bin;C:\go\gcc\mingw64\bin;" + $env:Path
    $env:GOARCH = "amd64"
    $env:GOOS = "windows"
    $env:CGO_ENABLED = 1
    if (-Not (Test-Path go.mod)) {
      go mod init {{mod_name}}
    }
    go mod tidy -go 1.25 -v
    if(-Not $?) { exit }
    if (-Not (Test-Path "{{dist}}")) { New-Item -ItemType Directory -Force -Path "{{dist}}" | Out-Null }
    Remove-Item -Force -ErrorAction SilentlyContinue -LiteralPath "{{dist}}\{{exe_name}}.exe","{{dist}}\{{exe_name}}_64.exe"
    go build -ldflags="{{ld_flags}}" -o "{{dist}}\{{exe_name}}_64.exe" ./cmd
    if(-Not $?) { exit }
    upx --force-overwrite -o {{dist}}\{{exe_name}}.exe {{dist}}\{{exe_name}}_64.exe
