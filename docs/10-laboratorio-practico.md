# Laboratorio Práctico: Clean Architecture en Go (1 hora)

## 🎯 Objetivos del Laboratorio

Al completar este laboratorio, habrás aprendido a:

1. **Identificar las capas de Clean Architecture** en un proyecto real
2. **Entender cómo fluye una solicitud** a través de las diferentes capas
3. **Implementar validaciones** que se propaguen correctamente por toda la arquitectura
4. **Mapear errores de dominio** a respuestas HTTP apropiadas
5. **Escribir tests** que validen el comportamiento de cada capa

---

## ⚡ Quick Start (5 minutos)

### Pre-requisitos
- Docker instalado
- Make instalado
- Git instalado
- Editor de código (VS Code recomendado)

### Levanta el proyecto

```bash
# 1. Clonar o descargar el proyecto
cd go-clean-arch

# 2. Copiar las variables de entorno
cp example.env .env

# 3. Levantar la base de datos y la aplicación
make dev-env  # Inicia MySQL en Docker (solo primera vez)
make up       # Inicia la app con hot-reload

# 4. En otra terminal, verifica que funciona
curl http://localhost:9090/articles
```

Si ves un JSON con artículos, ¡está funcionando! 🎉

---

## 📚 Conceptos Clave (10 minutos de presentación)

### ¿Qué es Clean Architecture?

Es un patrón de organización de código que **separa responsabilidades en capas**:

```
┌─────────────────────────────────┐
│     REST API (Echo)             │  ← Capa de Presentación
│    - Rutas                      │    (Cómo se presenta la info)
│    - Handlers                   │
│    - Mapeo a HTTP               │
└──────────────┬──────────────────┘
               │
┌──────────────▼──────────────────┐
│   Service (Casos de Uso)        │  ← Capa de Lógica de Negocio
│    - Orquestar operaciones      │    (Qué hacer con la info)
│    - Validar reglas             │
│    - Llamar a repositorios      │
└──────────────┬──────────────────┘
               │
┌──────────────▼──────────────────┐
│   Repository + MySQL            │  ← Capa de Datos
│    - Acceso a BD                │    (Cómo almacenar la info)
│    - Queries SQL                │
│    - Mapeo de datos             │
└──────────────┬──────────────────┘
               │
┌──────────────▼──────────────────┐
│   Domain (Entities + Errors)    │  ← Centro: Las reglas
│    - Estructura de datos        │    (Qué es un Artículo)
│    - Errores del negocio        │
└─────────────────────────────────┘
```

### Las 3 Capas Visibles en Este Proyecto

#### 1️⃣ **Domain** (`domain/`)
- **Responsabilidad**: Definir las *entidades* y *reglas de negocio*
- **Lo que vive aquí**:
  - Estructuras de datos: `Article`, `Author`
  - Errores del dominio: `ErrNotFound`, `ErrConflict`
  - **NO código de HTTP, NO SQL, NO detalles técnicos**
- **Por qué**: Las reglas del negocio son independientes de cómo se accede o presenta la info

#### 2️⃣ **Service** (`article/service.go`)
- **Responsabilidad**: *Orquestar* la lógica de negocio (casos de uso)
- **Lo que vive aquí**:
  - Métodos como `Store()`, `Delete()`, `Fetch()`
  - Validaciones y reglas
  - Llamadas a repositorios
  - **NO código de HTTP, NO acceso directo a BD**
- **Por qué**: Centraliza la lógica común; fácil testear sin BD

#### 3️⃣ **REST API** (`internal/rest/article.go`)
- **Responsabilidad**: *Exponer* el servicio como API HTTP
- **Lo que vive aquí**:
  - Handlers: `FetchArticle()`, `Store()`, `GetByID()`, `Delete()`
  - Mapeo: JSON ↔ Go structs, errores → códigos HTTP
  - **NO lógica de negocio directa**
- **Por qué**: El servicio podría exponerse vía gRPC, CLI, o eventos - solo cambias esta capa

### Ejemplo: Cómo fluye un request POST /articles

