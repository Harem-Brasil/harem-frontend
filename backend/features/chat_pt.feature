# language: pt
Funcionalidade: Endpoints de Chat
  Como um usuário autenticado
  Eu quero listar salas, criar salas, abrir uma sala e listar mensagens
  Para que eu possa conversar na plataforma

  Contexto:
    Dado que a API está em execução
    E eu estou autenticado como usuário "johndoe"

  Cenário: Listar salas de chat
    Quando eu enviar uma requisição GET para "/api/v1/chat/rooms"
    Então o código de status da resposta deve ser 200
    E a resposta deve ser um array
    E cada item da lista de salas deve conter os campos obrigatórios de sala

  Cenário: Criar uma nova sala de chat
    Quando eu enviar uma requisição POST para "/api/v1/chat/rooms" com:
      | name       | type   | description        |
      | Sala Geral | public | Discussão aberta   |
    Então o código de status da resposta deve ser 201
    E a resposta deve conter "id"
    E a resposta deve conter "name" com valor "Sala Geral"
    E a resposta deve conter "type" com valor "public"
    E a resposta deve conter "created_by"

  Cenário: Criar sala com tipo padrão público
    Quando eu enviar uma requisição POST para "/api/v1/chat/rooms" com:
      | name       |
      | Minha sala |
    Então o código de status da resposta deve ser 201
    E a resposta deve conter "type" com valor "public"

  Cenário: Criar sala sem nome deve falhar na validação
    Quando eu enviar uma requisição POST para "/api/v1/chat/rooms" com:
      | name | description |
      |      | Sem título    |
    Então o código de status da resposta deve ser 422

  Cenário: Criar sala com tipo inválido deve falhar na validação
    Quando eu enviar uma requisição POST para "/api/v1/chat/rooms" com:
      | name          | type              |
      | Sala inválida | tipo-inexistente  |
    Então o código de status da resposta deve ser 422
    E a resposta deve indicar erro de validação para o campo "type"

  Cenário: Obter uma sala existente retorna os dados
    Quando eu enviar uma requisição POST para "/api/v1/chat/rooms" com:
      | name       | type   | description    |
      | Sala Obter | public | Para teste GET |
    Então o código de status da resposta deve ser 201
    Quando eu enviar uma requisição GET para a sala usando o id da última resposta
    Então o código de status da resposta deve ser 200
    E a resposta deve conter "id"
    E a resposta deve conter "name" com valor "Sala Obter"
    E a resposta deve conter "type" com valor "public"
    E a resposta deve conter "created_by"
    E a resposta deve conter "member_count"

  Cenário: Obter sala inexistente retorna não encontrado
    Quando eu enviar uma requisição GET para "/api/v1/chat/rooms/00000000-0000-4000-8000-000000000099"
    Então o código de status da resposta deve ser 404

  Cenário: Listar mensagens de uma sala com paginação por cursor
    Quando eu enviar uma requisição POST para "/api/v1/chat/rooms" com:
      | name           | type   | description      |
      | Sala mensagens | public | Criada para GET  |
    Então o código de status da resposta deve ser 201
    Quando eu enviar uma requisição GET para as mensagens da sala usando o id da última resposta
    Então o código de status da resposta deve ser 200
    E a resposta deve conter "data"
    E a resposta deve conter "next_cursor"
    E a resposta deve conter "has_more"

  Cenário: Salas de chat exigem autenticação
    Dado que eu não estou autenticado
    Quando eu enviar uma requisição GET para "/api/v1/chat/rooms"
    Então o código de status da resposta deve ser 401
