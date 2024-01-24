integration:
	cd cmd/shortener && SERVER_PORT=$(random unused-port) && go build && cd ../.. && shortenertestbeta -test.v -test.run=^TestIteration5$ -binary-path=cmd/shortener/shortener -server-port=$SERVER_PORT
