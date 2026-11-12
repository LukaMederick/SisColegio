import { useEffect, useState } from 'react';
import { obtenerEstadoMatricula } from '../patterns/MatriculaState.js';

async function request(url, options) {
  const response = await fetch(url, options);
  if (!response.ok) {
    const body = await response.json().catch(() => ({}));
    throw new Error(body.error || 'No se pudo completar la operación.');
  }
  return response.status === 204 ? null : response.json();
}

export default function MatriculasCRUD() {
  const [matriculas, setMatriculas] = useState([]);
  const [estudiante, setEstudiante] = useState('');
  const [apellido, setApellido] = useState('');
  const [gradoCurso, setGradoCurso] = useState('');
  const [seccion, setSeccion] = useState('');
  const [fechaInscripcion, setFechaInscripcion] = useState('');
  const [tutorResponsable, setTutorResponsable] = useState('');
  const [estadoInicial, setEstadoInicial] = useState('Activa');
  const [edicion, setEdicion] = useState(null);
  const [error, setError] = useState('');
  const [mensaje, setMensaje] = useState('');

  async function cargar() {
    try {
      setMatriculas(await request('/api/matriculas'));
      setError('');
    } catch (err) {
      setError(err.message);
    }
  }

  useEffect(() => {
    cargar();
  }, []);

  async function crear(event) {
    event.preventDefault();
    try {
      const nueva = await request('/api/matriculas', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          estudiante,
          apellido,
          grado_curso: gradoCurso,
          seccion,
          fecha_inscripcion: fechaInscripcion,
          tutor_responsable: tutorResponsable,
          estado: estadoInicial,
        }),
      });
      setMatriculas((actuales) => [...actuales, nueva]);
      setEstudiante('');
      setApellido('');
      setGradoCurso('');
      setSeccion('');
      setFechaInscripcion('');
      setTutorResponsable('');
      setError('');
      setMensaje(`Matrícula creada correctamente con estado ${nueva.estado}.`);
    } catch (err) {
      setError(err.message);
      setMensaje('');
    }
  }

  async function cambiarEstado(matricula, siguienteEstado) {
    const estadoActual = obtenerEstadoMatricula(matricula.estado);
    if (!estadoActual.permite(siguienteEstado)) {
      setError(`Transición no permitida desde ${matricula.estado} hacia ${siguienteEstado}.`);
      return;
    }

    try {
      const res = await request(`/api/matriculas/${matricula.id}/estado`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ estado: siguienteEstado }),
      });
      setMatriculas((actuales) => actuales.map((item) => (
        item.id === matricula.id ? { ...item, estado: siguienteEstado } : item
      )));
      setError('');
      setMensaje(res?.mensaje ? `Estado de ${matricula.estudiante} actualizado a ${siguienteEstado}.` : 'Estado actualizado.');
    } catch (err) {
      setError(err.message);
      setMensaje('');
    }
  }

  async function guardarEdicion(event) {
    event.preventDefault();
    try {
      const actualizada = await request(`/api/matriculas/${edicion.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          estudiante: edicion.estudiante,
          apellido: edicion.apellido,
          grado_curso: edicion.grado_curso,
          seccion: edicion.seccion,
          fecha_inscripcion: edicion.fecha_inscripcion,
          tutor_responsable: edicion.tutor_responsable,
        }),
      });
      setMatriculas((actuales) => actuales.map((item) => (
        item.id === actualizada.id ? actualizada : item
      )));
      setEdicion(null);
      setError('');
      setMensaje('Datos del estudiante actualizados.');
    } catch (err) {
      setError(err.message);
      setMensaje('');
    }
  }

  async function eliminar(id) {
    try {
      await request(`/api/matriculas/${id}`, { method: 'DELETE' });
      setMatriculas((actuales) => actuales.filter((item) => item.id !== id));
      setError('');
      setMensaje('Matrícula eliminada.');
    } catch (err) {
      setError(err.message);
      setMensaje('');
    }
  }

  return (
    <section className="panel">
      <div className="panel-heading">
        <div>
          <p className="eyebrow">INSCRIPCIONES</p>
          <h2>Matrículas</h2>
        </div>
        <span className="pattern-tag">State</span>
      </div>
      <p className="panel-description">Gestiona matrículas controlando transiciones estrictas de estado (Activa, Finalizada, Anulada).</p>

      <form className="entry-form" onSubmit={crear}>
        <label>
          Nombre del alumno
          <input
            value={estudiante}
            onChange={(event) => setEstudiante(event.target.value)}
            required
            placeholder="Nombre del estudiante"
          />
        </label>
        <label>
          Apellido
          <input
            value={apellido}
            onChange={(event) => setApellido(event.target.value)}
            required
            placeholder="Apellido"
          />
        </label>
        <label>
          Grado / curso
          <input
            value={gradoCurso}
            onChange={(event) => setGradoCurso(event.target.value)}
            required
            placeholder="Ej. 3° primaria"
          />
        </label>
        <label>
          Sección
          <input
            value={seccion}
            onChange={(event) => setSeccion(event.target.value)}
            required
            placeholder="A"
          />
        </label>
        <label>
          Fecha de inscripción
          <input
            type="date"
            value={fechaInscripcion}
            onChange={(event) => setFechaInscripcion(event.target.value)}
            required
          />
        </label>
        <label>
          Tutor o responsable
          <input
            value={tutorResponsable}
            onChange={(event) => setTutorResponsable(event.target.value)}
            required
            placeholder="Nombre del tutor"
          />
        </label>
        <label>
          Estado inicial
          <select value={estadoInicial} onChange={(event) => setEstadoInicial(event.target.value)}>
            <option value="Activa">Activa</option>
            <option value="Pendiente">Pendiente</option>
          </select>
        </label>
        <button className="primary-button" type="submit">Registrar</button>
      </form>

      {mensaje && <p className="success-banner" role="status">{mensaje}</p>}
      {error && <p className="error-message" role="alert">{error}</p>}

      <div className="table-scroll">
        <table>
          <thead>
            <tr>
              <th>Nombre del alumno</th>
              <th>Apellido</th>
              <th>Grado / curso</th>
              <th>Sección</th>
              <th>Fecha de inscripción</th>
              <th>Tutor o responsable</th>
              <th>Estado</th>
              <th>Acciones</th>
            </tr>
          </thead>
          <tbody>
            {matriculas.map((matricula) => {
              const estado = obtenerEstadoMatricula(matricula.estado);
              return (
                <tr key={matricula.id}>
                  <td className="student-name">
                    {edicion?.id === matricula.id ? (
                      <form className="inline-edit-form" onSubmit={guardarEdicion}>
                        <input
                          aria-label="Nombre del alumno"
                          value={edicion.estudiante}
                          onChange={(event) => setEdicion({ ...edicion, estudiante: event.target.value })}
                          required
                        />
                        <button className="primary-button" type="submit">Guardar</button>
                        <button type="button" onClick={() => setEdicion(null)}>Cancelar</button>
                      </form>
                    ) : (matricula.estudiante || '—')}
                  </td>
                  <td>
                    {edicion?.id === matricula.id ? (
                      <input
                        aria-label="Apellido"
                        value={edicion.apellido || ''}
                        onChange={(event) => setEdicion({ ...edicion, apellido: event.target.value })}
                        required
                      />
                    ) : (matricula.apellido || '—')}
                  </td>
                  <td>
                    {edicion?.id === matricula.id ? (
                      <input
                        aria-label="Grado / curso"
                        value={edicion.grado_curso || ''}
                        onChange={(event) => setEdicion({ ...edicion, grado_curso: event.target.value })}
                        required
                      />
                    ) : (matricula.grado_curso || '—')}
                  </td>
                  <td>
                    {edicion?.id === matricula.id ? (
                      <input
                        aria-label="Sección"
                        value={edicion.seccion || ''}
                        onChange={(event) => setEdicion({ ...edicion, seccion: event.target.value })}
                        required
                      />
                    ) : (matricula.seccion || '—')}
                  </td>
                  <td>
                    {edicion?.id === matricula.id ? (
                      <input
                        type="date"
                        aria-label="Fecha de inscripción"
                        value={edicion.fecha_inscripcion || ''}
                        onChange={(event) => setEdicion({ ...edicion, fecha_inscripcion: event.target.value })}
                        required
                      />
                    ) : (matricula.fecha_inscripcion || '—')}
                  </td>
                  <td>
                    {edicion?.id === matricula.id ? (
                      <input
                        aria-label="Tutor o responsable"
                        value={edicion.tutor_responsable || ''}
                        onChange={(event) => setEdicion({ ...edicion, tutor_responsable: event.target.value })}
                        required
                      />
                    ) : (matricula.tutor_responsable || '—')}
                  </td>
                  <td><span className={`status-badge status-${matricula.estado.toLowerCase()}`}>{matricula.estado}</span></td>
                  <td className="action-cell">
                    {edicion?.id !== matricula.id && (
                      <>
                        <button type="button" onClick={() => setEdicion({
                          id: matricula.id,
                          estudiante: matricula.estudiante,
                          apellido: matricula.apellido || '',
                          grado_curso: matricula.grado_curso || '',
                          seccion: matricula.seccion || '',
                          fecha_inscripcion: matricula.fecha_inscripcion || '',
                          tutor_responsable: matricula.tutor_responsable || '',
                        })}>Editar</button>
                        <button className="danger-button" type="button" onClick={() => eliminar(matricula.id)}>Eliminar</button>
                        {estado.acciones().map((accion) => (
                          <button
                            key={accion}
                            type="button"
                            onClick={() => cambiarEstado(matricula, accion)}
                          >
                            Pasar a {accion}
                          </button>
                        ))}
                      </>
                    )}
                    {edicion?.id !== matricula.id && estado.acciones().length === 0 && <span className="terminal-state">Sin transiciones (Estado terminal)</span>}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
        {matriculas.length === 0 && <p className="empty-state">Aún no hay matrículas registradas.</p>}
      </div>
    </section>
  );
}
