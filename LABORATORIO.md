# LABORATORIO: Implementar Validaciones en Clean Architecture

**Duración**: 50 minutos  
**Tipo**: Exploratorio - Aprendizaje práctico (no es necesario entregar código)  
**Pre-lectura recomendada** (para hacer en casa): [docs/10-laboratorio-practico.md](docs/10-laboratorio-practico.md)

---

## 📋 Resumen de la Tarea

Vas a implementar **validaciones de artículos** que fluyan correctamente por todas las capas de Clean Architecture.

**Reglas a implementar:**
- ✅ Título: mínimo 5 caracteres
- ✅ Contenido: mínimo 20 caracteres

**Capas que van a cambiar:**
1. **Domain** - Definir validaciones + error
2. **Service** - Ejecutar validaciones
3. **REST** - Mapear error a HTTP 400

---

## 🎯 Paso 1: Domain Layer (5 minutos)

### 1.1 Agregar método `Validate()` a `Article`

**Archivo**: `domain/article.go`

Busca la struct `Article` y crea un método que valide:

```go
// Validate checks if the article has valid data according to business rules
func (a *Article) Validate() error {
	if len(a.Title) < 5 {
		return ErrValidationFailed // A definir en el paso 1.2
	}
	if len(a.Content) < 20 {
		return ErrValidationFailed
	}
	return nil
}
```

**Hints:**
- El método debe ser un receiver en `Article`
- Usar `len()` para calcular la longitud
- Retornar un error si las reglas no se cumplen

### 1.2 Crear error de validación en `domain/errors.go`

**Archivo**: `domain/errors.go`

Agrega un nuevo error después de los existentes:

```go
// ErrValidationFailed will throw if the article data does not meet business rules
ErrValidationFailed = errors.New("article data validation failed: title must be at least 5 characters and content must be at least 20 characters")
```

**Hints:**
- Sigue el patrón de los otros errores (comienzan con `Err` mayúscula)
- El mensaje debe ser claro para el cliente (será retornado en HTTP)

---

## 🎯 Paso 2: Service Layer (15 minutos)

### 2.1 Ejecutar validación en `article/service.go`

**Archivo**: `article/service.go`

Busca el método `Store()` y agrega la validación **ANTES** de guardar:

```go
func (a *Service) Store(ctx context.Context, ar *domain.Article) error {
	// AGREGAR AQUÍ: Validar el artículo
	if err := ar.Validate(); err != nil {
		return err
	}
	
	// El resto del código continúa igual
	// ... el método continúa con la lógica existente
}
```

**Hints:**
- La validación debe estar al principio del método
- Si retorna error, el método termina sin tocar la BD
- Si no retorna error, continúa normalmente

**Por qué va aquí y no en otra capa?**
- ✅ Service: Si alguien llamara el servicio desde CLI o gRPC, la validación se ejecutaría
- ❌ REST Handler: Si se llamara el servicio directamente, se saltaría la validación
- ❌ Repository: Las reglas de negocio no deben estar en acceso a datos

---

## 🎯 Paso 3: REST API Layer (15 minutos)

### 3.1 Mapear error de validación a HTTP 400

**Archivo**: `internal/rest/article.go`

Busca el método `Store()` en la struct `ArticleHandler`:

```go
func (a *ArticleHandler) Store(c echo.Context) error {
	// Parseo del JSON (ya está)
	var article domain.Article
	if err := c.BindAndValidate(&article); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	// Aquí va TU CÓDIGO: Llamar al servicio
	err := a.Service.Store(c.Request().Context(), &article)
	if err != nil {
		// IMPORTANTE: Mapear errores a HTTP status
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	// Si todo funciona, retorna 201 Created
	return c.JSON(http.StatusCreated, article)
}
```

**Hints:**
- Usa `getStatusCode(err)` que ya existe en el archivo
- La función `getStatusCode()` mapea `ErrValidationFailed` → 400
- Retorna un `ResponseError` con el mensaje de error

**Verifica que `getStatusCode()` maneja tu error:**
- Si llamas `getStatusCode(domain.ErrValidationFailed)`, debe retornar `http.StatusBadRequest` (400)
- Busca la función `getStatusCode()` en el mismo archivo para ver si lo hace

---

## 🎯 Paso 4: Tests (15 minutos)

### 4.1 Escribir tests en `article/service_test.go`

**Archivo**: `article/service_test.go`

Agrégalos junto a los tests existentes:

