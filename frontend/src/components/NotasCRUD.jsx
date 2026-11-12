import { useEffect, useState } from 'react';

async function request(url, options) {
  const response = await fetch(url, options);
  if (!response.ok) {
    const body = await response.json().catch(() => ({}));
    throw new Error(body.error || 'No se pudo completar la operación.');
  }
  return response.status === 204 ? null : response.json();
}

function HistorialModal({ nota, onClose, onRestaurada, reportarError, reportarMensaje }) {
  const [historial, setHistorial] = useState([]);
  const [cargando, setCargando] = useState(true);
  const [restaurandoId, setRestaurandoId] = useState(null);

  async function cargarHistorial() {
    setCargando(true);
    try {
      const data = await request(`/api/notas/${nota.id}/historial`);
      setHistorial(data || []);
    } catch (err) {
      reportarError(err.message);
    } finally {
      setCargando(false);
    }
  }

  useEffect(() => {
    if (nota) {
      cargarHistorial();
    }
  }, [nota?.id]);

  async function restaurar(version) {
    const confirmar = window.confirm(
      `¿Deseas restaurar la calificación a ${version.calificacion}? La versión actual quedará respaldada como un nuevo punto en el historial.`
    );
    if (!confirmar) return;

    setRestaurandoId(version.id);
    try {
      const notaActualizada = await request(`/api/notas/${nota.id}/restaurar`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ version_id: version.id }),
      });
      onRestaurada(notaActualizada);
      reportarMensaje(
        `Nota restaurada con éxito a ${notaActualizada.calificacion} (versión de ${version.fecha_modificacion}). El estado previo se guardó en el historial.`
      );
      await cargarHistorial();
    } catch (err) {
      reportarError(err.message);
    } finally {
      setRestaurandoId(null);
    }
  }

  if (!nota) return null;

  return (
    <div className="modal-overlay" onClick={onClose} role="dialog" aria-modal="true" aria-labelledby="historial-title">
      <div className="modal-dialog" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3 id="historial-title">Historial Memento: {nota.estudiante}</h3>
          <button className="modal-close-btn" type="button" onClick={onClose} aria-label="Cerrar modal">&times;</button>
        </div>

        <div className="modal-body">
          <p className="page-description" style={{ margin: 0 }}>
            Patrón Memento: Gestión de snapshots históricos con restauración bidireccional.
          </p>

          <div>
            <div className="memento-history-title">Estado Actual (Originator)</div>
            <div className="current-version-box">
              <div className="memento-meta">
                <span className="student-name">{nota.estudiante}</span>
                <span className="course-name">{nota.curso}</span>
              </div>
              <div style={{ textAlign: 'right' }}>
                <span className="badge-current">Versión actual</span>
                <div className="memento-grade" style={{ marginTop: '4px' }}>
                  {nota.calificacion} <small style={{ color: '#64748b', fontSize: '0.75rem' }}>/ 20</small>
                </div>
              </div>
            </div>
          </div>

          <div>
            <div className="memento-history-title">Versiones Anteriores Guardadas (Mementos)</div>
            {cargando ? (
              <p className="empty-state">Cargando versiones...</p>
            ) : historial.length === 0 ? (
              <p className="empty-state">
                No hay versiones previas. Se generará un snapshot automático cada vez que edites la calificación.
              </p>
            ) : (
              <ul className="memento-list">
                {historial.map((item, index) => (
                  <li key={item.id} className="memento-item">
                    <div className="memento-meta">
                      <span className="memento-date">
                        Modificado: <strong>{item.fecha_modificacion}</strong>
                      </span>
                      <span className="course-name">{item.curso}</span>
                    </div>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '14px' }}>
                      <span className="memento-grade">
                        {item.calificacion} <small style={{ color: '#64748b', fontSize: '0.72rem' }}>/ 20</small>
                      </span>
                      <button
                        className="restore-button"
                        type="button"
                        onClick={() => restaurar(item)}
                        disabled={restaurandoId === item.id}
                      >
                        {restaurandoId === item.id ? 'Restaurando...' : 'Restaurar'}
                      </button>
                    </div>
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>

        <div className="modal-footer">
          <button type="button" onClick={onClose}>Cerrar</button>
        </div>
      </div>
    </div>
  );
}

function NotaFila({ nota, onGuardada, onEliminar, onAbrirHistorial, reportarError, reportarMensaje }) {
  const [editando, setEditando] = useState(false);
  const [calificacionTemp, setCalificacionTemp] = useState(String(nota.calificacion));
  const [guardando, setGuardando] = useState(false);

  useEffect(() => {
    setCalificacionTemp(String(nota.calificacion));
  }, [nota.calificacion]);

  async function guardar() {
    const calificacion = Number(calificacionTemp);
    if (calificacionTemp === '' || !Number.isFinite(calificacion) || calificacion < 0 || calificacion > 20) {
      reportarError('La calificación debe estar entre 0 y 20.');
      return;
    }

    setGuardando(true);
    try {
      const actualizada = await request(`/api/notas/${nota.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          estudiante: nota.estudiante,
          curso: nota.curso,
          calificacion,
        }),
      });
      onGuardada(actualizada);
      setEditando(false);
      reportarError('');
      reportarMensaje(`Nota de ${nota.estudiante} actualizada a ${actualizada.calificacion}. Versión anterior respaldada en el historial.`);
    } catch (err) {
      reportarError(err.message);
    } finally {
      setGuardando(false);
    }
  }

  return (
    <li className="note-item">
      <div className="note-details">
        <span className="student-name">{nota.estudiante}</span>
        <span className="course-name">{nota.curso}</span>
      </div>
      {editando ? (
        <div className="note-actions note-edit-form">
          <label>
            Nueva calificación
            <input
              type="number"
              min="0"
              max="20"
              step="0.01"
              value={calificacionTemp}
              onChange={(event) => setCalificacionTemp(event.target.value)}
              autoFocus
            />
          </label>{' '}
          <button className="primary-button" type="button" onClick={guardar} disabled={guardando}>Guardar</button>{' '}
          <button type="button" onClick={() => { setCalificacionTemp(String(nota.calificacion)); setEditando(false); }} disabled={guardando}>
            Cancelar
          </button>
        </div>
      ) : (
        <>
          <span className="grade-value">{nota.calificacion}<small> / 20</small></span>
          <div className="note-actions">
            <button type="button" onClick={() => setEditando(true)}>Editar</button>{' '}
            <button className="history-button" type="button" onClick={() => onAbrirHistorial(nota)}>Historial</button>{' '}
            <button className="danger-button" type="button" onClick={() => onEliminar(nota.id)}>Eliminar</button>
          </div>
        </>
      )}
    </li>
  );
}

export default function NotasCRUD() {
  const [notas, setNotas] = useState([]);
  const [matriculas, setMatriculas] = useState([]);
  const [clases, setClases] = useState([]);
  const [formulario, setFormulario] = useState({ estudiante: '', curso: '', calificacion: '' });
  const [error, setError] = useState('');
  const [mensaje, setMensaje] = useState('');
  const [notaHistorial, setNotaHistorial] = useState(null);

  const matriculasAceptadas = matriculas.filter((item) => ['Activa', 'Finalizada'].includes(item.estado));

  async function cargar() {
    try {
      const [listaNotas, listaMatriculas, listaClases] = await Promise.all([
        request('/api/notas'),
        request('/api/matriculas'),
        request('/api/clases'),
      ]);
      setNotas(listaNotas);
      setMatriculas(listaMatriculas);
      setClases(listaClases);
      setError('');
    } catch (err) {
      setError(err.message);
    }
  }

  useEffect(() => {
    cargar();
  }, []);

  function actualizarCampo(event) {
    setFormulario({ ...formulario, [event.target.name]: event.target.value });
  }

  async function crear(event) {
    event.preventDefault();
    const nombreEstudiante = formulario.estudiante.trim();
    const tieneMatriculaValida = matriculasAceptadas.some((item) => item.estudiante.trim() === nombreEstudiante);
    if (!tieneMatriculaValida) {
      setError('Solo se pueden registrar notas para estudiantes con matrícula activa o finalizada.');
      setMensaje('');
      return;
    }

    try {
      const nueva = await request('/api/notas', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ ...formulario, estudiante: nombreEstudiante, calificacion: Number(formulario.calificacion) }),
      });
      setNotas((actuales) => [...actuales, nueva]);
      setFormulario({ estudiante: '', curso: '', calificacion: '' });
      setError('');
      setMensaje(`Nota creada exitosamente con calificación ${nueva.calificacion}.`);
    } catch (err) {
      setError(err.message);
      setMensaje('');
    }
  }

  async function eliminar(id) {
    try {
      await request(`/api/notas/${id}`, { method: 'DELETE' });
      setNotas((actuales) => actuales.filter((nota) => nota.id !== id));
      if (notaHistorial?.id === id) {
        setNotaHistorial(null);
      }
      setError('');
      setMensaje('Nota eliminada correctamente.');
    } catch (err) {
      setError(err.message);
      setMensaje('');
    }
  }

  function actualizarNota(actualizada) {
    setNotas((actuales) => actuales.map((nota) => (nota.id === actualizada.id ? actualizada : nota)));
    if (notaHistorial?.id === actualizada.id) {
      setNotaHistorial(actualizada);
    }
  }

  return (
    <section className="panel">
      <div className="panel-heading">
        <div>
          <p className="eyebrow">RENDIMIENTO</p>
          <h2>Notas</h2>
        </div>
        <span className="pattern-tag">Memento</span>
      </div>
      <p className="panel-description">
        Registra calificaciones con respaldo automático de versiones (Originator, Memento, Caretaker) y restauración bidireccional.
      </p>

      <form className="entry-form grades-form" onSubmit={crear}>
        <label>
          Estudiante
          <select
            name="estudiante"
            value={formulario.estudiante}
            onChange={actualizarCampo}
            required
            disabled={matriculasAceptadas.length === 0}
          >
            <option value="">{matriculasAceptadas.length === 0 ? 'No hay estudiantes aceptados' : 'Selecciona un estudiante'}</option>
            {matriculasAceptadas.map((matricula) => (
              <option key={matricula.id} value={matricula.estudiante}>{matricula.estudiante}</option>
            ))}
          </select>
        </label>{' '}
        <label>
          Curso
          <select
            name="curso"
            value={formulario.curso}
            onChange={actualizarCampo}
            required
            disabled={clases.length === 0}
          >
            <option value="">{clases.length === 0 ? 'No hay cursos disponibles' : 'Selecciona un curso'}</option>
            {clases.map((curso) => (
              <option key={curso.id} value={curso.nombre}>{curso.nombre}</option>
            ))}
          </select>
        </label>{' '}
        <label>
          Calificación
          <input
            name="calificacion"
            type="number"
            min="0"
            max="20"
            step="0.01"
            value={formulario.calificacion}
            onChange={actualizarCampo}
            required
            placeholder="0 - 20"
          />
        </label>{' '}
        <button className="primary-button" type="submit" disabled={matriculasAceptadas.length === 0 || clases.length === 0}>Crear nota</button>
      </form>

      {mensaje && <p className="success-banner" role="status">{mensaje}</p>}
      {error && <p className="error-message" role="alert">{error}</p>}

      <ul className="notes-list">
        {notas.map((nota) => (
          <NotaFila
            key={nota.id}
            nota={nota}
            onGuardada={actualizarNota}
            onEliminar={eliminar}
            onAbrirHistorial={(n) => setNotaHistorial(n)}
            reportarError={(msg) => { setError(msg); setMensaje(''); }}
            reportarMensaje={(msg) => { setMensaje(msg); setError(''); }}
          />
        ))}
      </ul>
      {notas.length === 0 && <p className="empty-state">Aún no hay calificaciones registradas.</p>}

      {notaHistorial && (
        <HistorialModal
          nota={notaHistorial}
          onClose={() => setNotaHistorial(null)}
          onRestaurada={actualizarNota}
          reportarError={(msg) => { setError(msg); setMensaje(''); }}
          reportarMensaje={(msg) => { setMensaje(msg); setError(''); }}
        />
      )}
    </section>
  );
}
