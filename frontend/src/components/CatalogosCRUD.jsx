import { useEffect, useState } from 'react';

const catalogos = {
  docentes: {
    endpoint: '/api/docentes',
    title: 'Docentes',
    eyebrow: 'EQUIPO ACADÉMICO',
    description: 'Administra los docentes que estarán a cargo de las clases.',
    fields: [
      { name: 'nombre', label: 'Nombre', placeholder: 'Nombre completo' },
      { name: 'especialidad', label: 'Especialidad', placeholder: 'Área o especialidad' },
      { name: 'correo', label: 'Correo', type: 'email', placeholder: 'docente@colegio.edu' },
    ],
    columns: [
      { key: 'nombre', label: 'Docente', className: 'student-name' },
      { key: 'especialidad', label: 'Especialidad' },
      { key: 'correo', label: 'Correo' },
    ],
    emptyMessage: 'Aún no hay docentes registrados.',
  },
  clases: {
    endpoint: '/api/clases',
    title: 'Clases',
    eyebrow: 'ORGANIZACIÓN',
    description: 'Asigna cada clase a un docente para habilitar el registro de asistencia.',
    dependencies: ['/api/docentes'],
    fields: [
      { name: 'nombre', label: 'Nombre de la clase', placeholder: 'Matemática' },
      { name: 'horario', label: 'Horario', placeholder: 'Lun. y Mié. · 08:00' },
      {
        name: 'docenteId',
        label: 'Docente',
        type: 'select',
        options: (resources) => resources.docentes.map((item) => ({
          value: item.id,
          label: `${item.nombre} · ${item.especialidad}`,
        })),
      },
    ],
    numericFields: ['docenteId'],
    columns: [
      { key: 'nombre', label: 'Clase', className: 'student-name' },
      { key: 'horario', label: 'Horario' },
      { key: 'docenteNombre', label: 'Docente' },
    ],
    emptyMessage: 'Aún no hay clases registradas.',
  },
  asistencias: {
    endpoint: '/api/asistencias',
    title: 'Asistencias',
    eyebrow: 'SEGUIMIENTO',
    description: 'Registra la asistencia de cada estudiante matriculado por clase y fecha.',
    dependencies: ['/api/matriculas', '/api/clases'],
    fields: [
      {
        name: 'matriculaId',
        label: 'Estudiante matriculado',
        type: 'select',
        options: (resources) => resources.matriculas
          .filter((item) => ['Activa', 'Finalizada'].includes(item.estado))
          .map((item) => ({
            value: item.id,
            label: `${item.estudiante} · ${item.estado}`,
          })),
      },
      {
        name: 'claseId',
        label: 'Clase',
        type: 'select',
        options: (resources) => resources.clases.map((item) => ({
          value: item.id,
          label: `${item.nombre} · ${item.docenteNombre}`,
        })),
      },
      { name: 'fecha', label: 'Fecha', type: 'date' },
      {
        name: 'estado',
        label: 'Estado',
        type: 'select',
        options: () => ['Presente', 'Tardanza', 'Ausente'].map((value) => ({ value, label: value })),
      },
    ],
    numericFields: ['matriculaId', 'claseId'],
    columns: [
      { key: 'fecha', label: 'Fecha' },
      { key: 'estudiante', label: 'Estudiante', className: 'student-name' },
      { key: 'clase', label: 'Clase' },
      { key: 'estado', label: 'Estado' },
    ],
    emptyMessage: 'Aún no hay asistencias registradas.',
  },
};

async function request(url, options) {
  const response = await fetch(url, options);
  if (!response.ok) {
    const body = await response.json().catch(() => ({}));
    throw new Error(body.error || 'No se pudo completar la operación.');
  }
  return response.status === 204 ? null : response.json();
}

function emptyForm(fields) {
  return Object.fromEntries(fields.map((field) => [
    field.name,
    field.name === 'fecha' ? localDate() : field.name === 'estado' ? 'Presente' : '',
  ]));
}

