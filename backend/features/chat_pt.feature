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

  Cenário: Obter sala inexistente retorna não encontrado
    Quando eu enviar uma requisição GET para "/api/v1/chat/rooms/00000000-0000-4000-8000-000000000099"
    Então o código de status da resposta deve ser 404

  Cenário: Listar mensagens de uma sala com paginação por cursor
    Quando eu enviar uma requisição GET para "/api/v1/chat/rooms/00000000-0000-4000-8000-000000000088/messages"
    Então o código de status da resposta deve ser 200
    E a resposta deve conter "data"
    E a resposta deve conter "next_cursor"
    E a resposta deve conter "has_more"

  Cenário: Salas de chat exigem autenticação
    Dado que eu não estou autenticado
    Quando eu enviar uma requisição GET para "/api/v1/chat/rooms"
    Então o código de status da resposta deve ser 401
