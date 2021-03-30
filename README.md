# Stock-Cipher-API 

---

## Descripción
Son dos APIs que interactúan en conjunto `stock-api` y `encryption-api` para obtener información de un stock y 
luego encriptarlo utilizando AES-256 en la encryption-api.

---

## ¿Como iniciar las aplicaciones?

Este proyecto usa dock-compose (es necesario tener docker) para iniciar las dos APIs, 
puedes ejecutar los siguientes comandos desde el root del proyecto:

    docker-compose build
    docker-compose up

---

### Stock-API

Cuenta con dos endpoints, un GET que se encarga de buscar la información del último día para un stock y lo encripta obteniendo el payload
y un token para poder descifrar (el token es one shoot) y un POST para descifrar este payload junto con el token. 

#### Recursos

**GET para stock encriptado**

Endpoint: /stock/:symbol

Ejemplo del request:

    curl -X GET http://localhost:8080/stock/IBM 
    
Response:

    {
        "token": "0e3b0862-4823-40ee-8d47-a259d732e012",
        "payload": "c0de60bbce8fe0f32e7d72d44ea5cda698f46ee235e332f0c9e60a923029c5624c542013c323dd5505a5bfd85826fcf4d474e6bed12d2285acb2ba7a0f987c78621e9c9d7da4614ad4052575014e7ae853ed9f705397c9af6c0ecb526ead11216cc7005efa4870626f89cff7555f93706d65075701fead9ea8d2ed4a40683a75e258804e769afb2a1420c3baf65c63a9"
    }    

**POST para descifrar el stock**

Endpoint: /stock/decrypt/:token

Ejemplo:

    curl -X POST \
      http://localhost:8080/stock/decrypt/0e3b0862-4823-40ee-8d47-a259d732e012 \
      -d '{
    	"payload": "c0de60bbce8fe0f32e7d72d44ea5cda698f46ee235e332f0c9e60a923029c5624c542013c323dd5505a5bfd85826fcf4d474e6bed12d2285acb2ba7a0f987c78621e9c9d7da4614ad4052575014e7ae853ed9f705397c9af6c0ecb526ead11216cc7005efa4870626f89cff7555f93706d65075701fead9ea8d2ed4a40683a75e258804e769afb2a1420c3baf65c63a9"
        }'
        
Response:
    
    {
        "date": "2021-03-29T00:00:00Z",
        "open": "135.9800",
        "high": "137.0700",
        "low": "135.5100",
        "close": "135.8600",
        "volume": "4622664"
    }