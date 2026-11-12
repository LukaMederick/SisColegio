import { useState, useEffect } from 'react';
import { Estados, puedeCambiar } from '../patterns/MatriculaState';

export default function MatriculasCRUD() {
  const [list, setList] = useState([]);
  const [form, setForm] = useState({ estudiante: '', curso: '' });
  const load = () => fetch('/api/matriculas').then(r => r.json()).then(setList);
  useEffect(() => { load(); }, []);

  const crear = async (e) => {
    e.preventDefault();
    if (!form.estudiante || !form.curso) return alert('Ingresa estudiante y curso');
    await fetch('/api/matriculas', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(form) });
    setForm({ estudiante: '', curso: '' });
    load();
  };

  const transicionar = async (id, actual, nuevo) => {
    if (!puedeCambiar(actual, nuevo)) return alert('Transición inválida');
    await fetch(`/api/matriculas/${id}`, { method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ estado: nuevo }) });
    load();
  };

  const eliminar = (id) => fetch(`/api/matriculas/${id}`, { method: 'DELETE' }).then(load);

  return (
    <div>
      <h2>Matrículas (Patrón State)</h2>
      <form onSubmit={crear} style={{ display: 'flex', gap: 8, marginBottom: 15 }}>
        <input placeholder="Nombre del Estudiante" value={form.estudiante} onChange={e => setForm({ ...form, estudiante: e.target.value })} required />
        <input placeholder="Nombre del Curso" value={form.curso} onChange={e => setForm({ ...form, curso: e.target.value })} required />
        <button type="submit" style={{ whiteSpace: 'nowrap' }}>+ Matricular</button>
      </form>
      <table width="100%">
        <thead><tr><th>Estudiante</th><th>Curso</th><th>Estado</th><th>Transiciones (State)</th><th>Acción</th></tr></thead>
        <tbody>
          {list.map(m => (
            <tr key={m.id}>
              <td>{m.estudiante}</td><td>{m.curso}</td><td><b>{m.estado}</b></td>
              <td>{Estados[m.estado]?.map(e => (
                <button key={e} onClick={() => transicionar(m.id, m.estado, e)} style={{ marginRight: 4 }}>→ {e}</button>
              ))}</td>
              <td><button onClick={() => eliminar(m.id)}>Eliminar</button></td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
