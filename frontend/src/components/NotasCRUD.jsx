import { useState, useEffect } from 'react';

export default function NotasCRUD() {
  const [list, setList] = useState([]);
  const [form, setForm] = useState({ estudiante: '', curso: '', calificacion: '' });
  const load = () => fetch('/api/notas').then(r => r.json()).then(setList);
  useEffect(() => { load(); }, []);

  const crear = async (e) => {
    e.preventDefault();
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
        <input placeholder="Estudiante" value={form.estudiante} onChange={e => setForm({ ...form, estudiante: e.target.value })} required />
        <input placeholder="Curso" value={form.curso} onChange={e => setForm({ ...form, curso: e.target.value })} required />
        <input placeholder="Nota (0-20)" type="number" step="0.1" value={form.calificacion} onChange={e => setForm({ ...form, calificacion: e.target.value })} required />
        <button type="submit">+ Guardar</button>
      </form>
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
