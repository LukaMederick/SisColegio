import { useState } from 'react';
import MatriculasCRUD from './components/MatriculasCRUD';
import NotasCRUD from './components/NotasCRUD';

export default function App() {
  const [tab, setTab] = useState('matriculas');

  return (
    <div style={{ maxWidth: 850, margin: '20px auto', padding: '0 16px' }}>
      <h1>Sistema Académico</h1>
      <div style={{ display: 'flex', gap: 10, marginBottom: 20 }}>
        <button
          onClick={() => setTab('matriculas')}
          style={{ fontWeight: tab === 'matriculas' ? 'bold' : 'normal', background: tab === 'matriculas' ? '#4263eb' : '', color: tab === 'matriculas' ? '#fff' : '' }}
        >
          1. Matrículas (State)
        </button>
        <button
          onClick={() => setTab('notas')}
          style={{ fontWeight: tab === 'notas' ? 'bold' : 'normal', background: tab === 'notas' ? '#4263eb' : '', color: tab === 'notas' ? '#fff' : '' }}
        >
          2. Notas (Memento)
        </button>
      </div>

      {tab === 'matriculas' ? <MatriculasCRUD /> : <NotasCRUD />}
    </div>
  );
}
