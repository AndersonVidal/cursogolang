# Desafio: Client-Server API com Context e Banco de Dados

Este projeto implementa dois sistemas em Go (client.go e server.go) que simulam uma comunicação cliente-servidor para consulta da cotação do dólar. O projeto utiliza requisições HTTP, manipulação de arquivos, contextos com timeout e persistência em banco de dados SQLite.

## 🧠 Objetivo
- O client.go deve fazer uma requisição HTTP para o server.go solicitando a cotação do dólar.
- O server.go deve:
    - Consultar a API pública `https://economia.awesomeapi.com.br/json/last/USD-BRL`
    - Retornar apenas o campo "bid" para o cliente
    - Persistir a cotação recebida em um banco de dados SQLite (applicationDB.db)
- Ambas as aplicações devem utilizar o pacote context com timeouts definidos.

## ⏱️ Regras de timeout
- `client.go` deve:
    - Ter um timeout de 300ms para receber a resposta do servidor
    - Salvar a cotação em um arquivo cotacao.txt no formato: Dólar: {valor}
- `server.go` deve:
    - Ter um timeout de 200ms para consultar a API de cotação
    - Ter um timeout de 10ms para persistir os dados no banco SQLite
- Todos os contextos devem registrar erros em caso de timeout

## 🚀 Como executar
1. Executar o servidor
```bash
cd client-server-api/server
go run server.go
```
Isso iniciará o servidor HTTP na porta :8080, expondo o endpoint /cotacao.

2. Executar o cliente
Em outro terminal:
```bash
cd client-server-api/client
go run client.go
```
O cliente fará uma requisição para `http://localhost:8080/cotacao`, receberá o valor da cotação e salvará no arquivo cotacao.txt.

## ✅ Exemplo de saída no arquivo cotacao.txt
`Dólar: 5.2473`

## 📁 Estrutura do projeto
```text
client-server-api/
├── client/
│   ├── client.go         # Cliente HTTP que consulta a cotação
│   ├── cotacao.txt       # Arquivo onde a cotação é salva
├── server/
│   ├── server.go         # Servidor HTTP com contexto, API e banco
│   ├── go.mod / go.sum   # Dependências Go
│   └── database/
│       └── applicationDB.db  # Banco SQLite com cotações persistidas
```

## 📌 Logs esperados (simulados)
```bash
[server] Buscando cotação na API externa...
[server] Persistindo cotação no banco...
[client] Cotação recebida: 5.2473
[client] Cotação salva em cotacao.txt
```
Em caso de timeout, os logs devem exibir erros como:
```bash
[server] Erro: Timeout ao consultar a API
[client] Erro: resposta do servidor demorou mais de 300ms
```
