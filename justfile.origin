# set shell := ["pwsh", "", "-CommandWithArgs"]
# set positional-arguments
shebang := 'pwsh.exe'
# Variables
exe_name := "korrectkm"
mod_name := "korrectkm"
css_input := "tailwind.css"
css_output := "embeded/root/assets/css/tailwind.css"
db_path := "./repo/selfdb/migrations/sqlite/zakaz.db"
migrations_dir := "././repo/selfdb/migrations/sqlite"
ld_flags :="-H=windowsgui -s -w -X 'korrectkm/config.Mode=production'"
# dev_port := "4000"
# browser_sync_port := "4001"

default:
  just --list

# Build production CSS
build:
    tailwindcss -i {{css_input}} -o {{css_output}} --minify

# Run database migrations
migrate cmnd="up":
    #!{{shebang}}
    if (-Not (Test-Path {{db_path}})) {
      Remove-Item {{db_path}}
    }
    sqlite3 {{db_path}} ""
    goose sqlite3 {{db_path}} {{cmnd}} --dir={{migrations_dir}}

migrateinit :
    #!{{shebang}}
    if (Test-Path {{db_path}}) {
      Remove-Item {{db_path}}
    }
    sqlite3 {{db_path}} ""
    goose sqlite3 {{db_path}} up --dir={{migrations_dir}}

win64:
    #!{{shebang}}
    $env:Path = "C:\Go\go.124\bin;C:\go\gcc\mingw64\bin;" + $env:Path
    $env:GOARCH = "amd64"
    $env:GOOS = "windows"
    $env:CGO_ENABLED = 1
    if (-Not (Test-Path go.mod)) {
      go mod init {{mod_name}}
    }
    go mod tidy -go 1.24 -v
    if(-Not $?) { exit }
    Remove-Item ..\dist\{{exe_name}}.exe, ..\dist\{{exe_name}}_64.exe 2>$null
    # go build -ldflags="{{ld_flags}}" -o ../dist/{{exe_name}}_64_.exe ./cmd
    go build -tags "production,desktop" -gcflags "all=-N -l" -ldflags="{{ld_flags}}" -o ../dist/{{exe_name}}_64.exe ./cmd
    if(-Not $?) { exit }
    upx --force-overwrite -o ../dist/{{exe_name}}.exe ../dist/{{exe_name}}_64.exe
