import { useState, useEffect } from 'react';
import { Estados, puedeCambiar } from '../patterns/MatriculaState';

export default function MatriculasCRUD() {
  const [list, setList] = useState([]), [form, setForm] = useState({ estudiante: '', curso: '' });
  const load = () => fetch('/api/matriculas').then(r => r.json()).then(setList);
  useEffect(() => { load(); }, []);

  const crear = async (e) => {
    e.preventDefault();
    await fetch('/api/matriculas', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(form) });
    setForm({ estudiante: '', curso: '' }); load();
  };
  const cambiar = (id, act, n) => puedeCambiar(act, n) && fetch(`/api/matriculas/${id}`, { method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ estado: n }) }).then(load);
  const eliminar = (id) => fetch(`/api/matriculas/${id}`, { method: 'DELETE' }).then(load);

  return (
    <div>
      <h2>Matrículas</h2>
      <form onSubmit={crear} style={{ display: 'flex', gap: 8, marginBottom: 15 }}>
        <input placeholder="Estudiante" value={form.estudiante} onChange={e => setForm({ ...form, estudiante: e.target.value })} required />
        <input placeholder="Curso" value={form.curso} onChange={e => setForm({ ...form, curso: e.target.value })} required />
        <button type="submit">+ Matricular</button>
      </form>
      <table width="100%">
        <thead><tr><th>Estudiante</th><th>Curso</th><th>Estado</th><th>Acciones</th></tr></thead>
        <tbody>
          {list.map(m => (
            <tr key={m.id}>
              <td>{m.estudiante}</td><td>{m.curso}</td><td><b>{m.estado}</b></td>
              <td>
                {Estados[m.estado]?.map(e => <button key={e} onClick={() => cambiar(m.id, m.estado, e)} style={{ marginRight: 4 }}>→ {e}</button>)}
                <button onClick={() => eliminar(m.id)}>Eliminar</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