function localDate() {
  const date = new Date();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${date.getFullYear()}-${month}-${day}`;
}

export default function CatalogosCRUD({ tipo }) {
  const config = catalogos[tipo];
  const [records, setRecords] = useState([]);
  const [resources, setResources] = useState({ docentes: [], matriculas: [], clases: [] });
  const [form, setForm] = useState(() => emptyForm(config.fields));
  const [editingId, setEditingId] = useState(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(true);

  async function load() {
    setLoading(true);
    try {
      const [items, ...related] = await Promise.all([
        request(config.endpoint),
        ...(config.dependencies || []).map((endpoint) => request(endpoint)),
      ]);
      const nextResources = { docentes: [], matriculas: [], clases: [] };
      (config.dependencies || []).forEach((endpoint, index) => {
        nextResources[endpoint.slice('/api/'.length)] = related[index];
      });
      setRecords(items);
      setResources(nextResources);
      setError('');
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
  }, [tipo]);

  async function guardar(event) {
    event.preventDefault();
    const payload = { ...form };
    (config.numericFields || []).forEach((field) => {
      payload[field] = Number(payload[field]);
    });

    try {
      await request(editingId ? `${config.endpoint}/${editingId}` : config.endpoint, {
        method: editingId ? 'PUT' : 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });
      setForm(emptyForm(config.fields));
      setEditingId(null);
      await load();
    } catch (err) {
      setError(err.message);
    }
  }

  function editar(record) {
    setEditingId(record.id);
    setForm(Object.fromEntries(config.fields.map((field) => [
      field.name,
      record[field.name] ?? '',
    ])));
    setError('');
  }

  async function eliminar(id) {
    try {
      await request(`${config.endpoint}/${id}`, { method: 'DELETE' });
      if (editingId === id) {
        setEditingId(null);
        setForm(emptyForm(config.fields));
      }
      await load();
    } catch (err) {
      setError(err.message);
    }
  }

  function cancelar() {
    setEditingId(null);
    setForm(emptyForm(config.fields));
    setError('');
  }

  return (
    <section className="panel">
      <div className="panel-heading">
        <div>
          <p className="eyebrow">{config.eyebrow}</p>
          <h2>{config.title}</h2>
        </div>
        <span className="record-count">{records.length} {records.length === 1 ? 'registro' : 'registros'}</span>
      </div>
      <p className="panel-description">{config.description}</p>

      <form className="entry-form catalog-form" onSubmit={guardar}>
        {config.fields.map((field) => {
          const options = field.options ? field.options(resources) : [];
          const hasOptions = field.type !== 'select' || options.length > 0;
          return (
            <label key={field.name}>
              {field.label}
              {field.type === 'select' ? (
                <select
                  name={field.name}
                  value={form[field.name]}
                  onChange={(event) => setForm({ ...form, [field.name]: event.target.value })}
                  required
                  disabled={!hasOptions}
                >
                  <option value="" disabled>
                    {hasOptions ? 'Selecciona una opción' : 'No hay opciones disponibles'}
                  </option>
                  {options.map((option) => (
                    <option key={option.value} value={option.value}>{option.label}</option>
                  ))}
                </select>
              ) : (
                <input
                  name={field.name}
                  type={field.type || 'text'}
                  value={form[field.name]}
                  onChange={(event) => setForm({ ...form, [field.name]: event.target.value })}
                  placeholder={field.placeholder}
                  required
                />
              )}
            </label>
          );
        })}
        <div className="form-actions">
          <button className="primary-button" type="submit">
            {editingId ? 'Guardar cambios' : `Crear ${config.title.toLowerCase().replace(/s$/, '')}`}
          </button>
          {editingId && <button type="button" onClick={cancelar}>Cancelar</button>}
        </div>
      </form>

      {error && <p className="error-message" role="alert">{error}</p>}
      <div className="table-scroll">
        <table>
          <thead>
            <tr>
              {config.columns.map((column) => <th key={column.key}>{column.label}</th>)}
              <th>Acciones</th>
            </tr>
          </thead>
          <tbody>
            {records.map((record) => (
              <tr key={record.id}>
                {config.columns.map((column) => (
                  <td key={column.key} className={column.className}>{record[column.key]}</td>
                ))}
                <td className="action-cell">
                  <button type="button" onClick={() => editar(record)}>Editar</button>
                  <button className="danger-button" type="button" onClick={() => eliminar(record.id)}>Eliminar</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {!loading && records.length === 0 && <p className="empty-state">{config.emptyMessage}</p>}
        {loading && <p className="empty-state">Cargando registros...</p>}
      </div>
    </section>
  );
}
