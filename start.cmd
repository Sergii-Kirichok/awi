del bin\\awi-service.exe
go build -tags="STA" -o bin\\awi-service.exe awi.go
bin\\awi-service.exe install
bin\\awi-service.exe debug
