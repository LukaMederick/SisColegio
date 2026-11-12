import { useState } from 'react';
import MatriculasCRUD from './components/MatriculasCRUD.jsx';
import NotasCRUD from './components/NotasCRUD.jsx';
import CatalogosCRUD from './components/CatalogosCRUD.jsx';

const modulos = [
  { id: 'docentes', label: '1. Docentes' },
  { id: 'clases', label: '2. Clases' },
  { id: 'matriculas', label: '3. Matrículas' },
  { id: 'asistencias', label: '4. Asistencias' },
  { id: 'notas', label: '5. Notas' },
];

export default function App() {
  const [modulo, setModulo] = useState('docentes');

  function contenido() {
    if (modulo === 'matriculas') return <MatriculasCRUD />;
    if (modulo === 'notas') return <NotasCRUD />;
    return <CatalogosCRUD key={modulo} tipo={modulo} />;
  }

  return (
    <div className="app-shell">
      <header className="page-header">
        <div className="brand-mark" aria-hidden="true">GE</div>
        <div>
          <p className="eyebrow">PANEL ACADÉMICO</p>
          <h1>Gestión educativa</h1>
          <p className="page-description">Matrículas, notas, asistencias, clases y docentes en un solo lugar.</p>
        </div>
      </header>
      <nav className="section-nav" aria-label="Secciones de gestión">
        {modulos.map((item) => (
          <button
            key={item.id}
            className={modulo === item.id ? 'nav-button active' : 'nav-button'}
            type="button"
            aria-current={modulo === item.id ? 'page' : undefined}
            onClick={() => setModulo(item.id)}
          >
            {item.label}
          </button>
        ))}
      </nav>
      <main className="content-grid">
        {contenido()}
      </main>
      <footer className="page-footer">Sistema de gestión educativa</footer>
    </div>
  );
}
