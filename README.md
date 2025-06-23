# Vorin - Web Directory & Admin Scanner

**Vorin** é uma ferramenta de varredura para diretórios e caminhos ocultos em aplicações web. Escrita em Go, com foco em desempenho, simplicidade e clareza nos resultados. Inspirada em ferramentas como Gobuster e FFUF, mas com estilo próprio.

## Características

- Scan rápido com múltiplas threads
- Suporte a wordlist customizada
- Detecção de diretórios comuns, páginas administrativas, arquivos sensíveis
- Saída colorida e limpa no terminal
- Fácil de compilar e usar em qualquer sistema

## Instalação

```bash
git clone https://github.com/JuaanReis/vorin.git
cd vorin
go build -o vorin
```

## Uso

```bash
./vorin -u http://example.com -w wordlist.txt -t 50
```

### Parâmetros

| Flag     | Descrição                                  |
|----------|--------------------------------------------|
| `-u`     | URL do alvo                                |
| `-w`     | Caminho da wordlist                        |
| `-t`     | Número de threads (padrão: 50)             |

## Exemplo de uso

Abaixo, um exemplo real de execução da ferramenta em ambiente de testes (como o site `testphp.vulnweb.com`), demonstrando a detecção de diretórios e arquivos sensíveis:

![Exemplo de Scan](assets/screenshots/showing.png)

> Todos os testes foram realizados em ambiente seguro e controlado, sem impactar sistemas reais.

## Wordlist

Você pode usar qualquer wordlist personalizada. A recomendação é começar com uma lista leve e ir aumentando conforme necessidade.

Exemplo de conteúdo:

```
admin
admin/login
.git
.htaccess
phpinfo.php
uploads
includes
```

## Segurança

Esta ferramenta foi criada para **uso educacional e profissional controlado**. O uso indevido contra sistemas sem autorização pode configurar crime conforme a legislação vigente.

## License

MIT License. Consulte o arquivo `LICENSE` para mais detalhes.