```
1. Cliente hace:
   curl -X POST http://localhost:9090/articles \
     -H "Content-Type: application/json" \
     -d '{"title":"Mi Artículo","content":"Contenido..."}'

2. REST Handler (internal/rest/article.go)
   ├─ Parsea el JSON
   ├─ Valida sintaxis básica
   └─ Llama Service.Store(ctx, article)

3. Service (article/service.go)
   ├─ Ejecuta validaciones de negocio (título no vacío, etc.)
   ├─ Verifica conflictos (¿ya existe?)
   └─ Llama Repository.Store(ctx, article)

4. Repository (internal/repository/mysql/article.go)
   ├─ Construye query SQL
   ├─ Ejecuta INSERT en MySQL
   └─ Retorna error o ID nuevo

5. Service retorna al Handler
   ├─ Si error: Handler mapea a HTTP (400, 404, 500, etc.)
   └─ Si éxito: Handler retorna 201 + JSON del artículo

6. Cliente recibe:
   HTTP/1.1 201 Created
   {"id":5,"title":"Mi Artículo",...}
```

### ¿Por qué separar en capas?

| Beneficio | Ejemplo |
|-----------|---------|
| **Testeable** | Tests del Service sin BD real (usa mocks) |
| **Reutilizable** | Service usable desde CLI, gRPC, API REST |
| **Mantenible** | Cambio en BD: solo modifica Repository |
| **Escalable** | Agregar caché: agrega lógica en Service, no en REST |

---

## 🧠 La Tarea (50 minutos prácticos)

### Descripción

Implementarás **validaciones de artículos** que se propaguen correctamente por todas las capas.

**Reglas de negocio a implementar:**
- Título debe tener mínimo 5 caracteres
- Contenido debe tener mínimo 20 caracteres

**Resultado esperado:**
- Si POST /articles con datos inválidos → HTTP 400 con mensaje claro
- Si POST /articles con datos válidos → HTTP 201 con artículo creado
- Tests que validen ambos casos

### Archivos que vas a tocar

Estos 5 archivos son los únicos que cambiarás:

1. **`domain/article.go`** - Definir la validación
2. **`domain/errors.go`** - Crear error de validación
3. **`article/service.go`** - Ejecutar validación
4. **`internal/rest/article.go`** - Mapear error a HTTP 400
5. **`article/service_test.go`** - Escribir tests

### Distribución de tiempo

| Paso | Tiempo | Qué Hacer |
|------|--------|-----------|
| 1. Definir validaciones en Domain | 5 min | Código en `domain/article.go` + error en `domain/errors.go` |
| 2. Ejecutar en Service | 15 min | Llamar validación en `article/service.go` método `Store()` |
| 3. Mapear a HTTP | 15 min | Capturar error en `internal/rest/article.go` método `Store()` |
| 4. Escribir tests | 15 min | Agregar casos en `article/service_test.go` |
| 5. Verificar | 5-10 min | `make tests` + `curl` con datos inválidos |

---

## 💡 Hints por Paso

### Paso 1: Definir la validación

**Pista 1a - Domain Layer (`domain/article.go`)**
- Busca dónde está la struct `Article`
- Crea un método `Validate()` que retorne `error`
- Este método debe chequear longitud de título y contenido
- ¿Cómo se calcula longitud en Go? (hint: `len(a.Title)`)

**Pista 1b - Domain Errors (`domain/errors.go`)**
- Observa cómo están definidos los otros errores: `ErrNotFound`, `ErrConflict`
- Crea uno nuevo: `ErrValidationFailed` o similar
- Piensa: ¿qué mensaje debe llevar para que el cliente entienda?

### Paso 2: Ejecutar la validación en Service

**Pista 2a - Service (`article/service.go`)**
- Busca el método `Store()` 
- Observa cómo recibe un `*domain.Article`
- **Antes** de llamar `a.articleRepo.Store()`, llama a `ar.Validate()`
- Si error: retorna el error
- Si OK: continúa normalmente

### Paso 3: Mapear a HTTP

