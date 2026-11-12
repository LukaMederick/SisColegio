// PATRÓN STATE (Frontend): Matriz de transiciones permitidas según el estado actual
export const Estados = {
  Pendiente: ['Activa', 'Anulada'], // De Pendiente solo pasa a Activa o Anulada
  Activa: ['Finalizada', 'Anulada'], // De Activa solo pasa a Finalizada o Anulada
  Finalizada: [], // Estado terminal
  Anulada: []     // Estado terminal
};

// Valida si el cambio solicitado es legal
export const puedeCambiar = (act, n) => Estados[act]?.includes(n);
