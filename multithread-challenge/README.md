# Desafio: Consulta de CEP com Multithreading

Este desafio consiste em realizar chamadas concorrentes para duas APIs distintas de consulta de CEP e exibir a resposta mais rápida.

## 🧠 Objetivo

* Chamar `BrasilAPI` e `ViaCEP` simultaneamente
* Exibir apenas a resposta mais rápida
* Retornar um erro de timeout caso nenhuma API responda em até 1 segundo

## 🔗 APIs utilizadas

* `https://brasilapi.com.br/api/cep/v1/{cep}`
* `https://viacep.com.br/ws/{cep}/json/`

## 🚀 Como executar

1. Dentro da pasta raiz do repositório, execute:

```bash
cd multithread-challenge
go run main.go
```

Isso irá iniciar a aplicação, instanciando um servidor local.

2. Com isso, a rota `http://localhost:8080/cep/{cep}` será exposta, e você poderá testá-la via navegador ou com o comando `curl`:

```bash
curl http://localhost:8080/cep/00000000
```

3. A resposta será um JSON retornado pela API que responder mais rapidamente. Um atributo adicional (`source`) é incluído na resposta para indicar qual API forneceu os dados:

```json
{
  "source": "ViaCEP",
  "data": {
    "cep": "xxxxx-xxx",
    "logradouro": "xxxxxxxxxxxxxxxx",
    "complemento": "xxxxx",
    "bairro": "xxxxxx",
    "localidade": "xxxxxxx",
    "uf": "xxxxx",
    "ibge": "xxxxxxx",
    "gia": "xxx",
    "ddd": "xxx",
    "siafi": "xxxx",
    "unidade": "xxx",
    "estado": "xxxxxxx",
    "regiao": "xxxxxxx"
  }
}
```

4. Além disso, um log será exibido no terminal onde o servidor estiver rodando, indicando a recepção da requisição:
```bash
2025/06/24 10:56:25 "GET http://localhost:8080/cep/xxxxxxxx HTTP/1.1" from [::1]:47172 - 200 327B in 240.249842ms
```