**Pista 3a - REST Handler (`internal/rest/article.go`)**
- Busca el método `Store()` en `ArticleHandler`
- Observa cómo llama `a.Service.Store(ctx, ...)`
- Captura el error retornado
- Usa la función `getStatusCode(err)` (ya existe) para mapear a HTTP status
- Construye un `ResponseError` con el mensaje
- Retorna `c.JSON(status, ResponseError{...})`
- **Prueba**: `curl -X POST http://localhost:9090/articles -H "Content-Type: application/json" -d '{"title":"Hi","content":"x"}'`
  - Debe retornar 400 con mensaje de validación

### Paso 4: Escribir Tests

**Pista 4a - Service Tests (`article/service_test.go`)**
- Abre el archivo y observa la estructura de tests existentes
- Crea 2 tests nuevos:
  - `TestStoreInvalidTitle` - título demasiado corto
  - `TestStoreInvalidContent` - contenido demasiado corto
- Cada test debe:
  1. Crear un mock `ArticleRepository`
  2. Crear el servicio: `svc := NewService(mockRepo, mockAuthorRepo)`
  3. Llamar `svc.Store(ctx, invalidArticle)`
  4. Assert que retorna `ErrValidationFailed` (o el error que creaste)
  5. Assert que **no** llamó a `mockRepo.Store()` (nunca llegó a BD)

**Pista 4b - Patrón Mockery**
- El proyecto usa `mockery` para generar mocks automáticamente
- Los mocks ya están listos en `article/mocks/`
- Observa cómo los tests existentes usan: `mockRepo.On("Store", ...).Return(...)`

### Paso 5: Verificar todo funciona

```bash
# Ejecutar tests
make tests

# Levantar el proyecto (si no está)
make up

# Probar con curl - CASO INVÁLIDO
curl -X POST http://localhost:9090/articles \
  -H "Content-Type: application/json" \
  -d '{"title":"Hi","content":"x","author":{"id":1}}'
# Esperado: 400 Bad Request con mensaje de validación

# Probar con curl - CASO VÁLIDO
curl -X POST http://localhost:9090/articles \
  -H "Content-Type: application/json" \
  -d '{"title":"Título válido con 5+ chars","content":"Este es un contenido con más de 20 caracteres","author":{"id":1}}'
# Esperado: 201 Created con el artículo creado
```

---

## 🤔 Si te quedas atascado

1. **Paso 1 - Validaciones**: ¿El método compila? ¿Retorna un error?
2. **Paso 2 - Service**: ¿Llamas al `Validate()` ANTES de guardar?
3. **Paso 3 - HTTP**: ¿La función `getStatusCode()` existe y maneja tu error?
4. **Paso 4 - Tests**: ¿Copiaste la estructura de un test existente?
5. **Paso 5 - Verificar**: ¿Ejecutaste `make tests` sin errores?

**Si aún no funciona**: Mira la rama `solution/validations` para ver una solución de referencia.

---

## ✅ Checklist Final

Cuando termines, verifica que:

- [ ] Código compila sin errores: `make build`
- [ ] Todos los tests pasan: `make tests`
- [ ] `curl` con datos inválidos retorna HTTP 400
- [ ] `curl` con datos válidos retorna HTTP 201
- [ ] Entiende cómo fluye el request por las 3 capas
- [ ] Puede explicar por qué la validación va en Service y no en Handler o Repository

---

## 🎓 Reflexión Final

Piensa en estas preguntas:

1. **¿Qué pasaría si la validación solo estuviera en el Handler (REST)?**
   - Respuesta: Si alguien llamara al Service directamente (ej: desde una CLI), la validación no se ejecutaría.

2. **¿Qué pasaría si la validación estuviera en el Repository?**
   - Respuesta: Sería una regla de datos, no de negocio. Además, queremos validar ANTES de tocar la BD.

3. **¿Por qué el Domain tiene los errores y no el Handler?**
   - Respuesta: Porque los errores son reglas del negocio, no detalles de HTTP. Mañana podrías usar gRPC y necesitarías distintos códigos de error.

¡Felicidades por completar el laboratorio! 🚀

---

## 📖 Referencias

- [Documentación: Conceptos de Clean Architecture](01-clean-architecture-conceptos.md)
- [Documentación: Mapa del Proyecto](02-mapa-del-proyecto.md)
- [Testing y Calidad](06-testing-y-calidad.md)
- [Troubleshooting](07-troubleshooting.md)
