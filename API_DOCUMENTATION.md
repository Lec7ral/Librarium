# Documentación de la API Librarium

Esta es la documentación de referencia para todos los endpoints disponibles en la API Librarium.

**URL Base:** `http://localhost:8080`

---

## Autenticación

La mayoría de los endpoints de escritura (`POST`, `PUT`, `DELETE`) están protegidos y requieren un token JWT.

**Flujo:**
1.  Regístrate con `POST /register`.
2.  Inicia sesión con `POST /login` para obtener un token.
3.  En las peticiones a endpoints protegidos, incluye la siguiente cabecera:
    - `Authorization: Bearer <tu_token_jwt>`

### `POST /register`

Registra un nuevo usuario.

**Cuerpo de la Petición (Body):** `application/json`
```json
{
    "username": "nuevo_usuario",
    "password": "una_contraseña_segura"
}
```

**Respuesta Exitosa:** `201 Created`

### `POST /login`

Inicia sesión y obtiene un token JWT.

**Cuerpo de la Petición (Body):** `application/json`
```json
{
    "username": "mi_usuario",
    "password": "mi_contraseña"
}
```

**Respuesta Exitosa:** `200 OK`
```json
{
    "token": "ey..."
}
```

---

## Libros (`/books`)

### `GET /books`

Obtiene una lista paginada, filtrada y ordenada de libros.

**Parámetros de Consulta (Query Params):**
-   `limit` (int, opcional, por defecto: 20): Número de resultados por página.
-   `page` (int, opcional, por defecto: 1): Número de la página a obtener.
-   `title` (string, opcional): Filtra los libros cuyo título contenga este texto.
-   `author` (string, opcional): Filtra los libros cuyo nombre de autor contenga este texto.
-   `sort` (string, opcional): Campo por el que ordenar. Valores permitidos: `title`, `author`, `published_date`, `stock`.
-   `order` (string, opcional, por defecto: `asc`): Dirección del orden. Valores permitidos: `asc`, `desc`.

**Ejemplo:** `GET /books?limit=5&page=2&sort=title&order=desc`

**Respuesta Exitosa:** `200 OK` (con metadatos de paginación)

### `GET /books/{id}`

Obtiene los detalles de un libro específico por su ID.

**Respuesta Exitosa:** `200 OK`

### `POST /books`

Crea un nuevo libro. **(Endpoint Protegido)**

**Cuerpo de la Petición (Body):** `application/json`
```json
{
    "title": "1984",
    "published_date": "1949-06-08",
    "isbn": "978-0451524935",
    "stock": 5,
    "author_id": 1
}
```

**Respuesta Exitosa:** `201 Created`

### `PUT /books/{id}`

Actualiza la información de un libro existente. **(Endpoint Protegido)**

**Cuerpo de la Petición (Body):** `application/json` (igual que en `POST /books`)

**Respuesta Exitosa:** `200 OK`

### `DELETE /books/{id}`

Elimina un libro por su ID. **(Endpoint Protegido)**

**Respuesta Exitosa:** `204 No Content`

---

## Autores (`/authors`)

### `GET /authors`

Obtiene una lista de todos los autores.

**Respuesta Exitosa:** `200 OK`

### `GET /authors/{id}`

Obtiene los detalles de un autor específico por su ID.

**Respuesta Exitosa:** `200 OK`

### `POST /authors`

Crea un nuevo autor. **(Endpoint Protegido)**

**Cuerpo de la Petición (Body):** `application/json`
```json
{
    "name": "George Orwell",
    "bio": "Autor de 1984 y Rebelión en la granja."
}
```

**Respuesta Exitosa:** `201 Created`
