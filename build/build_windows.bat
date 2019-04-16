SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=386
go build -o PresetConverterWin32.exe ../

SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o PresetConverterWin64.exe ../

SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -o PresetConverterMac ../