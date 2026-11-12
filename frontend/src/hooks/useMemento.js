import { useState } from 'react';

// Patrón Memento: permite guardar instantáneas y deshacer
export default function useMemento(initial = []) {
  const [data, setData] = useState(initial);
  const [history, setHistory] = useState([]);

  const setSnapshot = (val) => {
    setHistory((h) => [...h, data]);
    setData(val);
  };

  const undo = () => {
    if (!history.length) return;
    setData(history[history.length - 1]);
    setHistory((h) => h.slice(0, -1));
  };

  return { data, setData: setSnapshot, undo, canUndo: history.length > 0 };
}
