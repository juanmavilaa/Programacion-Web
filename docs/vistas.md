# Vistas

La aplicación se renderiza en el servidor con templ y usa HTMX para actualizar solo una parte de la página cuando se crea o se borra algo. No hay JavaScript propio ni respuestas JSON: el servidor siempre devuelve HTML, ya sea la página entera (`GET /`) o un fragmento.

## Componentes

Los componentes están en la carpeta `views/`:

- `layout.templ`: estructura base de la página. Incluye Pico.css, `static/styles.css` y HTMX.
- `entity_form.templ`: formularios para crear un cliente y un entrenamiento.
- `entity_list.templ`: tablas con los clientes y los entrenamientos, con el botón Borrar en cada fila.
- `fragments.templ`: los fragmentos que se devuelven después de crear algo y el mensaje de error.
- `index.templ`: página principal (`IndexPage`), que usa el layout y junta los formularios y las listas.

Los archivos `.templ` se compilan a código Go con `templ generate`.

## HTMX

HTMX se carga en el `<head>` del layout desde el CDN, con la versión 4.0.0 que se usó en clase.

### Alta

Los formularios tienen `hx-post` con la misma ruta de siempre (`/clients` y `/workouts`), `hx-target` con el id del contenedor de la lista (`#client-list` y `#workout-list`) y `hx-swap="outerHTML"`. El contenedor es el `<div>` que envuelve la tabla en `ClientList` y `WorkoutList`, así que cuando llega la respuesta se reemplaza entero por la lista nueva.

El handler ya no redirige. Después de insertar vuelve a consultar la lista y renderiza solo el fragmento (`ClientCreated` o `WorkoutCreated`). Además de la lista, el fragmento trae otros `<div>` con `id` y `hx-swap-oob="true"`, que HTMX reemplaza en la página por id aunque no sean el target:

- el `#client-form`, para que el formulario quede vacío;
- el `#workout-form`, para que el `<select>` tenga los clientes actualizados (si no había ningún cliente, pasa del aviso al formulario);
- el `#mensaje`, para borrar un error anterior.

Esos tres `<div>` con id también están en `index.templ`, envolviendo los formularios y el lugar donde se muestran los errores. Así HTMX encuentra de dónde sacar el contenido viejo.

### Borrado

Cada fila tiene un botón con `hx-delete="/clients/{id}"` (o `/workouts/{id}`), `hx-target="closest tr"`, `hx-swap="outerHTML"` y `hx-confirm` para que el navegador pregunte antes de borrar.

El handler lee el id de la URL, llama a `DeleteClient` o `DeleteWorkout` y responde 200 sin cuerpo. Como no hay nada para poner en lugar de la fila, `outerHTML` la saca de la página. Se responde 200 y no 204 porque HTMX no hace ningún reemplazo cuando recibe un 204.

### Errores

Las validaciones son las mismas de siempre. Si algo no es válido el handler responde 422 (datos incorrectos) o 409 (conflicto, como un email repetido o un cliente con entrenamientos), con el header `HX-Reswap: none` y un cuerpo que es solo el `#mensaje` como contenido OOB. HTMX 4 también reemplaza el contenido cuando la respuesta es un 4xx, por eso hace falta el header: se muestra el mensaje y la lista queda como estaba.

La función que arma esa respuesta es `mostrarError`, en `handlers/helpers.go`.

## Handlers

Los handlers están en la carpeta `handlers/`:

- `index.go`: atiende `GET /`. Busca los clientes y los entrenamientos con sqlc y renderiza `IndexPage`.
- `clients.go`: atiende `POST /clients` y `DELETE /clients/{id}`.
- `workouts.go`: atiende `POST /workouts` y `DELETE /workouts/{id}`.
- `helpers.go`: `mostrarError`.

Para las rutas con id, el handler separa la URL con `strings.Split` y convierte el id con `strconv.ParseInt`. Si la URL no tiene la forma `/clients/{id}` o el id no es un número positivo que entre en un `int32`, responde 400. Cualquier otro método sobre esas rutas da 405.

Las validaciones son estas: el nombre no puede estar vacío, el email tiene que ser válido y no repetirse, el cliente de un entrenamiento tiene que existir, las series y repeticiones tienen que ser mayores a 0 y el peso tiene que ser un número entre 0 y 9999.99. Un cliente que tiene entrenamientos no se puede borrar.

Borrar algo que ya no existe no da error: la fila se saca igual, porque en la base ya no está.

## Qué cambió respecto al TP5

En el TP5 los formularios hacían un POST, el servidor redirigía a `/` (Post/Redirect/Get) y el navegador volvía a cargar la página entera. El borrado se hacía con un POST a `/clients/delete` y `/workouts/delete` porque un formulario HTML no puede mandar DELETE. Ahora esas rutas no existen y el error ya no viaja en `/?error=`.

## Limitaciones

Como el borrado responde vacío, después de borrar un cliente el `<select>` de entrenamientos sigue mostrándolo hasta que se recarga la página. Si se intenta cargar un entrenamiento para ese cliente el servidor lo rechaza con el mensaje "no existe el cliente". Por la misma razón, un mensaje de error queda en pantalla hasta el próximo alta o hasta recargar.

## Generación del código

`make generate` ejecuta `sqlc generate` y `templ generate`.

templ se ejecuta con `go run` en la misma versión que figura en `go.mod` (v0.3.1020). Hace falta porque el código que genera cada versión de templ solo funciona con el runtime de esa misma versión. Si se cambia la versión en `go.mod`, hay que cambiarla también en el `Makefile`.
