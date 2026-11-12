import { useState } from 'react';
import MatriculasCRUD from './components/MatriculasCRUD';
import NotasCRUD from './components/NotasCRUD';

export default function App() {
  const [tab, setTab] = useState('matriculas');
  return (
    <div style={{ maxWidth: 800, margin: '20px auto', padding: '0 16px' }}>
      <h1>Sistema Académico</h1>
      <div style={{ display: 'flex', gap: 10, marginBottom: 20 }}>
        <button onClick={() => setTab('matriculas')} style={{ fontWeight: tab === 'matriculas' ? 700 : 400 }}>1. Matrículas</button>
        <button onClick={() => setTab('notas')} style={{ fontWeight: tab === 'notas' ? 700 : 400 }}>2. Notas</button>
      </div>
      {tab === 'matriculas' ? <MatriculasCRUD /> : <NotasCRUD />}
    </div>
  );
}
