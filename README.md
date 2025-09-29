# netscan
to start 
cd /apps/server
go mod tidy          
go run ./cmd/netscan-server -oui ./oui.csv -addr :8080

then cd /apps/web
npm run dev