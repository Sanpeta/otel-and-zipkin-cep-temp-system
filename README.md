# Projeto de Tracing Distribuído com OpenTelemetry e Zipkin

Este projeto implementa o conceito de tracing distribuído entre dois serviços (Serviço A e Serviço B) utilizando OpenTelemetry e Zipkin. O objetivo é rastrear e monitorar a comunicação entre os serviços, coletando dados de performance e análise de tempo de resposta. A comunicação entre os serviços é feita via HTTP, e o tracing é exportado para o Zipkin, que pode ser utilizado para visualização dos dados de tracing.

## Estrutura do Projeto

- **Serviço A**: Um servidor HTTP que recebe uma requisição com um CEP, faz o processamento e chama o Serviço B para obter mais informações sobre o CEP.
- **Serviço B**: Um servidor HTTP que recebe o CEP, faz uma busca (simulada) e retorna a resposta com os dados solicitados.
- **OpenTelemetry**: Usado para capturar dados de tracing (spans) durante a comunicação entre os serviços.
- **Zipkin**: Usado como coletor de tracing. O OpenTelemetry envia os spans para o Zipkin, onde os dados podem ser visualizados.

## Funcionalidades

1. **Serviço A**:

   - Recebe um CEP via requisição HTTP.
   - Envia esse CEP para o Serviço B via requisição HTTP.
   - Utiliza OpenTelemetry para criar um trace distribuído, medindo o tempo de resposta da requisição para o Serviço B.
   - Envia a resposta do Serviço B de volta para o cliente.

2. **Serviço B**:

   - Recebe um CEP via requisição HTTP.
   - Simula uma busca de dados baseados no CEP.
   - Retorna os dados simulados para o Serviço A.

3. **Tracing**:
   - Cada requisição entre os serviços é rastreada utilizando spans do OpenTelemetry.
   - Os spans são exportados para o Zipkin, onde podem ser visualizados.
   - O tempo de resposta de cada serviço é monitorado para ajudar a identificar possíveis gargalos.

## Tecnologias Utilizadas

- **Go**: Linguagem de programação para ambos os serviços.
- **OpenTelemetry**: Framework para rastreamento distribuído.
- **Zipkin**: Sistema de coleta e visualização de tracing distribuído.
- **HTTP**: Protocolo de comunicação entre os serviços.
- **JSON**: Formato de troca de dados entre os serviços.

## Como Executar o Projeto

### Pré-requisitos

1. **Docker**: Para executar o Zipkin como coletor de traces.
2. **Go**: Para compilar e executar os serviços em Go.

### Passos para Execução

#### 1. Executar o Zipkin (Coletor de Tracing)

O Zipkin pode ser executado rapidamente utilizando o Docker. Execute o seguinte comando para iniciar o contêiner do Zipkin:

```bash
docker run -d -p 9411:9411 openzipkin/zipkin
```

Isso fará o Zipkin ficar disponível em `http://localhost:9411`, onde você poderá visualizar os traces.

#### 2. Executar o Serviço A ou Teste A

O Serviço A é responsável por receber o CEP, enviar para o Serviço B e retornar a resposta para o cliente. Para executá-lo:

- Clone o repositório do projeto.
- Navegue até o diretório do Serviço A.
- Execute o seguinte comando para rodar o servidor:

```bash
go run main.go
```

O Serviço A estará disponível em `http://localhost:8080`.

#### 3. Executar o Serviço B ou Teste B

O Serviço B é responsável por simular uma busca de dados com base no CEP enviado pelo Serviço A. Para executá-lo:

- Navegue até o diretório do Serviço B.
- Execute o seguinte comando para rodar o servidor:

```bash
go run main.go
```

O Serviço B estará disponível em `http://localhost:8081`.

#### 4. Fazer uma Requisição

Agora, para testar o sistema, faça uma requisição para o Serviço A. Você pode usar uma ferramenta como o **Postman** ou **curl**.

Exemplo de requisição com **curl**:

```bash
curl -X POST http://localhost:8080/cep -d '{"cep": "01001010"}' -H "Content-Type: application/json"
```

Isso enviará um CEP para o Serviço A, que irá chamar o Serviço B e retornar a resposta com os dados simulados.

#### 5. Visualizar os Traces no Zipkin

Após fazer a requisição, acesse `http://localhost:9411` para visualizar o tracing distribuído no Zipkin. Você poderá ver os detalhes sobre os spans gerados durante a comunicação entre o Serviço A e o Serviço B.

### Estrutura de Diretórios

```
/servico-a
    main.go         # Código do Serviço A
/servico-b
    main.go         # Código do Serviço B
/teste-a
    main.go         # Código do Teste A
/teste-b
    main.go         # Código do Teste B
```

### Detalhamento do Código

- **Serviço A**: O serviço inicia um span quando recebe a requisição com o CEP e faz uma chamada HTTP para o Serviço B. Ele usa o OpenTelemetry para rastrear o tempo de resposta da chamada HTTP para o Serviço B.
- **Serviço B**: O Serviço B simplesmente simula a busca do CEP e retorna os dados ao Serviço A. O tempo de resposta do Serviço B também será monitorado no tracing.
- **Tracing**: O OpenTelemetry cria spans, que representam unidades de trabalho, e os exporta para o Zipkin. O tempo de cada span é registrado, permitindo medir o desempenho entre os serviços.

### Como Personalizar

- **Alterar o Zipkin URL**: Se o seu Zipkin estiver em um servidor diferente, altere a URL do coletor no código de ambos os serviços para apontar para o novo local.
- **Configuração do OpenTelemetry**: O OpenTelemetry pode ser configurado para enviar spans a diferentes tipos de coletores, como o Jaeger, dependendo das necessidades do seu projeto.

## Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para mais detalhes.
