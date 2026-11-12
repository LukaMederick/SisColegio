export const Estados = {
  Pendiente: ['Activa', 'Anulada'],
  Activa: ['Finalizada', 'Anulada'],
  Finalizada: [], Anulada: []
};
export const puedeCambiar = (act, n) => Estados[act]?.includes(n);
