import { useCallback, useEffect, useState } from 'react';

// El historial guarda mementos inmutables del valor anterior a cada edición.
export default function useMemento(initialValue) {
  const [value, setValue] = useState(initialValue);
  const [history, setHistory] = useState([]);

  useEffect(() => {
    setValue(initialValue);
    setHistory([]);
  }, [initialValue]);

  const update = useCallback((nextValue) => {
    setHistory((previous) => [...previous, value]);
    setValue(nextValue);
  }, [value]);

  const undo = useCallback(() => {
    if (history.length === 0) return value;
    const previousValue = history[history.length - 1];
    setHistory((previous) => previous.slice(0, -1));
    setValue(previousValue);
    return previousValue;
  }, [history, value]);

  const commit = useCallback((savedValue) => {
    setValue(savedValue);
    setHistory([]);
  }, []);

  return { value, update, undo, commit, canUndo: history.length > 0 };
}