```go
func TestStoreWithInvalidTitle(t *testing.T) {
	// Arrange: Crear mocks
	mockArticleRepo := new(mocks.ArticleRepository)
	mockAuthorRepo := new(mocks.AuthorRepository)
	
	article := &domain.Article{
		Title:   "Hi",     // Menos de 5 caracteres: INVÁLIDO
		Content: "Este es un contenido válido con más de 20 caracteres",
	}

	svc := NewService(mockArticleRepo, mockAuthorRepo)

	// Act
	err := svc.Store(context.Background(), article)

	// Assert
	assert.Equal(t, domain.ErrValidationFailed, err)
	// Verificar que nunca intentó guardar en BD
	mockArticleRepo.AssertNotCalled(t, "Store")
}

func TestStoreWithInvalidContent(t *testing.T) {
	// Arrange
	mockArticleRepo := new(mocks.ArticleRepository)
	mockAuthorRepo := new(mocks.AuthorRepository)
	
	article := &domain.Article{
		Title:   "Título válido",
		Content: "Corto", // Menos de 20 caracteres: INVÁLIDO
	}

	svc := NewService(mockArticleRepo, mockAuthorRepo)

	// Act
	err := svc.Store(context.Background(), article)

	// Assert
	assert.Equal(t, domain.ErrValidationFailed, err)
	mockArticleRepo.AssertNotCalled(t, "Store")
}

func TestStoreWithValidData(t *testing.T) {
	// Arrange
	mockArticleRepo := new(mocks.ArticleRepository)
	mockAuthorRepo := new(mocks.AuthorRepository)
	
	article := &domain.Article{
		Title:   "Título válido con 5+ caracteres",
		Content: "Este es un contenido válido con más de 20 caracteres",
		Author:  domain.Author{ID: 1},
	}

	// Mock: Esperar que se llame a Store
	mockArticleRepo.On("Store", context.Background(), article).Return(nil)
	// Mock: Esperar que se llame a GetByID para llenar datos del autor
	mockAuthorRepo.On("GetByID", context.Background(), int64(1)).Return(
		domain.Author{ID: 1, Name: "John Doe"},
		nil,
	)

	svc := NewService(mockArticleRepo, mockAuthorRepo)

	// Act
	err := svc.Store(context.Background(), article)

	// Assert
	assert.NoError(t, err)
	mockArticleRepo.AssertCalled(t, "Store", context.Background(), article)
}
```

**Hints:**
- Usa `assert` del package `github.com/stretchr/testify/assert`
- `On()` y `AssertCalled()` ya están configurados (busca ejemplos en el archivo)
- Los mocks están en `article/mocks/` (generados automáticamente)
- Un test debe fallar si intentas guardar sin validar

---

## 🎯 Paso 5: Verificación (5-10 minutos)

### 5.1 Ejecutar tests

```bash
make tests
```

✅ Todos deben pasar sin errores.

### 5.2 Levantar la aplicación

```bash
# Si no está corriendo
make up

# En otra terminal
curl http://localhost:9090/articles
```

### 5.3 Probar casos

**Caso A: Datos inválidos → HTTP 400**

```bash
curl -X POST http://localhost:9090/articles \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Hi",
    "content": "x",
    "author": {"id": 1}
  }'
```

**Resultado esperado:**
```json
{
  "message": "article data validation failed: title must be at least 5 characters..."
}
```
HTTP Status: `400 Bad Request`

**Caso B: Datos válidos → HTTP 201**

```bash
curl -X POST http://localhost:9090/articles \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Un artículo con título válido",
    "content": "Este es un contenido con más de veinte caracteres.",
    "author": {"id": 1}
  }'
```

**Resultado esperado:**
```json
{
  "id": 5,
  "title": "Un artículo con título válido",
  "content": "Este es un contenido con más de veinte caracteres.",
  "author": {"id": 1, "name": "..."},
  "created_at": "...",
  "updated_at": "..."
}
```
HTTP Status: `201 Created`

---

## 🤔 Troubleshooting

| Problema | Solución |
|----------|----------|
| `ErrValidationFailed` no existe | Asegúrate de crear el error en `domain/errors.go` |
| Tests no compilan | ¿Importaste `assert` de `testify`? ¿Los mocks están en `article/mocks/`? |
| HTTP 500 en lugar de 400 | Verifica que `getStatusCode()` maneja correctamente `ErrValidationFailed` |
| `make tests` falla | Verifica que el código compila: `make build` primero |
| La validación no se ejecuta | Asegúrate de llamar `ar.Validate()` en `Service.Store()` ANTES de guardar |

---

## 🎓 Reflexión

**Después de completar, piensa:**

1. ¿Por qué la validación está en Service y no en Handler?
   - Respuesta: Porque puede haber otras formas de acceder al servicio (CLI, gRPC, eventos), y queremos que la validación sea consistente.

2. ¿Qué sucedería si la validación solo estuviera en el Handler?
   - Respuesta: Si alguien llamara `Service.Store()` directamente, se saltaría la validación.

3. ¿Por qué el error se define en Domain?
   - Respuesta: Porque es una regla de negocio. En el futuro podrías cambiar a gRPC y necesitarías un código de error diferente.

---

## ✅ Checklist Final

Cuando termines, verifica:

- [ ] `make build` sin errores
- [ ] `make tests` sin fallos
- [ ] HTTP 400 cuando título < 5 caracteres
- [ ] HTTP 400 cuando contenido < 20 caracteres
- [ ] HTTP 201 cuando datos son válidos
- [ ] Tests nuevos pasan todos
- [ ] Entiendo por qué cada validación está donde está

---

## 📚 Archivos Que Vas a Editar

```
domain/
├── article.go          ← Agregar método Validate()
└── errors.go           ← Agregar ErrValidationFailed

article/
├── service.go          ← Llamar Validate() en Store()
└── service_test.go     ← Agregar 3 tests nuevos

internal/rest/
└── article.go          ← Mapear error a HTTP 400
```

**Total**: 5 archivos, cambios localizados en cada uno.

---

## 🚀 Bono (si terminas antes)

Si terminas con tiempo, intenta:

1. **Agregar validaciones adicionales** - Ej: `ID > 0`, `UpdatedAt <= ahora`
2. **Escribir más tests** - Casos edge (strings con espacios, Unicode, etc.)
3. **Explorar el código existente** - Cómo funciona el mapeo de otros errores (`ErrNotFound`, `ErrConflict`)
4. **Revisar mocks** - Abre `article/mocks/ArticleRepository.go` para ver cómo se generan

---

¡Buena suerte! 🎯
