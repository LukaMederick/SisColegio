import { useState, useEffect } from 'react';

export default function NotasCRUD() {
  const [list, setList] = useState([]);
  const [matriculas, setMatriculas] = useState([]);
  const [form, setForm] = useState({ estudiante: '', curso: '', calificacion: '' });

  const load = () => {
    fetch('/api/notas').then(r => r.json()).then(setList);
    fetch('/api/matriculas').then(r => r.json()).then(setMatriculas);
  };
  useEffect(() => { load(); }, []);

  // Alumnos y cursos dinámicos según las matrículas creadas por el usuario
  const alumnos = [...new Set(matriculas.map(m => m.estudiante))];
  const cursos = form.estudiante
    ? [...new Set(matriculas.filter(m => m.estudiante === form.estudiante).map(m => m.curso))]
    : [...new Set(matriculas.map(m => m.curso))];

  const crear = async (e) => {
    e.preventDefault();
    if (!form.estudiante || !form.curso) return alert('Selecciona estudiante y curso');
    await fetch('/api/notas', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ ...form, calificacion: parseFloat(form.calificacion) || 0 }) });
    setForm({ estudiante: '', curso: '', calificacion: '' });
    load();
  };

  const editar = async (n) => {
    const val = prompt('Nueva calificación (se guardará Memento para deshacer):', n.calificacion);
    if (val === null || val === '') return;
    await fetch(`/api/notas/${n.id}`, { method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ ...n, calificacion: parseFloat(val) }) });
    load();
  };

  const undo = (id) => fetch(`/api/notas/${id}/undo`, { method: 'POST' }).then(load);
  const eliminar = (id) => fetch(`/api/notas/${id}`, { method: 'DELETE' }).then(load);

  return (
    <div>
      <h2>Notas (Patrón Memento)</h2>
      <form onSubmit={crear} style={{ display: 'flex', gap: 8, marginBottom: 15 }}>
        <select value={form.estudiante} onChange={e => setForm({ ...form, estudiante: e.target.value, curso: '' })} required>
          <option value="">-- Estudiante Matriculado --</option>
          {alumnos.map(a => <option key={a} value={a}>{a}</option>)}
        </select>
        <select value={form.curso} onChange={e => setForm({ ...form, curso: e.target.value })} required>
          <option value="">-- Curso --</option>
          {cursos.map(c => <option key={c} value={c}>{c}</option>)}
        </select>
        <input placeholder="Nota (0-20)" type="number" step="0.1" min="0" max="20" value={form.calificacion} onChange={e => setForm({ ...form, calificacion: e.target.value })} required style={{ maxWidth: 120 }} />
        <button type="submit" style={{ whiteSpace: 'nowrap' }}>+ Guardar</button>
      </form>
      {alumnos.length === 0 && <p style={{ color: '#888' }}>* No hay alumnos aún. Agrega primero en "1. Matrículas".</p>}
      <table width="100%">
        <thead><tr><th>Estudiante</th><th>Curso</th><th>Nota</th><th>Acciones</th></tr></thead>
        <tbody>
          {list.map(n => (
            <tr key={n.id}>
              <td>{n.estudiante}</td><td>{n.curso}</td><td><b>{n.calificacion}</b></td>
              <td>
                <button onClick={() => editar(n)} style={{ marginRight: 4 }}>Editar</button>
                <button onClick={() => undo(n.id)} style={{ marginRight: 4 }}>↩ Deshacer (Memento)</button>
                <button onClick={() => eliminar(n.id)}>Eliminar</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
