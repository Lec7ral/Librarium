# Librarium API - Roadmap

Este documento traza el progreso y los futuros pasos para el desarrollo de la API Librarium.

---

## ✅ Fase 1: Funcionalidad y Arquitectura Base (Completado)

- [x] **Configuración Inicial del Proyecto:** Estructura de directorios para una API Go escalable (`cmd`, `internal`, `configs`).
- [x] **Base de Datos:** Inicialización de una base de datos SQLite y creación del esquema inicial.
- [x] **CRUD Básico de Libros:** Implementación de los cinco endpoints (GET, GET by ID, POST, PUT, DELETE) para la gestión de libros.
- [x] **Refactorización de la Base de Datos:** Eliminación de la creación de pools de conexión en cada petición mediante Inyección de Dependencias (patrón `Env`).
- [x] **Enrutamiento Profesional:** Reemplazo del enrutador manual por la librería estándar de la industria `gorilla/mux`.
- [x] **Gestión de Configuración:** Externalización de la configuración (puerto, DSN de la BD) para que se cargue desde variables de entorno.
- [x] **Middleware de Logging:** Implementación de un middleware para registrar los detalles de cada petición entrante (método, ruta, estado, duración).

---

## ✅ Fase 2: Seguridad y Calidad del Producto (Completado)

- [x] **Autenticación con JWT:**
    - [x] Creación de endpoints `/register` y `/login`.
    - [x] Protección de las rutas de escritura (POST, PUT, DELETE) mediante un middleware de autenticación JWT.
    - [x] Hasheo seguro de contraseñas con `bcrypt`.
- [x] **Validación Avanzada:** Implementación de validación de datos de entrada a nivel de struct tags con `go-playground/validator`.
- [x] **Apagado Controlado (Graceful Shutdown):** Implementación de un mecanismo para que el servidor termine las peticiones en curso antes de apagarse.
- [x] **Estandarización de Respuestas:**
    - [x] Creación de helpers (`respondWithError`, `respondWithJSON`) en un paquete `web` para asegurar respuestas JSON consistentes.
    - [x] Refactorización de todos los handlers y middlewares para usar los helpers de respuesta.

---

## ✅ Fase 3: Expansión del Dominio y Refactorización (Completado)

- [x] **Modelo de Datos Relacional:**
    - [x] Creación de la entidad `Author` con su propia tabla en la base de datos.
    - [x] Modificación de la tabla `books` para usar una clave externa (`author_id`) en lugar de un campo de texto.
- [x] **CRUD de Autores:** Implementación de los endpoints para crear y leer autores.
- [x] **Refactorización a Patrón Repositorio:**
    - [x] Creación de una capa de abstracción de datos (`BookRepository`, `AuthorRepository`, `UserRepository`).
    - [x] Refactorización de todos los handlers para que dependan de las interfaces de los repositorios en lugar de la base de datos directamente.
- [x] **API de Consulta Avanzada:**
    - [x] Implementación de **paginación** (`limit`, `page`) en el endpoint de lista de libros.
    - [x] Implementación de **filtrado dinámico** (`title`, `author`) en el endpoint de lista de libros.
    - [x] Implementación de **ordenamiento dinámico** (`sort`, `order`) en el endpoint de lista de libros.
    - [x] Enriquecimiento de la respuesta de la lista de libros con **metadatos de paginación** (total_records, total_pages, etc.).
- [x] **Pruebas Unitarias (Testing):**
    - [x] Introducción a las pruebas unitarias con `go-sqlmock`.
    - [x] Creación de una suite de pruebas para los repositorios, cubriendo casos de éxito y de error.

---

## 🚀 Fase 4: Próximos Pasos hacia un Producto Completo

- [x] **Gestión de Inventario y Préstamos (Lógica de Negocio Compleja):**
    - [x] Añadir un campo `stock` a la tabla `books`.
    - [x] Crear una nueva tabla `loans` (préstamos) que relacione un `user_id` con un `book_id`, con fechas de préstamo y devolución.
    - [x] Implementar endpoints para que un usuario pueda "tomar prestado" un libro (crear un `loan` y decrementar el `stock`) y "devolverlo" (eliminar el `loan` y aumentar el `stock`).

- [x] **Roles y Permisos (Autorización Avanzada):**
    - [x] Añadir un campo `role` a la tabla `users` (p. ej., "miembro" y "bibliotecario").
    - [x] Modificar la emisión de tokens JWT para que incluyan el rol del usuario.
    - [x] Crear un nuevo middleware de autorización que restrinja ciertas acciones (como crear autores o añadir stock) solo a los usuarios con el rol de "bibliotecario".

- [x] **Mejora de la Experiencia del Desarrollador (DX):**
    - [x] **Documentación de la API con Swagger/OpenAPI:** Generar documentación interactiva de la API a partir del código para que otros desarrolladores puedan descubrir y probar los endpoints fácilmente.

- [ ] **Optimización y Rendimiento:**
    - [ ] **Resolver el Problema N+1:** Analizar las consultas y optimizar la obtención de datos relacionados para evitar múltiples viajes a la base de datos.
    - [ ] **Caching:** Introducir una capa de caché (p. ej., con Redis) para las consultas frecuentes, como la lista de autores o los detalles de un libro popular.

- [ ] **Cobertura Total de Pruebas:**
    - [ ] Escribir pruebas unitarias para todos los métodos de los repositorios que aún no están cubiertos.
    - [ ] Introducir **Pruebas de Integración**, que prueben el flujo completo desde la petición HTTP hasta la base de datos (usando una base de datos de prueba real).
