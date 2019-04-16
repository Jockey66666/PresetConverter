SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=386
go build -o ./bin/Windows/x86/PresetConverterWin.exe ./src/

SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o ./bin/Windows/x64/PresetConverterWin.exe ./src/

SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -o ./bin/Mac/PresetConverterMac ./src/

cp ./src/thumbnail.png ./bin/Windows/
cp ./src/thumbnail.png ./bin/Mac/