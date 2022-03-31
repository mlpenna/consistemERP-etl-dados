# RPA e ETL para Consistem ERP
Um pequeno projeto em Go que implementa um RPA para extração de dados de algumas rotinas do [Consistem ERP](https://consistem.com.br/) e alguns ELTs para transportá-los a um BD PostgreSQL.

O projeto possui três componentes principais:
### Agendador de tarefas

Utilizando a biblioteca cron disponível [aqui](https://github.com/robfig/cron), implementa um gerenciador e agendador de execução das funçoes de RPA e ELT do programa
para atender aos diferentes requisitos de carregamento dos dados pela área de dados.

### RPA

Utiliza a biblioteca [go-vgo/robotgo](https://github.com/go-vgo/robotgo) para extrair relatórios do sistema em formato .csv. Isto é feito através da interface web do ERP
e automação de comandos na tela.

### ETL

Após extraídos, os dados presentes nos arquivos .csv são lidos, tratados, e carregados em um banco de dados PostgreSQL para serem consumidos pela área de dados.

