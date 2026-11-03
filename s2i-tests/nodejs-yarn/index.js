const http = require('node:http');

const port = Number(process.env.PORT) || 3000;

const server = http.createServer((req, res) => {
  if (req.url === '/health') {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ status: 'UP' }));
    return;
  }
  res.writeHead(200, { 'Content-Type': 'text/plain' });
  res.end('Hello from nodejs-yarn\n');
});

server.listen(port, () => console.log(`nodejs-yarn listening on port ${port}`));

process.on('SIGTERM', () => server.close(() => process.exit(0)));
