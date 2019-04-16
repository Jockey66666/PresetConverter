SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=386
go build -o ./Windows/x86/PresetConverterWin.exe ../src/

SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o ./Windows/x64/PresetConverterWin.exe ../src/

SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -o ./Mac/PresetConverterMac ../src/

cp ../src/thumbnail.png ./Windows/
cp ../src/thumbnail.png ./Mac/